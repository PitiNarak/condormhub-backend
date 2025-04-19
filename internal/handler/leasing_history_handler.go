package handler

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingHistoryHandler struct {
	service     ports.LeasingHistoryService
	dormService ports.DormService
}

func NewLeasingHistoryHandler(service ports.LeasingHistoryService, dormService ports.DormService) ports.LeasingHistoryHandler {
	return &LeasingHistoryHandler{service: service, dormService: dormService}
}

// SetEndTimestamp godoc
// @Summary Set end date of a leasing history
// @Description Set end date of a leasing history in the database
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
	leasingHistoryID, err := parseIdParam(c)
	if err != nil {
		return err
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
		urls := h.service.GetImageUrl(v.Images)
		resData[i] = v.ToDTO(urls)
		resData[i].Images = h.dormService.GetImageUrl(v.Dorm.Images)
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

	leasingHistory, totalPage, totalRows, err := h.service.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	resData := make([]dto.LeasingHistory, len(leasingHistory))
	for i, v := range leasingHistory {
		urls := h.service.GetImageUrl(v.Images)
		resData[i] = v.ToDTO(urls)
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID godoc
// @Summary Get a leasing history by historyid
// @Description Retrieve a leasing history by historyid
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "HistoryID"
// @Success 200 {object} dto.LeasingHistory "Retrive history successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID"
// @Router /history/{id} [get]
func (h *LeasingHistoryHandler) GetByID(c *fiber.Ctx) error {
	historyID, err := parseIdParam(c)
	if err != nil {
		return err
	}
	leasingHistory, err := h.service.GetByID(historyID)
	if err != nil {
		return err
	}

	urls := h.service.GetImageUrl(leasingHistory.Images)
	resData := leasingHistory.ToDTO(urls)
	resData.Images = h.dormService.GetImageUrl(leasingHistory.Dorm.Images)

	return c.Status(fiber.StatusOK).JSON(resData)
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
	leasingHistoryID, err := parseIdParam(c)
	if err != nil {
		return err
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
	dormID, err := parseIdParam(c)
	if err != nil {
		return err
	}

	leasingHistory, err := h.service.Create(userID, dormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return err
	}
	urls := h.service.GetImageUrl(leasingHistory.Images)

	res := dto.Success(leasingHistory.ToDTO(urls))

	return c.Status(fiber.StatusCreated).JSON(res)
}

// CreateReview godoc
// @Summary Create a new review
// @Description Add a new review to the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "HistoryID"
// @Param user body dto.ReviewCreateRequestBody true "review information"
// @Success 201 {object} dto.SuccessResponse[dto.Review]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/{id}/review [post]
func (h *LeasingHistoryHandler) CreateReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	body := new(dto.ReviewCreateRequestBody)
	err := c.BodyParser(&body)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	historyID, err := parseIdParam(c)

	if err != nil {
		return err
	}
	history, err := h.service.GetByID(historyID)
	if err != nil {
		return err
	}
	review, err := h.service.CreateReview(user, historyID, body.Message, int(body.Rate))
	if err != nil {
		return err
	}
	urls := h.service.GetImageUrl(history.Images)
	res := dto.Success(review.ToDTO(urls))
	return c.Status(fiber.StatusOK).JSON(res)
}

// UploadReviewImage godoc
// @Summary Upload multiple images for a review
// @Description Upload multiple images for a specific review by its historyID, by attaching the images as value for the key field name "image", as a multipart form-data
// @Tags history
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "historyID"
// @Param image formData file true "ReviewImage" collectionFormat "multi"
// @Success 200 {object} dto.SuccessResponse[dto.ReviewImageUploadResponseBody] "Successful image upload"
// @Failure 400 {object} dto.ErrorResponse "Invalid Request"
// @Failure 403 {object} dto.ErrorResponse "unauthorized to upload image to review"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "History not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to upload review image"
// @Router /history/{id}/review/images [post]
func (h *LeasingHistoryHandler) UploadReviewImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	historyID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return apperror.BadRequestError(err, "Invalid multipart form data")
	}

	files := form.File["image"]

	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}
	isAdmin := user.Role == domain.AdminRole

	urls := []string{}
	for _, file := range files {
		fileData, err := file.Open()
		if err != nil {
			return apperror.InternalServerError(err, "error opening file")
		}
		defer fileData.Close()

		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return apperror.BadRequestError(errors.New("uploaded file is not an image"), "uploaded file is not an image")
		}

		url, err := h.service.UploadReviewImage(c.Context(), historyID, file.Filename, contentType, fileData, userID, isAdmin)
		if err != nil {
			return err
		}

		urls = append(urls, url)
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(dto.ReviewImageUploadResponseBody{ImageURL: urls}))
}

// DeleteReviewImageByURL godoc
// @Summary Delete a review image by its url
// @Description Deletes a review image using its percent encoded url from bucket storage. Encode URL using encodeURIComponent().
// @Tags history
// @Security Bearer
// @Accept json
// @Produce json
// @Param url path string true "Percent encoded URL"
// @Success 204 "Image deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 403 {object} dto.ErrorResponse "You do not have permission to delete this review image"
// @Failure 404 {object} dto.ErrorResponse "Image not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete image"
// @Router /history/review/image/{url} [delete]
func (h *LeasingHistoryHandler) DeleteReviewImageByURL(c *fiber.Ctx) error {
	decodedURL, err := url.PathUnescape(c.Params("url"))
	if err != nil {
		return apperror.BadRequestError(err, "Invalid URL")
	}

	userID := c.Locals("userID").(uuid.UUID)
	user := c.Locals("user").(*domain.User)
	if user.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}

	isAdmin := user.Role == domain.AdminRole

	if err := h.service.DeleteImageByURL(c.Context(), decodedURL, userID, isAdmin); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateReview godoc
// @Summary Update a new review
// @Description Update a review in the database
// @Tags history
// @Security Bearer
// @Produce json
// @Param id path string true "HistoryID"
// @Param user body dto.ReviewUpdateRequestBody true "review information"
// @Success 201 {object} dto.SuccessResponse[dto.Review]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format or your request is invalid or review not exist"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "Can not parse UUID or failed to save leasing history to database"
// @Router /history/{id}/review [patch]
func (h *LeasingHistoryHandler) UpdateReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	body := new(dto.ReviewUpdateRequestBody)
	err := c.BodyParser(&body)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	historyID, err := parseIdParam(c)
	if err != nil {
		return err
	}
	history, err := h.service.GetByID(historyID)
	if err != nil {
		return err
	}
	review, err := h.service.UpdateReview(user, historyID, body.Message, int(body.Rate))
	if err != nil {
		return err
	}
	urls := h.service.GetImageUrl(history.Images)
	res := dto.Success(review.ToDTO(urls))
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
// @Router /history/{id}/review [delete]
func (h *LeasingHistoryHandler) DeleteReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	historyID, err := parseIdParam(c)
	if err != nil {
		return err
	}
	err = h.service.DeleteReview(user, historyID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LeasingHistoryHandler) GetReportedReviews(c *fiber.Ctx) error {
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

	leasingHistory, totalPage, totalRows, err := h.service.GetReportedReviews(limit, page)
	if err != nil {
		return err
	}

	resData := make([]dto.LeasingHistory, len(leasingHistory))
	for i, v := range leasingHistory {
		urls := h.service.GetImageUrl(v.Images)
		resData[i] = v.ToDTO(urls)
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

func parseIdParam(c *fiber.Ctx) (uuid.UUID, error) {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return uuid.Nil, apperror.BadRequestError(err, "Incorrect UUID format")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return uuid.Nil, err
		}
		return uuid.Nil, apperror.InternalServerError(err, "Can not parse UUID")
	}
	return parsedID, nil
}
