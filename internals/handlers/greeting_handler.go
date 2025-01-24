package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

type GreetingHandler struct {
}

func NewGreetingHandler() *GreetingHandler {
	return &GreetingHandler{}
}

func (e *GreetingHandler) Greeting(c *fiber.Ctx) error {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "unknown"
	}
	return c.JSON(fiber.Map{"name": "CondormHub API", "env": env})
}
