package error_handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ErrorHandler struct {
	Code    int
	Message string
	Err     error
}

func NewErrorHandler(code int, message string, err error) *ErrorHandler {
	return &ErrorHandler{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *ErrorHandler) Error() string {
	// handle pass nil error
	if e.Err == nil {
		log.Warnf("your error is nil. Passed error message: %s", e.Message)
		return e.Message
	}
	return e.Err.Error()
}

func InternalServerError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusInternalServerError, msg, err)
}

func BadRequestError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusBadRequest, msg, err)
}

func UnauthorizedError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusUnauthorized, msg, err)
}

func ForbiddenError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusForbidden, msg, err)
}

func NotFoundError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusNotFound, msg, err)
}

func ConflictError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusConflict, msg, err)
}

func UnprocessableEntityError(err error, msg string) *ErrorHandler {
	return NewErrorHandler(fiber.StatusUnprocessableEntity, msg, err)
}

func Handler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *ErrorHandler
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else {
		message = err.Error()
	}

	if e != nil && e.Err != nil {
		log.Errorf("Error: %v, Code: %d, Message: %s", e.Error(), code, message)
	} else {
		log.Errorf("Error: %s, Code: %d, Message: %s", err.Error(), code, message)
	}

	return c.Status(code).JSON(&http_response.HttpResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
