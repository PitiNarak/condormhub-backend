package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
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
// @Success 201  {object}  httpResponse.HttpResponse{data=domain.LeasingHistory,pagination=nil} "Dorm successfully created"
// @Failure 400  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "Incorrect UUID format"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "Dorm not found or leasing history not found"
// @Failure 500  {object}  httpResponse.HttpResponse{data=nil,pagination=nil} "Can not parse UUID or failed to save leasing history to database"
// @Router /history/{id} [post]
func (h *LeasingHistoryHandler) Create(c *fiber.Ctx) error {
	reqBody := new(dto.LeasingHistoryCreateRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Your request is invalid")
	}
	if err := c.BodyParser(reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Your request is invalid")
	}
	userID := c.Locals("userID").(uuid.UUID)
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return errorHandler.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return errorHandler.InternalServerError(err, "Can not parse UUID")
	}
	leasingHistory, err := h.service.Create(userID, dormID, reqBody.Price)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(httpResponse.SuccessResponse("Leasing history successfully created", leasingHistory))
}
