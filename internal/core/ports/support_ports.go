package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SupportRepository interface {
	Create(support *domain.SupportRequest) error
	GetAll(limit int, page int, userID uuid.UUID, isAdmin bool) ([]domain.SupportRequest, int, int, error)
	GetByID(id uuid.UUID) (*domain.SupportRequest, error)
	UpdateStatus(id uuid.UUID, status domain.SupportStatus) error
}

type SupportService interface {
	Create(support *domain.SupportRequest) error
	GetAll(limit int, page int, userID uuid.UUID, isAdmin bool) ([]domain.SupportRequest, int, int, error)
	GetByID(id uuid.UUID) (*domain.SupportRequest, error)
	UpdateStatus(id uuid.UUID, status domain.SupportStatus) error
}

type SupportHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
}
