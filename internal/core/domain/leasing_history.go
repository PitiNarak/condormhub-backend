package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeasingHistory struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DormID     uuid.UUID `gorm:"type:uuid;not null"`
	Dorm       Dorm      `gorm:"foreignKey:DormID;references:ID"`
	LesseeID   uuid.UUID `gorm:"type:uuid;not null"`
	Lessee     User      `gorm:"foreignKey:LesseeID;references:ID"`
	Orders     []Order   `gorm:"foreignKey:LeasingHistoryID"`
	Start      time.Time
	End        time.Time `gorm:"default:null"`
	Price      float64
	ReviewFlag bool
	Review     Review `gorm:"embedded"`
}

func (l *LeasingHistory) ToDTO() dto.LeasingHistory {
	orders := make([]dto.OrderResponseBody, len(l.Orders))
	for i, v := range l.Orders {
		orders[i] = v.ToDTO()
	}
	var review dto.Review
	if l.ReviewFlag {
		review = l.Review.ToDTO()
	}
	return dto.LeasingHistory{
		ID:         l.ID,
		Dorm:       l.Dorm.ToDTO(),
		Lessee:     l.Lessee.ToDTO(),
		Orders:     orders,
		Start:      l.Start,
		End:        l.End,
		Price:      l.Price,
		Review:     review,
		ReviewFlag: l.ReviewFlag,
	}
}

type Review struct {
	Message  string     `gorm:"default:null"`
	Rate     int        `gorm:"default:null"`
	CreateAt *time.Time `gorm:"autoUpdateTime;default:null"`
}

func (r *Review) ToDTO() dto.Review {
	return dto.Review{
		Message:  r.Message,
		Rate:     r.Rate,
		CreateAt: *r.CreateAt,
	}
}

func (l *LeasingHistory) AfterUpdate(tx *gorm.DB) (err error) {
	var avgRating float64

	err = tx.Model(&LeasingHistory{}).Select("COALESCE(avg(rate), 0)").Where("dorm_id = ?", l.DormID).Scan(&avgRating).Error
	if err != nil {
		return err
	}

	return tx.Model(Dorm{}).Where("id = ?", l.DormID).Update("rating", avgRating).Error
}
