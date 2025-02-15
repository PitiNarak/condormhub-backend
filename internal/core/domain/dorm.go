package domain

import (
	"time"

	"github.com/google/uuid"
)

type Dorm struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt    time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt    time.Time `json:"updateAt" gorm:"autoUpdateTime"`
	Name        string    `json:"name" validate:"required"`
	OwnerID     uuid.UUID `validate:"required"`
	Owner       User
	Size        float64 `json:"size" validate:"required,gt=0"`
	Bedrooms    int     `json:"bedrooms" validate:"required,gte=0"`
	Bathrooms   int     `json:"bathrooms" validate:"required,gte=0"`
	Province    string  `json:"province" validate:"required"`
	District    string  `json:"district" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Rating      float64 `json:"rating" gorm:"default:0" validate:"gte=0,lte=5"`
	Description string  `json:"description" gorm:"type:text"`
}
