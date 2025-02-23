package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
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

func (h *LeasingHistoryHandler) Create(c *fiber.Ctx) error {
	reqBody := new(dto.LeasingHistoryCreateRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Your request body is invalid")
	}
	userID := c.Locals("userID").(uuid.UUID)
	dormIDstr := reqBody.DormID
	dormID, err := uuid.Parse(dormIDstr)
	if err != nil {
		return errorHandler.BadRequestError(err, "Can not parse dormID")
	}
	leasingHistory, err := h.service.Create(userID, dormID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Leasing history successfully deleted", leasingHistory))
}
func (h *LeasingHistoryHandler) Update(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) Delete(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) GetByUserID(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) PatchEndTimestamp(c *fiber.Ctx) error {
	return nil
}
