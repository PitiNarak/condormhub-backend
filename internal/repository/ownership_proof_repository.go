package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type OwnershipProofRepository struct {
	db *database.Database
}

func NewOwnershipProofRepository(db *database.Database) ports.OwnershipProofRepository {
	return &OwnershipProofRepository{db: db}
}

func (o *OwnershipProofRepository) Create(ownershipProof *domain.OwnershipProof) error {
	if err := o.db.Create(ownershipProof).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save ownership proof to database")
	}
	return nil
}

func (o *OwnershipProofRepository) Delete(dormID uuid.UUID) error {
	if err := o.db.Delete(&domain.OwnershipProof{}, dormID).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete ownership proof")
	}
	return nil
}

func (o *OwnershipProofRepository) GetByDormID(dormID uuid.UUID) (*domain.OwnershipProof, error) {
	ownershipProof := new(domain.OwnershipProof)
	if err := o.db.First(ownershipProof, dormID).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Ownership proof not found")
	}
	return ownershipProof, nil
}

func (o *OwnershipProofRepository) UpdateDocument(dormID uuid.UUID, updateDocumentRequestBody *dto.UpdateOwnerShipProofRequestBody) error {
	existingOwnershipProof, err := o.GetByDormID(dormID)
	if err != nil {
		return apperror.NotFoundError(err, "Ownership proof not found")
	}
	updateErr := o.db.Model(existingOwnershipProof).Where("dorm_id = ?", dormID).Updates(updateDocumentRequestBody).Error
	if updateErr != nil {
		return apperror.InternalServerError(updateErr, "failed to update document")
	}

	return nil
}

func (o *OwnershipProofRepository) UpdateStatus(dormID uuid.UUID, updateStatusRequestBody *dto.UpdateOwnerShipProofStatusRequestBody) error {
	existingOwnershipProof, err := o.GetByDormID(dormID)
	if err != nil {
		return apperror.NotFoundError(err, "Ownership proof not found")
	}
	updateErr := o.db.Model(existingOwnershipProof).Where("dorm_id =  ?", dormID).Updates(updateStatusRequestBody).Error
	if updateErr != nil {
		return apperror.InternalServerError(updateErr, "failed to update status")
	}

	return nil
}
