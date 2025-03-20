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

func (ct *ContractService) DeleteContract(contractID uuid.UUID) error {
	return ct.contractRepo.Delete(contractID)
}

func (ct *ContractService) GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error) {
	return ct.contractRepo.GetContractByContractID(contractID)
}

func (ct *ContractService) UpdateStatus(contractID uuid.UUID, status domain.ContractStatus, userID uuid.UUID) error {
	//check valid lessor ID

	user, userErr := ct.userRepo.GetUserByID(userID)
	if userErr != nil {
		return userErr
	}
	if user == nil || user.Role == "" {
		return apperror.BadRequestError(errors.New("invalid user"), "user not found or role is missing")
	}
	if user.Role == domain.AdminRole {
		return apperror.BadRequestError(errors.New("invalid user"), "role mismatch")
	}

	if user.Role == domain.LessorRole {
		if err := ct.contractRepo.UpdateLessorStatus(contractID, status); err != nil {
			return err
		}
	} else {
		if err := ct.contractRepo.UpdateLesseeStatus(contractID, status); err != nil {
			return err
		}
	}

	contract, err := ct.contractRepo.GetContractByContractID(contractID)
	if err != nil {
		return err
	}

	if contract.LesseeStatus == domain.Signed && contract.LessorStatus == domain.Signed {
		if err := ct.contractRepo.UpdateContractStatus(contractID, domain.Signed); err != nil {
			return err
		}
	}

	if contract.LesseeStatus == domain.Cancelled || contract.LessorStatus == domain.Cancelled {
		if err := ct.contractRepo.UpdateContractStatus(contractID, domain.Cancelled); err != nil {
			return err
		}
	}
	return nil
}
