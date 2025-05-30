package server

import (
	"context"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/middleware"
	"github.com/PitiNarak/condormhub-backend/pkg/email"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/PitiNarak/condormhub-backend/pkg/redis"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/yokeTH/go-pkg/apperror"
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
	jwtUtils       *jwt.JWTUtils
	authMiddleware *middleware.AuthMiddleware
	redis          *redis.Redis
	db             *database.Database
	smtpConfig     *email.SMTPConfig
	stripeConfig   *stripe.Config
	stripe         *stripe.Stripe
	handler        *handler
	service        *service
	repository     *repository
}

func NewServer(config Config, smtpConfig email.SMTPConfig, jwtConfig jwt.JWTConfig, storageConfig storage.Config, stripeConfig stripe.Config, redis *redis.Redis, db *database.Database) *Server {

	app := fiber.New(fiber.Config{
		AppName:               config.Name,
		BodyLimit:             config.MaxBodyLimitMB * 1024 * 1024,
		CaseSensitive:         true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		ErrorHandler:          apperror.ErrorHandler,
	})

	jwtUtils := jwt.NewJWTUtils(&jwtConfig, redis)
	storage := storage.NewStorage(storageConfig)
	stripe := stripe.New(stripeConfig)

	return &Server{
		app:          app,
		config:       config,
		storage:      storage,
		jwtUtils:     jwtUtils,
		db:           db,
		redis:        redis,
		smtpConfig:   &smtpConfig,
		stripeConfig: &stripeConfig,
		stripe:       stripe,
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
	s.authMiddleware = middleware.NewAuthMiddleware(s.jwtUtils, s.repository.user)
}
