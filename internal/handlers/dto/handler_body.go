package dto

type ResetPasswordRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyRequestBody struct {
	Token string `json:"token" validate:"required"`
}
