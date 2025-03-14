package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type Status string

const (
	RequestPending  Status = "PENDING"
	RequestAccepted Status = "ACCEPT"
	RequestRejected Status = "REJECT"
	RequestCanceled Status = "CANCELED"
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
	Start    time.Time
	End      time.Time `gorm:"default:null"`
}

func (l *LeasingRequest) ToDTO() dto.LeasingRequest {
	return dto.LeasingRequest{
		ID:     l.ID,
		Status: string(*l.Status),
		Dorm:   l.Dorm.ToDTO(),
		Lessee: l.Lessee.ToDTO(),
		Lessor: l.Lessor.ToDTO(),
		Start:  l.Start,
		End:    l.End,
	}
}
