package dto

import (
	"time"

	"github.com/google/uuid"
)

type LeasingRequest struct {
	ID     uuid.UUID        `json:"id"`
	Status string           `json:"status"`
	Dorm   DormResponseBody `json:"dorm"`
	Lessee UserResponse     `json:"lessee"`
	Lessor UserResponse     `json:"lessor"`
	Create time.Time        `json:"create"`
	End    time.Time        `json:"end"`
}
