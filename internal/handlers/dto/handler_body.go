package dto

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

type ResponseResetPasswordBody struct {
	Password string `json:"password" validate:"required"`
}

type RegisterRequestBody struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required"`
}
