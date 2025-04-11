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
	// Calculate the new average rating after an update to leasingHistory
	var avgRating float64
	err = tx.Model(&LeasingHistory{}).Select("COALESCE(avg(rate), 0)").Where("dorm_id = ?", l.DormID).Scan(&avgRating).Error
	if err != nil {
		return err
	}

	// Update the rating field in the dorms table
	err = tx.Model(Dorm{}).Where("id = ?", l.DormID).Update("rating", avgRating).Error
	if err != nil {
		return err
	}

	// Find the dorm the is related to the updated leasingHistory to find the lessor related
	dorm := new(Dorm)
	err = tx.First(dorm, l.DormID).Error
	if err != nil {
		return err
	}

	// Count the amount of review a lessor has after an update to leasingHistory
	var count int64
	err = tx.Model(&LeasingHistory{}).Joins("JOIN dorms ON dorms.id = leasing_histories.dorm_id").Where("dorms.owner_id = ? AND review_flag = true", dorm.OwnerID).Count(&count).Error
	if err != nil {
		return err
	}

	// Update the review count in the users table
	err = tx.Model(&User{}).Where("id = ?", dorm.OwnerID).Update("review_count", count).Error
	if err != nil {
		return nil
	}

	return nil
}
