package handlers

import (
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetByDormID godoc
// @Summary Get all leasing history by userid
// @Description Retrieve a list of all leasing history by userid
// @Tags history
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object} httpResponse.HttpResponse{data=[]domain.LeasingHistory} "Retrive history successfully"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil} "Incorrect UUID format"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil} "leasing history not found"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil} "Can not parse UUID"
// @Router /history/bydorm/{id} [get]
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return errorHandler.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return errorHandler.InternalServerError(err, "Can not parse UUID")
	}
	leasingHistory, err := h.service.GetByDormID(dormID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Retrive history successfully", leasingHistory))
}
