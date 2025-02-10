package dto

import "time"

type UserRequestBody struct {
	Username        string    `json:"username" gorm:"unique"`
	Password        string    `json:"password" validate:"omitempty,min=8"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	NationalID      string    `json:"nationalID"`
	Gender          string    `json:"gender"`
	BirthDate       time.Time `json:"birthDate"`
	StudentEvidence string    `json:"studentEvidence"`
}
