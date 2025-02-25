package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/gofiber/fiber/v2"
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
	CreateOrder(leasingHistoryID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetOrderByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler)
	GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler)
	UpdateOrder(order *domain.Order) *errorHandler.ErrorHandler
	DeleteOrder(orderID uuid.UUID) *errorHandler.ErrorHandler
}

type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	GetOrderByID(c *fiber.Ctx) error
	// GetUnpaidOrderByUserID(c *fiber.Ctx) ([]domain.Order, int, int, *errorHandler.ErrorHandler)
	// UpdateOrder(c *fiber.Ctx) *errorHandler.ErrorHandler
	// DeleteOrder(c *fiber.Ctx) *errorHandler.ErrorHandler
}
