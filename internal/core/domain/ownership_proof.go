package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OwnershipProofStatus string

const (
	Pending  OwnershipProofStatus = "Pending"
	Approved OwnershipProofStatus = "Approved"
	Rejected OwnershipProofStatus = "Rejected"
)

type OwnershipProof struct {
	DormID    uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AdminID   uuid.UUID            `gorm:"type:uuid;default:null"`
	Status    OwnershipProofStatus `gorm:"default:Pending"`
	FileKey   string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// func (o *OwnershipProof) ToDTO(ctx context.Context, storage *storage.Storage) (dto.OwnershipProofResponseBody, error) {
// 	url, err := storage.GetSignedUrl(ctx, o.FileKey, 60*time.Minute)
// 	if err != nil {
// 		return dto.OwnershipProofResponseBody{}, err
// 	}
// 	return dto.OwnershipProofResponseBody{
// 		Url:     url,
// 		DormID:  o.DormID,
// 		AdminID: o.AdminID,
// 		Status:  dto.OwnershipProofStatus(o.Status),
// 	}, nil
// }
