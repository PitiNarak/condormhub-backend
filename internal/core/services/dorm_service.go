package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
)

type DormService struct {
	dormRepo ports.DormRepository
}

func NewDormService(repo ports.DormRepository) ports.DormService {
	return &DormService{dormRepo: repo}
}

func (s *DormService) Create(dorm *domain.Dorm) error {
	return s.dormRepo.Create(dorm)
}

func (s *DormService) GetAll(limit int, page int) ([]domain.Dorm, int, int, error) {
	dorms, totalPages, totalRows, err := s.dormRepo.GetAll(limit, page)
	if err != nil {
		return nil, totalPages, totalRows, err
	}
	return dorms, totalPages, totalRows, nil
}

func (s *DormService) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	return s.dormRepo.GetByID(id)
}

func (s *DormService) Update(userID uuid.UUID, isAdmin bool, dormID uuid.UUID, updateData *dto.DormRequestBody) (*domain.Dorm, error) {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return nil, err
	}

	if dorm.OwnerID != userID && !isAdmin {
		return nil, apperror.ForbiddenError(errors.New("unauthorized to update this room"), "unauthorized to update this room")
	}

	dorm.Name = updateData.Name
	dorm.Size = updateData.Size
	dorm.Bedrooms = updateData.Bedrooms
	dorm.Bathrooms = updateData.Bathrooms
	dorm.Address = updateData.Address
	dorm.Price = updateData.Price
	dorm.Description = updateData.Description

	if err := s.dormRepo.Update(dormID, dorm); err != nil {
		return nil, err
	}

	return dorm, nil
}

func (s *DormService) Delete(id uuid.UUID) error {
	return s.dormRepo.Delete(id)
}
