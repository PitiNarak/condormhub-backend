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
	Config *config.AppConfig
}

func NewEmailService(config *config.AppConfig) ports.EmailServicePort {
	return &EmailService{Config: config}
}

func (e *EmailService) SendVerificationEmail(email string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, &e.Config.JWT)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", e.Config.SMTP.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")

	verLink := fmt.Sprintf("http://localhost:3000/verify/%s", token)
	body := fmt.Sprintf("<html><body><p>Click the link to verify your account: </p><a href='%s'>verify</a></body></html>", verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.Config.SMTP.Host, e.Config.SMTP.Port, e.Config.SMTP.Email, e.Config.SMTP.Password)

	return dailer.DialAndSend(message)
}

func (e *EmailService) SendResetPasswordEmail(email string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, &e.Config.JWT)
	if err != nil {
		return err
	}
	message := gomail.NewMessage()
	message.SetHeader("From", e.Config.SMTP.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Reset Password")

	verLink := fmt.Sprintf("http://localhost:3000/resetpassword/%s", token)
	body := fmt.Sprintf("<html><body><p>Click the link to reset your password: </p><a href='%s'></a></body></html>", verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.Config.SMTP.Host, e.Config.SMTP.Port, e.Config.SMTP.Email, e.Config.SMTP.Password)

	return dailer.DialAndSend(message)
}
