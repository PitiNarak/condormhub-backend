package services

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
)

type SupportService struct {
	repo ports.SupportRepository
}

func NewSupportRepository(repo ports.SupportRepository) ports.SupportService {
	return &SupportService{repo: repo}
}

func (s *SupportService) Create(support *domain.SupportRequest) error {
	return s.repo.Create(support)
}
