package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
)

type ContractService struct {
	contractRepo          ports.ContractRepository
	userRepo              ports.UserRepository
	dormRepo              ports.DormRepository
	leasingHistoryService ports.LeasingHistoryService
	dormService           ports.DormService
}

func NewContractService(contractRepo ports.ContractRepository, userRepo ports.UserRepository, dormRepo ports.DormRepository, leasingHistoryService ports.LeasingHistoryService, dormService ports.DormService) ports.ContractService {
	return &ContractService{
		contractRepo:          contractRepo,
		userRepo:              userRepo,
		dormRepo:              dormRepo,
		leasingHistoryService: leasingHistoryService,
		dormService:           dormService,
	}
}

func (ct *ContractService) DeleteContract(contractID uuid.UUID) error {
	return ct.contractRepo.Delete(contractID)
}

func (ct *ContractService) GetContractByContractID(contractID uuid.UUID) (*dto.ContractResponseBody, error) {
	contract, err := ct.contractRepo.GetContractByContractID(contractID)
	if err != nil {
		return nil, err
	}
	urls := ct.dormService.GetImageUrl(contract.Dorm.Images)
	contractResponse := contract.ToDTO(urls)
	return &contractResponse, nil
}

func (ct *ContractService) GetByUserID(userID uuid.UUID, limit, page int) (*[]dto.ContractResponseBody, int, int, error) {
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
	contracts, totalPage, totalRows, err := ct.getContractsByRole(user.Role, userID, limit, page)

	resData := make([]dto.ContractResponseBody, len(*contracts))
	for i, v := range *contracts {
		urls := ct.dormService.GetImageUrl(v.Dorm.Images)
		resData[i] = v.ToDTO(urls)
	}
	return &resData, totalPage, totalRows, err
}
func (ct *ContractService) getContractsByRole(role domain.Role, userID uuid.UUID, limit, page int) (*[]domain.Contract, int, int, error) {
	if role == domain.LessorRole {
		return ct.contractRepo.GetContractByLessorID(userID, limit, page)
	}
	return ct.contractRepo.GetContractByLesseeID(userID, limit, page)

}

func (ct *ContractService) GetByDormID(dormID uuid.UUID, limit, page int) (*[]dto.ContractResponseBody, int, int, error) {
	contracts, totalPage, totalRows, err := ct.contractRepo.GetContractByLessorID(dormID, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	resData := make([]dto.ContractResponseBody, len(*contracts))
	for i, v := range *contracts {
		urls := ct.dormService.GetImageUrl(v.Dorm.Images)
		resData[i] = v.ToDTO(urls)
	}
	return &resData, totalPage, totalRows, err
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
