package domain

type ResetPasswordBody struct {
	Email string `json:"email"`
}

type RespondResetPasswordBody struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
