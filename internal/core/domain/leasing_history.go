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
	Review     Review        `gorm:"embedded"`
	Images     []ReviewImage `gorm:"foreignKey:HistoryID"` // Link to ReviewImage
	ReportFlag bool          `gorm:"default:false"`
}

func (l *LeasingHistory) ToDTO(urls []string) dto.LeasingHistory {
	orders := make([]dto.OrderResponseBody, len(l.Orders))
	for i, v := range l.Orders {
		orders[i] = v.ToDTO()
	}
	var review dto.Review
	if l.ReviewFlag {
		review = l.Review.ToDTO(urls)
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

type ReviewImage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt  time.Time `json:"createAt" gorm:"autoCreateTime"`
	HistoryID uuid.UUID `gorm:"type:uuid;not null"`
	ImageKey  string    `gorm:"type:text;not null"`
}

func (r *Review) ToDTO(urls []string) dto.Review {

	return dto.Review{
		Message:  r.Message,
		Rate:     r.Rate,
		CreateAt: *r.CreateAt,
		Images:   urls,
	}
}

func (r *Review) ToReportedReviewDTO(urls []string, reviewer dto.UserResponse) dto.ReportedReview {
	return dto.ReportedReview{
		Message:  r.Message,
		Rate:     r.Rate,
		Reviewer: reviewer,
		CreateAt: *r.CreateAt,
		Images:   urls,
	}
}

func (l *LeasingHistory) AfterUpdate(tx *gorm.DB) (err error) {
	// Calculate the new average rating after an update to leasingHistory
	var avgRating float64
	err = tx.Model(&LeasingHistory{}).Select("COALESCE(avg(rate), 0)").Where("dorm_id = ?", l.DormID).Where("review_flag = ?", true).Scan(&avgRating).Error
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
	err = tx.Model(&LeasingHistory{}).Joins("JOIN dorms ON dorms.id = leasing_histories.dorm_id").Where("dorms.owner_id = ?", dorm.OwnerID).Where("leasing_histories.review_flag = ?", true).Count(&count).Error
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

func updateDormsLeasedCount(tx *gorm.DB, lesseeID uuid.UUID) error {
	var count int64
	if err := tx.Model(&LeasingHistory{}).Where("lessee_id = ?", lesseeID).Count(&count).Error; err != nil {
		return err
	}

	if err := tx.Model(&User{}).Where("id = ?", lesseeID).Update("dorms_leased", count).Error; err != nil {
		return err
	}

	return nil
}

func (l *LeasingHistory) AfterCreate(tx *gorm.DB) (err error) {
	return updateDormsLeasedCount(tx, l.LesseeID)
}
