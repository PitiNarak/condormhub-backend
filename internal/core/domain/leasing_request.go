package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	RequestPending  Status = "PENDING"
	RequestAccepted Status = "ACCEPT"
	RequestRejected Status = "REJECT"
)

type LeasingRequest struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Status   *Status   `gorm:"default:null"`
	DormID   uuid.UUID `gorm:"type:uuid;not null"`
	Dorm     Dorm      `gorm:"foreignKey:DormID;references:ID"`
	LesseeID uuid.UUID `gorm:"type:uuid;not null"`
	Lessee   User      `gorm:"foreignKey:LesseeID;references:ID"`
	LessorID uuid.UUID `gorm:"type:uuid;not null"`
	Lessor   User      `gorm:"foreignKey:LessorID;references:ID"`
	Create   time.Time
	End      time.Time `gorm:"default:null"`
}
