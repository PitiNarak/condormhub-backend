package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Message    string    `json:"message" validate:"required"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"update_at"`
	UserName   string    `json:"userName" gorm:"unique"`
	Password   string    `json:"password"`
	Email      string    `json:"email" gorm:"unique"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	NationalID string    `json:"nationalID" `
	Gender     string    `json:"gender"`
	BirthDate  string    `json:"birthDate"`
	IsVerified bool      `json:"isVerified"`
	Role       string    `json:"role"`
	// studentEvidence
	FilledPersonalInfo bool `json:"filledPersonalInfo"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
