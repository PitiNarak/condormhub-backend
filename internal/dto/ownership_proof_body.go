package dto

import (
	"github.com/google/uuid"
)

type OwnershipProofStatus string

const (
	Pending  OwnershipProofStatus = "Pending"
	Approved OwnershipProofStatus = "Approved"
	Rejected OwnershipProofStatus = "Rejected"
)

type UpdateOwnerShipProofStatusRequestBody struct {
	Status  OwnershipProofStatus `json:"status"`
	AdminID uuid.UUID            `json:"adminId"`
}

type OwnershipProofResponseBody struct {
	Url     string               `json:"url"`
	DormID  uuid.UUID            `json:"dormId"`
	AdminID uuid.UUID            `json:"adminId"`
	Status  OwnershipProofStatus `json:"status"`
}
