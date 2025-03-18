package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type ContractRepository struct {
	db *database.Database
}

func NewContractRepository(db *database.Database) ports.ContractRepository {
	return &ContractRepository{db: db}
}

func (ct *ContractRepository) Create(contract *domain.Contract) error {
	if err := ct.db.Create(contract).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save contract to database")
	}
	return nil
}

func (ct *ContractRepository) GetByDormID(dormID uuid.UUID) (*domain.Contract, error) {
	contract := new(domain.Contract)
	if err := ct.db.First(contract, dormID).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contract not found")
	}
	return contract, nil
}
