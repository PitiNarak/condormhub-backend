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
	contract       ports.ContractRepository
	leasingRequest ports.LeasingRequestRepository
	receipt        ports.ReceiptRepository
	support        ports.SupportRepository
}

func (s *Server) initRepository() {
	user := repository1.NewUserRepo(s.db)
	dorm := repository1.NewDormRepository(s.db)
	leasingHistory := repository1.NewLeasingHistoryRepository(s.db)
	order := repository1.NewOrderRepository(s.db)
	tsx := repository1.NewTransactionRepository(s.db)
	ownershipProof := repository1.NewOwnershipProofRepository(s.db)
	contract := repository1.NewContractRepository(s.db)
	leasingRequest := repository1.NewLeasingRequestRepository(s.db)
	receipt := repository1.NewReceiptRepository(s.db)
	support := repository1.NewSupportRepository(s.db)

	s.repository = &repository{
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
