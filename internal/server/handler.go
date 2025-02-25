package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
)

type handler struct {
	greeting       *handlers.GreetingHandler
	user           ports.UserHandler
	exampleUpload  *handlers.TestUploadHandler
	dorm           ports.DormHandler
	leasingHistory ports.LeasingHistoryHandler
	order          ports.OrderHandler
}

func (s *Server) initHandler() {
	greeting := handlers.NewGreetingHandler()
	user := handlers.NewUserHandler(s.service.user)
	exampleUpload := handlers.NewTestUploadHandler(s.storage)
	dorm := handlers.NewDormHandler(s.service.dorm)
	leasingHistory := handlers.NewLeasingHistoryHandler(s.service.leasingHistory)
	order := handlers.NewOrderHandler(s.service.order)

	s.handler = &handler{
		greeting:       greeting,
		user:           user,
		exampleUpload:  exampleUpload,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
	}
}
