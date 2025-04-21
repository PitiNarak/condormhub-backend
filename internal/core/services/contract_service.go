package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
)

type ContractService struct {
	contractRepo          ports.ContractRepository
	userRepo              ports.UserRepository
	dormRepo              ports.DormRepository
	leasingHistoryService ports.LeasingHistoryService
}

func NewContractService(contractRepo ports.ContractRepository, userRepo ports.UserRepository, dormRepo ports.DormRepository, leasingHistoryService ports.LeasingHistoryService) ports.ContractService {
	return &ContractService{
		contractRepo:          contractRepo,
		userRepo:              userRepo,
		dormRepo:              dormRepo,
		leasingHistoryService: leasingHistoryService,
	}
}

func (ct *ContractService) DeleteContract(contractID uuid.UUID) error {
	return ct.contractRepo.Delete(contractID)
}

func (ct *ContractService) GetContractByContractID(contractID uuid.UUID) (*domain.Contract, error) {
	return ct.contractRepo.GetContractByContractID(contractID)
}

func (ct *ContractService) GetByUserID(userID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	user, userErr := ct.userRepo.GetUserByID(userID)
	if userErr != nil {
		return nil, 0, 0, userErr
	}
	if user == nil || user.Role == "" {
		return nil, 0, 0, apperror.BadRequestError(errors.New("invalid user"), "user not found or role is missing")
	}
	if user.Role == domain.AdminRole {
		return nil, 0, 0, apperror.BadRequestError(errors.New("invalid user"), "role mismatch")
	}
	return ct.getContractsByRole(user.Role, userID, limit, page)
}
func (ct *ContractService) getContractsByRole(role domain.Role, userID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	if role == domain.LessorRole {
		return ct.contractRepo.GetContractByLessorID(userID, limit, page)
	}
	return ct.contractRepo.GetContractByLesseeID(userID, limit, page)
}

func (ct *ContractService) GetByDormID(dormID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	contracts, totalPage, totalRows, err := ct.contractRepo.GetContractByLessorID(dormID, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return contracts, totalPage, totalRows, nil
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

	if err := ct.contractRepo.UpdateStatus(contractID, status, &user.Role); err != nil {
		return err
	}

	contract, err := ct.contractRepo.GetContractByContractID(contractID)
	if err != nil {
		return err
	}

	if contract.LesseeStatus == domain.Signed && contract.LessorStatus == domain.Signed {
		if err := ct.contractRepo.UpdateStatus(contractID, domain.Signed, nil); err != nil {
			return err
		}
		if _, err := ct.leasingHistoryService.Create(contract.LesseeID, contract.DormID); err != nil {
			return err
		}
	}

	if contract.LesseeStatus == domain.Cancelled || contract.LessorStatus == domain.Cancelled {
		if err := ct.contractRepo.UpdateStatus(contractID, domain.Cancelled, nil); err != nil {
			return err
		}
	}
	return nil
}
