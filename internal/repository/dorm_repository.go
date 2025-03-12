package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type DormRepository struct {
	db *database.Database
}

func NewDormRepository(db *database.Database) ports.DormRepository {
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

func (d *DormRepository) GetAll(limit int, page int) ([]domain.Dorm, int, int, error) {
	var dorms []domain.Dorm
	query := d.db.Preload("Owner")

	totalPages, totalRows, err := d.db.Paginate(&dorms, query, limit, page, "create_at DESC")
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve dorms")
	}

	return dorms, totalPages, totalRows, nil
}

func (d *DormRepository) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	dorm := new(domain.Dorm)
	if err := d.db.Preload("Owner").First(dorm, id).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Dorm not found")
	}
	return dorm, nil
}

func (d *DormRepository) Update(id uuid.UUID, dorm dto.DormUpdateRequestBody) error {
	updateData := make(map[string]interface{})
	if dorm.Name != "" {
		updateData["name"] = dorm.Name
	}
	if dorm.Size != 0 {
		updateData["size"] = dorm.Size
	}
	if dorm.Bedrooms != 0 {
		updateData["bedrooms"] = dorm.Bedrooms
	}
	if dorm.Bathrooms != 0 {
		updateData["bathrooms"] = dorm.Bathrooms
	}
	if dorm.Price != 0 {
		updateData["price"] = dorm.Price
	}
	if dorm.Description != "" {
		updateData["description"] = dorm.Description
	}
	if dorm.Address.District != "" {
		updateData["district"] = dorm.Address.District
	}
	if dorm.Address.Subdistrict != "" {
		updateData["subdistrict"] = dorm.Address.Subdistrict
	}
	if dorm.Address.Province != "" {
		updateData["province"] = dorm.Address.Province
	}
	if dorm.Address.Zipcode != "" {
		updateData["zipcode"] = dorm.Address.Zipcode
	}

	res := d.db.Model(&domain.Dorm{}).Where("id = ?", id).Updates(updateData)
	if res.Error != nil {
		return apperror.InternalServerError(res.Error, "Failed to update room")
	}

	return nil
}

func (d *DormRepository) SaveDormImage(dormImage *domain.DormImage) error {
	if err := d.db.Create(dormImage).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to save dorm's image to database")
	}
	return nil
}
