package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResetPasswordCreateRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequestBody struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type VerifyRequestBody struct {
	Token string `json:"token" validate:"required"`
}

type UserInformationRequestBody struct {
	Username        string    `json:"username,omitempty" validate:"omitempty,min=2"`
	Password        string    `json:"password,omitempty" validate:"omitempty,min=8"`
	Firstname       string    `json:"firstname,omitempty" validate:"omitempty,min=2"`
	Lastname        string    `json:"lastname,omitempty" validate:"omitempty,min=2"`
	NationalID      string    `json:"nationalID,omitempty" validate:"omitempty,len=13"`
	Gender          string    `json:"gender,omitempty"`
	BirthDate       time.Time `json:"birthDate,omitempty"`
	StudentEvidence string    `json:"studentEvidence,omitempty"`
	Lifestyles      []string  `json:"lifestyles,omitempty" validate:"omitempty,lifestyle"`
	PhoneNumber     string    `json:"phoneNumber,omitempty" validate:"omitempty,phoneNumber"`
}

type UserResponse struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	Gender             string    `json:"gender"`
	BirthDate          time.Time `json:"birthDate"`
	IsVerified         bool      `json:"isVerified"`
	Role               string    `json:"role"`
	FilledPersonalInfo bool      `json:"filledPersonalInfo"`
	Lifestyles         []string  `json:"lifestyles"`
	PhoneNumber        string    `json:"phoneNumber"`
	StudentEvidence    string    `json:"studentEvidence"`
	IsStudentVerified  bool      `json:"isStudentVerified"`
}
