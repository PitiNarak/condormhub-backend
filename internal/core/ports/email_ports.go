package ports

import "github.com/google/uuid"

type EmailServicePort interface {
	SendVerificationEmail(email, name string, token string) error
	SendResetPasswordEmail(email, name string, userID uuid.UUID) error
}
