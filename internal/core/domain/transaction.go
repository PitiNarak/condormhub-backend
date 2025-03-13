package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type CheckoutStatus string

const (
	StatusOpen     CheckoutStatus = "open"
	StatusComplete CheckoutStatus = "complete"
	StatusExpired  CheckoutStatus = "expired"
)

type Transaction struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	SessionStatus CheckoutStatus `json:"status" gorm:"default:open"`
	CreateAt      time.Time      `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt      time.Time      `json:"updateAt" gorm:"autoUpdateTime"`
	Price         int64          `json:"price"`
	Order         Order          `json:"-" gorm:"foreignKey:OrderID"`
	OrderID       uuid.UUID      `json:"-"`
}

func (t *Transaction) ToDTO() dto.TransactionResponse {
	if t == nil {
		return dto.TransactionResponse{}
	}
	return dto.TransactionResponse{
		ID:            t.ID,
		SessionStatus: string(t.SessionStatus),
		CreateAt:      t.CreateAt,
		UpdateAt:      t.UpdateAt,
		Price:         t.Price,
	}
}
