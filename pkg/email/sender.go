package email

import (
	"fmt"
	"os"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/go-gomail/gomail"
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
	jwtUtils    *jwt.JWTUtils
}

func NewEmailService(emailConfig *SMTPConfig, jwtUtils *jwt.JWTUtils) ports.EmailServicePort {
	return &EmailService{emailConfig: emailConfig, jwtUtils: jwtUtils}
}

func (e *EmailService) SendVerificationEmail(email, name string, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")
	cwd, _ := os.Getwd()
	verLink := fmt.Sprintf(e.emailConfig.LinkHostname+"/verify?token=%s", token)
	html, _ := readTemplate(cwd + "/pkg/html_template/verify-compress.html")
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.emailConfig.Host, e.emailConfig.Port, e.emailConfig.Email, e.emailConfig.Password)

	return dailer.DialAndSend(message)
}

func (e *EmailService) SendResetPasswordEmail(email, name string, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@condormhub.xyz")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Reset Password")

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.InternalServerError(err, "cannot get current path")
	}
	verLink := fmt.Sprintf(e.emailConfig.LinkHostname+"/newpassword/token=%s", token)
	html, err := readTemplate(cwd + "/pkg/html_template/reset-compress.html")
	if err != nil {
		return apperror.InternalServerError(err, "cannot load html template")
	}
	body := fmt.Sprintf(html, name, verLink, verLink)
	message.SetBody("text/html", body)

	dailer := gomail.NewDialer(e.emailConfig.Host, e.emailConfig.Port, e.emailConfig.Email, e.emailConfig.Password)
	err = dailer.DialAndSend(message)
	if err != nil {
		return apperror.InternalServerError(err, "cannot sent email")
	}

	return nil
}

func readTemplate(path string) (string, error) {
	byteContent, err := os.ReadFile(path)
	if err != nil { //many people wrap this into a function
		return "", err
	}
	return string(byteContent), nil
}
