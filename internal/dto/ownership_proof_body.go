package dto

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type DormIDForOwnershipProofRequestBody struct {
	DormID uuid.UUID `json:"dormId" `
}

type UpdateOwnerShipProofRequestBody struct {
	FileKey string `json:"filekey" `
}

type UpdateOwnerShipProofStatusRequestBody struct {
	Status  domain.OwnershipProofStatus `json:"status"`
	AdminID uuid.UUID                   `json:"adminId"`
}
