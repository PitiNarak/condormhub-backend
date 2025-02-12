package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
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
// @Success 201  {object}  http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "user successfully registered"
// @Failure 400  {object}  http_response.HttpResponse{data=nil} "your request is invalid"
// @Failure 500  {object}  http_response.HttpResponse{data=nil} "system cannot register user"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequestBody)
	err := c.BodyParser(&user)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	gormUser := &domain.User{
		Email:    user.Email,
		Username: user.UserName,
		Password: user.Password,
	}

	token, err := h.userService.Create(gormUser)
	if err != nil {
		return err
	}

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     token,
		UserInformation: *gormUser,
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("user successfully registered", response))

}
