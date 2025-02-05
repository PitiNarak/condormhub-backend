package services

import (
	"github.com/go-gomail/gomail"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (e *EmailService) SendPasswordResetEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Password Reset")

	// Set the verification URL with the token
	message.SetBody("text/plain", "Click the link to verify your account: http://localhost:3000/verify/"+token)

	// Set up the dialer with SMTP server credentials
	dailer := gomail.NewDialer(e.Host, e.Port, e.Email, e.Password)

	// Send the email
	return dailer.DialAndSend(message)
}

type User struct {
	gorm.Model
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type reset_password_body struct {
	Email string `json:"email"`
}

func GetUser(db *gorm.DB, email string) *User {
	var user User
	result := db.First(&user, "email = ?", email)
	if result.Error != nil {
		// log.Fatalf("Error finding user: %v", result.Error)
		return nil
	}
	return &user
}

func ResetPassword(c *fiber.Ctx, db *gorm.DB) error {
	body := new(reset_password_body)
	emailService := &EmailService{
		Email:    "pmanoret@gmail.com",  // Your email address
		Host:     "smtp.gmail.com",      // SMTP server host
		Port:     587,                   // SMTP server port
		Password: "zpby yooa zfje ddvk", // SMTP email password
	}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user := GetUser(db, body.Email)
	if user != nil { //we will return oke even though user is not found
		err := emailService.SendPasswordResetEmail(user.Email, "sample-verification-token")
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
