package services

import (
	"errors"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type LeasingRequestService struct {
	requestRepo  ports.LeasingRequestRepository
	dormRepo     ports.DormRepository
	contractRepo ports.ContractRepository
}

func NewLeasingRequestService(requestRepo ports.LeasingRequestRepository, dormRepo ports.DormRepository, contractRepo ports.ContractRepository) ports.LeasingRequestService {
	return &LeasingRequestService{requestRepo: requestRepo, dormRepo: dormRepo, contractRepo: contractRepo}
}

func (s *LeasingRequestService) Create(leeseeID uuid.UUID, dormID uuid.UUID, message string) (*domain.LeasingRequest, error) {
	leasingRequest := &domain.LeasingRequest{
		Status:   domain.RequestPending,
		DormID:   dormID,
		LesseeID: leeseeID,
		Message:  message,
	}
	err := s.requestRepo.Create(leasingRequest)
	if err != nil {
		return nil, err
	}
	leasingRequest, err = s.requestRepo.GetByID(leasingRequest.ID)
	if err != nil {
		return nil, err
	}
	return leasingRequest, nil
}
func (s *LeasingRequestService) Delete(id uuid.UUID) error {
	err := s.requestRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (s *LeasingRequestService) GetByUserID(id uuid.UUID, role domain.Role, limit, page int) ([]domain.LeasingRequest, int, int, error) {
	leasingRequest, totalPage, totalRows, err := s.requestRepo.GetByUserID(id, limit, page, role)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingRequest, totalPage, totalRows, nil
}
func (s *LeasingRequestService) Approve(id, userId uuid.UUID, isAdmin bool) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if leasingRequest.Status != domain.RequestPending {
		return apperror.BadRequestError(errors.New("request is not in the pending status"), "request is not in the pending status")
	}
	if userId != leasingRequest.Dorm.OwnerID && !isAdmin {
		return apperror.UnauthorizedError(errors.New("user is unauthorized"), "user is unauthorized")
	}
	leasingRequest.End = time.Now()
	leasingRequest.Status = domain.RequestAccepted
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	contract := &domain.Contract{LesseeID: leasingRequest.LesseeID, DormID: leasingRequest.DormID}
	err = s.contractRepo.Create(contract)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingRequestService) Reject(id, userId uuid.UUID, isAdmin bool) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if leasingRequest.Status != domain.RequestPending {
		return apperror.BadRequestError(errors.New("request is not in the pending status"), "request is not in the pending status")
	}
	if userId != leasingRequest.Dorm.OwnerID && !isAdmin {
		return apperror.UnauthorizedError(errors.New("user is unauthorized"), "user is unauthorized")
	}
	leasingRequest.End = time.Now()
	leasingRequest.Status = domain.RequestRejected
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingRequestService) Cancel(id, userId uuid.UUID, isAdmin bool) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if leasingRequest.Status != domain.RequestPending {
		return apperror.BadRequestError(errors.New("request is not in the pending status"), "request is not in the pending status")
	}
	if userId != leasingRequest.LesseeID && !isAdmin {
		return apperror.UnauthorizedError(errors.New("user is unauthorized"), "user is unauthorized")
	}
	leasingRequest.End = time.Now()
	leasingRequest.Status = domain.RequestCanceled
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingRequestService) GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error) {
	leasingRequest, totalPage, totalRows, err := s.requestRepo.GetByDormID(id, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingRequest, totalPage, totalRows, nil
}
