package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DormRepository struct {
	db *gorm.DB
}

func NewDormRepository(db *gorm.DB) ports.DormRepository {
	return &DormRepository{db: db}
}

func (d *DormRepository) Create(dorm *domain.Dorm) error {
	if err := d.db.Create(dorm).Error; err != nil {
		return errorHandler.InternalServerError(err, "Failed to save dorm to database")
	}
	return nil
}

func (d *DormRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.Dorm{}, id).Error; err != nil {
		return errorHandler.InternalServerError(err, "Failed to delete dorm")
	}
	return nil
}

func (d *DormRepository) GetAll() ([]domain.Dorm, error) {
	var dorms []domain.Dorm
	if err := d.db.Preload("Owner").Find(&dorms).Error; err != nil {
		return nil, errorHandler.InternalServerError(err, "Failed to retrieve dorms")
	}
	return dorms, nil
}

func (d *DormRepository) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	dorm := new(domain.Dorm)
	if err := d.db.Preload("Owner").First(dorm, id).Error; err != nil {
		return nil, errorHandler.NotFoundError(err, "Dorm not found")
	}
	return dorm, nil
}

func (d *DormRepository) Update(id uuid.UUID, dorm *domain.Dorm) error {
	existingDorm, err := d.GetByID(id)
	if err != nil {
		return errorHandler.NotFoundError(err, "Dorm not found")
	}

	err = d.db.Model(existingDorm).Updates(dorm).Error
	if err != nil {
		return errorHandler.InternalServerError(err, "Failed to update room")
	}

	return nil
}
