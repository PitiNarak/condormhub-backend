package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type GreetingHandler struct {
}

func NewGreetingHandler() *GreetingHandler {
	return &GreetingHandler{}
}

func (e *GreetingHandler) Greeting(c *fiber.Ctx) error {
	err := c.Send([]byte("Hello, World!"))
	if err != nil {
		return err
	}
	return nil
}
