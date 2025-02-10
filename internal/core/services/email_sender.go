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
	Host         string `env:"HOST,required"`
	Port         int    `env:"PORT,required"`
	Email        string `env:"EMAIL,required"`
	Password     string `env:"PASSWORD,required"`
	LinkHostname string `env:"LINK_HOSTNAME,required"`
}

type EmailService struct {
	emailConfig *SMTPConfig
	jwtUtils    *utils.JWTUtils
}

func NewEmailService(emailConfig *SMTPConfig, jwtUtils *utils.JWTUtils) ports.EmailServicePort {
	return &EmailService{emailConfig: emailConfig, jwtUtils: jwtUtils}
}

func (e *EmailService) SendVerificationEmail(email, name string, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")
	cwd, _ := os.Getwd()
	verLink := fmt.Sprintf(e.emailConfig.LinkHostname+"/verify?token=%s", token)
	html, _ := utils.ReadTemplate(cwd + "/pkg/html_template/verify-compress.html")
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.emailConfig.Host, e.emailConfig.Port, e.emailConfig.Email, e.emailConfig.Password)

	return dailer.DialAndSend(message)
}

func (e *EmailService) SendResetPasswordEmail(email, name string, userID uuid.UUID) error {
	token, generateErr := e.jwtUtils.GenerateJWT(userID)
	if generateErr != nil {
		return generateErr
	}
	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Reset Password")

	cwd, _ := os.Getwd()
	verLink := fmt.Sprintf(e.emailConfig.LinkHostname+"/newpassword/token=%s", token)
	html, _ := utils.ReadTemplate(cwd + "/pkg/html_template/reset-compress.html")
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.emailConfig.Host, e.emailConfig.Port, e.emailConfig.Email, e.emailConfig.Password)

	return dailer.DialAndSend(message)
}
