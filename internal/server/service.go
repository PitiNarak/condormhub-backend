package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/pkg/email"
)

type service struct {
	user           ports.UserService
	dorm           ports.DormService
	leasingHistory ports.LeasingHistoryService
	order          ports.OrderService
	tsx            ports.TransactionService
	ownershipProof ports.OwnershipProofService
}

func (s *Server) initService() {
	email := email.NewEmailService(s.smtpConfig, s.jwtUtils)
	user := services.NewUserService(s.repository.user, email, s.jwtUtils)
	dorm := services.NewDormService(s.repository.dorm)
	leasingHistory := services.NewLeasingHistoryService(s.repository.leasingHistory, s.repository.dorm)
	order := services.NewOrderService(s.repository.order, s.repository.leasingHistory)
	tsx := services.NewTransactionService(s.repository.tsx, s.repository.order, s.stripe)
	ownershipProof := services.NewOwnershipProofService(s.repository.ownershipProof, s.repository.user)

	s.service = &service{
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
		ownershipProof: ownershipProof,
	}
}
