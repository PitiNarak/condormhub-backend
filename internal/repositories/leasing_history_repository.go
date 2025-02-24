package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/pagination"
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
	existingHistory, err := d.GetByID(LeasingHistory.ID)
	if err != nil {
		return errorHandler.NotFoundError(err, "History not found")
	}
	err = d.db.Model(existingHistory).Updates(LeasingHistory).Error
	if err != nil {
		return errorHandler.InternalServerError(err, "Failed to update leasing history")
	}
	return nil
}
func (d *LeasingHistoryRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.LeasingHistory{}, id).Error; err != nil {
		return errorHandler.InternalServerError(err, "Failed to delete leasing history")
	}
	return nil
}
func (d *LeasingHistoryRepository) GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	var leasingHistory []domain.LeasingHistory
	scope, totalPage, totalRows, err := pagination.Paginate(leasingHistory, d.db, limit, page, "")
	if err != nil {
		return nil, 0, 0, err
	}
	err = d.db.Scopes(scope).Preload("Dorm").Preload("Lessee").Preload("Orders").Preload("Dorm.Owner").Where("lessee_id = ?", id).Find(&leasingHistory).Error
	if err != nil {
		return nil, 0, 0, errorHandler.NotFoundError(err, "leasing history not found")
	}
	return leasingHistory, totalPage, totalRows, nil
}
func (d *LeasingHistoryRepository) GetByDormID(id uuid.UUID) ([]domain.LeasingHistory, error) {
	var leasingHistory []domain.LeasingHistory
	err := d.db.Preload("Dorm").Preload("Dorm.Owner").
		Where("dorm_id = ?", id).Find(&leasingHistory).Error

	if err != nil {
		return nil, errorHandler.NotFoundError(err, "leasing history not found")
	}
	return leasingHistory, nil
}
