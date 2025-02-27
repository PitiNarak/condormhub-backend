package domain

import (
	"time"

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
