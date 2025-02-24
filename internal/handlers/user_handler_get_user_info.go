package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/gofiber/fiber/v2"
)

// GetUserInfo godoc
// @Summary Get user information
// @Description Get user information
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} httpResponse.HttpResponse{data=domain.User, pagination=nil} "get user information successfully"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil} "system cannot get user information"
// @Router /user/me [get]
func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("get user information successfully", user))

}
