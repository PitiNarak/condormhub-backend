package dto

import (
	"time"

	"github.com/google/uuid"
)

type LeasingHistory struct {
	ID     uuid.UUID           `json:"id"`
	Dorm   DormResponseBody    `json:"dorm"`
	Lessee UserResponse        `json:"lessee"`
	Orders []OrderResponseBody `json:"orders"`
	Start  time.Time           `json:"start"`
	End    time.Time           `json:"end"`
	Price  float64             `json:"price"`
	Review Review              `json:"review"`
}

type Review struct {
	ID       uuid.UUID `json:"id"`
	Message  string    `json:"message"`
	Rate     int       `json:"rate"`
	CreateAt time.Time `json:"crateAt"`
}

type ReviewRequestBody struct {
	Message string    `json:"message" validate:"omitempty"`
	Rate    int       `json:"rate" validate:"omitempty,gte=0,lte=5"`
	ID      uuid.UUID `json:"id" validate:"omitempty"`
}
