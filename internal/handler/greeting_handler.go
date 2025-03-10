package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
)

type GreetingHandler struct {
}

func NewGreetingHandler() *GreetingHandler {
	return &GreetingHandler{}
}

func (e *GreetingHandler) Greeting(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(dto.Success("Hello from CondromHub Api."))
}
