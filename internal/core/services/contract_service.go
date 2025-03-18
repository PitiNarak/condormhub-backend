package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type ContractService struct {
	contractRepo ports.ContractRepository
	userRepo     ports.UserRepository
	dormRepo     ports.DormRepository
}

func NewContractService(contractRepo ports.ContractRepository, userRepo ports.UserRepository, dormRepo ports.DormRepository) ports.ContractService {
	return &ContractService{
		contractRepo: contractRepo,
		userRepo:     userRepo,
		dormRepo:     dormRepo,
	}
}

func (ct *ContractService) Create(contract *domain.Contract) error {
	lessor, lessorErr := ct.userRepo.GetUserByID(contract.LessorID)
	if lessorErr != nil {
		return lessorErr
	}
	if lessor == nil || lessor.Role == "" {
		return apperror.BadRequestError(errors.New("invalid lessor"), "lessor not found or role is missing")
	}

	if lessor.Role != domain.LessorRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an lessor")
	}

	lessee, lesseeErr := ct.userRepo.GetUserByID(contract.LesseeID)
	if lesseeErr != nil {
		return lesseeErr
	}
	if lessee == nil || lessee.Role == "" {
		return apperror.BadRequestError(errors.New("invalid lessee"), "lessee not found or role is missing")
	}

	if lessee.Role != domain.LesseeRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an lessee")
	}

	_, dormErr := ct.dormRepo.GetByID(contract.DormID)
	if dormErr != nil {
		return dormErr
	}

	if err := ct.contractRepo.Create(contract); err != nil {
		return err
	}
	return nil
}

func (ct *ContractService) GetByDormID(dormID uuid.UUID) (*domain.Contract, error) {
	return ct.contractRepo.GetByDormID(dormID)
}
