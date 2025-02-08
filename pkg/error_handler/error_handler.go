package error_handler

import "github.com/gofiber/fiber/v2"

type ErrorHandler struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     string `json:"error"`
}

func NewErrorHandler(code int, message string, err error) *ErrorHandler {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return &ErrorHandler{
		Code:    code,
		Message: message,
		Err:     errMsg,
	}
}

func (e *ErrorHandler) Error() string {
	return e.Message
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
