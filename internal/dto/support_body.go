package dto

import (
	"time"

	"github.com/google/uuid"
)

type SupportRequestBody struct {
	Message string `json:"message" validate:"required"`
}

type UpdateStatusRequestBody struct {
	Status string `json:"status" validate:"required"`
}

type SupportResponseBody struct {
	ID       uuid.UUID `json:"id"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
	UserID   uuid.UUID `json:"userID"`
	Message  string    `json:"message"`
	Status   string    `json:"status"`
}
