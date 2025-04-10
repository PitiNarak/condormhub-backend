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
	contract       ports.ContractService
	leasingRequest ports.LeasingRequestService
	receipt        ports.ReceiptService
	support        ports.SupportService
}

func (s *Server) initService() {
	email := email.NewEmailService(s.smtpConfig, s.jwtUtils)
	user := services.NewUserService(s.repository.user, email, s.jwtUtils, s.storage)
	dorm := services.NewDormService(s.repository.dorm, s.storage)
	leasingHistory := services.NewLeasingHistoryService(s.repository.leasingHistory, s.repository.dorm)
	order := services.NewOrderService(s.repository.order, s.repository.leasingHistory)
	ownershipProof := services.NewOwnershipProofService(s.repository.ownershipProof, s.repository.user, s.storage)
	contract := services.NewContractService(s.repository.contract, s.repository.user, s.repository.dorm, leasingHistory)
	leasingRequest := services.NewLeasingRequestService(s.repository.leasingRequest, s.repository.dorm, s.repository.contract)
	receipt := services.NewReceiptService(s.repository.receipt, s.repository.user, s.repository.tsx, s.repository.order, s.repository.leasingHistory, s.repository.dorm, s.storage)
	tsx := services.NewTransactionService(s.repository.tsx, s.repository.order, s.stripe, s.repository.leasingHistory, receipt)
	support := services.NewSupportService(s.repository.support)

	s.service = &service{
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
		ownershipProof: ownershipProof,
		contract:       contract,
		leasingRequest: leasingRequest,
		receipt:        receipt,
		support:        support,
	}
}
