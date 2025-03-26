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

type LeasingRequestRepository struct {
	db *database.Database
}

func NewLeasingRequestRepository(db *database.Database) ports.LeasingRequestRepository {
	return &LeasingRequestRepository{db: db}
}

func (d *LeasingRequestRepository) Create(LeasingRequest *domain.LeasingRequest) error {
	if err := d.db.Create(LeasingRequest).Error; err != nil {
		return apperror.InternalServerError(err, "failed to save leasing request to database")
	}

	return nil
}

func (d *LeasingRequestRepository) GetByID(id uuid.UUID) (*domain.LeasingRequest, error) {
	leasingRequest := new(domain.LeasingRequest)
	if err := d.db.Preload("Dorm").Preload("Lessee").Preload("Dorm.Owner").First(leasingRequest, id).Error; err != nil {
		return nil, apperror.NotFoundError(err, "leasing request not found")
	}
	return leasingRequest, nil
}

func (d *LeasingRequestRepository) Update(LeasingRequest *domain.LeasingRequest) error {
	existingRequest, err := d.GetByID(LeasingRequest.ID)
	if err != nil {
		return apperror.NotFoundError(err, "leasing request not found")
	}
	err = d.db.Model(existingRequest).Updates(LeasingRequest).Error
	if err != nil {
		return apperror.InternalServerError(err, "failed to update leasing request")
	}
	return nil
}
func (d *LeasingRequestRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.LeasingRequest{}, id).Error; err != nil {
		return apperror.InternalServerError(err, "failed to delete leasing request")
	}
	return nil
}
func (d *LeasingRequestRepository) GetByUserID(id uuid.UUID, limit, page int, role domain.Role) ([]domain.LeasingRequest, int, int, error) {
	var leasingRequest []domain.LeasingRequest
	var query *gorm.DB
	if role == domain.LesseeRole {
		query = d.db.Preload("Dorm").
			Preload("Lessee").
			Preload("Dorm.Owner").
			Where("lessee_id = ?", id)
	} else if role == domain.LessorRole {
		query = d.db.Preload("Dorm").
			Preload("Lessee").
			Preload("Dorm.Owner").
			Joins("JOIN dorms ON dorms.id = leasing_requests.dorm_id").
			Where("owner_id = ?", id)
	} else {
		query = d.db.Preload("Dorm").
			Preload("Lessee").
			Preload("Dorm.Owner").
			Joins("LEFT JOIN dorms ON dorms.id = leasing_requests.dorm_id").
			Where("lessee_id = ? OR owner_id = ?", id, id)
	}
	totalPage, totalRows, err := d.db.Paginate(&leasingRequest, query, limit, page, "start DESC")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "leasing request not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get leasing request")
	}

	return leasingRequest, totalPage, totalRows, nil
}

func (d *LeasingRequestRepository) GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingRequest, int, int, error) {
	var leasingRequest []domain.LeasingRequest
	query := d.db.Preload("Dorm").Preload("Dorm.Owner").Preload("Lessee").Where("dorm_id = ?", id)
	totalPage, totalRows, err := d.db.Paginate(&leasingRequest, query, limit, page, "start")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "leasing history not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get leasing history")
	}

	return leasingRequest, totalPage, totalRows, nil
}
