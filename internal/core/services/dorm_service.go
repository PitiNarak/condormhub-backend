package services

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
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

func (s *DormService) GetAll(limit, page int) ([]domain.Dorm, int, int, error) {
	dorms, totalPages, totalRows, err := s.dormRepo.GetAll(limit, page)
	if err != nil {
		return nil, totalPages, totalRows, err
	}
	return dorms, totalPages, totalRows, nil
}

func (s *DormService) GetByID(id uuid.UUID) (*domain.Dorm, error) {
	return s.dormRepo.GetByID(id)
}

func (s *DormService) Update(id uuid.UUID, dorm *domain.Dorm) error {
	return s.dormRepo.Update(id, dorm)
}

func (s *DormService) Delete(id uuid.UUID) error {
	return s.dormRepo.Delete(id)
}
