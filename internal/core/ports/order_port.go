package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(orderID uuid.UUID) (*domain.Order, error)
	GetUnpaidByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, error)
	Update(order *domain.Order) error
	Delete(orderID uuid.UUID) error
}

type OrderService interface {
	CreateOrder(leasingHistoryID uuid.UUID) (*domain.Order, error)
	GetOrderByID(orderID uuid.UUID) (*domain.Order, error)
	GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(orderID uuid.UUID) error
}

type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	GetOrderByID(c *fiber.Ctx) error
	GetUnpaidOrderByUserID(c *fiber.Ctx) error
	GetMyUnpaidOrder(c *fiber.Ctx) error
	// UpdateOrder(c *fiber.Ctx) error
	// DeleteOrder(c *fiber.Ctx) error
}
