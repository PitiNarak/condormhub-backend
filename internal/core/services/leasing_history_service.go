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
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	leasingHistory := &domain.LeasingHistory{DormID: dormID, LesseeID: userID, Start: createTime, Price: dorm.Price}
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
	err := s.historyRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (s *LeasingHistoryService) GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	leasingHistory, totalPage, totalRows, err := s.historyRepo.GetByUserID(id, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingHistory, totalPage, totalRows, nil
}
func (s *LeasingHistoryService) GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	leasingHistory, totalPage, totalRows, err := s.historyRepo.GetByDormID(id, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingHistory, totalPage, totalRows, nil
}
func (s *LeasingHistoryService) SetEndTimestamp(id uuid.UUID) error {
	leasingHistory, err := s.historyRepo.GetByID(id)
	if err != nil {
		return err
	}
	leasingHistory.End = time.Now()
	err = s.historyRepo.Update(leasingHistory)
	if err != nil {
		return err
	}
	return nil
}
