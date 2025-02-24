package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary Refresh user
// @Description Refresh user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RefreshTokenRequestBody true "user information"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.TokenResponseBody} "user successfully Refresh in"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil} "user not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil} "system cannot refresh user"
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return errorHandler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return errorHandler.BadRequestError(err, "your request body is incorrect")
	}

	accessToken, refreshToken, loginErr := h.userService.RefreshToken(c.Context(), req.RefreshToken)
	if loginErr != nil {
		return loginErr
	}

	response := dto.TokenResponseBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("refresh successful", response))
}
