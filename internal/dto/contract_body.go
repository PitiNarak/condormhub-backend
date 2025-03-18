package dto

import "github.com/google/uuid"

type ContractRequestBody struct {
	LessorID uuid.UUID `json:"lessorId"`
	LesseeID uuid.UUID `json:"lesseeId"`
	DormID   uuid.UUID `json:"dormId"`
}
