package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
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

func (o *OwnershipProof) ConvertToDTOWithFile(url string, expires time.Time) dto.OwnershipProofWithFileResponseBody {
	ownershipProofWithFileResponseBody := dto.OwnershipProofWithFileResponseBody{
		Url:     url,
		Expires: expires,
		DormID:  ownershipProof.DormID,
		AdminID: ownershipProof.AdminID,
		Status:  ownershipProof.Status,
	}

	return ownershipProofWithFileResponseBody
}
