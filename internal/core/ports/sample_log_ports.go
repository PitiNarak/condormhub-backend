package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SampleLogService interface {
	Save(message string) error
	Delete(id uuid.UUID) error
	GetAll() ([]domain.SampleLog, error)
	EditMessage(id uuid.UUID, message string) error
}

type SampleLogRepository interface {
	Save(message string) error
	Delete(id uuid.UUID) error
	GetAll() ([]domain.SampleLog, error)
	EditMessage(id uuid.UUID, message string) error
}

type SampleLogHandler interface {
	Save(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	EditMessage(c *fiber.Ctx) error
}
