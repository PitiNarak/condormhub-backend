package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckoutStatus string

const (
	StatusOpen     CheckoutStatus = "open"
	StatusComplete CheckoutStatus = "complete"
	StatusExpired  CheckoutStatus = "expired"
)

type Transaction struct {
	ID            string         `gorm:"primaryKey"`
	SessionStatus CheckoutStatus `gorm:"default:open"`
	CreateAt      time.Time      `gorm:"autoCreateTime"`
	UpdateAt      time.Time      `gorm:"autoUpdateTime"`
	Price         int64
	Order         Order `gorm:"foreignKey:OrderID"`
	OrderID       uuid.UUID
	DeletedAt     gorm.DeletedAt `gorm:"index"`
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
