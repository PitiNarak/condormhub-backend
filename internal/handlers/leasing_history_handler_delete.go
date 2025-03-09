package handlers

import (
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Delete godoc
// @Summary Delete a leasing history
// @Description Delete a leasing history in the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingHistoryId"
// @Success 200  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "Delete successfully"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "Incorrect UUID format"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "Can not parse UUID or Failed to delete leasing history"
// @Router /history/{id} [delete]
func (h *LeasingHistoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingHistoryID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Delete(leasingHistoryID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
