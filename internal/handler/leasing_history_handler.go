package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/go-playground/validator"
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
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingHistoryID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.SetEndTimestamp(leasingHistoryID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
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
// @Param limit query int false "Number of dorms to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Success 200 {object} dto.PaginationResponse[dto.LeasingHistory] "Retrive history successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Router /history/me [get]
func (h *LeasingHistoryHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	leasingHistory, totalPage, totalRows, err := h.service.GetByUserID(userID, limit, page)
	if err != nil {
		return err
	}

	resData := make([]dto.LeasingHistory, len(leasingHistory))
	for i, v := range leasingHistory {
		resData[i] = v.ToDTO()
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
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
// @Param limit query int false "Number of dorms to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Success 200 {object} dto.PaginationResponse[dto.LeasingHistory] "Retrive history successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect or page exceeded"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID"
// @Router /history/bydorm/{id} [get]
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	id := c.Params("id")
	dormID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	leasingHistory, totalPage, totalRows, err := h.service.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	resData := make([]dto.LeasingHistory, len(leasingHistory))
	for i, v := range leasingHistory {
		resData[i] = v.ToDTO()
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
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
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Delete(leasingHistoryID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Create godoc
// @Summary Create a new leasing history
// @Description Add a new leasing history to the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 201 {object} dto.SuccessResponse[dto.LeasingHistory] "Dorm successfully created"
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
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	leasingHistory, err := h.service.Create(userID, dormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}

	res := dto.Success(leasingHistory.ToDTO())

	return c.Status(fiber.StatusCreated).JSON(res)
}

// CreateReview godoc
// @Summary Create a new review
// @Description Add a new review to the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param user body dto.ReviewRequestBody true "review information"
// @Success 201 {object} dto.SuccessResponse[dto.Review]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/review/ [post]
func (h *LeasingHistoryHandler) CreateReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	body := new(dto.ReviewRequestBody)
	err := c.BodyParser(&body)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	review, err := h.service.CreateReview(user, body.ID, body.Message, int(body.Rate))
	if err != nil {
		return err
	}
	res := dto.Success(review.ToDTO())
	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateReview godoc
// @Summary Update a new review
// @Description Update a review in the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param user body dto.ReviewRequestBody true "review information"
// @Success 201 {object} dto.SuccessResponse[dto.Review]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/review/ [patch]
func (h *LeasingHistoryHandler) UpdateReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	body := new(dto.ReviewRequestBody)
	err := c.BodyParser(&body)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	review, err := h.service.UpdateReview(user, body.ID, body.Message, int(body.Rate))
	if err != nil {
		return err
	}
	res := dto.Success(review.ToDTO())
	return c.Status(fiber.StatusCreated).JSON(res)
}

// Delete godoc
// @Summary Delete a review
// @Description Delete a review
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "HistoryID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found or leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/delete/ [delete]
func (h *LeasingHistoryHandler) DeleteReview(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*domain.User)
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingHistoryID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.DeleteReview(user, leasingHistoryID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
