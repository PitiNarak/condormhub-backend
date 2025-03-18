package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(contract *domain.Contract) error
	GetByDormID(dormID uuid.UUID) (*domain.Contract, error)
}

type ContractService interface {
	Create(contract *domain.Contract) error
	GetByDormID(dormID uuid.UUID) (*domain.Contract, error)
}

type ContractHandler interface {
	Create(c *fiber.Ctx) error
}
