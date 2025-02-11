package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole  Role = "ADMIN"
	LesseeRole Role = "LESSEE"
	LessorRole Role = "LESSOR"
)

type User struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt           time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt           time.Time `json:"updateAt" gorm:"autoUpdateTime"`
	Username           string    `json:"username" gorm:"unique" validate:"required"`
	Password           string    `json:"-" validate:"required,min=8"`
	Email              string    `json:"email" gorm:"unique" validate:"required,email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	NationalID         string    `json:"nationalID" `
	Gender             string    `json:"gender"`
	BirthDate          time.Time `json:"birthDate" gorm:"type:DATE;default:null"`
	IsVerified         bool      `json:"isVerified" gorm:"default:false"`
	Role               Role      `json:"role" gorm:"default:null"`
	FilledPersonalInfo bool      `json:"filledPersonalInfo" gorm:"default:false"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `json:"isStudentVerified" gorm:"default:false"`
}
