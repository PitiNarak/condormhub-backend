package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	repository1 "github.com/PitiNarak/condormhub-backend/internal/repository"
)

type repository struct {
	user           ports.UserRepository
	dorm           ports.DormRepository
	leasingHistory ports.LeasingHistoryRepository
	order          ports.OrderRepository
	tsx            ports.TransactionRepository
	ownershipProof ports.OwnershipProofRepository
	leasingRequest ports.LeasingRequestRepository
}

func (s *Server) initRepository() {
	user := repository1.NewUserRepo(s.db)
	dorm := repository1.NewDormRepository(s.db)
	leasingHistory := repository1.NewLeasingHistoryRepository(s.db)
	order := repository1.NewOrderRepository(s.db)
	tsx := repository1.NewTransactionRepository(s.db)
	ownershipProof := repository1.NewOwnershipProofRepository(s.db)
	leasingRequest := repository1.NewLeasingRequestRepository(s.db)

	s.repository = &repository{
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
		ownershipProof: ownershipProof,
		leasingRequest: leasingRequest,
	}
}
