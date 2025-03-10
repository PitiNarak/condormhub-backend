package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(UserService ports.UserService) ports.UserHandler {
	return &UserHandler{userService: UserService}
}
