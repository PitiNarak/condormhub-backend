package dto

type LeasingHistoryCreateRequestBody struct {
	Price int `json:"price" validate:"required,gt=100"`
}
