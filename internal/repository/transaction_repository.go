package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/yokeTH/go-pkg/apperror"
)

type TransactionRepository struct {
	db *database.Database
}

func NewTransactionRepository(db *database.Database) ports.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tsx *domain.Transaction) error {
	err := r.db.Create(tsx).Error
	if err != nil {
		return apperror.InternalServerError(err, "Failed to create order")
	}
	return nil
}

func (r *TransactionRepository) Update(tsx *domain.Transaction) error {
	err := r.db.Model(&tsx).Where("id = ?", tsx.ID).Updates(tsx).Error
	if err != nil {
		return apperror.InternalServerError(err, "Failed to create order")
	}
	return nil
}

func (r *TransactionRepository) GetByID(id string) (domain.Transaction, error) {
	var tsx domain.Transaction
	err := r.db.Where("id = ?", id).First(&tsx).Error
	if err != nil {
		return tsx, apperror.NotFoundError(err, "Transaction not found")
	}
	return tsx, nil
}
