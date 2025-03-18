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
}
