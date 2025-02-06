package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/PitiNarak/condormhub-backend/internals/core/utils"
	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
)

type EmailService struct {
	Config *config.Config
}

func NewEmailService(config *config.Config) ports.EmailServicePort {
	return &EmailService{Config: config}
}

func (e *EmailService) SendVerificationEmail(email string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, e.Config)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", e.Config.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")

	// have not yet implement API for /verify
	verLink := fmt.Sprintf("http://localhost:3000/verify/%s", token)
	body := fmt.Sprintf("<html><body><p>Click the link to verify your account: </p><a href='%s'></a></body></html>", verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.Config.Host, e.Config.Port, e.Config.Email, e.Config.Password)

	return dailer.DialAndSend(message)
}
