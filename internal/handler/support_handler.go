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

	return c.Status(fiber.StatusCreated).JSON(dto.Success(support))
}
