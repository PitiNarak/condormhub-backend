package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type SupportRepository interface {
	Create(support *domain.SupportRequest) error
}

type SupportService interface {
	Create(support *domain.SupportRequest) error
}

type SupportHandler interface {
	Create(c *fiber.Ctx) error
}
