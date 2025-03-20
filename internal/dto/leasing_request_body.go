package dto

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	RequestPending  Status = "PENDING"
	RequestAccepted Status = "ACCEPT"
	RequestRejected Status = "REJECT"
	RequestCanceled Status = "CANCELED"
)

type LeasingRequest struct {
	ID     uuid.UUID        `json:"id"`
	Status Status           `json:"status"`
	Dorm   DormResponseBody `json:"dorm"`
	Lessee UserResponse     `json:"lessee"`
	Lessor UserResponse     `json:"lessor"`
	Start  time.Time        `json:"start"`
	End    time.Time        `json:"end"`
}
