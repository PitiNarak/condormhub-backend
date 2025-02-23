package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/repositories"
)

type repository struct {
	user ports.UserRepository
	dorm ports.DormRepository
}

func (s *Server) initRepository() {
	user := repositories.NewUserRepo(s.db)
	dorm := repositories.NewDormRepository(s.db)
	s.repository = &repository{
		user: user,
		dorm: dorm,
	}
}
