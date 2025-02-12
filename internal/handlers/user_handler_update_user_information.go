package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/go-playground/validator"
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
// @Success 200 {object} http_response.HttpResponse{data=domain.User} "user successfully updated account information"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot update your account information"
// @Router /user [patch]
func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var requestBody *dto.UserInformationRequestBody

	user := c.Locals("user").(*domain.User)
	if user == nil {
		return error_handler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := c.BodyParser(&requestBody)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(requestBody); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.userService.UpdateInformation(user.ID, *requestBody)

	if err != nil {
		return error_handler.InternalServerError(err, "system cannot update your account information")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("user successfully updated account information", userInfo))

}
