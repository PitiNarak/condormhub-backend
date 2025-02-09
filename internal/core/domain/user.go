package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreateAt   time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt   time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	UserName   string    `json:"userName" gorm:"unique" validate:"required"`
	Password   string    `json:"password" validate:"required,min=8"`
	Email      string    `json:"email" gorm:"unique" validate:"required,email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	NationalID string    `json:"nationalID" `
	Gender     string    `json:"gender"`
	BirthDate  string    `json:"birthDate"`
	IsVerified bool      `gorm:"default:false" json:"isAccountVerified"`
	Role       string    `json:"role"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `gorm:"default:false" json:"isStudentVerified"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateInfo struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	NationalID      string `json:"nationalID" `
	Gender          string `json:"gender"`
	BirthDate       string `json:"birthDate"`
	Role            string `json:"role"`
	StudentEvidence string `json:"studentEvidence"`
}
