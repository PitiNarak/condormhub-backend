package domain

import "github.com/google/uuid"

type ContractStatus string

const (
	Waiting   ContractStatus = "Waiting"
	Signed    ContractStatus = "Signed"
	Cancelled ContractStatus = "Cancelled"
)

type Contract struct {
	LessorID     uuid.UUID      `gorm:"primaryKey"`
	LesseeID     uuid.UUID      `gorm:"primaryKey"`
	DormID       uuid.UUID      `gorm:"primaryKey"`
	LessorStatus ContractStatus `gorm:"default:Waiting"`
	LesseeStatus ContractStatus `gorm:"default:Waiting"`
	Status       ContractStatus `gorm:"default:Waiting"`
}
