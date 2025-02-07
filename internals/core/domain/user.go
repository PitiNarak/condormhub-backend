package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	Password   string    `json:"password" validate:"required"`
	Email      string    `gorm:"unique" json:"email" validate:"required,email"`
	Name       string    `json:"name" validate:"required"`
	IsVerified bool      `gorm:"default:false" json:"isVerified"`
}
