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
	OwnerID     uuid.UUID `json:"ownerId" validate:"required"`
	Owner       User      `json:"-"`
	Size        float64   `json:"size" validate:"required,gt=0"`
	Bedrooms    int       `json:"bedrooms" validate:"required,gte=0"`
	Bathrooms   int       `json:"bathrooms" validate:"required,gte=0"`
	Address     Address   `json:"address" gorm:"embedded" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Rating      float64   `json:"rating" gorm:"default:0" validate:"gte=0,lte=5"`
	Description string    `json:"description" gorm:"type:text"`
}

type Address struct {
	District    string `json:"district" validate:"required"`
	Subdistrict string `json:"subdistrict" validate:"required"`
	Province    string `json:"province" validate:"required"`
	Zipcode     string `json:"zipcode" validate:"required,numeric,len=5"`
}
