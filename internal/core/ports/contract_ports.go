package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(contract *domain.Contract) error
	GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error)
	GetContractByLessorID(LessorID uuid.UUID) (*[]domain.Contract, error)
	GetContractByLesseeID(LesseeID uuid.UUID) (*[]domain.Contract, error)
	GetContractByDormID(DormID uuid.UUID) (*[]domain.Contract, error)
	Delete(contractID uuid.UUID) error
	UpdateLessorStatus(contractID uuid.UUID, lessorStatus domain.ContractStatus) error
	UpdateLesseeStatus(contractID uuid.UUID, lesseeStatus domain.ContractStatus) error
	UpdateContractStatus(contractID uuid.UUID, status domain.ContractStatus) error
}

type ContractService interface {
	Create(contract *domain.Contract) error
	GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error)
	DeleteContract(contractID uuid.UUID) error
	UpdateStatus(contractID uuid.UUID, lessorStatus domain.ContractStatus, userID uuid.UUID) error
}

type ContractHandler interface {
	Create(c *fiber.Ctx) error
	SignContract(c *fiber.Ctx) error
	CancelContract(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
