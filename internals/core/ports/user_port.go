package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUser(userID uuid.UUID) (*domain.User, error)
	UpdateUser(user domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}

type UserService interface {
	Create(user *domain.User) error
	Update(user domain.User) error
	Login(email string, password string) (string, error)
	VerifyUser(token string) error
	ResetPasswordCreate(email string) error
	ResetPasswordResponse(token string, password string) error
}

type UserHandler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	ResetPasswordResponse(c *fiber.Ctx) error
}
