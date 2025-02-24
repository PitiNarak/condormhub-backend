package handlers

import (
	"errors"
	"strconv"

	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetByDormID godoc
// @Summary Get all leasing history by userid
// @Description Retrieve a list of all leasing history by userid
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Param limit query string true "Number of history to be retirved"
// @Param page query string true "Page to retrive"
// @Success 200 {object} httpResponse.HttpResponse{data=[]domain.LeasingHistory} "Retrive history successfully"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil} "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil} "leasing history not found"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil} "Can not parse UUID"
// @Router /history/bydorm/{id} [get]
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	id := c.Params("id")
	dormID, err := uuid.Parse(id)
	if err != nil {
		return errorHandler.InternalServerError(err, "Can not parse UUID")
	}
	params := c.Queries()
	limitStr, ok := params["limit"]
	if !ok {
		return errorHandler.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return errorHandler.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	pageStr, ok := params["page"]
	if !ok {
		return errorHandler.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return errorHandler.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	leasingHistory, totalPage, totalRows, err := h.service.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	response := dto.GetLeasingHistoryResponseBody{
		Currentpage:    page,
		Lastpage:       totalPage,
		Limit:          limit,
		Total:          totalRows,
		LeasingHistory: leasingHistory,
	}
	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Retrive history successfully", response))
}
