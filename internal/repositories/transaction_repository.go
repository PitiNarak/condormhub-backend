package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ports.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tsx *domain.Transaction) *errorHandler.ErrorHandler {
	err := r.db.Create(tsx).Error
	if err != nil {
		return errorHandler.InternalServerError(err, "Failed to create order")
	}
	return nil
}

func (r *TransactionRepository) Update(tsx *domain.Transaction) *errorHandler.ErrorHandler {
	err := r.db.Model(tsx).Where("id = ?", tsx.ID).Updates(tsx).Error
	if err != nil {
		return errorHandler.InternalServerError(err, "Failed to create order")
	}
	return nil
}
