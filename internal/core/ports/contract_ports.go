package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(contract *domain.Contract) error
	GetContract(LesseeID uuid.UUID, DormID uuid.UUID) (*[]domain.Contract, error)
	GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error)
	GetContractByLessorID(LessorID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error)
	GetContractByLesseeID(LesseeID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error)
	GetContractByDormID(DormID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error)
	Delete(contractID uuid.UUID) error
	UpdateStatus(contractID uuid.UUID, status domain.ContractStatus, role *domain.Role) error
}

type ContractService interface {
	GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error)
	GetByUserID(userID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error)
	GetByDormID(lesseeID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error)
	DeleteContract(contractID uuid.UUID) error
	UpdateStatus(contractID uuid.UUID, lessorStatus domain.ContractStatus, userID uuid.UUID) error
}

type ContractHandler interface {
	GetContractByContractID(c *fiber.Ctx) error
	GetContractByUserID(c *fiber.Ctx) error
	GetContractByDormID(c *fiber.Ctx) error
	SignContract(c *fiber.Ctx) error
	CancelContract(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
