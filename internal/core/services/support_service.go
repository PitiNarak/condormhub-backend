package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
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

func (s *SupportService) UpdateStatus(id uuid.UUID, status domain.SupportStatus) (*domain.SupportRequest, error) {
	support, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NotFoundError(err, "Support request not found")
	}

	if status != domain.ProblemOpen && status != domain.ProblemInProgress && status != domain.ProblemResolved {
		return nil, apperror.UnprocessableEntityError(errors.New("invalid status value"), "Invalid status value")
	}

	support.Status = status

	return support, s.repo.UpdateStatus(id, status)
}
