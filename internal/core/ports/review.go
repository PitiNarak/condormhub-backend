package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReviewRepository interface {
	Create(Review *domain.Review) error
	Update(Review *domain.Review) error
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*domain.Review, error)
}

type ReviewService interface {
	Create(Massage string, Rate int) (*domain.Review, error)
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) ([]domain.Review, int, int, error)
	Update(Massage string, Rate int) (*domain.Review, error)
}

type ReviewHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}
