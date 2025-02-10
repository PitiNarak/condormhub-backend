package error_handler

import "github.com/gofiber/fiber/v2"

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
