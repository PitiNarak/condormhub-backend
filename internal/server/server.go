package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/internal/middlewares"
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
	app            *fiber.App
	config         Config
	storage        *storage.Storage
	jwtUtils       *utils.JWTUtils
	authMiddleware *middlewares.AuthMiddleware
	db             *gorm.DB
	smtpConfig     *services.SMTPConfig
	handler        *handler
	service        *service
	repository     *repository
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

	jwtUtils := utils.NewJWTUtils(&jwtConfig)
	storage := storage.NewStorage(storageConfig)

	return &Server{
		app:        app,
		config:     config,
		storage:    storage,
		jwtUtils:   jwtUtils,
		db:         db,
		smtpConfig: &smtpConfig,
	}
}

func (s *Server) Start(ctx context.Context, stop context.CancelFunc) {

	s.initServerMiddleware()
	s.initRepository()
	s.initAuthMiddleware()

	s.initService()
	s.initHandler()
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

func (s *Server) initServerMiddleware() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.CorsAllowOrigins,
		AllowMethods:     s.config.CorsAllowMethods,
		AllowHeaders:     s.config.CorsAllowHeaders,
		AllowCredentials: s.config.CorsAllowCredentials,
	}))

	s.app.Use(requestid.New())
	s.app.Use(logger.New(logger.Config{
		DisableColors: true,
	}))

}

func (s *Server) initAuthMiddleware() {
	s.authMiddleware = middlewares.NewAuthMiddleware(s.jwtUtils, s.repository.user)
}
