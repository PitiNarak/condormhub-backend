package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type DormRepository interface {
	Create(dorm *domain.Dorm) error
	GetAll() ([]domain.Dorm, error)
	GetByID(id uuid.UUID) (*domain.Dorm, error)
	Update(dorm *domain.Dorm) error
	Delete(id uuid.UUID) error
}
