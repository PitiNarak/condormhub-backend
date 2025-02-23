package services

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
)

type LeasingHistoryService struct {
	repo ports.LeasingHistoryRepository
}

func NewLeasingHistoryService(repo ports.LeasingHistoryRepository) ports.LeasingHistoryService {
	return &LeasingHistoryService{repo: repo}
}

func (s *LeasingHistoryService) Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error) {
	createTime := time.Now()
	leasingHistory := &domain.LeasingHistory{DormID: dormID, LesseeID: userID, Start: createTime}
	err := s.repo.Create(leasingHistory)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	leasingHistory, err = s.repo.GetByID(leasingHistory.ID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) Delete(id uuid.UUID) error {
	return nil
}
func (s *LeasingHistoryService) GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	leasingHistory, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	leasingHistory, err := s.repo.GetByDormID(id)
	if err != nil {
		return nil, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) SetEndTimestamp(id uuid.UUID) error {
	return nil
}
