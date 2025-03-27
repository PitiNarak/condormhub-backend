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
	GetByID(id uuid.UUID) (*domain.LeasingRequest, error)
	GetByUserID(id uuid.UUID, limit, page int, role domain.Role) ([]domain.LeasingRequest, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
	GetID(dormID, userID uuid.UUID) (uuid.UUID, error)
}

type LeasingRequestService interface {
	Create(leeseeID uuid.UUID, dormID uuid.UUID, message string) (*domain.LeasingRequest, error)
	Delete(id uuid.UUID) error
	GetByUserID(id uuid.UUID, role domain.Role, limit, page int) ([]domain.LeasingRequest, int, int, error)
	Approve(id, userId uuid.UUID, isAdmin bool) error
	Reject(id, userId uuid.UUID, isAdmin bool) error
	Cancel(id, userId uuid.UUID, isAdmin bool) error
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error)
}

type LeasingRequestHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
}
