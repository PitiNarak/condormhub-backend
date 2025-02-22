package server

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
)

type service struct {
	email ports.EmailServicePort
	user  ports.UserService
	dorm  ports.DormService
}
