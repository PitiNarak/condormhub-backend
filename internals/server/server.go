package server

import (
	"log"
	"os"

	"github.com/PitiNarak/condormhub-backend/internals/core/services"
	"github.com/PitiNarak/condormhub-backend/internals/handlers"
	"github.com/PitiNarak/condormhub-backend/internals/repositories"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"gorm.io/gorm"
)

type Server struct {
	app              *fiber.App
	greetingHandler  *handlers.GreetingHandler
	sampleLogHandler *handlers.SampleLogHandler
	userHandler      *handlers.UserHandler
}

func NewServer(db *gorm.DB) *Server {
	sampleLogRepository := repositories.NewSampleLogRepository(db)
	userRepository := repositories.NewUserRepo(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	return &Server{
		app:              fiber.New(),
		greetingHandler:  handlers.NewGreetingHandler(),
		sampleLogHandler: handlers.NewSampleLogHandler(sampleLogRepository),
		userHandler:      userHandler,
	}
}

func (s *Server) Start(port string) {

	sampleLogRoutes := s.app.Group("/log")
	sampleLogRoutes.Get("/", s.sampleLogHandler.GetAll)
	sampleLogRoutes.Post("/", s.sampleLogHandler.Save)
	sampleLogRoutes.Delete("/:id", s.sampleLogHandler.Delete)
	sampleLogRoutes.Patch("/:id", s.sampleLogHandler.EditMessage)

	userRoutes := s.app.Group("/user")
	userRoutes.Post("/register", s.userHandler.Create)

	userRoutes.Post("/login", s.userHandler.Login)
	s.app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	userRoutes.Put("/update", s.userHandler.Update)

	s.app.All("/", s.greetingHandler.Greeting)

	err := s.app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Listen Failed: %v", err)
	}
}
