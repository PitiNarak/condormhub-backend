package services

import (
	"errors"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type OwnershipProofService struct {
	ownershipProofRepo ports.OwnershipProofRepository
	userRepo           ports.UserRepository
}

func NewOwnershipProofService(ownershipProofRepo ports.OwnershipProofRepository, userRepo ports.UserRepository) ports.OwnershipProofService {
	return &OwnershipProofService{
		ownershipProofRepo: ownershipProofRepo,
		userRepo:           userRepo,
	}
}

func (o *OwnershipProofService) Create(ownershipProof *domain.OwnershipProof) error {
	if err := o.ownershipProofRepo.Create(ownershipProof); err != nil {
		return err
	}
	return nil
}
func (o *OwnershipProofService) Delete(dormID uuid.UUID) error {
	if err := o.ownershipProofRepo.Delete(dormID); err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) GetByDormID(dormID uuid.UUID) (*domain.OwnershipProof, error) {
	ownershipProof, err := o.ownershipProofRepo.GetByDormID(dormID)
	if err != nil {
		return nil, err
	}
	return ownershipProof, nil
}

func (o *OwnershipProofService) UpdateDocument(dormID uuid.UUID, fileKey string) error {
	if err := o.ownershipProofRepo.UpdateDocument(dormID, fileKey); err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) UpdateStatus(dormID uuid.UUID, adminID uuid.UUID, status domain.OwnershipProofStatus) error {
	admin, err := o.userRepo.GetUserByID(adminID)
	if err != nil {
		return err
	}

	if admin == nil || admin.Role == nil {
		return apperror.BadRequestError(errors.New("invalid admin"), "Admin not found or role is missing")
	}

	if *admin.Role != domain.AdminRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an admin")
	}

	updateStatusRequestBody := new(dto.UpdateOwnerShipProofStatusRequestBody)
	updateStatusRequestBody.Status = status
	updateStatusRequestBody.AdminID = adminID

	if err := o.ownershipProofRepo.UpdateStatus(dormID, updateStatusRequestBody); err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) ConvertToDTOWithFile(ownershipProof domain.OwnershipProof, url string, expires time.Time) dto.OwnershipProofWithFileResponseBody {
	ownershipProofWithFileResponseBody := dto.OwnershipProofWithFileResponseBody{
		Url:     url,
		Expires: expires,
		DormID:  ownershipProof.DormID,
		AdminID: ownershipProof.AdminID,
		Status:  ownershipProof.Status,
	}

	return ownershipProofWithFileResponseBody
}

func (o *OwnershipProofService) ConvertToDTO(ownershipProof domain.OwnershipProof) dto.OwnershipProofResponseBody {
	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		DormID:  ownershipProof.DormID,
		AdminID: ownershipProof.AdminID,
		Status:  ownershipProof.Status,
	}

	return ownershipProofResponseBody
}
