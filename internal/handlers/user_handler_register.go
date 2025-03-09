package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequestBody true "user information"
// @Success 201  {object}  httpResponse.HttpResponse{data=dto.TokenWithUserInformationResponseBody,pagination=nil} "user successfully registered"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "system cannot register user"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequestBody)
	err := c.BodyParser(&user)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	gormUser := &domain.User{
		Email:    user.Email,
		Username: user.UserName,
		Password: user.Password,
	}

	accessToken, refreshToken, err := h.userService.Create(c.Context(), gormUser)
	if err != nil {
		return err
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: *gormUser,
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusCreated).JSON(res)

}
