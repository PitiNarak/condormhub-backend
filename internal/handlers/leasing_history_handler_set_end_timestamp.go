package handlers

import (
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SetEndTimestamp godoc
// @Summary Delete a leasing history
// @Description Delete a leasing history in the database
// @Tags history
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "LeasingHistoryId"
// @Success 200  {object}  httpResponse.HttpResponse{data=nil} "Set end timestamp successfully"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil} "Incorrect UUID format"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil} "leasing history not found"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil} "Can not parse UUID or Failed to update leasing history"
// @Router /history/end/{id} [post]
func (h *LeasingHistoryHandler) SetEndTimestamp(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return errorHandler.BadRequestError(err, "Incorrect UUID format")
	}

	leasingHistoryID, err := uuid.Parse(id)
	if err != nil {
		return errorHandler.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.SetEndTimestamp(leasingHistoryID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Set end timestamp successfully", nil))
}
