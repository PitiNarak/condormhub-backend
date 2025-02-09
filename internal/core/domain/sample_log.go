package domain

import (
	"time"

	"github.com/google/uuid"
)

type SampleLog struct {
	Id       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Message  string    `json:"message" validate:"required"`
	CreateAt time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt time.Time `gorm:"autoUpdateTime" json:"update_at"`
}

func NewSampleLog(id uuid.UUID, message string) *SampleLog {
	return &SampleLog{
		Id:      id,
		Message: message,
	}
}

func (l *SampleLog) GetId() uuid.UUID {
	return l.Id
}

func (l *SampleLog) GetMessage() string {
	return l.Message
}
