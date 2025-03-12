package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type OrderType string

const (
	InsuranceOrderType   OrderType = "insurance"
	MonthlyBillOrderType OrderType = "monthly_bill"
)

type Order struct {
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt          time.Time       `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt          time.Time       `json:"updateAt" gorm:"autoUpdateTime"`
	Type              OrderType       `json:"type"`
	Price             int64           `json:"price"`
	Transactions      []*Transaction  `json:"transactions" gorm:"foreignKey:OrderID"`
	PaidTransaction   *Transaction    `json:"paidTransaction" gorm:"foreignKey:OrderID;default:null"`
	PaidTransactionID string          `json:"paidTransactionID" gorm:"default:null"`
	LeasingHistory    *LeasingHistory `json:"leasingHistory" gorm:"foreignKey:LeasingHistoryID"`
	LeasingHistoryID  uuid.UUID       `json:"leasingHistoryID"`
}

func (o *Order) ToDTO() dto.OrderResponseBody {
	return dto.OrderResponseBody{
		ID:    o.ID,
		Type:  string(o.Type),
		Price: o.Price,
		PaidTransaction: &dto.TransactionResponse{
			ID:            o.PaidTransaction.ID,
			SessionStatus: string(o.PaidTransaction.SessionStatus),
			CreateAt:      o.PaidTransaction.CreateAt,
			UpdateAt:      o.PaidTransaction.UpdateAt,
			Price:         o.PaidTransaction.Price,
		},
	}
}
