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
