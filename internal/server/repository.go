package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/repositories"
)

type repository struct {
	user           ports.UserRepository
	dorm           ports.DormRepository
	leasingHistory ports.LeasingHistoryRepository
	order          ports.OrderRepository
	tsx            ports.TransactionRepository
}

func (s *Server) initRepository() {
	user := repositories.NewUserRepo(s.db)
	dorm := repositories.NewDormRepository(s.db)
	leasingHistory := repositories.NewLeasingHistoryRepository(s.db)
	order := repositories.NewOrderRepository(s.db)
	tsx := repositories.NewTransactionRepository(s.db)

	s.repository = &repository{
		user:           user,
		dorm:           dorm,
		leasingHistory: leasingHistory,
		order:          order,
		tsx:            tsx,
	}
}
