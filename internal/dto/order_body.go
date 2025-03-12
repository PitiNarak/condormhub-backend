package dto

import (
	"github.com/google/uuid"
)

type OrderRequestBody struct {
	LeasingHistoryID uuid.UUID `json:"leasingHistoryId" validate:"required"`
}

type OrderResponseBody struct {
	ID              uuid.UUID            `json:"id"`
	Type            string               `json:"type"`
	Price           int64                `json:"price"`
	PaidTransaction *TransactionResponse `json:"paidTransaction"`
}
