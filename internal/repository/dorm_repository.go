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
	// TODO: Cascade delete for all field that reference to dorm
	if err := d.db.Select("Images").Delete(&domain.Dorm{ID: id}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete dorm")
	}
	return nil
}

func (d *DormRepository) GetAll(
	limit int, page int,
	search string,
	minPrice int, maxPrice int,
	district string,
	subdistrict string,
	province string,
	zipcode string,
) ([]domain.Dorm, int, int, error) {
	var dorms []domain.Dorm
	query := d.db.Preload("Owner").Preload("Images")

	if search != "" {
		regex := "%" + search + "%"
		query.Where("name ILIKE ? OR province ILIKE ? OR district ILIKE ? OR subdistrict ILIKE ? OR zipcode ILIKE ?", regex, regex, regex, regex, regex)
	}

	if minPrice != -1 {
		query.Where("price >= ?", minPrice)
	}

	if maxPrice != -1 {
		query.Where("price <= ?", maxPrice)
	}

	if district != "" {
		query.Where("district ILIKE ?", "%"+district+"%")
	}

	if subdistrict != "" {
		query.Where("subdistrict ILIKE ?", "%"+subdistrict+"%")
	}

	if province != "" {
		query.Where("province ILIKE ?", "%"+province+"%")
	}

	if zipcode != "" {
		query.Where("zipcode ILIKE ?", "%"+zipcode+"%")
	}

	totalPages, totalRows, err := d.db.Paginate(&dorms, query, limit, page, "create_at DESC")
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve dorms")
	}

	return dorms, totalPages, totalRows, nil
}

func (d *DormRepository) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	dorm := new(domain.Dorm)
	if err := d.db.Preload("Owner").Preload("Images").First(dorm, id).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Dorm not found")
	}
	return dorm, nil
}

func (d *DormRepository) Update(id uuid.UUID, dorm dto.DormUpdateRequestBody) error {
	updatedDorm := domain.Dorm{
		Name:        dorm.Name,
		Size:        dorm.Size,
		Bedrooms:    dorm.Bedrooms,
		Bathrooms:   dorm.Bathrooms,
		Address:     domain.Address(dorm.Address),
		Price:       dorm.Price,
		Description: dorm.Description,
	}

	res := d.db.Model(&domain.Dorm{}).Where("id = ?", id).Updates(updatedDorm)
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

func (d *DormRepository) GetByOwnerID(ownerID uuid.UUID, limit int, page int) ([]domain.Dorm, int, int, error) {
	var dorms []domain.Dorm
	query := d.db.Preload("Owner").Preload("Images").Where("owner_id = ?", ownerID)

	totalPages, totalRows, err := d.db.Paginate(&dorms, query, limit, page, "create_at DESC")
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to retrieve dorms")
	}

	return dorms, totalPages, totalRows, nil
}

func (d *DormRepository) DeleteImageByKey(imageKey string) error {
	if err := d.db.Where("image_key = ?", imageKey).Delete(&domain.DormImage{}).Error; err != nil {
		return apperror.InternalServerError(err, "Failed to delete image")
	}
	return nil
}

func (d *DormRepository) GetImageByKey(imageKey string) (*domain.DormImage, error) {
	dormImage := new(domain.DormImage)
	if err := d.db.Where("image_key = ?", imageKey).First(dormImage).Error; err != nil {
		return nil, apperror.NotFoundError(err, "Image not found")
	}
	return dormImage, nil
}
