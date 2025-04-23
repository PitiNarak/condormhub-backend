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
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt     time.Time      `gorm:"autoCreateTime"`
	LesseeID     uuid.UUID      `gorm:"type:uuid;not null"`
	Lessee       User           `gorm:"foreignKey:LesseeID;references:ID"`
	DormID       uuid.UUID      `gorm:"type:uuid;not null"`
	Dorm         Dorm           `gorm:"foreignKey:DormID;references:ID"`
	LessorStatus ContractStatus `gorm:"default:WAITING"`
	LesseeStatus ContractStatus `gorm:"default:WAITING"`
	Status       ContractStatus `gorm:"default:WAITING"`
}

func (ct *Contract) ToDTO(urls []string) dto.ContractResponseBody {
	dormResponse := ct.Dorm.ToDTO()
	dormResponse.Images = urls
	return dto.ContractResponseBody{
		ID:             ct.ID,
		Lessee:         ct.Lessee.ToDTO(),
		Dorm:           dormResponse,
		LessorStatus:   dto.ContractStatus(ct.LessorStatus),
		LesseeStatus:   dto.ContractStatus(ct.LesseeStatus),
		ContractStatus: dto.ContractStatus(ct.Status),
	}
}
