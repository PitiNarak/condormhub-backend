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
	if err := ct.db.Where("contract_id = ?", contractID).Delete(&domain.Contract{}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete contract")
	}
	return nil
}

func (ct *ContractRepository) GetContract(lessorID uuid.UUID, lesseeID uuid.UUID, dormID uuid.UUID) (*[]domain.Contract, error) {
	var contracts []domain.Contract
	if err := ct.db.Where("lessor_id = ? AND lessee_id = ? AND dorm_id = ?", lessorID, lesseeID, dormID).Find(&contracts).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contracts not found")
	}
	return &contracts, nil
}

func (ct *ContractRepository) GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error) {
	contract := new(domain.Contract)
	if err := ct.db.Where("contract_id = ? ", contractID).Find(&contract).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contract not found")
	}
	return contract, nil
}

func (ct *ContractRepository) GetContractByLessorID(lessorID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.Where("lessor_id = ? ", lessorID).Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) GetContractByLesseeID(lesseeID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.Where("lessee_id = ? ", lesseeID).Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) GetContractByDormID(dormID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	var contracts []domain.Contract
	query := ct.db.Where("dorm_id = ? ", dormID).Find(&contracts)

	totalPage, totalRows, err := ct.db.Paginate(&contracts, query, limit, page, "create_at DESC")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "contract not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get contract")
	}

	return &contracts, totalPage, totalRows, nil
}

func (ct *ContractRepository) UpdateLessorStatus(contractID uuid.UUID, LessorStatus domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("contract_id = ?",
		contractID).Updates(
		map[string]interface{}{"lessor_status": LessorStatus}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update lessor status")
	}
	return nil
}

func (ct *ContractRepository) UpdateLesseeStatus(contractID uuid.UUID, LesseeStatus domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("contract_id = ?",
		contractID).Updates(
		map[string]interface{}{"lessee_status": LesseeStatus}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update lessee status")
	}
	return nil
}

func (ct *ContractRepository) UpdateContractStatus(contractID uuid.UUID, status domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("contract_id = ?",
		contractID).Updates(
		map[string]interface{}{"status": status}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update contract status")
	}
	return nil
}
