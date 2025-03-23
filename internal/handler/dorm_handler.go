package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
// @Param dorm body dto.DormCreateRequestBody true "Dorm information"
// @Success 201 {object} dto.SuccessResponse[dto.DormResponseBody] "Dorm successfully created"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 403 {object} dto.ErrorResponse "You do not have permission to create a dorm"
// @Failure 400 {object} dto.ErrorResponse "Your request is invalid"
// @Failure 500 {object} dto.ErrorResponse "Failed to save dorm"
// @Router /dorms [post]
func (d *DormHandler) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	userRole := user.Role
	if userRole == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}

	reqBody := new(dto.DormCreateRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	userID := c.Locals("userID").(uuid.UUID)

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request body is invalid")
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

	if err := d.dormService.Create(userRole, dorm); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create dorm error")
	}

	data, err := d.dormService.GetByID(dorm.ID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get dorm error")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(data))
}

// Delete godoc
// @Summary Delete a dorm
// @Description Removes a dorm from the database based on the give ID
// @Tags dorms
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 204 "Dorm successfully deleted"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 403 {object} dto.ErrorResponse "You do not have permission to delete this dorm"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete dorm"
// @Router /dorms/{id} [delete]
func (d *DormHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}

	isAdmin := user.Role == domain.AdminRole

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	if err := d.dormService.Delete(c.Context(), userID, isAdmin, dormID); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "delete dorm error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetAll godoc
// @Summary Get all dorms by a search string
// @Description Retrieve a list of all dorms filtered by a search query. If no query is provided, all dorms are returned.
// @Tags dorms
// @Param search query string false "Search query"
// @Param minPrice query int false "Filter min price"
// @Param maxPrice query int false "Filter max price"
// @Param district query string false "Filter district price"
// @Param subdistrict query string false "Filter subdistrict price"
// @Param province query string false "Filter province price"
// @Param zipcode query string false "Filter zipcode price"
// @Param limit query int false "Number of dorms to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.DormResponseBody] "All dorms retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve dorms"
// @Router /dorms [get]
func (d *DormHandler) GetAll(c *fiber.Ctx) error {
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

	search := c.Query("search")
	minPrice := c.QueryInt("minPrice", -1)
	maxPrice := c.QueryInt("maxPrice", -1)

	if minPrice > maxPrice {
		err := errors.New("min price cannot more than max price")
		return apperror.BadRequestError(err, err.Error())
	}

	district := c.Query("district")
	subdistrict := c.Query("subdistrict")
	province := c.Query("province")
	zipcode := c.Query("zipcode")

	dorms, totalPages, totalRows, err := d.dormService.GetAll(limit, page, search, minPrice, maxPrice, district, subdistrict, province, zipcode)
	if err != nil {
		return err
	}

	res := dto.SuccessPagination(dorms, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPages,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID godoc
// @Summary Get a dorm by ID
// @Description Retrieve a specific dorm based on its ID
// @Tags dorms
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object} dto.SuccessResponse[dto.DormResponseBody] "Dorm data successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to retrieve dorm"
// @Router /dorms/{id} [get]
func (d *DormHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	dorm, err := d.dormService.GetByID(dormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get dorm error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(dorm))
}

// Update godoc
// @Summary Update an existing dorm
// @Description Modifies an existing room's details based on the given ID
// @Tags dorms
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "DormID"
// @Param dorm body dto.DormUpdateRequestBody true "Updated Room Data"
// @Success 200 {object} dto.SuccessResponse[dto.DormResponseBody] "Dorm data updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid Request"
// @Failure 403 {object} dto.ErrorResponse "unauthorized to update this dorm"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to update dorm"
// @Router /dorms/{id} [patch]
func (d *DormHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)

	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}

	isAdmin := user.Role == domain.AdminRole

	id := c.Params("id")
	updateReqBody := new(dto.DormUpdateRequestBody)
	if err := c.BodyParser(updateReqBody); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	validate := validator.New()
	if err := validate.Struct(updateReqBody); err != nil {
		return apperror.BadRequestError(err, "Your request body is invalid")
	}

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	updatedDorm, err := d.dormService.Update(userID, isAdmin, dormID, updateReqBody)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "update dorm error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(updatedDorm))
}

// UploadDormImage godoc
// @Summary Upload an image for a dorm
// @Description Upload an image for a specific dorm by its ID, by attaching the image as a value for the key field name "image", as a multipart form-data
// @Tags dorms
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "DormID"
// @Param image formData file true "DormImage"
// @Success 200 {object} dto.SuccessResponse[dto.DormImageUploadResponseBody] "Successful image upload"
// @Failure 400 {object} dto.ErrorResponse "Invalid Request"
// @Failure 403 {object} dto.ErrorResponse "unauthorized to upload image to dorm"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Dorm not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to upload dorm image"
// @Router /dorms/{id}/images [post]
func (d *DormHandler) UploadDormImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}
	isAdmin := user.Role == domain.AdminRole

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	contentType := file.Header.Get("Content-Type")
	url, err := d.dormService.UploadDormImage(c.Context(), dormID, file.Filename, contentType, fileData, userID, isAdmin)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(dto.DormImageUploadResponseBody{ImageURL: url}))
}

// GetByOwnerID godoc
// @Summary Get dorms by owner ID
// @Description Retrieve all dorms of a specific owner ID
// @Tags dorms
// @Param id path string true "OwnerID"
// @Param limit query int false "Number of dorms to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.DormResponseBody] "All dorms retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve dorms"
// @Router /dorms/owner/{id} [get]
func (d *DormHandler) GetByOwnerID(c *fiber.Ctx) error {
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

	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	ownerID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	dorms, totalPages, totalRows, err := d.dormService.GetByOwnerID(ownerID, limit, page)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get dorms error")
	}

	res := dto.SuccessPagination(dorms, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPages,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

func (d *DormHandler) DeleteDormImageByURL(c *fiber.Ctx) error {
	reqBody := new(dto.DormImageDeleteRequestBody)
	if err := c.BodyParser(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	if err := d.dormService.DeleteImageByURL(c.Context(), reqBody.ImageURL); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
