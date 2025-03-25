package ports

import (
	"context"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReceiptRepository interface {
	Create(receipt *domain.Receipt) error
}

type ReceiptService interface {
	Create(ctx context.Context, ownerID uuid.UUID, transactionID uuid.UUID) (*domain.Receipt, string, error)
}

type ReceiptHandler interface {
	Create(c *fiber.Ctx) error
}
