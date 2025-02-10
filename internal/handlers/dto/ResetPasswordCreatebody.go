package dto

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}
