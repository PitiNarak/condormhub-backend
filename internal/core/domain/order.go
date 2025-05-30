package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderType string

const (
	InsuranceOrderType   OrderType = "insurance"
	MonthlyBillOrderType OrderType = "monthly_bill"
)

type Order struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt          time.Time `gorm:"autoCreateTime"`
	UpdateAt          time.Time `gorm:"autoUpdateTime"`
	Type              OrderType
	Price             int64
	Transactions      []*Transaction `gorm:"foreignKey:OrderID"`
	PaidTransaction   *Transaction   `gorm:"foreignKey:OrderID;default:null"`
	PaidTransactionID string         `gorm:"default:null"`
	LeasingHistory    LeasingHistory `gorm:"foreignKey:LeasingHistoryID"`
	LeasingHistoryID  uuid.UUID
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (o *Order) ToDTO() dto.OrderResponseBody {
	return dto.OrderResponseBody{
		ID:              o.ID,
		Type:            string(o.Type),
		Price:           o.Price,
		PaidTransaction: o.PaidTransaction.ToDTO(),
	}
}

func (o *Order) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Model(&Transaction{}).Where("order_id = ?", o.ID).Delete(&Transaction{}).Error
}
