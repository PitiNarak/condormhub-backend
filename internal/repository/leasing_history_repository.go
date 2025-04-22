package repository

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
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
	if err := d.db.
		Preload("Dorm").
		Preload("Lessee").
		Preload("Orders").
		Preload("Dorm.Owner").
		Preload("Images").
		First(leasingHistory, id).Error; err != nil {
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
	return nil
}

func (d *LeasingHistoryRepository) SaveReviewImage(reviewImage *domain.ReviewImage) error {
	if err := d.db.Create(reviewImage).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save review's image to database")
	}
	return nil
}

func (d *LeasingHistoryRepository) DeleteImageByKey(imageKey string) error {
	if err := d.db.Where("image_key = ?", imageKey).Delete(&domain.ReviewImage{}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete image")
	}
	return nil
}

func (d *LeasingHistoryRepository) GetImageByKey(imageKey string) (*domain.ReviewImage, error) {
	reviewImage := new(domain.ReviewImage)
	if err := d.db.Where("image_key = ?", imageKey).First(reviewImage).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Image not found")
	}
	return reviewImage, nil
}

func (d *LeasingHistoryRepository) DeleteReview(leasingHistory *domain.LeasingHistory) error {
	err := d.db.Model(leasingHistory).Update("review_flag", false).Error
	if err != nil {
		return apperror.InternalServerError(err, "Failed to delete review")
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
	query := d.db.Preload("Dorm").
		Preload("Lessee").
		Preload("Orders").
		Preload("Dorm.Owner").
		Preload("Images").
		Where("lessee_id = ?", id)
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
	query := d.db.Preload("Dorm").
		Preload("Dorm.Owner").
		Preload("Images").
		Where("dorm_id = ?", id)
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

func (d *LeasingHistoryRepository) GetReviewByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	var reviews []domain.LeasingHistory
	query := d.db.Preload("Lessee").
		Preload("Dorm").
		Preload("Images").
		Where("review_flag = ?", true).
		Where("dorm_id = ?", id)
	totalPage, totalRows, err := d.db.Paginate(&reviews, query, limit, page, "start")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "Review not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get reviews")
	}

	return reviews, totalPage, totalRows, nil
}

func (d *LeasingHistoryRepository) GetReportedReviews(limit int, page int) ([]domain.LeasingHistory, int, int, error) {
	var reviews []domain.LeasingHistory
	query := d.db.Preload("Lessee").Preload("Images").Where("report_flag = ?", true).Where("review_flag = ?", true)
	totalPages, totalRows, err := d.db.Paginate(&reviews, query, limit, page, "id")
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve reviews")
	}
	return reviews, totalPages, totalRows, nil
}
