package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
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

func (o *OwnershipProofRepository) UpdateDocument(dormID uuid.UUID, fileKey string) error {
	if err := o.db.Model(&domain.OwnershipProof{}).Where("dorm_id = ?", dormID).Updates(map[string]interface{}{"file_key": fileKey}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update document")
	}
	return nil
}

func (o *OwnershipProofRepository) UpdateStatus(dormID uuid.UUID, updateStatusRequestBody *dto.UpdateOwnerShipProofStatusRequestBody) error {
	if err := o.db.Model(&domain.OwnershipProof{}).Where("dorm_id = ?", dormID).Updates(updateStatusRequestBody).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update status")
	}
	return nil
}
