package server

import (
	"log"

	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/PitiNarak/condormhub-backend/internals/core/services"
	"github.com/PitiNarak/condormhub-backend/internals/handlers"
	"github.com/PitiNarak/condormhub-backend/internals/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	app              *fiber.App
	greetingHandler  *handlers.GreetingHandler
	sampleLogHandler *handlers.SampleLogHandler
	userHandler      *handlers.UserHandler
	config           *config.AppConfig
}

func NewServer(db *gorm.DB) *Server {
	config := config.Load()

	sampleLogRepository := repositories.NewSampleLogRepository(db)
	userRepository := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepository)
	emailService := services.NewEmailService(&config.SMTP, &config.JWT)
	userHandler := handlers.NewUserHandler(userService, emailService, config)

	return &Server{
		app:              fiber.New(),
		greetingHandler:  handlers.NewGreetingHandler(),
		sampleLogHandler: handlers.NewSampleLogHandler(sampleLogRepository),
		userHandler:      userHandler,
		config:           config,
	}
}

func (s *Server) Start(port string) {

	sampleLogRoutes := s.app.Group("/log")
	sampleLogRoutes.Get("/", s.sampleLogHandler.GetAll)
	sampleLogRoutes.Post("/", s.sampleLogHandler.Save)
	sampleLogRoutes.Delete("/:id", s.sampleLogHandler.Delete)
	sampleLogRoutes.Patch("/:id", s.sampleLogHandler.EditMessage)

	s.app.All("/", s.greetingHandler.Greeting)

	s.app.Post("/register", s.userHandler.Create)
	s.app.Get("/verify/:token", s.userHandler.VerifyEmail)
	s.app.Get("/resetPassword", s.userHandler.ResetPasswordCreate)
	s.app.Get("/resetPasswordRespond", s.userHandler.ResetPasswordRespond)

	err := s.app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Listen Failed: %v", err)
	}
}
