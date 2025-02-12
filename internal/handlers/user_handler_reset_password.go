package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordRequestBody true "token"
// @Success 200 {object} http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "password reset successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot reset password"
// @Router /user/newpassword [post]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordRequestBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	tokenString := body.Token
	if tokenString == "" {
		return error_handler.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	user, err := h.userService.ResetPassword(tokenString, body.Password)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("password reset successfully", fiber.Map{
		"userInformation": user,
		"accessToken":     tokenString,
	}))
}
