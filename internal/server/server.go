package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
	"github.com/PitiNarak/condormhub-backend/internal/middlewares"
	"github.com/PitiNarak/condormhub-backend/internal/repositories"
	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

type Config struct {
	Name                 string `env:"NAME"`
	Port                 int    `env:"PORT"`
	Env                  string `env:"ENV"`
	MaxBodyLimitMB       int    `env:"MAX_BODY_LIMIT_MB"`
	CorsAllowOrigins     string `env:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods     string `env:"CORS_ALLOW_METHODS"`
	CorsAllowHeaders     string `env:"CORS_ALLOW_HEADERS"`
	CorsAllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS"`
}

type Server struct {
	app             *fiber.App
	config          Config
	greetingHandler *handlers.GreetingHandler
	userHandler     ports.UserHandler
	// orderHandler      ports.TransactionHandler
	testUploadHandler *handlers.TestUploadHandler
	storage           *storage.Storage
	jwtUtils          *utils.JWTUtils
	authMiddleware    *middlewares.AuthMiddleware
	dormHandler       ports.DormHandler
}

func NewServer(config Config, smtpConfig services.SMTPConfig, jwtConfig utils.JWTConfig, storageConfig storage.Config, stripeConfig stripe.Config, db *gorm.DB) *Server {

	app := fiber.New(fiber.Config{
		AppName:       config.Name,
		BodyLimit:     config.MaxBodyLimitMB * 1024 * 1024,
		CaseSensitive: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var e *error_handler.ErrorHandler
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			} else {
				message = err.Error()
			}

			if e != nil && e.Err != nil {
				log.Printf("Error: %v, Code: %d, Message: %s", e.Error(), code, message)
			} else {
				log.Printf("Error: %s, Code: %d, Message: %s", err.Error(), code, message)
			}

			return c.Status(code).JSON(&http_response.HttpResponse{
				Success: false,
				Message: message,
				Data:    nil,
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.CorsAllowOrigins,
		AllowMethods:     config.CorsAllowMethods,
		AllowHeaders:     config.CorsAllowHeaders,
		AllowCredentials: config.CorsAllowCredentials,
	}))

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		DisableColors: true,
	}))

	jwtUtils := utils.NewJWTUtils(&jwtConfig)
	storage := storage.NewStorage(storageConfig)

	userRepository := repositories.NewUserRepo(db)

	emailService := services.NewEmailService(&smtpConfig, jwtUtils)
	userService := services.NewUserService(userRepository, emailService, jwtUtils)
	userHandler := handlers.NewUserHandler(userService)
	testUploadHandler := handlers.NewTestUploadHandler(storage)

	// stripe := stripe.New(stripeConfig)
	dormRepository := repositories.NewDormRepository(db)
	dormService := services.NewDormService(dormRepository)
	dormHandler := handlers.NewDormHandler(dormService)

	// tsxRepo := repositories.NewTransactionRepository(db)
	// tsxService := services.NewTransactionService(tsxRepo, dormRepository, stripe)
	// tsxHandler := handlers.NewTransactionHandler(tsxService, &stripeConfig)

	authMiddleware := middlewares.NewAuthMiddleware(jwtUtils, userRepository)
	return &Server{
		app:             app,
		greetingHandler: handlers.NewGreetingHandler(),
		userHandler:     userHandler,
		// orderHandler:      tsxHandler,
		config:            config,
		testUploadHandler: testUploadHandler,
		storage:           storage,
		jwtUtils:          jwtUtils,
		authMiddleware:    authMiddleware,
		dormHandler:       dormHandler,
	}
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc, jwtConfig utils.JWTConfig) {

	// init routes
	s.initRoutes()

	// start server
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Panicf("Failed to start server: %v\n", err)
			stop()
		}
	}()

	// shutdown server at the end
	defer func() {
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Failed to shutdown server: %v\n", err)
		}
		log.Println("Server stopped")
	}()

	<-ctx.Done()

	log.Println("Server is shutting down...")
}

func (s *Server) initRoutes() {
	// greeting
	s.app.Get("/", s.greetingHandler.Greeting)

	// swagger
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	// upload file example
	s.app.Post("/upload/public", s.testUploadHandler.UploadToPublicBucketHandler)
	s.app.Post("/upload/private", s.testUploadHandler.UploadToPrivateBucketHandler)
	s.app.Get("/signedurl/*", s.testUploadHandler.GetSignedUrlHandler)

	// user
	userRoutes := s.app.Group("/user")

	userRoutes.Get("/me", s.authMiddleware.Auth, s.userHandler.GetUserInfo)

	userRoutes.Post("/verify", s.userHandler.VerifyEmail)
	userRoutes.Post("/resetpassword", s.userHandler.ResetPasswordCreate)
	userRoutes.Post("/newpassword", s.authMiddleware.Auth, s.userHandler.ResetPassword)
	userRoutes.Patch("/", s.authMiddleware.Auth, s.userHandler.UpdateUserInformation)
	userRoutes.Delete("/", s.authMiddleware.Auth, s.userHandler.DeleteAccount)

	authRoutes := s.app.Group("/auth")
	authRoutes.Post("/register", s.userHandler.Register)
	authRoutes.Post("/login", s.userHandler.Login)

	// order
	// orderRoutes := s.app.Group("/order")
	// orderRoutes.Post("/", s.authMiddleware.Auth, s.orderHandler.CreateOrder)
	// orderRoutes.Post("/webhook", s.orderHandler.Webhook)
	// dorm
	dormRoutes := s.app.Group("/dorms")
	dormRoutes.Post("/", s.authMiddleware.Auth, s.dormHandler.Create)
	dormRoutes.Get("/", s.dormHandler.GetAll)
	dormRoutes.Get("/:id", s.dormHandler.GetByID)
	dormRoutes.Patch("/:id", s.authMiddleware.Auth, s.dormHandler.Update)
	dormRoutes.Delete("/:id", s.authMiddleware.Auth, s.dormHandler.Delete)
}
