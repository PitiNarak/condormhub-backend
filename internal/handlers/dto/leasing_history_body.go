package dto

type LeasingHistoryCreateRequestBody struct {
	Price int `json:"password" validate:"required,gt=100"`
}
