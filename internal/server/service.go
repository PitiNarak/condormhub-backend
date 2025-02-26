package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/core/services"
)

type service struct {
	email          ports.EmailServicePort
	user           ports.UserService
	dorm           ports.DormService
	leasingHistory ports.LeasingHistoryService
	order          ports.OrderService
	tsx            ports.TransactionService
}

func (s *Server) initService() {
	email := services.NewEmailService(s.smtpConfig, s.jwtUtils)
	user := services.NewUserService(s.repository.user, email, s.jwtUtils)
	dorm := services.NewDormService(s.repository.dorm)
	leasingHistory := services.NewLeasingHistoryService(s.repository.leasingHistory, s.repository.dorm)
	order := services.NewOrderService(s.repository.order, s.repository.leasingHistory)
	tsx := services.NewTransactionService(s.repository.tsx, s.repository.order, s.stripe)

	s.service = &service{
		email:          email,
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
	}
}
