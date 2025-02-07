package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUser(userID uuid.UUID) (domain.User, error)
	UpdateUser(user domain.User) error
	GetUserByEmail(email string) (domain.User, error)
}

type UserService interface {
	Create(user *domain.User) (*domain.User, error)
	VerifyUser(userID uuid.UUID) error
	ResetPasswordCreate(email string) (domain.User, error)
	ResetPasswordResponse(userID uuid.UUID, password string) error
}

type UserHandler interface {
	Create(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	ResetPasswordResponse(c *fiber.Ctx) error
}
