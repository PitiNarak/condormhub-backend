package dto

import "github.com/google/uuid"

type DormRequestBody struct {
	Name        string    `json:"name" validate:"required"`
	OwnerID     uuid.UUID `json:"ownerId" validate:"required"`
	Size        float64   `json:"size" validate:"required,gt=0"`
	Bedrooms    int       `json:"bedrooms" validate:"required,gte=0"`
	Bathrooms   int       `json:"bathrooms" validate:"required,gte=0"`
	Province    string    `json:"province" validate:"required"`
	District    string    `json:"district" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Description string    `json:"description" gorm:"type:text"`
}
