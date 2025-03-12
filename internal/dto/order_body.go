package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderRequestBody struct {
	LeasingHistoryID uuid.UUID `json:"leasingHistoryId" validate:"required"`
}

type OrderResponseBody struct {
	ID                uuid.UUID              `json:"id"`
	CreateAt          time.Time              `json:"-"`
	UpdateAt          time.Time              `json:"-"`
	Type              string                 `json:"type"`
	Price             int64                  `json:"price"`
	Transactions      []*TransactionResponse `json:"-"`
	PaidTransaction   *TransactionResponse   `json:"paidTransaction"`
	PaidTransactionID string                 `json:"-"`
	LeasingHistory    *LeasingHistory        `json:"-"`
	LeasingHistoryID  uuid.UUID              `json:"-"`
}
