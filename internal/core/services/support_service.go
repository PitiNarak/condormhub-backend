package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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

func (s *SupportService) UpdateStatus(id uuid.UUID, status domain.SupportStatus) error {
	if status != domain.ProblemOpen && status != domain.ProblemInProgress && status != domain.ProblemResolved {
		return apperror.BadRequestError(errors.New("invalid status value"), "Invalid status value")
	}
	return s.repo.UpdateStatus(id, status)
}
