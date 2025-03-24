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
	// Validate input contract
	if err := ct.validateContract(contract); err != nil {
		return err
	}
	// Check and validate lessor
	if _, err := ct.getUserAndValidateRole(contract.LessorID, domain.LessorRole); err != nil {
		return err
	}
	// Check and validate lessee
	if _, err := ct.getUserAndValidateRole(contract.LesseeID, domain.LesseeRole); err != nil {
		return err
	}
	// Check if dorm exists
	if _, err := ct.dormRepo.GetByID(contract.DormID); err != nil {
		return err
	}
	// Check for existing active contract
	if err := ct.checkForExistingActiveContract(contract); err != nil {
		return err
	}
	// Create the contract
	if err := ct.contractRepo.Create(contract); err != nil {
		return err
	}
	return nil
}
func (ct *ContractService) validateContract(contract *domain.Contract) error {
	if contract.LessorID == uuid.Nil || contract.LesseeID == uuid.Nil || contract.DormID == uuid.Nil {
		return apperror.BadRequestError(errors.New("invalid contract"), "missing required fields")
	}
	return nil
}
func (ct *ContractService) getUserAndValidateRole(userID uuid.UUID, role domain.Role) (*domain.User, error) {
	user, err := ct.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Role == "" {
		return nil, apperror.BadRequestError(errors.New("invalid user"), "user not found or role is missing")
	}
	if user.Role != role {
		return nil, apperror.BadRequestError(errors.New("role mismatch"), "user role mismatch")
	}
	return user, nil
}
func (ct *ContractService) checkForExistingActiveContract(contract *domain.Contract) error {
	contracts, err := ct.contractRepo.GetContract(contract.LessorID, contract.LesseeID, contract.DormID)
	if err != nil {
		return err
	}
	for _, c := range *contracts {
		if c.Status == domain.Waiting {
			return apperror.BadRequestError(errors.New("contract already exists"), "active contract already exists")
		}
	}
	return nil
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
	}

	if contract.LesseeStatus == domain.Cancelled || contract.LessorStatus == domain.Cancelled {
		if err := ct.contractRepo.UpdateStatus(contractID, domain.Cancelled, nil); err != nil {
			return err
		}
	}
	return nil
}
