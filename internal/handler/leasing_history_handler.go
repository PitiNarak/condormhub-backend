package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingHistoryHandler struct {
	service ports.LeasingHistoryService
}

func NewLeasingHistoryHandler(service ports.LeasingHistoryService) ports.LeasingHistoryHandler {
	return &LeasingHistoryHandler{service: service}
}

// SetEndTimestamp godoc
// @Summary Delete a leasing history
// @Description Delete a leasing history in the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingHistoryId"
// @Success 204 "Set end timestamp successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to update leasing history"
// @Router /history/{id} [patch]
func (h *LeasingHistoryHandler) SetEndTimestamp(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingHistoryID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.SetEndTimestamp(leasingHistoryID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetByUserID godoc
// @Summary Get all leasing history by userid
// @Description Retrieve a list of all leasing history by userid
// @Tags history
// @Security Bearer
// @Produce json
// @Param limit query string true "Number of history to be retirved"
// @Param page query string true "Page to retrive"
// @Success 200 {object} dto.PaginationResponse[domain.LeasingHistory] "Retrive history successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Router /history/me [get]
func (h *LeasingHistoryHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	leasingHistory, totalPage, totalRows, err := h.service.GetByUserID(userID, limit, page)
	if err != nil {
		return err
	}

	res := dto.SuccessPagination(leasingHistory, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByDormID godoc
// @Summary Get all leasing history by userid
// @Description Retrieve a list of all leasing history by userid
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Param limit query string true "Number of history to be retirved"
// @Param page query string true "Page to retrive"
// @Success 200 {object} dto.PaginationResponse[domain.LeasingHistory] "Retrive history successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID"
// @Router /history/bydorm/{id} [get]
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	id := c.Params("id")
	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	leasingHistory, totalPage, totalRows, err := h.service.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	// response := dto.PaginationResponseBody{
	// 	Currentpage: page,
	// 	Lastpage:    totalPage,
	// 	Limit:       limit,
	// 	Total:       totalRows,
	// }

	res := dto.SuccessPagination(leasingHistory, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// Delete godoc
// @Summary Delete a leasing history
// @Description Delete a leasing history in the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingHistoryId"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to delete leasing history"
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
