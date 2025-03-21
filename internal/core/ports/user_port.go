package ports

import (
	"context"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUserByID(userID uuid.UUID) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error
	UpdateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	DeleteAccount(userID uuid.UUID) error
}

type UserService interface {
	Create(ctx context.Context, user *domain.User) (string, string, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error)
	Login(context.Context, string, string) (*domain.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	VerifyUser(ctx context.Context, token string) (*domain.User, string, string, error)
	ResetPasswordCreate(context.Context, string) error
	ResetPassword(context.Context, string, string) (*domain.User, string, string, error)
	DeleteAccount(userID uuid.UUID) error
	ResendVerificationEmailService(ctx context.Context, email string) error
}

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	UpdateUserInformation(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	DeleteAccount(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	ResendVerificationEmailHandler(c *fiber.Ctx) error
}
