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

// Register godoc
// @Summary Create a new dorm
// @Description Add a new room to the database with the given details
// @Tags dorms
// @Security Bearer
// @Accept json
// @Produce json
// @Param dorm body dto.DormRequestBody true "Dorm information"
// @Success 201  {object}  http_response.HttpResponse{data=domain.Dorm} "Dorm successfully created"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 400  {object}  http_response.HttpResponse{data=nil} "Your request is invalid"
// @Failure 500  {object}  http_response.HttpResponse{data=nil} "Failed to save dorm"
// @Router /dorms [post]
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
		Name:      reqBody.Name,
		OwnerID:   userID,
		Size:      reqBody.Size,
		Bedrooms:  reqBody.Bedrooms,
		Bathrooms: reqBody.Bathrooms,
		Address: domain.Address{
			District:    reqBody.Address.District,
			Subdistrict: reqBody.Address.Subdistrict,
			Province:    reqBody.Address.Province,
			Zipcode:     reqBody.Address.Zipcode,
		},
		Price:       reqBody.Price,
		Description: reqBody.Description,
	}

	if err := d.dormService.Create(dorm); err != nil {
		return err
	}

	res, err := d.dormService.GetByID(dorm.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("Dorm successfully created", res))
}

// Delete godoc
// @Summary Delete a dorm
// @Description Removes a dorm from the database based on the give ID
// @Tags dorms
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object} http_response.HttpResponse{data=nil} "Dorm successfully deleted"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "Incorrect UUID format"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} http_response.HttpResponse{data=nil} "Dorm not found"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "Failed to delete dorm"
// @Router /dorms/{id} [delete]
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

// GetAll godoc
// @Summary Get all dorms
// @Description Retrieve a list of all dorms
// @Tags dorms
// @Produce json
// @Success 200 {object} http_response.HttpResponse{data=[]domain.Dorm} "All dorms retrieved successfully"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "Failed to retrieve dorms"
// @Router /dorms [get]
func (d *DormHandler) GetAll(c *fiber.Ctx) error {
	dorms, err := d.dormService.GetAll()
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("All dorms retrieved successfully", dorms))
}

// GetByID godoc
// @Summary Get a dorm by ID
// @Description Retrieve a specific dorm based on its ID
// @Tags dorms
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object} http_response.HttpResponse{data=domain.Dorm} "Dorm data successfully retrieved"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "Incorrect UUID format"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} http_response.HttpResponse{data=nil} "Dorm not found"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "Server failed to retrieve dorm"
// @Router /dorms/{id} [get]
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

// Update godoc
// @Summary Update an existing dorm
// @Description Modifies an existing room's details based on the given ID
// @Tags dorms
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "DormID"
// @Param dorm body dto.DormRequestBody true "Updated Room Data"
// @Success 200 {object} http_response.HttpResponse{data=domain.Dorm} "Dorm data updated successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "Invalid Request"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} http_response.HttpResponse{data=nil} "Dorm not found"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "Server failed to update dorm"
// @Router /dorms/{id} [patch]
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
		Name:      reqBody.Name,
		OwnerID:   userID,
		Size:      reqBody.Size,
		Bedrooms:  reqBody.Bedrooms,
		Bathrooms: reqBody.Bathrooms,
		Address: domain.Address{
			District:    reqBody.Address.District,
			Subdistrict: reqBody.Address.Subdistrict,
			Province:    reqBody.Address.Province,
			Zipcode:     reqBody.Address.Zipcode,
		},
		Price:       reqBody.Price,
		Description: reqBody.Description,
	}

	err = d.dormService.Update(dormID, dorm)
	if err != nil {
		return err
	}

	res, err := d.dormService.GetByID(dormID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Dorm data updated successfully", res))
}
