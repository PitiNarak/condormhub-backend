package ports

import (
	"context"
	"io"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OwnershipProofRepository interface {
	Create(ownershipProof *domain.OwnershipProof) error
	Delete(dormID uuid.UUID) error
	GetByDormID(dormID uuid.UUID) (*domain.OwnershipProof, error)
	UpdateDocument(dormID uuid.UUID, fileKey string) error
	UpdateStatus(dormID uuid.UUID, updateStatusRequestBody *dto.UpdateOwnerShipProofStatusRequestBody) error
}

type OwnershipProofService interface {
	Create(*domain.OwnershipProof) error
	Delete(dormID uuid.UUID) error
	GetByDormID(dormID uuid.UUID) (*domain.OwnershipProof, error)
	UpdateDocument(dormID uuid.UUID, fileKey string) error
	UpdateStatus(dormID uuid.UUID, adminID uuid.UUID, status domain.OwnershipProofStatus) error
	ConvertToDTO(ownershipProof domain.OwnershipProof, url string, expires time.Time) dto.OwnershipProofResponseBody
	UploadFile(ctx context.Context, dormID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error)
}

type OwnershipProofHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
}
