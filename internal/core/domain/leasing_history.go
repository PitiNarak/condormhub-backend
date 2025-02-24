package domain

import (
	"time"

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
}
