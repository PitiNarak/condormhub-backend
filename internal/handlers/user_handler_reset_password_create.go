package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resend verification email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordCreateRequestBody true "token"
// @Success 200 {object} http_response.HttpResponse{data=nil} "email is sent to user successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot resend verification email"
// @Router /user/resetpassword [post]
func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordCreateRequestBody)

	if err := c.BodyParser(body); err != nil {
		return errorHandler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return errorHandler.BadRequestError(err, "your request body is incorrect")
	}

	err := h.userService.ResetPasswordCreate(body.Email)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("email is sent to user successfully", nil))
}
