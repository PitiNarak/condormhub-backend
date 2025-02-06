package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/go-gomail/gomail"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type reset_password_body struct {
	Email string `json:"email"`
}

type repond_reset_password_body struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (e *EmailService) SendPasswordResetEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Config.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Password Reset")

	// Set the verification URL with the token
	message.SetBody("text/plain", "Click the link to verify your account: http://localhost:3000/verify/"+token)

	// Set up the dialer with SMTP server credentials
	dailer := gomail.NewDialer(e.Config.Host, e.Config.Port, e.Config.Email, e.Config.Password)

	// Send the email
	return dailer.DialAndSend(message)
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

func ResetPassword(db *gorm.DB, id int, newPasswod string) error {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		// log.Fatalf("Error finding user: %v", result.Error)
		return result.Error
	}
	pass := []byte("mypassword")
	hp, err := bcrypt.GenerateFromPassword(pass, 0)
	if err != nil {
		return err
	}
	user.Password = string(hp)
	result = db.Save(user)
	if result.Error != nil {
		// log.Fatalf("Error finding user: %v", result.Error)
		return result.Error
	}
	return nil
}

func SendingResetPasswordEmail(c *fiber.Ctx, db *gorm.DB) error {
	body := new(reset_password_body)
	emailService := &EmailService{Config: config.Load()}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user := GetUser(db, body.Email)
	if user != nil { //we will return oke even though user is not found
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		secretKey := "secret" // Fix later
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID

		t, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		err = emailService.SendPasswordResetEmail(user.Email, t)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func RespondHandler(c *fiber.Ctx, db *gorm.DB) error {
	body := new(repond_reset_password_body)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	tokenString := body.Token
	secretKey := "secret"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(float64)
		ResetPassword(db, int(userID), body.Password)
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token claims")
}
