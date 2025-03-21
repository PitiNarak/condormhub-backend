package dto

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID       uuid.UUID `json:"id"`
	Message  string    `json:"message"`
	Rate     int       `json:"rate"`
	CreateAt time.Time `json:"crateAt"`
}

type ReviewRequestBody struct {
	Message string    `json:"message" validate:"omitempty"`
	Rate    float64   `json:"rate" validate:"omitempty,gte=0,lte=5"`
	ID      uuid.UUID `json:"id" validate:"omitempty"`
}
