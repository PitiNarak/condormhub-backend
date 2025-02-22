package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/internal/middlewares"
	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
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
<<<<<<< HEAD
	app                   *fiber.App
	config                Config
	greetingHandler       *handlers.GreetingHandler
	sampleLogHandler      *handlers.SampleLogHandler
	userHandler           ports.UserHandler
	testUploadHandler     *handlers.TestUploadHandler
	storage               *storage.Storage
	jwtUtils              *utils.JWTUtils
	authMiddleware        *middlewares.AuthMiddleware
	dormHandler           ports.DormHandler
	leasingHistoryHandler ports.LeasingHistoryHandler
=======
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
>>>>>>> dev
}

func NewServer(config Config, smtpConfig services.SMTPConfig, jwtConfig utils.JWTConfig, storageConfig storage.Config, db *gorm.DB) *Server {

	app := fiber.New(fiber.Config{
		AppName:               config.Name,
		BodyLimit:             config.MaxBodyLimitMB * 1024 * 1024,
		CaseSensitive:         true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler.Handler,
	})

	jwtUtils := utils.NewJWTUtils(&jwtConfig)
	storage := storage.NewStorage(storageConfig)

<<<<<<< HEAD
	sampleLogRepository := repositories.NewSampleLogRepository(db)
	userRepository := repositories.NewUserRepo(db)

	emailService := services.NewEmailService(&smtpConfig, jwtUtils)
	userService := services.NewUserService(userRepository, emailService, jwtUtils)
	userHandler := handlers.NewUserHandler(userService)
	testUploadHandler := handlers.NewTestUploadHandler(storage)

	dormRepository := repositories.NewDormRepository(db)
	dormService := services.NewDormService(dormRepository)
	dormHandler := handlers.NewDormHandler(dormService)

	leasingHistoryRepository := repositories.NewLeasingHistoryRepository(db)
	leasingHistoryService := services.NewLeasingHistoryService(leasingHistoryRepository)
	leasingHistoryHandler := handlers.NewLeasingHistoryHandler(leasingHistoryService)

	authMiddleware := middlewares.NewAuthMiddleware(jwtUtils, userRepository)
	return &Server{
		app:                   app,
		greetingHandler:       handlers.NewGreetingHandler(),
		sampleLogHandler:      handlers.NewSampleLogHandler(sampleLogRepository),
		userHandler:           userHandler,
		config:                config,
		testUploadHandler:     testUploadHandler,
		storage:               storage,
		jwtUtils:              jwtUtils,
		authMiddleware:        authMiddleware,
		dormHandler:           dormHandler,
		leasingHistoryHandler: leasingHistoryHandler,
=======
	return &Server{
		app:        app,
		config:     config,
		storage:    storage,
		jwtUtils:   jwtUtils,
		db:         db,
		smtpConfig: &smtpConfig,
>>>>>>> dev
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
