package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(contract *domain.Contract) error
	GetContract(lessorID uuid.UUID, lesseeID uuid.UUID, dormID uuid.UUID) (*domain.Contract, error)
	Delete(lessorID uuid.UUID, lesseeID uuid.UUID, dormID uuid.UUID) error
	UpdateLessorStatus(contractRequestBody dto.ContractRequestBody, lessorStatus domain.ContractStatus) error
	UpdateLesseeStatus(contractRequestBody dto.ContractRequestBody, lesseeStatus domain.ContractStatus) error
	UpdateContractStatus(contractRequestBody dto.ContractRequestBody, status domain.ContractStatus) error
}

type ContractService interface {
	Create(contract *domain.Contract) error
	GetContract(lessorID uuid.UUID, lesseeID uuid.UUID, dormID uuid.UUID) (*domain.Contract, error)
	UpdateStatus(contractRequestBody dto.ContractRequestBody, lessorStatus domain.ContractStatus, userID uuid.UUID) error
}

type ContractHandler interface {
	Create(c *fiber.Ctx) error
	SignContract(c *fiber.Ctx) error
}
