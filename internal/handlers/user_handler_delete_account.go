package handlers

import (
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteAccount godoc
// @Summary Delete a user account
// @Description Delete a user account
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} httpResponse.HttpResponse{data=nil} "account successfully deleted"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil} "cannot parse uuid or cannot delete user"
// @Router /user/ [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.userService.DeleteAccount(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("account successfully deleted", nil))
}
