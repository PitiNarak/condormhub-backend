package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type LeasingHistory struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DormID   uuid.UUID `json:"dorm_id" gorm:"type:uuid;not null"`
	Dorm     Dorm      `json:"dorm" gorm:"foreignKey:DormID;references:ID"`
	LesseeID uuid.UUID `json:"lessee_id" gorm:"type:uuid;not null"`
	Lessee   User      `json:"lessee" gorm:"foreignKey:LesseeID;references:ID"`
	Orders   []Order   `json:"orders" gorm:"foreignKey:LeasingHistoryID"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end" gorm:"default:null"`
	Price    float64   `json:"price"`
}

func (l *LeasingHistory) ToDTO() dto.LeasingHistory {
	orders := make([]dto.OrderResponseBody, len(l.Orders))
	for i, v := range l.Orders {
		orders[i] = v.ToDTO()
	}

	return dto.LeasingHistory{
		ID:       l.ID,
		Dorm:     l.Dorm.ToDTO(),
		LesseeID: l.LesseeID,
		Lessee:   l.Lessee.ToDTO(),
		Orders:   orders,
		Start:    l.Start,
		End:      l.End,
		Price:    l.Price,
	}
}
