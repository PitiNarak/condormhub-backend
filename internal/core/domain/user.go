package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole  Role = "ADMIN"
	UserRole   Role = "USER"
	LesseeRole Role = "LESSEE"
	LesserRole Role = "LESSOR"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreateAt           time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt           time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	Username           string    `json:"username" gorm:"unique" validate:"required"`
	Password           string    `json:"-" validate:"required,min=8"`
	Email              string    `json:"email" gorm:"unique" validate:"required,email"`
	FirstName          string    `json:"firstName"`
	LastName           string    `json:"lastName"`
	NationalID         string    `json:"nationalID" `
	Gender             string    `json:"gender"`
	BirthDate          string    `json:"birthDate"`
	IsVerified         bool      `gorm:"default:false" json:"isVerified"`
	Role               Role      `gorm:"default:USER" json:"role"`
	FilledPersonalInfo bool      `gorm:"default:false" json:"filledPersonalInfo"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `gorm:"default:false" json:"isStudentVerified"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GetInfoRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserInfo struct {
	UserName           string `json:"username" gorm:"unique" validate:"required"`
	Email              string `json:"email" gorm:"unique" validate:"required,email"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	NationalID         string `json:"nationalID" `
	Gender             string `json:"gender"`
	BirthDate          string `json:"birthDate"`
	IsVerified         bool   `gorm:"default:false" json:"isVerified"`
	Role               string `json:"role"`
	StudentEvidence    string `json:"studentEvidence"`
	IsStudentVerified  bool   `gorm:"default:false" json:"isStudentVerified"`
	FilledPersonalInfo bool   `gorm:"default:false" json:"filledPersonalInfo"`
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
