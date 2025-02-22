package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderType string

const (
	InsuranceOrderType   OrderType = "insurance"
	MonthlyBillOrderType OrderType = "monthly_bill"
)

type Order struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt          time.Time      `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt          time.Time      `json:"updateAt" gorm:"autoUpdateTime"`
	Type              OrderType      `json:"type"`
	Transactions      []*Transaction `json:"transactions" gorm:"foreignKey:OrderID"`
	PaidTransaction   *Transaction   `json:"paidTransaction" gorm:"foreignKey:OrderID;default:null"`
	PaidTransactionID string         `json:"paidTransactionID" gorm:"default:null"`
}
