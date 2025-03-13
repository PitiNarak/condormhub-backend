package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type Dorm struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt    time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt    time.Time `json:"updateAt" gorm:"autoUpdateTime"`
	Name        string    `json:"name" validate:"required"`
	OwnerID     uuid.UUID `json:"ownerId" validate:"required"`
	Owner       User      `json:"owner"`
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

func (d *Dorm) ToDTO() dto.DormResponseBody {
	return dto.DormResponseBody{
		ID:          d.ID,
		CreateAt:    d.CreateAt,
		UpdateAt:    d.UpdateAt,
		Name:        d.Name,
		Owner:       d.Owner.ToDTO(),
		Size:        d.Size,
		Bedrooms:    d.Bedrooms,
		Bathrooms:   d.Bathrooms,
		Address:     d.Address.ToDTO(),
		Price:       d.Price,
		Rating:      d.Rating,
		Description: d.Description,
	}
}

func (a *Address) ToDTO() dto.Address {
	return dto.Address{
		District:    a.District,
		Subdistrict: a.Subdistrict,
		Province:    a.Province,
		Zipcode:     a.Zipcode,
	}
}
