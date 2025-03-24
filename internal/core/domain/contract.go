package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type ContractStatus string

const (
	Waiting   ContractStatus = "WAITING"
	Signed    ContractStatus = "SIGNED"
	Cancelled ContractStatus = "CANCELLED"
)

type Contract struct {
	ContractID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt     time.Time      `gorm:"autoCreateTime"`
	LessorID     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LesseeID     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	DormID       uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LessorStatus ContractStatus `gorm:"default:WAITING"`
	LesseeStatus ContractStatus `gorm:"default:WAITING"`
	Status       ContractStatus `gorm:"default:WAITING"`
}

func (ct *Contract) ToDTO() dto.ContractResponseBody {
	return dto.ContractResponseBody{
		ContractID:     ct.ContractID,
		LessorID:       ct.LessorID,
		LesseeID:       ct.LesseeID,
		DormID:         ct.DormID,
		LessorStatus:   string(ct.LessorStatus),
		LesseeStatus:   string(ct.LesseeStatus),
		ContractStatus: string(ct.Status),
	}
}
