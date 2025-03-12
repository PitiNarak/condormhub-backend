package dto

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type DormIDForOwnershipProofRequestBody struct {
	DormID uuid.UUID `json:"dormId" `
}

type UpdateOwnerShipProofStatusRequestBody struct {
	Status  domain.OwnershipProofStatus `json:"status"`
	AdminID uuid.UUID                   `json:"adminId"`
}

type OwnershipProofWithFileResponseBody struct {
	Url     string                      `json:"url"`
	Expires time.Time                   `json:"expires"`
	DormID  uuid.UUID                   `json:"dormId"`
	AdminID uuid.UUID                   `json:"adminId"`
	Status  domain.OwnershipProofStatus `json:"status" gorm:"default:Pending"`
}

type OwnershipProofResponseBody struct {
	DormID  uuid.UUID                   `json:"dormId"`
	AdminID uuid.UUID                   `json:"adminId"`
	Status  domain.OwnershipProofStatus `json:"status" gorm:"default:Pending"`
}
