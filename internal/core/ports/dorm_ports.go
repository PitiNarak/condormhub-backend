package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DormRepository interface {
	Create(dorm *domain.Dorm) error
	GetAll(limit int, page int) ([]domain.Dorm, int, int, error)
	GetByID(id uuid.UUID) (*domain.Dorm, error)
	Update(id uuid.UUID, dorm *domain.Dorm) error
	Delete(id uuid.UUID) error
}

type DormService interface {
	Create(userRole domain.Role, dorm *domain.Dorm) error
	GetAll(limit int, page int) ([]domain.Dorm, int, int, error)
	GetByID(id uuid.UUID) (*domain.Dorm, error)
	Update(userID uuid.UUID, isAdmin bool, dormID uuid.UUID, dorm *dto.DormRequestBody) (*domain.Dorm, error)
	Delete(userID uuid.UUID, isAdmin bool, dormID uuid.UUID) error
}

type DormHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
