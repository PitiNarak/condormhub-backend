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
}

func (s *Server) initService() {
	email := services.NewEmailService(s.smtpConfig, s.jwtUtils)
	user := services.NewUserService(s.repository.user, email, s.jwtUtils)
	dorm := services.NewDormService(s.repository.dorm)
	leasingHistory := services.NewLeasingHistoryService(s.repository.leasingHistory, s.repository.dorm)

	s.service = &service{
		email:          email,
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
	}
}
