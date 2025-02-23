package repositories

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
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
	if err := d.db.Create(LeasingHistory).Error; err != nil {
		return errorHandler.InternalServerError(err, "failed to save leasing history to database")
	}
	return nil
}

func (d *LeasingHistoryRepository) GetByID(id uuid.UUID) (*domain.LeasingHistory, error) {
	leasingHistory := new(domain.LeasingHistory)
	if err := d.db.Preload("Dorm").Preload("Lessee").Preload("Orders").Preload("Dorm.Owner").First(leasingHistory, id).Error; err != nil {
		return nil, errorHandler.NotFoundError(err, "leasing history not found")
	}
	return leasingHistory, nil
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
func (d *LeasingHistoryRepository) PatchEndTimestamp(id uuid.UUID, endTime time.Time) error {
	return nil
}
