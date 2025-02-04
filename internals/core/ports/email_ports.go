package ports

type EmailServicePort interface {
	SendVerificationEmail(email, token string) error
}
