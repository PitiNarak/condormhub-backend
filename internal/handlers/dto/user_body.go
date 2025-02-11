package dto

import "time"

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

type ResponseResetPasswordBody struct {
	Password string `json:"password" validate:"required"`
}

type UserInformationRequestBody struct {
	Username        string    `json:"username" gorm:"unique"`
	Password        string    `json:"password" validate:"omitempty,min=8"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	NationalID      string    `json:"nationalID"`
	Gender          string    `json:"gender"`
	BirthDate       time.Time `json:"birthDate"`
	StudentEvidence string    `json:"studentEvidence"`
}
