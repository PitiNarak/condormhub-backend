package domain

import (
	"github.com/google/uuid"
)

type OwnershipProofStatus string

const (
	Pending  OwnershipProofStatus = "Pending"
	Approved OwnershipProofStatus = "Approved"
	Rejected OwnershipProofStatus = "Rejected"
)

type OwnershipProof struct {
	DormID  uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID uuid.UUID            `gorm:"type:uuid;default:null"`
	Status  OwnershipProofStatus `gorm:"default:Pending"`
	FileKey string
}
