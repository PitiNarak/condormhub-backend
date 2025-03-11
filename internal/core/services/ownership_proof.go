package services

import (
	"errors"

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
	err := o.ownershipProofRepo.Create(ownershipProof)
	if err != nil {
		return err
	}
	return nil
}
func (o *OwnershipProofService) Delete(lessorID uuid.UUID) error {
	err := o.ownershipProofRepo.Delete(lessorID)
	if err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) GetByLessorID(lessorID uuid.UUID) (*domain.OwnershipProof, error) {
	ownershipProof, err := o.ownershipProofRepo.GetByLessorID(lessorID)
	if err != nil {
		return nil, err
	}
	return ownershipProof, nil
}

func (o *OwnershipProofService) UpdateDocument(lessorID uuid.UUID, updateDocumentRequestBody *dto.UpdateOwnerShipProofRequestBody) error {
	err := o.ownershipProofRepo.UpdateDocument(lessorID, updateDocumentRequestBody)
	if err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) UpdateStatus(lessorID uuid.UUID, adminID uuid.UUID, status domain.OwnershipProofStatus) error {
	admin, admin_err := o.userRepo.GetUserByID(adminID)
	if admin_err != nil {
		return admin_err
	}

	if *admin.Role != domain.Role("ADMIN") {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an admin")
	}

	updateStatusRequestBody := new(dto.UpdateOwnerShipProofStatusRequestBody)
	updateStatusRequestBody.Status = status
	updateStatusRequestBody.AdminID = adminID

	err := o.ownershipProofRepo.UpdateStatus(lessorID, updateStatusRequestBody)
	if err != nil {
		return err
	}
	return nil
}
