package dto

import (
	"github.com/google/uuid"
)

type ContractRequestBody struct {
	ContractID uuid.UUID `json:"contractId"`
	LessorID   uuid.UUID `json:"lessorId"`
	LesseeID   uuid.UUID `json:"lesseeId"`
	DormID     uuid.UUID `json:"dormId"`
}

type ContractResponseBody struct {
	ContractID     uuid.UUID        `json:"contractId"`
	Lessor         UserResponse     `json:"lessor"`
	Lessee         UserResponse     `json:"lessee"`
	Dorm           DormResponseBody `json:"dorm"`
	LessorStatus   string           `json:"lessorStatus"`
	LesseeStatus   string           `json:"lesseeStatus"`
	ContractStatus string           `json:"contractStatus"`
}
