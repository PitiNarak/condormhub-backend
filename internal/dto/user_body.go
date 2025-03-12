package dto

import (
	"time"
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
