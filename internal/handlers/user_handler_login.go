package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
// @Success 200 {object} dto.SuccessResponse[dto.TokenWithUserInformationResponseBody] "user successfully logged in"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "user not found"
// @Failure 500 {object} dto.ErrorResponse "system cannot login user"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	user, accessToken, refreshToken, loginErr := h.userService.Login(c.Context(), req.Email, req.Password)
	if loginErr != nil {
		return loginErr
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: *user,
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}
