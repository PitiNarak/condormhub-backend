package server

import (
	"log"

	"github.com/PitiNarak/condormhub-backend/internals/handlers"
	"github.com/PitiNarak/condormhub-backend/internals/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	app              *fiber.App
	greetingHandler  *handlers.GreetingHandler
	sampleLogHandler *handlers.SampleLogHandler
}

func NewServer(db *gorm.DB) *Server {
	sampleLogRepository := repositories.NewSampleLogRepository(db)

	return &Server{
		app:              fiber.New(),
		greetingHandler:  handlers.NewGreetingHandler(),
		sampleLogHandler: handlers.NewSampleLogHandler(sampleLogRepository),
	}
}

func (s *Server) Start(port string) {

	sampleLogRoutes := s.app.Group("/log")
	sampleLogRoutes.Get("/", s.sampleLogHandler.GetAll)
	sampleLogRoutes.Post("/", s.sampleLogHandler.Save)
	sampleLogRoutes.Delete("/:id", s.sampleLogHandler.Delete)
	sampleLogRoutes.Patch("/:id", s.sampleLogHandler.EditMessage)

	s.app.All("/", s.greetingHandler.Greeting)

	err := s.app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Listen Failed: %v", err)
	}
}
