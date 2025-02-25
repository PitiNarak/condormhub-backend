package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type Transaction struct {
	ID            string                       `json:"id" gorm:"primaryKey"`
	SessionStatus stripe.CheckoutSessionStatus `json:"status" gorm:"default:open"`
	CreateAt      time.Time                    `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt      time.Time                    `json:"updateAt" gorm:"autoUpdateTime"`
	Price         int64                        `json:"price"`
	Order         Order                        `json:"order" gorm:"foreignKey:OrderID"`
	OrderID       uuid.UUID                    `json:"orderID"`
}
