package dto

import (
	"time"

	"github.com/google/uuid"
)

type TransactionRequestBody struct {
	OrderID uuid.UUID `json:"orderID"`
}

type CreateTransactionResponseBody struct {
	CheckoutUrl string `json:"checkoutUrl"`
}

type TransactionResponse struct {
	ID            string            `json:"id"`
	SessionStatus string            `json:"status"`
	CreateAt      time.Time         `json:"createAt"`
	UpdateAt      time.Time         `json:"updateAt"`
	Price         int64             `json:"price"`
	Order         OrderResponseBody `json:"-"`
	OrderID       uuid.UUID         `json:"-"`
}
