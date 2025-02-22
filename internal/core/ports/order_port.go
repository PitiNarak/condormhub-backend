package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *domain.Order) *error_handler.ErrorHandler
	GetByID(orderID uuid.UUID) (*domain.Order, *error_handler.ErrorHandler)
	GetAll() ([]*domain.Order, *error_handler.ErrorHandler)
	Update(order *domain.Order) *error_handler.ErrorHandler
	Delete(orderID uuid.UUID) *error_handler.ErrorHandler
}
