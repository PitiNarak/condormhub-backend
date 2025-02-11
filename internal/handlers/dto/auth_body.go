package dto

type RegisterRequestBody struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required"`
}

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
