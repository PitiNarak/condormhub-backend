package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingRequestRepository interface {
	Create(LeasingRequest *domain.LeasingRequest) error
	Update(LeasingRequest *domain.LeasingRequest) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.LeasingHistory, error)
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
}

type LeasingRequestService interface {
	Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingRequest, error)
	Delete(id uuid.UUID) error
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
	Approve(id uuid.UUID) error
	Reject(id uuid.UUID) error
}

type LeasingRequestHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
}
