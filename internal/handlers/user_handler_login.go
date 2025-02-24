package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequestBody true "user information"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.TokenWithUserInformationResponseBody, pagination=nil} "user successfully logged in"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil} "user not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil} "system cannot login user"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return errorHandler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return errorHandler.BadRequestError(err, "your request body is incorrect")
	}

	user, accessToken, refreshToken, loginErr := h.userService.Login(c.Context(), req.Email, req.Password)
	if loginErr != nil {
		return loginErr
	}

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: *user,
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Login successful", response, nil))
}
