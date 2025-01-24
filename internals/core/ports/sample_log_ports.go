package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SampleLogService interface {
	Save(message string) error
	Delete(id uuid.UUID) error
	GetAll() ([]domain.SampleLog, error)
}

type SampleLogRepository interface {
	Save(message string) error
	Delete(id uuid.UUID) error
	GetAll() ([]domain.SampleLog, error)
}

type SampleLogHandler interface {
	Save(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}
