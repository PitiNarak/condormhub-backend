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
	LessorID uuid.UUID             `json:"lessorId" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID  uuid.UUID             `json:"adminId" gorm:"type:uuid;default:null"`
	Lessor   User                  `json:"lessor" gorm:"foreignKey:LessorID;references:ID"`
	Admin    User                  `json:"admin" gorm:"foreignKey:AdminID;references:ID"`
	Status   *OwnershipProofStatus `json:"status"`
	FileKey  string                `json:"file_key"`
}
