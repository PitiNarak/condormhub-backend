package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService ports.UserService
}

func NewUserHandler(UserService ports.UserService) *UserHandler {
	return &UserHandler{UserService: UserService}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user domain.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	create_err := h.UserService.Create(user)
	if create_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": create_err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true})

}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	token, loginErr := h.UserService.Login(req.Email, req.Password)
	if loginErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(loginErr.Error())
	}
	return c.Status(200).JSON(fiber.Map{"token": token})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user domain.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	updatedUser, updateErr := h.UserService.Update(user)
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(updateErr.Error())
	}

	return c.Status(200).JSON(fiber.Map{"user": updatedUser})
}
