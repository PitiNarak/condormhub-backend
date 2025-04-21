package repository

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
	"gorm.io/gorm"
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

func (ct *ContractRepository) Delete(contractID uuid.UUID) error {
	if err := ct.db.Where("id = ?", contractID).Delete(&domain.Contract{}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete contract")
	}
	return nil
}

func (ct *ContractRepository) GetContract(lesseeID uuid.UUID, dormID uuid.UUID) (*[]domain.Contract, error) {
	var contracts []domain.Contract
	if err := ct.db.
		Preload("Lessee").
		Preload("Dorm").
		Where("lessee_id = ? AND dorm_id = ?", lesseeID, dormID).
		Find(&contracts).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contracts not found")
	}
	return &contracts, nil
}

func (ct *ContractRepository) GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error) {
	contract := new(domain.Contract)
	if err := ct.db.
		Preload("Lessee").
		Preload("Dorm").
		Where("id = ? ", contractID).
		Find(&contract).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contract not found")
	}
	return contract, nil
}

func (ct *ContractRepository) GetContractByLessorID(lessorID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.
		Joins("JOIN dorms ON dorms.id = contracts.dorm_id").
		Where("dorms.owner_id = ?", lessorID).
		Preload("Lessee").
		Preload("Dorm").
		Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) GetContractByLesseeID(lesseeID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.
		Preload("Lessee").
		Preload("Dorm").
		Where("lessee_id = ? ", lesseeID).
		Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) GetContractByDormID(dormID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.
		Preload("Lessee").
		Preload("Dorm").
		Where("dorm_id = ? ", dormID).
		Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) UpdateStatus(contractID uuid.UUID, status domain.ContractStatus, role *domain.Role) error {
	if role == nil {
		contract := domain.Contract{
			Status: status,
		}
		if err := ct.db.Model(&domain.Contract{}).Where("id = ?",
			contractID).Updates(contract).Error; err != nil {
			return apperror.InternalServerError(err, "failed to update contract status")
		}
	} else if *role == domain.LessorRole {
		contract := domain.Contract{
			LessorStatus: status,
		}
		if err := ct.db.Model(&domain.Contract{}).Where("id = ?",
			contractID).Updates(contract).Error; err != nil {
			return apperror.InternalServerError(err, "failed to update lessor status")
		}
	} else if *role == domain.LesseeRole {
		contract := domain.Contract{
			LesseeStatus: status,
		}
		if err := ct.db.Model(&domain.Contract{}).Where("id = ?",
			contractID).Updates(contract).Error; err != nil {
			return apperror.InternalServerError(err, "failed to update lessee status")
		}
	} else {
		return apperror.InternalServerError(errors.New("invalid role"), "Invalid role")
	}

	return nil
}
