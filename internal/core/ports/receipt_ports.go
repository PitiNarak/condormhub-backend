package ports

import (
	"context"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReceiptRepository interface {
	Create(receipt *domain.Receipt) error
	GetByUserID(userID uuid.UUID, limit int, page int) ([]domain.Receipt, int, int, error)
}

type ReceiptService interface {
	Create(c context.Context, ownerID uuid.UUID, transaction domain.Transaction) error
	GetByUserID(userID uuid.UUID, limit int, page int) ([]domain.Receipt, int, int, error)
	GetUrl(c context.Context, receipt domain.Receipt) (string, error)
}

type ReceiptHandler interface {
	GetByUserID(c *fiber.Ctx) error
}
