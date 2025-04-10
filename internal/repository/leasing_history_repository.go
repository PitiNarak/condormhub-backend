package repository

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeasingHistoryRepository struct {
	db *database.Database
}

func NewLeasingHistoryRepository(db *database.Database) ports.LeasingHistoryRepository {
	return &LeasingHistoryRepository{db: db}
}

func (d *LeasingHistoryRepository) Create(LeasingHistory *domain.LeasingHistory) error {
	if err := d.db.Create(LeasingHistory).Error; err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "failed to save leasing history to database")
	}
	return nil
}

func (d *LeasingHistoryRepository) GetByID(id uuid.UUID) (*domain.LeasingHistory, error) {
	leasingHistory := new(domain.LeasingHistory)
	if err := d.db.Preload("Dorm").Preload("Lessee").Preload("Orders").Preload("Dorm.Owner").First(leasingHistory, id).Error; err != nil {
		if apperror.IsAppError(err) {
			return nil, err
		}
		return nil, apperror.NotFoundError(err, "leasing history not found")
	}
	return leasingHistory, nil
}

func (d *LeasingHistoryRepository) Update(LeasingHistory *domain.LeasingHistory) error {
	existingHistory, err := d.GetByID(LeasingHistory.ID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.NotFoundError(err, "History not found")
	}
	err = d.db.Model(existingHistory).Updates(LeasingHistory).Error
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Failed to update leasing history")
	}
	if existingHistory.ReviewFlag != LeasingHistory.ReviewFlag {
		err = d.db.Model(existingHistory).UpdateColumn("ReviewFlag", LeasingHistory.ReviewFlag).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *LeasingHistoryRepository) Delete(id uuid.UUID) error {
	// TODO: Cascade delete?
	if err := d.db.Delete(&domain.LeasingHistory{}, id).Error; err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Failed to delete leasing history")
	}
	return nil
}
func (d *LeasingHistoryRepository) GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	var leasingHistory []domain.LeasingHistory
	query := d.db.Preload("Dorm").Preload("Lessee").Preload("Dorm.Images").Preload("Orders").Preload("Dorm.Owner").Where("lessee_id = ?", id)
	totalPage, totalRows, err := d.db.Paginate(&leasingHistory, query, limit, page, "start")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "leasing history not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get leasing history")
	}

	return leasingHistory, totalPage, totalRows, nil
}
func (d *LeasingHistoryRepository) GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	var leasingHistory []domain.LeasingHistory
	query := d.db.Preload("Dorm").Preload("Dorm.Images").Preload("Dorm.Owner").Where("dorm_id = ?", id)
	totalPage, totalRows, err := d.db.Paginate(leasingHistory, query, limit, page, "start")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "leasing history not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get leasing history")
	}

	return leasingHistory, totalPage, totalRows, nil
}
