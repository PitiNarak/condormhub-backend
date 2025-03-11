package repository

import (
	"fmt"

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
		fmt.Println("DB Error:", err)
		return apperror.InternalServerError(err, "Failed to save ownership proof to database")
	}
	return nil
}

func (o *OwnershipProofRepository) Delete(lessorID uuid.UUID) error {
	if err := o.db.Delete(&domain.OwnershipProof{}, lessorID).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete ownership proof")
	}
	return nil
}

func (o *OwnershipProofRepository) GetByLessorID(lessorID uuid.UUID) (*domain.OwnershipProof, error) {
	ownershipProof := new(domain.OwnershipProof)
	if err := o.db.First(ownershipProof, lessorID).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Ownership proof not found")
	}
	return ownershipProof, nil
}

func (o *OwnershipProofRepository) UpdateDocument(lessorID uuid.UUID, updateDocumentRequestBody *dto.UpdateOwnerShipProofRequestBody) error {
	existingOwnershipProof, err := o.GetByLessorID(lessorID)
	if err != nil {
		return apperror.NotFoundError(err, "Ownership proof not found")
	}
	updateErr := o.db.Model(existingOwnershipProof).Where("lessor_id = ?", lessorID).Updates(updateDocumentRequestBody).Error
	if updateErr != nil {
		return apperror.InternalServerError(err, "failed to update document")
	}

	return nil
}

func (o *OwnershipProofRepository) UpdateStatus(lessorID uuid.UUID, updateStatusRequestBody *dto.UpdateOwnerShipProofStatusRequestBody) error {
	existingOwnershipProof, err := o.GetByLessorID(lessorID)
	if err != nil {
		return apperror.NotFoundError(err, "Ownership proof not found")
	}
	updateErr := o.db.Model(existingOwnershipProof).Where("lessorId = ?", lessorID).Updates(updateStatusRequestBody).Error
	if updateErr != nil {
		return apperror.InternalServerError(err, "failed to update status")
	}

	return nil
}
