package services

import (
	"fmt"
	"os"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
)

type SMTPConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	Email    string `env:"EMAIL,required"`
	Password string `env:"PASSWORD,required"`
	Hostname string `env:"LINK_HOSTNAME,required"`
}

type EmailService struct {
	EmailConfig *SMTPConfig
	JWTConfig   *utils.JWTConfig
}

func NewEmailService(emailConfig *SMTPConfig, jwtConfig *utils.JWTConfig) ports.EmailServicePort {
	return &EmailService{EmailConfig: emailConfig, JWTConfig: jwtConfig}
}

func (e *EmailService) SendVerificationEmail(email, name string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, e.JWTConfig)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")
	cwd, _ := os.Getwd()
	verLink := fmt.Sprintf(e.EmailConfig.Hostname+"/verify?token=%s", token)
	html, _ := utils.ReadTemplate(cwd + "/pkg/html_template/verify-compress.html")
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.EmailConfig.Host, e.EmailConfig.Port, e.EmailConfig.Email, e.EmailConfig.Password)

	return dailer.DialAndSend(message)
}

func (e *EmailService) SendResetPasswordEmail(email, name string, userID uuid.UUID) error {
	token, err := utils.GenerateJWT(userID, e.JWTConfig)
	if err != nil {
		return err
	}
	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Reset Password")

	cwd, _ := os.Getwd()
	verLink := fmt.Sprintf(e.EmailConfig.Hostname+"/newpassword/token=%s", token)
	html, _ := utils.ReadTemplate(cwd + "/pkg/html_template/reset-compress.html")
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.EmailConfig.Host, e.EmailConfig.Port, e.EmailConfig.Email, e.EmailConfig.Password)

	return dailer.DialAndSend(message)
}
