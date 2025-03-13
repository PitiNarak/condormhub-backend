package dto

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type UpdateOwnerShipProofStatusRequestBody struct {
	Status  domain.OwnershipProofStatus `json:"status"`
	AdminID uuid.UUID                   `json:"adminId"`
}

type OwnershipProofResponseBody struct {
	Url     string                      `json:"url"`
	DormID  uuid.UUID                   `json:"dormId"`
	AdminID uuid.UUID                   `json:"adminId"`
	Status  domain.OwnershipProofStatus `json:"status"`
}
