package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
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

func (ct *ContractRepository) Delete(contractID uuid.UUID) error {
	if err := ct.db.Where("contract_id = ?", contractID).Delete(&domain.Contract{}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete contract")
	}
	return nil
}

func (ct *ContractRepository) GetContractByLessorID(lessorID uuid.UUID) (*[]domain.Contract, error) {
	var contracts []domain.Contract
	if err := ct.db.Where("lessor_id = ? ", lessorID).Find(&contracts).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contracts not found")
	}
	return &contracts, nil
}

func (ct *ContractRepository) GetContractByLesseeID(lesseeID uuid.UUID) (*[]domain.Contract, error) {
	var contracts []domain.Contract
	if err := ct.db.Where("lessee_id = ? ", lesseeID).Find(&contracts).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contracts not found")
	}
	return &contracts, nil
}

func (ct *ContractRepository) GetContractByDormID(DormID uuid.UUID) (*[]domain.Contract, error) {
	var contracts []domain.Contract
	if err := ct.db.Where("dorm_id = ? ", DormID).Find(&contracts).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Contracts not found")
	}
	return &contracts, nil
}

func (ct *ContractRepository) UpdateLessorStatus(contractRequestBody dto.ContractRequestBody, LessorStatus domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("dorm_id = ? AND lessor_id = ? AND lessee_id = ?",
		contractRequestBody.DormID, contractRequestBody.LessorID, contractRequestBody.LesseeID).Updates(
		map[string]interface{}{"lessor_status": LessorStatus}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update lessor status")
	}
	return nil
}

func (ct *ContractRepository) UpdateLesseeStatus(contractRequestBody dto.ContractRequestBody, LesseeStatus domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("dorm_id = ? AND lessor_id = ? AND lessee_id = ?",
		contractRequestBody.DormID, contractRequestBody.LessorID, contractRequestBody.LesseeID).Updates(
		map[string]interface{}{"lessee_status": LesseeStatus}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update lessee status")
	}
	return nil
}

func (ct *ContractRepository) UpdateContractStatus(contractRequestBody dto.ContractRequestBody, status domain.ContractStatus) error {
	if err := ct.db.Model(&domain.Contract{}).Where("dorm_id = ? AND lessor_id = ? AND lessee_id = ?",
		contractRequestBody.DormID, contractRequestBody.LessorID, contractRequestBody.LesseeID).Updates(
		map[string]interface{}{"status": status}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update contract status")
	}
	return nil
}
