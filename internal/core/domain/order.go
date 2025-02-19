package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type OrderType string

const (
	InsuranceOrder OrderType = "insurance"
	MonthlyBill    OrderType = "monthly_bill"
)

type Order struct {
	ID            uuid.UUID                    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt      time.Time                    `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt      time.Time                    `json:"updateAt" gorm:"autoUpdateTime"`
	SessionID     string                       `json:"sessionId"`
	Price         int64                        `json:"price"`
	SessionStatus stripe.CheckoutSessionStatus `json:"status" gorm:"default:'open'"`
	Type          OrderType                    `json:"type"`
	LessorID      uuid.UUID                    `json:"lessorID"`
	Lessor        *User                        `json:"lessor" gorm:"foreignKey:LessorID"`
	LesseeID      uuid.UUID                    `json:"lesseeID"`
	Lessee        *User                        `json:"lessee" gorm:"foreignKey:LesseeID"`
	DormitoryID   uuid.UUID                    `json:"dormitoryID"` // TODO: add dormitory struct
	CheckoutUrl   string                       `json:"checkoutUrl"`
}
