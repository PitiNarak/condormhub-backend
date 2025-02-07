package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	Create(user domain.User) error
	Update(user domain.User) (domain.User, error)
	GetUserViaEmail(email string) (domain.User, error)
}

type UserService interface {
	Create(user domain.User) error
	Update(user domain.User) (domain.User, error)
	Login(email string, password string) (string, error)
}

type UserHandler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}
