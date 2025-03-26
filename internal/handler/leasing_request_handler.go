package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingRequestHandler struct {
	service ports.LeasingRequestService
}

func NewLeasingRequestHandler(service ports.LeasingRequestService) ports.LeasingRequestHandler {
	return &LeasingRequestHandler{service: service}
}

// Approve godoc
// @Summary Approve a leasing request
// @Description Approve a leasing request
// @Tags request
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingRequestId"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized or lessee cannot approve a request"
// @Failure 404 {object} dto.ErrorResponse "leasing request not found or request is not in the pending status"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to update leasing request"
// @Router /request/{id}/approve [patch]
func (h *LeasingRequestHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*domain.User)
	if user.Role == domain.LesseeRole {
		return apperror.UnauthorizedError(errors.New("user is a lessee"), "lessee cannot approve a request")
	}
	if err := uuid.Validate(id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingRequestID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Approve(leasingRequestID, user.ID, user.Role == domain.AdminRole)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Reject godoc
// @Summary Reject a leasing request
// @Description Reject a leasing request
// @Tags request
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingRequestId"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse "leasing request not found or request is not in the pending status"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized or lessee cannot reject a request"
// @Failure 404 {object} dto.ErrorResponse "leasing request not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to update leasing request"
// @Router /request/{id}/reject [patch]
func (h *LeasingRequestHandler) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*domain.User)
	if user.Role == domain.LesseeRole {
		return apperror.UnauthorizedError(errors.New("user is a lessee"), "lessee cannot reject a request")
	}
	if err := uuid.Validate(id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingRequestID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Reject(leasingRequestID, user.ID, user.Role == domain.AdminRole)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Cancel godoc
// @Summary Cancel a leasing request
// @Description Cancel a leasing request
// @Tags request
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingRequestId"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse "leasing request not found or request is not in the pending status"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized or lessor cannot cancel a request"
// @Failure 404 {object} dto.ErrorResponse "leasing request not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to update leasing request"
// @Router /request/{id}/cancel [patch]
func (h *LeasingRequestHandler) Cancel(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(*domain.User)
	if user.Role == domain.LessorRole {
		return apperror.UnauthorizedError(errors.New("user is a lessor"), "lessor cannot cancel a request")
	}
	if err := uuid.Validate(id); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingRequestID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Cancel(leasingRequestID, user.ID, user.Role == domain.AdminRole)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetByUserID godoc
// @Summary Get all leasing request by userid
// @Description Retrieve a list of all leasing request by userid
// @Tags request
// @Security Bearer
// @Produce json
// @Param limit query int false "Number of dorms to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Success 200 {object} dto.PaginationResponse[dto.LeasingRequest]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or limit parameter is incorrect or page parameter is incorrect"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing request not found"
// @Router /request/me [get]
func (h *LeasingRequestHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
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
	leasingHistory, totalPage, totalRows, err := h.service.GetByUserID(userID, user.Role, limit, page)
	if err != nil {
		return err
	}

	resData := make([]dto.LeasingRequest, len(leasingHistory))
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
// @Summary Delete a leasing request
// @Description Delete a leasing request in the database
// @Tags request
// @Security Bearer
// @Produce json
// @Param id path string true "LeasingRequestId"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized or only admin can delete a request"
// @Failure 404 {object} dto.ErrorResponse "leasing request not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or Failed to delete leasing request"
// @Router /request/{id} [delete]
func (h *LeasingRequestHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	user := c.Locals("user").(*domain.User)
	if user.Role != domain.AdminRole {
		return apperror.UnauthorizedError(errors.New("user is not an admin"), "only admin can delete a request")
	}

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	leasingRequestID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	err = h.service.Delete(leasingRequestID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Create godoc
// @Summary Create a new leasing request
// @Description Add a new leasing request to the database
// @Tags request
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Param user body dto.LeasingRequestCreateRequestBody true "request information"
// @Success 201 {object} dto.SuccessResponse[dto.LeasingRequest] "Dorm successfully created"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found or leasing request not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing request to database"
// @Router /request/{id} [post]
func (h *LeasingRequestHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == domain.LessorRole {
		return apperror.UnauthorizedError(errors.New("user is a lessor"), "lessor cannot create a request")
	}
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	body := new(dto.LeasingRequestCreateRequestBody)
	if err := c.BodyParser(&body); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	leasingRequest, err := h.service.Create(userID, dormID, body.Message)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}

	res := dto.Success(leasingRequest.ToDTO())

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *LeasingRequestHandler) GetByDormID(c *fiber.Ctx) error {
	dormID, err := parseIdParam(c)
	if err != nil {
		return err
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

	leasingRequest, totalPage, totalRows, err := h.service.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	resData := make([]dto.LeasingRequest, len(leasingRequest))
	for i, v := range leasingRequest {
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
