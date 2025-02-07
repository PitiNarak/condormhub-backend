package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserService  ports.UserService
	EmailService ports.EmailServicePort
	Config       *config.AppConfig
}

func NewUserHandler(UserService ports.UserService, emailService ports.EmailServicePort, config *config.AppConfig) *UserHandler {
	return &UserHandler{UserService: UserService, EmailService: emailService, Config: config}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	user := new(domain.User)
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, create_err := h.UserService.Create(user)
	if create_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": create_err.Error()})
	}

	err = h.EmailService.SendVerificationEmail(user.Email, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email verification"})
	}

	return nil
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Params("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"erorr": "Invalid claims"})
	}

	userIDstr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err)
	}

	if err := h.UserService.VerifyUser(userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("User not found")
	}

	return nil
}

func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(domain.Reset_password_body)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error()) //Json
	}
	user, err := h.UserService.ResetPasswordCreate(body.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("User not found")
	}
	err = h.EmailService.SendResetPasswordEmail(user.Email, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send reset password email"})
	}
	return nil
}

func (h *UserHandler) ResetPasswordRespond(c *fiber.Ctx) error {
	body := new(domain.Repond_reset_password_body)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.JWTSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"erorr": "Invalid claims"})
	}
	userIDstr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse userId"})
	}
	err = h.UserService.ResetPasswordRespond(userID, body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}
	return nil
}
