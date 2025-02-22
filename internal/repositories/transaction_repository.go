package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tsx *domain.Transaction) *error_handler.ErrorHandler {
	err := r.db.Create(tsx).Error
	if err != nil {
		return error_handler.InternalServerError(err, "Failed to create order")
	}
	return nil
}

func (r *TransactionRepository) Update(tsx *domain.Transaction) *error_handler.ErrorHandler {
	err := r.db.Model(tsx).Where("id = ?", tsx.ID).Updates(tsx).Error
	if err != nil {
		return error_handler.InternalServerError(err, "Failed to create order")
	}
	return nil
}
