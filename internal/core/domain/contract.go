package domain

import "github.com/google/uuid"

type ContractStatus string

const (
	Waiting   ContractStatus = "Waiting"
	Signed    ContractStatus = "Signed"
	Cancelled ContractStatus = "Cancelled"
)

type Contract struct {
	ContractID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	LessorID     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LesseeID     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	DormID       uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LessorStatus ContractStatus `gorm:"default:Waiting"`
	LesseeStatus ContractStatus `gorm:"default:Waiting"`
	Status       ContractStatus `gorm:"default:Waiting"`
}
