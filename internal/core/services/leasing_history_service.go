package services

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
)

type LeasingHistoryService struct {
	historyRepo ports.LeasingHistoryRepository
	dormRepo    ports.DormRepository
}

func NewLeasingHistoryService(historyRepo ports.LeasingHistoryRepository, dormRepo ports.DormRepository) ports.LeasingHistoryService {
	return &LeasingHistoryService{historyRepo: historyRepo, dormRepo: dormRepo}
}

func (s *LeasingHistoryService) Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error) {
	_, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	createTime := time.Now()
	leasingHistory := &domain.LeasingHistory{DormID: dormID, LesseeID: userID, Start: createTime}
	err = s.historyRepo.Create(leasingHistory)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	leasingHistory, err = s.historyRepo.GetByID(leasingHistory.ID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) Delete(id uuid.UUID) error {
	return nil
}
func (s *LeasingHistoryService) GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	leasingHistory, err := s.historyRepo.GetByUserID(id)
	if err != nil {
		return nil, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	leasingHistory, err := s.historyRepo.GetByDormID(id)
	if err != nil {
		return nil, err
	}
	return leasingHistory, nil
}
func (s *LeasingHistoryService) SetEndTimestamp(id uuid.UUID) error {
	return nil
}
