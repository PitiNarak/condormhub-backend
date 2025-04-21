package domain

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dorm struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt    time.Time      `gorm:"autoCreateTime"`
	UpdateAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `validate:"required"`
	OwnerID     uuid.UUID      `validate:"required"`
	Owner       User
	Size        float64 `validate:"required,gt=0"`
	Bedrooms    int     `validate:"required,gte=0"`
	Bathrooms   int     `validate:"required,gte=0"`
	Address     Address `gorm:"embedded" validate:"required"`
	Price       float64 `validate:"required,gt=0"`
	Rating      float64 `gorm:"default:0" validate:"gte=0,lte=5"`
	Description string  `gorm:"type:text"`
	Images      []DormImage
}

type Address struct {
	District    string `validate:"required"`
	Subdistrict string `validate:"required"`
	Province    string `validate:"required"`
	Zipcode     string `validate:"required,numeric,len=5"`
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

type DormImage struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt time.Time `json:"createAt" gorm:"autoCreateTime"`
	DormID   uuid.UUID `gorm:"type:uuid;not null"`
	ImageKey string    `gorm:"type:text;not null"`
}

func updateDormsOwnedCount(tx *gorm.DB, lessorID uuid.UUID) error {
	var count int64
	if err := tx.Model(&Dorm{}).Where("owner_id = ?", lessorID).Count(&count).Error; err != nil {
		return err
	}

	if err := tx.Model(&User{}).Where("id = ?", lessorID).Update("dorms_owned", count).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dorm) AfterCreate(tx *gorm.DB) (err error) {
	return updateDormsOwnedCount(tx, d.OwnerID)
}

func (d *Dorm) AfterDelete(tx *gorm.DB) (err error) {
	return updateDormsOwnedCount(tx, d.OwnerID)
}

func (d *Dorm) BeforeDelete(tx *gorm.DB) (err error) {
	// Mark Pending leasing request as Rejected if the dorm related is deleted
	err = tx.Model(&LeasingRequest{}).Where("dorm_id = ? AND status = ?", d.ID, RequestPending).Update("status", RequestRejected).Error
	if err != nil {
		return err
	}

	return nil
}
