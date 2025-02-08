package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internals/core/services"
	"github.com/PitiNarak/condormhub-backend/internals/core/utils"
	"github.com/PitiNarak/condormhub-backend/internals/handlers"
	"github.com/PitiNarak/condormhub-backend/internals/repositories"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
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
	app              *fiber.App
	config           Config
	greetingHandler  *handlers.GreetingHandler
	sampleLogHandler *handlers.SampleLogHandler
	userHandler      *handlers.UserHandler
}

func NewServer(config Config, smtpConfig services.SMTPConfig, jwtConfig utils.JWTConfig, db *gorm.DB) *Server {

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
			}

			log.Printf("Error: %v, Code: %d, Message: %s", e.Error(), code, message)

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

	sampleLogRepository := repositories.NewSampleLogRepository(db)
	userRepository := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepository)
	emailService := services.NewEmailService(&smtpConfig, &jwtConfig)
	userHandler := handlers.NewUserHandler(userService, emailService, &jwtConfig)

	return &Server{
		app:              app,
		greetingHandler:  handlers.NewGreetingHandler(),
		sampleLogHandler: handlers.NewSampleLogHandler(sampleLogRepository),
		userHandler:      userHandler,
		config:           config,
	}
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {
	sampleLogRoutes := s.app.Group("/log")
	sampleLogRoutes.Get("/", s.sampleLogHandler.GetAll)
	sampleLogRoutes.Post("/", s.sampleLogHandler.Save)
	sampleLogRoutes.Delete("/:id", s.sampleLogHandler.Delete)
	sampleLogRoutes.Patch("/:id", s.sampleLogHandler.EditMessage)

	s.app.All("/", s.greetingHandler.Greeting)

	s.app.Post("/register", s.userHandler.Create)
	s.app.Get("/verify/:token", s.userHandler.VerifyEmail)
	s.app.Get("/resetpassword", s.userHandler.ResetPasswordCreate)
	s.app.Post("/resetpasswordresponse", s.userHandler.ResetPasswordResponse)

	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			log.Panicf("Failed to start server: %v\n", err)
			stop()
		}
	}()

	defer func() {
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Failed to shutdown server: %v\n", err)
		}
		log.Println("Server stopped")
	}()

	<-ctx.Done()

	log.Println("Server is shutting down...")
}
