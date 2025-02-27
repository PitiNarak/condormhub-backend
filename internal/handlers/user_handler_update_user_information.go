package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
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
// @Success 200 {object} httpResponse.HttpResponse{data=domain.User,pagination=nil} "user successfully updated account information"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "system cannot update your account information"
// @Router /user [patch]
func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var requestBody *dto.UserInformationRequestBody

	user := c.Locals("user").(*domain.User)
	if user == nil {
		return errorHandler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := c.BodyParser(&requestBody)
	if err != nil {
		return errorHandler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()
	lifestyle_err := validate.RegisterValidation("lifestyle", utils.ValidateLifestyles)
	if lifestyle_err != nil {
		return errorHandler.BadRequestError(err, "your lifestyle-tag is incorrect format")
	}

	phoneNumber_err := validate.RegisterValidation("phoneNumber", utils.ValidatePhone)
	if phoneNumber_err != nil {
		return errorHandler.BadRequestError(phoneNumber_err, "your phone number is incorrect format")
	}

	if err := validate.Struct(requestBody); err != nil {
		return errorHandler.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.userService.UpdateInformation(user.ID, *requestBody)

	if err != nil {
		return errorHandler.InternalServerError(err, "system cannot update your account information")
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("user successfully updated account information", userInfo))

}
