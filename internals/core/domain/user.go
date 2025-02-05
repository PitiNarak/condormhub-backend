package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Message    string    `json:"message" validate:"required"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"update_at"`
	Password   string    `json:"password"`
	Email      string    `gorm:"unique" json:"email"`
	Name       string    `json:"name"`
	IsVerified bool      `json:"isVerified"`
}

type Lesser struct {
	Id         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Message    string    `json:"message" validate:"required"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"update_at"`
	UserName   string    `json:"userName"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	NationalID string    `json:"nationalID"`
	Gender     string    `json:"gender"`
	BirthDate  string    `json:"birthDate"`
	IsVerified bool      `json:"isVerified"`
}
