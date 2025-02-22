package services

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
)

type LeasingHistoryService struct {
	repo *ports.DormRepository
}

func NewLeasingHistoryService(repo *ports.DormRepository) ports.LeasingHistoryService {
	return &LeasingHistoryService{repo: repo}
}

func (s *LeasingHistoryService) Create(LeasingHistory *domain.LeasingHistory) error
func (s *LeasingHistoryService) Update(LeasingHistory *domain.LeasingHistory) error
func (s *LeasingHistoryService) Delete(id uuid.UUID) error
func (s *LeasingHistoryService) GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error)
func (s *LeasingHistoryService) GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error)
func (s *LeasingHistoryService) AddNewOrder(id uuid.UUID, order *domain.Order) error
func (s *LeasingHistoryService) PatchEndTimestamp(id uuid.UUID, endTime time.Time) error
