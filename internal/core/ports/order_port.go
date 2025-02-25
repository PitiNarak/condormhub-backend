package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *domain.Order) *errorHandler.ErrorHandler
	GetByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetAll() ([]*domain.Order, *errorHandler.ErrorHandler)
	Update(order *domain.Order) *errorHandler.ErrorHandler
	Delete(orderID uuid.UUID) *errorHandler.ErrorHandler
}
