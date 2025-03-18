package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
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

func (ct *ContractService) GetContract(lessorID uuid.UUID, lesseeID uuid.UUID, dormID uuid.UUID) (*domain.Contract, error) {
	return ct.contractRepo.GetContract(lessorID, lesseeID, dormID)
}

func (ct *ContractService) UpdateStatus(contractRequestBody dto.ContractRequestBody, status domain.ContractStatus, userID uuid.UUID) error {
	//check valid lessor ID

	lessor, lessorErr := ct.userRepo.GetUserByID(contractRequestBody.LessorID)
	if lessorErr != nil {
		return lessorErr
	}
	if lessor == nil || lessor.Role == "" {
		return apperror.BadRequestError(errors.New("invalid lessor"), "lessor not found or role is missing")
	}

	if lessor.Role != domain.LessorRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an lessor")
	}

	//check valid lessee ID
	lessee, lesseeErr := ct.userRepo.GetUserByID(contractRequestBody.LesseeID)
	if lesseeErr != nil {
		return lesseeErr
	}
	if lessee == nil || lessee.Role == "" {
		return apperror.BadRequestError(errors.New("invalid lessee"), "lessee not found or role is missing")
	}

	if lessee.Role != domain.LesseeRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an lessee")
	}

	//check valid dorm ID
	_, dormErr := ct.dormRepo.GetByID(contractRequestBody.DormID)
	if dormErr != nil {
		return dormErr
	}

	if contractRequestBody.LesseeID == userID {
		if err := ct.contractRepo.UpdateLessorStatus(contractRequestBody, status); err != nil {
			return err
		}
	} else {
		if err := ct.contractRepo.UpdateLesseeStatus(contractRequestBody, status); err != nil {
			return err
		}
	}

	contract, err := ct.contractRepo.GetContract(contractRequestBody.LessorID, contractRequestBody.LesseeID, contractRequestBody.DormID)
	if err != nil {
		return err
	}

	if contract.LesseeStatus == domain.Signed && contract.LessorStatus == domain.Signed {
		if err := ct.contractRepo.UpdateContractStatus(contractRequestBody, domain.Signed); err != nil {
			return err
		}
	}
	return nil
}
