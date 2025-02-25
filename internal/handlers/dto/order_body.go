package dto

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type OrderRequestBody struct {
	LeasingHistoryID uuid.UUID `json:"leasingHistoryId" validate:"required"`
}

type OrderResponseBody struct {
	ID                uuid.UUID              `json:"id"`
	CreateAt          time.Time              `json:"-"`
	UpdateAt          time.Time              `json:"-"`
	Type              domain.OrderType       `json:"type"`
	Price             int64                  `json:"price"`
	Transactions      []*domain.Transaction  `json:"transactions"`
	PaidTransaction   *domain.Transaction    `json:"paidTransaction"`
	PaidTransactionID string                 `json:"-"`
	LeasingHistory    *domain.LeasingHistory `json:"-"`
	LeasingHistoryID  uuid.UUID              `json:"-"`
}
