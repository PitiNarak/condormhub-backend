package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
)

type repository struct {
	user ports.UserRepository
	dorm ports.DormRepository
}
