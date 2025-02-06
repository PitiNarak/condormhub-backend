package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserService  ports.UserService
	EmailService ports.EmailServicePort
	Config       *config.Config
}

func NewUserHandler(UserService ports.UserService, emailService ports.EmailServicePort, config *config.Config) *UserHandler {
	return &UserHandler{UserService: UserService, EmailService: emailService, Config: config}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	user := new(domain.User)
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, create_err := h.UserService.Create(*user)
	if create_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": create_err.Error()})
	}

	err = h.EmailService.SendVerificationEmail(user.Email, user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email verification"})
	}

	return nil
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Params("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWTSecretKey), nil
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
