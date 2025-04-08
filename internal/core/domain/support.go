package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type SupportStatus string

const (
	ProblemOpen       SupportStatus = "OPEN"
	ProblemInProgress SupportStatus = "IN-PROGRESS"
	ProblemResolved   SupportStatus = "RESOLVED"
)

type SupportRequest struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdateAt time.Time `gorm:"autoUpdateTime"`
	UserID   uuid.UUID
	Message  string        `gorn:"type:text;not null"`
	Status   SupportStatus `gorm:"default:'OPEN'"`
}

func (s *SupportRequest) ToDTO() dto.SupportResponseBody {
	return dto.SupportResponseBody{
		ID:       s.ID,
		CreateAt: s.CreateAt,
		UpdateAt: s.UpdateAt,
		UserID:   s.UserID,
		Message:  s.Message,
		Status:   string(s.Status),
	}
}
