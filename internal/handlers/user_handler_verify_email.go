package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.VerifyRequestBody true "token"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.TokenWithUserInformationResponseBody,pagination=nil} "email is verified successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "system cannot verify your email"
// @Router /user/verify [post]
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	body := new(dto.VerifyRequestBody)

	if err := c.BodyParser(body); err != nil {
		return errorHandler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return errorHandler.BadRequestError(err, "your request body is incorrect")
	}
	accessToken, user, err := h.userService.VerifyUser(c.Context(), body.Token)
	if err != nil {
		return err
	}

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		UserInformation: *user,
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("email is verified successfully", response))
}
