package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUser(userID uuid.UUID) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error
	UpdateUser(user domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}

type UserService interface {
	Create(user *domain.User) (string, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error)
	Login(email string, password string) (string, error)
	VerifyUser(token string) (string, *domain.User, error)
	ResetPasswordCreate(email string) error
	ResetPasswordResponse(token string, password string) error
}

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	UpdateUserInformation(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	ResetPasswordResponse(c *fiber.Ctx) error
}
