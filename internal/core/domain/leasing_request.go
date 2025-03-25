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
	Status   Status    `gorm:"default:null"`
	DormID   uuid.UUID `gorm:"type:uuid;not null"`
	Dorm     Dorm      `gorm:"foreignKey:DormID;references:ID"`
	LesseeID uuid.UUID `gorm:"type:uuid;not null"`
	Lessee   User      `gorm:"foreignKey:LesseeID;references:ID"`
	Start    time.Time `gorm:"autoCreateTime"`
	End      time.Time `gorm:"default:null"`
	Message  string
}

func (l *LeasingRequest) ToDTO() dto.LeasingRequest {
	return dto.LeasingRequest{
		ID:      l.ID,
		Status:  dto.Status(l.Status),
		Dorm:    l.Dorm.ToDTO(),
		Lessee:  l.Lessee.ToDTO(),
		Start:   l.Start,
		End:     l.End,
		Message: l.Message,
	}
}
