package ports

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/gofiber/fiber/v2"
)

type SampleLogService interface {
	Save(message string) error
	Delete(id string) error
	GetAll() ([]domain.SampleLog, error)
}

type SampleLogRepository interface {
	Save(message string) error
	Delete(id int) error
	GetAll() ([]domain.SampleLog, error)
}

type SampleLogHandler interface {
	Save(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}
