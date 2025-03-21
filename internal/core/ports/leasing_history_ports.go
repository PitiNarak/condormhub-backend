package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingHistoryRepository interface {
	Create(LeasingHistory *domain.LeasingHistory) error
	Update(LeasingHistory *domain.LeasingHistory) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.LeasingHistory, error)
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
}

type LeasingHistoryService interface {
	Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error)
	CreateReview(user domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error)
	UpdateReview(user domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error)
	DeleteReview(user domain.User, id uuid.UUID) error
	Delete(id uuid.UUID) error
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	SetEndTimestamp(id uuid.UUID) error
}

type LeasingHistoryHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	CreateReview(c *fiber.Ctx) error
	UpdateReview(c *fiber.Ctx) error
	DeleteReview(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
	SetEndTimestamp(c *fiber.Ctx) error
}
