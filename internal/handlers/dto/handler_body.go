package dto

type ResetPasswordRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}
