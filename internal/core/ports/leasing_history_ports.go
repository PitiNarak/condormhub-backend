package ports

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingHistoryRepository interface {
	Create(LeasingHistory *domain.LeasingHistory) error
	Update(LeasingHistory *domain.LeasingHistory) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.LeasingHistory, error)
	GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error)
	GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error)
	PatchEndTimestamp(id uuid.UUID, endTime time.Time) error
}

type LeasingHistoryService interface {
	Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error)
	Update(LeasingHistory *domain.LeasingHistory) error
	Delete(id uuid.UUID) error
	GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error)
	GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error)
	PatchEndTimestamp(id uuid.UUID, endTime time.Time) error
}

type LeasingHistoryHandler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
	PatchEndTimestamp(c *fiber.Ctx) error
}
