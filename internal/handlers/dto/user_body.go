package dto

import "time"

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

type ResponseResetPasswordBody struct {
	Password string `json:"password" validate:"required"`
}

type ResponseVerifyBody struct {
	Username           string `json:"username" gorm:"unique" validate:"required"`
	Email              string `json:"email" gorm:"unique" validate:"required,email"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	NationalID         string `json:"nationalID" `
	Gender             string `json:"gender"`
	BirthDate          string `json:"birthDate"`
	IsVerified         bool   `gorm:"default:false" json:"isVerified"`
	Role               string `json:"role"`
	FilledPersonalInfo bool   `gorm:"default:false" json:"filledPersonalInfo"`
	StudentEvidence    string `json:"studentEvidence"`
	IsStudentVerified  bool   `gorm:"default:false" json:"isStudentVerified"`
}

type UserInformationRequestBody struct {
	Username        string    `json:"username" gorm:"unique"`
	Password        string    `json:"password" validate:"omitempty,min=8"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	NationalID      string    `json:"nationalID"`
	Gender          string    `json:"gender"`
	BirthDate       time.Time `json:"birthDate"`
	StudentEvidence string    `json:"studentEvidence"`
}
