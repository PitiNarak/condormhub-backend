package dto

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type OrderRequestBody struct {
	LessorID    uuid.UUID `json:"lessorID"` // will be removed after dormitory is implemented
	LesseeID    uuid.UUID `json:"lesseeID"`
	DormitoryID uuid.UUID `json:"dormitoryID"` // TODO: add dormitory struct
}

type CreateOrderResponseBody struct {
	Order       domain.Order `json:"order"`
	CheckoutUrl string       `json:"checkoutUrl"`
}
