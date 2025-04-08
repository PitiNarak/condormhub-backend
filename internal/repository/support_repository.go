package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
)

type SupportRepository struct {
	db *database.Database
}

func NewSupportRepository(db *database.Database) ports.SupportRepository {
	return &SupportRepository{db: db}
}

func (s *SupportRepository) Create(support *domain.SupportRequest) error {
	if err := s.db.Create(support).Error; err != nil {
		return apperror.InternalServerError(err, "Could not submit support request")
	}
	return nil
}
