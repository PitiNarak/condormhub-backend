package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type LeasingHistory struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DormID   uuid.UUID `gorm:"type:uuid;not null"`
	Dorm     Dorm      `gorm:"foreignKey:DormID;references:ID"`
	LesseeID uuid.UUID `gorm:"type:uuid;not null"`
	Lessee   User      `gorm:"foreignKey:LesseeID;references:ID"`
	Orders   []Order   `gorm:"foreignKey:LeasingHistoryID"`
	Start    time.Time
	End      time.Time `gorm:"default:null"`
	Price    float64
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
