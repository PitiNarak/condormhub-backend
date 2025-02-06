package ports

import "github.com/google/uuid"

type EmailServicePort interface {
	SendVerificationEmail(email string, userID uuid.UUID) error
}
