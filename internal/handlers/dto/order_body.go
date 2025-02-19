package dto

import (
	"github.com/google/uuid"
)

type OrderBody struct {
	LessorID    uuid.UUID `json:"lessorID"` // will be removed after dormitory is implemented
	LesseeID    uuid.UUID `json:"lesseeID"`
	DormitoryID uuid.UUID `json:"dormitoryID"` // TODO: add dormitory struct
}
