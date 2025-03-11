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
	DormID  uuid.UUID             `json:"dormId" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID uuid.UUID             `json:"adminId" gorm:"type:uuid;default:null"`
	Status  *OwnershipProofStatus `json:"status" gorm:"default:Pending"`
	FileKey string                `json:"file_key"`
}
