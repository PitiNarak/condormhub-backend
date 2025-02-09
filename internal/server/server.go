package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
	"github.com/PitiNarak/condormhub-backend/internal/repositories"
	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	app               *fiber.App
	config            Config
	greetingHandler   *handlers.GreetingHandler
	sampleLogHandler  *handlers.SampleLogHandler
	userHandler       *handlers.UserHandler
	testUploadHandler *handlers.TestUploadHandler
	storage           *storage.Storage
}

func NewServer(config Config, smtpConfig services.SMTPConfig, jwtConfig utils.JWTConfig, storageConfig storage.Config, db *gorm.DB) *Server {

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
			}
			log.Printf("Error: %s, Code: %d, Message: %s", err.Error(), code, message)

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

	storage := storage.NewStorage(storageConfig)

	sampleLogRepository := repositories.NewSampleLogRepository(db)
	userRepository := repositories.NewUserRepo(db)

	emailService := services.NewEmailService(&smtpConfig, &jwtConfig)
	userService := services.NewUserService(userRepository, emailService, &jwtConfig)
	userHandler := handlers.NewUserHandler(userService)
	testUploadHandler := handlers.NewTestUploadHandler(storage)

	return &Server{
		app:               app,
		greetingHandler:   handlers.NewGreetingHandler(),
		sampleLogHandler:  handlers.NewSampleLogHandler(sampleLogRepository),
		userHandler:       userHandler,
		config:            config,
		testUploadHandler: testUploadHandler,
		storage:           storage,
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

	// test upload
	s.app.Post("/upload", s.testUploadHandler.UploadHandler)

	// sample log
	sampleLogRoutes := s.app.Group("/log")
	sampleLogRoutes.Get("/", s.sampleLogHandler.GetAll)
	sampleLogRoutes.Post("/", s.sampleLogHandler.Save)
	sampleLogRoutes.Delete("/:id", s.sampleLogHandler.Delete)
	sampleLogRoutes.Patch("/:id", s.sampleLogHandler.EditMessage)

	// user
	userRoutes := s.app.Group("/user")
	userRoutes.Post("/register", s.userHandler.Create)
	userRoutes.Post("/login", s.userHandler.Login)
	userRoutes.Patch("/update", s.userHandler.UpdateUserInformation)
	userRoutes.Get("/verify/", s.userHandler.VerifyEmail)
	userRoutes.Get("/resetpassword", s.userHandler.ResetPasswordCreate)
	userRoutes.Patch("/newpassword", s.userHandler.ResetPasswordResponse)
}
