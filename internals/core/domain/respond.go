package domain

type Reset_password_body struct {
	Email string `json:"email"`
}

type Repond_reset_password_body struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
