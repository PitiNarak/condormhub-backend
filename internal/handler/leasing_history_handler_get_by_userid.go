package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
