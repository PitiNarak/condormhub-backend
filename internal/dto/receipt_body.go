package dto

import (
	"github.com/google/uuid"
)

type ReceiptResponseBody struct {
	ID          uuid.UUID           `json:"receiptId"`
	Owner       UserResponse        `json:"owner"`
	Transaction TransactionResponse `json:"transaction"`
	Url         string              `json:"url"`
}
