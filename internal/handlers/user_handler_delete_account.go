package handlers

import (
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
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
// @Success 200 {object} http_response.HttpResponse{data=nil} "account successfully deleted"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "cannot parse uuid or cannot delete user"
// @Router /user/ [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userIDstr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return error_handler.InternalServerError(err, "cannot parse uuid")
	}

	err = h.userService.DeleteAccount(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("account successfully deleted", nil))
}
