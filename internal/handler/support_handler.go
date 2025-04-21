package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
)

type SupportHandler struct {
	service ports.SupportService
}

func NewSupportHandler(service ports.SupportService) ports.SupportHandler {
	return &SupportHandler{service: service}
}

// Create godoc
// @Summary Submit a support request
// @Description Let a user send a message to the admin
// @Tags support
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body dto.SupportRequestBody true "Support Request"
// @Success 201 {object} dto.SuccessResponse[dto.SupportResponseBody] "Support request submitted successfully"
// @Failure 400 {object} dto.ErrorResponse "Your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Could not submit support request"
// @Router /support [post]
func (h *SupportHandler) Create(c *fiber.Ctx) error {
	reqBody := new(dto.SupportRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request body is invalid")
	}

	userID := c.Locals("userID").(uuid.UUID)
	support := &domain.SupportRequest{
		UserID:  userID,
		Message: reqBody.Message,
	}
	if err := h.service.Create(support); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(support.ToDTO()))
}

// GetAll godoc
// @Summary Get all support requests
// @Description Retrieve a list of all support requests.
// @Tags support
// @Security Bearer
// @Param limit query int false "Number of support requests to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.SupportResponseBody] "All support requests retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Could not fetch support requests"
// @Router /support [get]
func (h *SupportHandler) GetAll(c *fiber.Ctx) error {
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

	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}
	isAdmin := user.Role == domain.AdminRole

	supports, totalPages, totalRows, err := h.service.GetAll(limit, page, userID, isAdmin)
	if err != nil {
		return err
	}

	resData := make([]dto.SupportResponseBody, len(supports))
	for i, support := range supports {
		resData[i] = support.ToDTO()
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPages,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateStatus godoc
// @Summary Update a support request's status
// @Description Modifies an existing support request's status based on the given ID. The possible status are OPEN, IN-PROGRESS, and RESOLVED.
// @Tags support
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "SupportRequestID"
// @Param dorm body dto.UpdateStatusRequestBody true "Status update"
// @Success 200 {object} dto.SuccessResponse[dto.SupportResponseBody] "Status updated"
// @Failure 400 {object} dto.ErrorResponse "Invalid Request"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 403 {object} dto.ErrorResponse "You do not have permission to update support request status"
// @Failure 404 {object} dto.ErrorResponse "Support request not found"
// @Failure 422 {object} dto.ErrorResponse "Invalid status value"
// @Failure 500 {object} dto.ErrorResponse "Server failed to update support request status"
// @Router /support/{id} [patch]
func (h *SupportHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	supportID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	req := new(dto.UpdateStatusRequestBody)
	if err = c.BodyParser(req); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		return apperror.BadRequestError(err, "Validation failed")
	}

	status := domain.SupportStatus(req.Status)
	updatedSupport, err := h.service.UpdateStatus(supportID, status)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(updatedSupport.ToDTO()))
}
