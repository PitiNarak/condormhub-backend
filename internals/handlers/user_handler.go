package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService ports.UserService
}

func NewUserHandler(UserService ports.UserService) *UserHandler {
	return &UserHandler{UserService: UserService}
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

	err = h.UserService.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return nil
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Params("token")

	if err := h.UserService.VerifyUser(tokenString); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(domain.ResetPasswordBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error()) //Json
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.UserService.ResetPasswordCreate(body.Email)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return nil
}

func (h *UserHandler) ResetPasswordResponse(c *fiber.Ctx) error {
	body := new(domain.ResponseResetPasswordBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.UserService.ResetPasswordResponse(body.Token, body.Password)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return nil
}
