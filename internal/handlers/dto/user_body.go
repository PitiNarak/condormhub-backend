package dto

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
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
	Username        string                `json:"username" gorm:"unique"`
	Password        string                `json:"password" validate:"omitempty,min=8"`
	Firstname       string                `json:"firstname"`
	Lastname        string                `json:"lastname"`
	NationalID      string                `json:"nationalID"`
	Gender          string                `json:"gender"`
	BirthDate       time.Time             `json:"birthDate"`
	StudentEvidence string                `json:"studentEvidence"`
	Lifestyles      domain.LifestyleArray `json:"lifestyles" validate:"lifestyle" gorm:"type:lifestyle_tag[]"`
}
