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
	ContractID uuid.UUID `json:"contractId"`
	LessorID   uuid.UUID `json:"lessorId"`
	LesseeID   uuid.UUID `json:"lesseeId"`
	DormID     uuid.UUID `json:"dormId"`
}
