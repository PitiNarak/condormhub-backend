package handlers

import (
	"os"

	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
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
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Hello, welcome to CondormHub!", map[string]string{"env": env}))
	// return error_handler.InternalServerError(errors.New("error from system"), "your error message")
}
