package repositories

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeasingHistoryRepository struct {
	db *gorm.DB
}

func NewLeasingHistoryRepository(db *gorm.DB) ports.LeasingHistoryRepository {
	return &LeasingHistoryRepository{db: db}
}

func (d *LeasingHistoryRepository) Create(LeasingHistory *domain.LeasingHistory) error {
	return nil
}
func (d *LeasingHistoryRepository) Update(LeasingHistory *domain.LeasingHistory) error {
	return nil
}
func (d *LeasingHistoryRepository) Delete(id uuid.UUID) error {
	return nil
}
func (d *LeasingHistoryRepository) GetByUserID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	return []domain.LeasingHistory{}, nil
}
func (d *LeasingHistoryRepository) GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	return []domain.LeasingHistory{}, nil
}
func (d *LeasingHistoryRepository) AddNewOrder(id uuid.UUID, order *domain.Order) error {
	return nil
}
func (d *LeasingHistoryRepository) PatchEndTimestamp(id uuid.UUID, endTime time.Time) error {
	return nil
}
