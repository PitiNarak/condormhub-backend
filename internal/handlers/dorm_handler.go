package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DormHandler struct {
	dormService ports.DormService
}

func NewDormHandler(service ports.DormService) ports.DormHandler {
	return &DormHandler{dormService: service}
}

func (d *DormHandler) Create(c *fiber.Ctx) error {
	reqBody := new(dto.DormRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return error_handler.BadRequestError(err, "Your request is invalid")
	}

	userIDstr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return error_handler.InternalServerError(err, "Can not parse UUID")
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return error_handler.BadRequestError(err, "Your request body is invalid")
	}

	dorm := &domain.Dorm{
		Name:        reqBody.Name,
		OwnerID:     userID,
		Size:        reqBody.Size,
		Bedrooms:    reqBody.Bedrooms,
		Bathrooms:   reqBody.Bathrooms,
		Province:    reqBody.Province,
		District:    reqBody.District,
		Price:       reqBody.Price,
		Description: reqBody.Description,
	}

	if err := d.dormService.Create(dorm); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("Dorm successfully created", dorm))
}

func (d *DormHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return error_handler.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return error_handler.InternalServerError(err, "Can not parse UUID")
	}

	_, err = d.dormService.GetByID(dormID)
	if err != nil {
		return err
	}

	err = d.dormService.Delete(dormID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Dorm successfully deleted", nil))
}

func (d *DormHandler) GetAll(c *fiber.Ctx) error {
	dorms, err := d.dormService.GetAll()
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("All dorms retrieved successfully", dorms))
}

func (d *DormHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return error_handler.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return error_handler.InternalServerError(err, "Can not parse UUID")
	}

	dorm, err := d.dormService.GetByID(dormID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Dorm data successfully retrieved", dorm))
}

func (d *DormHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	reqBody := new(dto.DormRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return error_handler.BadRequestError(err, "Your request is invalid")
	}

	userIDstr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return error_handler.InternalServerError(err, "cannot parse uuid")
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return error_handler.BadRequestError(err, "Your request body is invalid")
	}

	if err := uuid.Validate(id); err != nil {
		return error_handler.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return error_handler.InternalServerError(err, "Can not parse UUID")
	}

	dorm := &domain.Dorm{
		Name:        reqBody.Name,
		OwnerID:     userID,
		Size:        reqBody.Size,
		Bedrooms:    reqBody.Bedrooms,
		Bathrooms:   reqBody.Bathrooms,
		Province:    reqBody.Province,
		District:    reqBody.District,
		Price:       reqBody.Price,
		Description: reqBody.Description,
	}

	err = d.dormService.Update(dormID, dorm)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Dorm data updated successfully", dorm))
}
