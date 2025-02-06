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
