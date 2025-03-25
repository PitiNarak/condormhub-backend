package repository

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *ReceiptRepository) GetByUserID(userID uuid.UUID, limit int, page int) ([]domain.Receipt, int, int, error) {
	var receipts []domain.Receipt
	query := r.db.Preload("Owner").
		Preload("Transaction").
		Where("owner_id = ?", userID).
		Find(&receipts)

	totalPage, totalRows, err := r.db.Paginate(&receipts, query, limit, page, "create_at DESC")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "receipt not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to receipts")
	}

	return receipts, totalPage, totalRows, nil
}
