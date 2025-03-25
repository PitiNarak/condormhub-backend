package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
)

type ReceiptRepository struct {
	db *database.Database
}

func NewReceiptRepository(db *database.Database) ports.ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) Create(receipt *domain.Receipt) error {
	if err := r.db.Create(receipt).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create receipt")
	}

	if err := r.db.Preload("Owner").Preload("Transaction").First(receipt, receipt.ID).Error; err != nil {
		return apperror.InternalServerError(err, "failed to preload receipt")
	}

	return nil
}
