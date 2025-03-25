package dto

import (
	"github.com/google/uuid"
)

type ContractStatus string

const (
	Waiting   ContractStatus = "WAITING"
	Signed    ContractStatus = "SIGNED"
	Cancelled ContractStatus = "CANCELLED"
)

type ContractRequestBody struct {
	LesseeID uuid.UUID `json:"lesseeId"`
	DormID   uuid.UUID `json:"dormId"`
}

type ContractResponseBody struct {
	ID             uuid.UUID        `json:"id"`
	Lessor         UserResponse     `json:"lessor"`
	Lessee         UserResponse     `json:"lessee"`
	Dorm           DormResponseBody `json:"dorm"`
	LessorStatus   ContractStatus   `json:"lessorStatus"`
	LesseeStatus   ContractStatus   `json:"lesseeStatus"`
	ContractStatus ContractStatus   `json:"contractStatus"`
}
