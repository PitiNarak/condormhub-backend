package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OwnershipProofRepository interface {
	Create(ownershipProof *domain.OwnershipProof) error
	Delete(lessorID uuid.UUID) error
	GetByLessorID(lessorID uuid.UUID) (*domain.OwnershipProof, error)
	UpdateDocument(lessorID uuid.UUID, updateDocumentRequestBody *dto.UpdateOwnerShipProofRequestBody) error
	UpdateStatus(lessorID uuid.UUID, updateStatusRequestBody *dto.UpdateOwnerShipProofStatusRequestBody) error
}

type OwnershipProofService interface {
	Create(*domain.OwnershipProof) error
	Delete(lessorID uuid.UUID) error
	GetByLessorID(lessorID uuid.UUID) (*domain.OwnershipProof, error)
	UpdateDocument(lessorID uuid.UUID, updateDocumentRequestBody *dto.UpdateOwnerShipProofRequestBody) error
	UpdateStatus(lessorID uuid.UUID, adminID uuid.UUID, status domain.OwnershipProofStatus) error
}

type OwnershipProofHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	// GetByLessorID(c *fiber.Ctx) error
}
