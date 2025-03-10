package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UpdateUserInformation godoc
// @Summary Update user information
// @Description Update user information
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body dto.UserInformationRequestBody true "user information"
// @Success 200 {object} dto.SuccessResponse[domain.User] "user successfully updated account information"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "system cannot update your account information"
// @Router /user [patch]
func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var requestBody *dto.UserInformationRequestBody

	user := c.Locals("user").(*domain.User)
	if user == nil {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := c.BodyParser(&requestBody)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()
	lifeStyleErr := validate.RegisterValidation("lifestyle", utils.ValidateLifestyles)
	if lifeStyleErr != nil {
		return apperror.BadRequestError(lifeStyleErr, "your lifestyle-tag is incorrect format")
	}

	phoneNumberErr := validate.RegisterValidation("phoneNumber", utils.ValidatePhone)
	if phoneNumberErr != nil {
		return apperror.BadRequestError(phoneNumberErr, "your phone number is incorrect format")
	}

	if err := validate.Struct(requestBody); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.userService.UpdateInformation(user.ID, *requestBody)

	if err != nil {
		return apperror.InternalServerError(err, "system cannot update your account information")
	}

	res := dto.Success(userInfo)

	return c.Status(fiber.StatusOK).JSON(res)

}
