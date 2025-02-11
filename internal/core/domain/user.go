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
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreateAt           time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt           time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	Username           string    `json:"username" gorm:"unique" validate:"required"`
	Password           string    `json:"-" validate:"required,min=8"`
	Email              string    `json:"email" gorm:"unique" validate:"required,email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	NationalID         string    `json:"nationalID" `
	Gender             string    `json:"gender"`
	BirthDate          time.Time `json:"birthDate"`
	IsVerified         bool      `gorm:"default:false" json:"isVerified"`
	Role               Role      `gorm:"default:null" json:"role"`
	FilledPersonalInfo bool      `gorm:"default:false" json:"filledPersonalInfo"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `gorm:"default:false" json:"isStudentVerified"`
}
