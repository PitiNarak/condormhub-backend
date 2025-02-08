package domain

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

type ResponseResetPasswordBody struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserBody struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"userName" validate:"required"`
}
