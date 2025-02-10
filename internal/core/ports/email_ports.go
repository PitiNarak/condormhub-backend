package ports

type EmailServicePort interface {
	SendVerificationEmail(email, name string, token string) error
	SendResetPasswordEmail(email, name string, token string) error
}
