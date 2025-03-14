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
	requestRepo ports.LeasingRequestRepository
	dormRepo    ports.DormRepository
}

func NewLeasingRequestService(requestRepo ports.LeasingRequestRepository, dormRepo ports.DormRepository) ports.LeasingRequestService {
	return &LeasingRequestService{requestRepo: requestRepo, dormRepo: dormRepo}
}

func (s *LeasingRequestService) Create(leeseeID uuid.UUID, dormID uuid.UUID) (*domain.LeasingRequest, error) {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return &domain.LeasingRequest{}, err
	}
	createTime := time.Now()
	requestPending := domain.RequestPending
	leasingRequest := &domain.LeasingRequest{Status: &requestPending, DormID: dormID, LesseeID: leeseeID, LessorID: dorm.OwnerID, Start: createTime}
	err = s.requestRepo.Create(leasingRequest)
	if err != nil {
		return &domain.LeasingRequest{}, err
	}
	leasingRequest, err = s.requestRepo.GetByID(leasingRequest.ID)
	if err != nil {
		return &domain.LeasingRequest{}, err
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
func (s *LeasingRequestService) Approve(id, userId uuid.UUID) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if userId != leasingRequest.LessorID {
		return apperror.UnauthorizedError(errors.New("id not in the leasing request"), "id not in the leasing request")
	}
	leasingRequest.End = time.Now()
	requestAccept := domain.RequestAccepted
	leasingRequest.Status = &requestAccept
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingRequestService) Reject(id, userId uuid.UUID) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if userId != leasingRequest.LessorID {
		return apperror.UnauthorizedError(errors.New("id not in the leasing request"), "id not in the leasing request")
	}
	leasingRequest.End = time.Now()
	requestRejected := domain.RequestRejected
	leasingRequest.Status = &requestRejected
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingRequestService) Cancel(id, userId uuid.UUID) error {
	leasingRequest, err := s.requestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if userId != leasingRequest.LesseeID {
		return apperror.UnauthorizedError(errors.New("id not in the leasing request"), "id not in the leasing request")
	}
	leasingRequest.End = time.Now()
	requestCanceled := domain.RequestCanceled
	leasingRequest.Status = &requestCanceled
	err = s.requestRepo.Update(leasingRequest)
	if err != nil {
		return err
	}
	return nil
}
