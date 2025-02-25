package dto

import "github.com/google/uuid"

type OrderRequestBody struct {
	LeasingHistoryID uuid.UUID `json:"leasingHistoryId" validate:"required"`
}
