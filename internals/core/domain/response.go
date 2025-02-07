package domain

type ResetPasswordBody struct {
	Email string `json:"email" validate:"required"`
}

type ResponseResetPasswordBody struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}
