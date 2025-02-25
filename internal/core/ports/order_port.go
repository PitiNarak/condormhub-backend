package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *domain.Order) *errorHandler.ErrorHandler
	GetByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetUnpaidByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler)
	Update(order *domain.Order) *errorHandler.ErrorHandler
	Delete(orderID uuid.UUID) *errorHandler.ErrorHandler
}

type OrderService interface {
	CreateOrder(leasingHistoryID uuid.UUID) *errorHandler.ErrorHandler
	GetOrderByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler)
	UpdateOrder(order *domain.Order) *errorHandler.ErrorHandler
	DeleteOrder(orderID uuid.UUID) *errorHandler.ErrorHandler
}

type OrderHandler interface {
	CreateOrder(order *domain.Order) *errorHandler.ErrorHandler
	GetOrder(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler)
	UpdateOrder(order *domain.Order) *errorHandler.ErrorHandler
	DeleteOrder(orderID uuid.UUID) *errorHandler.ErrorHandler
}
