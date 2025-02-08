package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
)

type SMTPConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	Email    string `env:"EMAIL,required"`
	Password string `env:"PASSWORD,required"`
}

type EmailService struct {
	EmailConfig *SMTPConfig
	JWTConfig   *utils.JWTConfig
}

func NewEmailService(emailConfig *SMTPConfig, jwtConfig *utils.JWTConfig) ports.EmailServicePort {
	return &EmailService{EmailConfig: emailConfig, JWTConfig: jwtConfig}
}

func (e *EmailService) SendVerificationEmail(email string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, e.JWTConfig)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", e.EmailConfig.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")

	verLink := fmt.Sprintf("http://localhost:3000/verify/%s", token)
	body := fmt.Sprintf("<html><body><p>Click the link to verify your account: </p><a href='%s'>verify</a></body></html>", verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.EmailConfig.Host, e.EmailConfig.Port, e.EmailConfig.Email, e.EmailConfig.Password)

	return dailer.DialAndSend(message)
}

func (e *EmailService) SendResetPasswordEmail(email string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, e.JWTConfig)
	if err != nil {
		return err
	}
	message := gomail.NewMessage()
	message.SetHeader("From", e.EmailConfig.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Reset Password")

	verLink := fmt.Sprintf("http://localhost:3000/resetpassword/%s", token)
	body := fmt.Sprintf("<html><body><p>Click the link to reset your password: </p><a href='%s'></a></body></html>", verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.EmailConfig.Host, e.EmailConfig.Port, e.EmailConfig.Email, e.EmailConfig.Password)

	return dailer.DialAndSend(message)
}
