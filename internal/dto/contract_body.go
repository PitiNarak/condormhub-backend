package dto

import (
	"github.com/google/uuid"
)

type ContractStatus string

const (
	Waiting   ContractStatus = "Waiting"
	Signed    ContractStatus = "Signed"
	Cancelled ContractStatus = "Cancelled"
)

type ContractRequestBody struct {
	LessorID uuid.UUID `json:"lessorId"`
	LesseeID uuid.UUID `json:"lesseeId"`
	DormID   uuid.UUID `json:"dormId"`
}

type UpdateContractStatusRequestBody struct {
	Status  ContractStatus `json:"status"`
	AdminID uuid.UUID      `json:"adminId"`
}
