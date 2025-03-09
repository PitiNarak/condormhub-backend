package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/databases"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type DormRepository struct {
	db *databases.Database
}

func NewDormRepository(db *databases.Database) ports.DormRepository {
	return &DormRepository{db: db}
}

func (d *DormRepository) Create(dorm *domain.Dorm) error {
	if err := d.db.Create(dorm).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save dorm to database")
	}
	return nil
}

func (d *DormRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.Dorm{}, id).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete dorm")
	}
	return nil
}

func (d *DormRepository) GetAll() ([]domain.Dorm, error) {
	var dorms []domain.Dorm
	if err := d.db.Preload("Owner").Find(&dorms).Error; err != nil {
		return nil, apperror.InternalServerError(err, "Failed to retrieve dorms")
	}
	return dorms, nil
}

func (d *DormRepository) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	dorm := new(domain.Dorm)
	if err := d.db.Preload("Owner").First(dorm, id).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Dorm not found")
	}
	return dorm, nil
}

func (d *DormRepository) Update(id uuid.UUID, dorm *domain.Dorm) error {
	existingDorm, err := d.GetByID(id)
	if err != nil {
		return apperror.NotFoundError(err, "Dorm not found")
	}

	err = d.db.Model(existingDorm).Updates(dorm).Error
	if err != nil {
		return apperror.InternalServerError(err, "Failed to update room")
	}

	return nil
}
