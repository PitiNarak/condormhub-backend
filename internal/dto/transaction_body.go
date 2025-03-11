package dto

import (
	"github.com/google/uuid"
)

type TransactionRequestBody struct {
	OrderID uuid.UUID `json:"orderID"`
}

type CreateTransactionResponseBody struct {
	CheckoutUrl string `json:"checkoutUrl"`
}
