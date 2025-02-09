package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUser(userID uuid.UUID) (*domain.User, error)
	UpdateUser(user domain.User) error
	Update(email string, updateInfo domain.UpdateInfo) error
	GetUserByEmail(email string) (*domain.User, error)
}

type UserService interface {
	Create(user *domain.User) error
	GetUser(email string) (*domain.User, error)
	Update(user domain.User, updateInfo domain.UpdateInfo) error
	Login(email string, password string) (string, error)
	VerifyUser(token string) error
	ResetPasswordCreate(email string) error
	ResetPasswordResponse(token string, password string) error
}

type UserHandler interface {
	Create(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	ResetPasswordResponse(c *fiber.Ctx) error
}
