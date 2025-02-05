package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	Create(user domain.User) error
}

type UserService interface {
	Create(user domain.User) error
}

type UserHandler interface {
	Create(c *fiber.Ctx) error
}
