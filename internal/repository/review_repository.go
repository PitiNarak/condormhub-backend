package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *database.Database
}

func NewReviewRepository(db *database.Database) ports.ReviewRepository {
	return &ReviewRepository{db: db}
}

func (d *ReviewRepository) Create(Review *domain.Review) error {
	if err := d.db.Create(Review).Error; err != nil {
		return apperror.InternalServerError(err, "failed to save review to database")
	}

	return nil
}

func (d *ReviewRepository) Get(id uuid.UUID) (*domain.Review, error) {
	leasingRequest := new(domain.Review)
	if err := d.db.First(leasingRequest, id).Error; err != nil {
		return nil, apperror.NotFoundError(err, "review not found")
	}
	return leasingRequest, nil
}

func (d *ReviewRepository) Update(Review *domain.Review) error {
	existingRequest, err := d.Get(Review.ID)
	if err != nil {
		return apperror.NotFoundError(err, "review not found")
	}
	err = d.db.Model(existingRequest).Updates(Review).Error
	if err != nil {
		return apperror.InternalServerError(err, "failed to update review")
	}
	return nil
}
func (d *ReviewRepository) Delete(id uuid.UUID) error {
	if err := d.db.Delete(&domain.LeasingRequest{}, id).Error; err != nil {
		return apperror.InternalServerError(err, "failed to delete review")
	}
	return nil
}
