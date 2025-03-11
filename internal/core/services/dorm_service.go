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

func checkPermission(ownerID uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	if ownerID != userID && !isAdmin {
		return errors.New("unauthorized action")
	}
	return nil
}

func (s *DormService) Create(userRole domain.Role, dorm *domain.Dorm) error {
	if userRole != domain.AdminRole && userRole != domain.LessorRole {
		return apperror.ForbiddenError(errors.New("unauthorized action"), "You do not have permission to create a dorm")
	}
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

	if err = checkPermission(dorm.OwnerID, userID, isAdmin); err != nil {
		return nil, apperror.ForbiddenError(err, "You do not have permission to update this dorm")
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

func (s *DormService) Delete(userID uuid.UUID, isAdmin bool, dormID uuid.UUID) error {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return err
	}

	if err := checkPermission(dorm.OwnerID, userID, isAdmin); err != nil {
		return apperror.ForbiddenError(err, "You do not have permission to delete this dorm")
	}

	return s.dormRepo.Delete(dormID)
}
