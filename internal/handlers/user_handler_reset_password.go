package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
// @Success 200 {object} httpResponse.HttpResponse{data=dto.TokenWithUserInformationResponseBody,pagination=nil} "password reset successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "system cannot reset password"
// @Router /user/newpassword [post]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordRequestBody)

	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	tokenString := body.Token
	if tokenString == "" {
		return apperror.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	user, err := h.userService.ResetPassword(c.Context(), tokenString, body.Password)
	if err != nil {
		return err
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     tokenString,
		RefreshToken:    "",
		UserInformation: *user,
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}
