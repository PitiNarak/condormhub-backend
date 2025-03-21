package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type Review struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Message  string
	Rate     int
	CreateAt time.Time `gorm:"autoCreateTime"`
}

func (r *Review) ToDTO() dto.Review {
	return dto.Review{
		ID:       r.ID,
		Message:  r.Message,
		Rate:     r.Rate,
		CreateAt: r.CreateAt,
	}
}
