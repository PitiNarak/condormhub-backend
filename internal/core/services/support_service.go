package services

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
)

type SupportService struct {
	repo ports.SupportRepository
}

func NewSupportService(repo ports.SupportRepository) ports.SupportService {
	return &SupportService{repo: repo}
}

func (s *SupportService) Create(support *domain.SupportRequest) error {
	return s.repo.Create(support)
}

func (s *SupportService) GetAll(limit int, page int, userID uuid.UUID, isAdmin bool) ([]domain.SupportRequest, int, int, error) {
	return s.repo.GetAll(limit, page, userID, isAdmin)
}
