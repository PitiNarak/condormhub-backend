package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Message    string    `json:"message" validate:"required"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	UserName   string    `json:"userName" gorm:"unique" validate:"required"`
	Password   string    `json:"password" validate:"required"`
	Email      string    `json:"email" gorm:"unique" validate:"required,email"`
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
