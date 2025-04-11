package services

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
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

func (s *LeasingHistoryService) CreateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error) {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateUserForReview(user, history, true)
	if err != nil {
		return nil, err
	}
	review := domain.Review{
		Message: Message,
		Rate:    Rate,
	}
	history.Review = review
	history.ReviewFlag = true
	err = s.historyRepo.Update(history)
	if err != nil {
		return nil, err
	}
	history, err = s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &history.Review, nil
}

func (s *LeasingHistoryService) UpdateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error) {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateUserForReview(user, history, false)
	if err != nil {
		return nil, err
	}
	review := domain.Review{
		Message: Message,
		Rate:    Rate,
	}
	history.Review = review
	err = s.historyRepo.Update(history)
	if err != nil {
		return nil, err
	}
	history, err = s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &history.Review, nil
}

func (s *LeasingHistoryService) DeleteReview(user *domain.User, id uuid.UUID) error {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return err
	}
	err = utils.ValidateUserForReview(user, history, false)
	if err != nil {
		return err
	}
	history.ReviewFlag = false
	err = s.historyRepo.DeleteReview(history)
	if err != nil {
		return err
	}
	return nil
}
