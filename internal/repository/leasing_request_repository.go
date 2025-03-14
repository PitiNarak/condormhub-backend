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
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "failed to save leasing request to database")
	}
	return nil
}

func (d *LeasingRequestRepository) GetByID(id uuid.UUID) (*domain.LeasingRequest, error) {
	leasingRequest := new(domain.LeasingRequest)
	if err := d.db.Preload("Dorm").Preload("Lessee").Preload("Orders").Preload("Dorm.Owner").First(leasingRequest, id).Error; err != nil {
		if apperror.IsAppError(err) {
			return nil, err
		}
		return nil, apperror.NotFoundError(err, "leasing request not found")
	}
	return leasingRequest, nil
}

func (d *LeasingRequestRepository) Update(LeasingRequest *domain.LeasingRequest) error {
	existingRequest, err := d.GetByID(LeasingRequest.ID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.NotFoundError(err, "leasing request not found")
	}
	err = d.db.Model(existingRequest).Updates(LeasingRequest).Error
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "failed to update leasing request")
	}
	return nil
}
func (d *LeasingRequestRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.LeasingRequest{}, id).Error; err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "failed to delete leasing request")
	}
	return nil
}
func (d *LeasingRequestRepository) GetByUserID(id uuid.UUID, limit, page int, role domain.Role) ([]domain.LeasingRequest, int, int, error) {
	var leasingRequest []domain.LeasingRequest
	query := new(gorm.DB)
	if role == domain.LesseeRole {
		query = d.db.Preload("Dorm").Preload("Lessee").Preload("Lessor").Where("lessee_id = ?", id)
	} else {
		query = d.db.Preload("Dorm").Preload("Lessee").Preload("Lessor").Where("lessor_id = ?", id)
	}
	totalPage, totalRows, err := d.db.Paginate(&leasingRequest, query, limit, page, "start")

	if err != nil {
		if apperror.IsAppError(err) {
			return nil, 0, 0, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "leasing request not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get leasing request")
	}

	return leasingRequest, totalPage, totalRows, nil
}
