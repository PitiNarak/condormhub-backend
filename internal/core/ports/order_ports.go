package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	Webhook(c *fiber.Ctx) error
}

type OrderService interface {
	CreateOrder(orderType domain.OrderType, dormitoryID uuid.UUID, lessorID uuid.UUID, lesseeID uuid.UUID) (*domain.Order, *error_handler.ErrorHandler)
	UpdateOrderStatus(event stripe.Event) *error_handler.ErrorHandler
}

type OrderRepository interface {
	Create(order *domain.Order) *error_handler.ErrorHandler
	Update(order *domain.Order) *error_handler.ErrorHandler
	// GetById(orderId uuid.UUID) (*domain.Order, *error_handler.ErrorHandler)
	// GetByLessorId(lessorId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler)
	// GetByLesseeId(lesseeId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler)
	// GetByDormitoryId(dormitoryId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler)
	// GetAll() ([]domain.Order, *error_handler.ErrorHandler)
}
