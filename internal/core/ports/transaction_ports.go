package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	Webhook(c *fiber.Ctx) error
}

type TransactionService interface {
	CreateTransaction(orderID uuid.UUID) (*domain.Transaction, *string, error)
	UpdateTransactionStatus(event stripe.Event) error
}

type TransactionRepository interface {
	Create(order *domain.Transaction) error
	GetByID(id string) (domain.Transaction, error)
	Update(order *domain.Transaction) error
}
