package services

import (
	"github.com/PitiNarak/condormhub-backend/internals/config"
	_ "github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/go-gomail/gomail"
)

type EmailService struct {
	Config *config.Config
}

func NewEmailService(config *config.Config) *EmailService {
	return &EmailService{Config: config}
}

func (e *EmailService) SendVerificationEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Config.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")

	// have not yet implement API for /verify
	message.SetBody("text/plain", "Click the link to verify your account: http://localhost:3000/verify/"+token)

	dailer := gomail.NewDialer(e.Config.Host, e.Config.Port, e.Config.Email, e.Config.Password)

	return dailer.DialAndSend(message)
}
