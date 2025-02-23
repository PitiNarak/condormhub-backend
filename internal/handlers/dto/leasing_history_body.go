package dto

type LeasingHistoryCreateRequestBody struct {
	DormID string `json:"dormid" validate:"required"`
}
