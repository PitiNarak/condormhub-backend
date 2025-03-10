package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Create godoc
// @Summary Create a new leasing history
// @Description Add a new leasing history to the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 201 {object} dto.SuccessResponse[domain.LeasingHistory] "Dorm successfully created"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found or leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/{id} [post]
func (h *LeasingHistoryHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	leasingHistory, err := h.service.Create(userID, dormID)
	if err != nil {
		return err
	}

	res := dto.Success(leasingHistory)

	return c.Status(fiber.StatusCreated).JSON(res)
}
