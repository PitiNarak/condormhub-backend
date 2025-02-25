package ports

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	Webhook(c *fiber.Ctx) error
}

type TransactionService interface {
	CreateTransaction(orderID uuid.UUID) (*domain.Transaction, *string, *errorHandler.ErrorHandler)
	UpdateTransactionStatus(event stripe.Event) *errorHandler.ErrorHandler
}

type TransactionRepository interface {
	Create(order *domain.Transaction) *errorHandler.ErrorHandler
	Update(order *domain.Transaction) *errorHandler.ErrorHandler
}
