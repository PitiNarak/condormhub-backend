package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"

	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/google/uuid"
)

type OwnershipProofHandler struct {
	storage               *storage.Storage
	ownershipProofService ports.OwnershipProofService
}

func NewOwnershipProofHandler(OwnershipProofService ports.OwnershipProofService, storage *storage.Storage) ports.OwnershipProofHandler {
	return &OwnershipProofHandler{ownershipProofService: OwnershipProofService, storage: storage}
}

// Create godoc
// @Summary Upload new ownership proof
// @Description Upload a new file as ownership proof for a dorm
// @Tags ownership
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Ownership proof file"
// @Param dormId formData string true "Dorm ID (UUID format)"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ownership/create [post]
func (o *OwnershipProofHandler) Create(c *fiber.Ctx) error {

	//extract file from http
	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	dormIDStr := c.FormValue("dormId")
	dormID, dorm_err := uuid.Parse(dormIDStr)

	if dorm_err != nil {
		return apperror.BadRequestError(err, "Dorm id is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	//prepare element for upload to storage
	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	file_uuid := uuid.New().String()
	fileKey := fmt.Sprintf("ownership-proof/%s-%s", filename, file_uuid)

	//upload to storage
	err = o.storage.UploadFile(c.Context(), fileKey, contentType, fileData, storage.PrivateBucket)
	if err != nil {
		return apperror.InternalServerError(err, "error uploading file")
	}

	//get key file
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	ownershipProof := &domain.OwnershipProof{
		DormID:  dormID,
		FileKey: fileKey,
	}

	db_err := o.ownershipProofService.Create(ownershipProof)

	if db_err != nil {
		return db_err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)})
}

// Delete godoc
// @Summary Delete ownership proof
// @Description Delete an ownership proof file
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ownership [delete]
func (o *OwnershipProofHandler) Delete(c *fiber.Ctx) error {

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if body_err := c.BodyParser(dormReqBody); body_err != nil {
		return apperror.BadRequestError(body_err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	ownershipProof, err := o.ownershipProofService.GetByDormID(dormID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	storage_err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if storage_err != nil {
		return apperror.InternalServerError(err, "error deleting file")
	}

	err = o.ownershipProofService.Delete(dormID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success("Ownership proof successfully deleted"))
}

// Update godoc
// @Summary Update ownership proof
// @Description Replace an existing ownership proof file
// @Tags ownership
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "New ownership proof file"
// @Param dormId formData string true "Dorm ID (UUID format)"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ownership/update [post]
func (o *OwnershipProofHandler) Update(c *fiber.Ctx) error {

	//extract file from http
	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	dormIDStr := c.FormValue("dormId")
	dormID, dorm_err := uuid.Parse(dormIDStr)

	if dorm_err != nil {
		return apperror.BadRequestError(err, "Dorm id is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	ownershipProof, err := o.ownershipProofService.GetByDormID(dormID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	storage_err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if storage_err != nil {
		return apperror.InternalServerError(err, "error deleting file")
	}

	//prepare element for upload to storage
	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	file_uuid := uuid.New().String()
	newFileKey := fmt.Sprintf("ownership-proof/%s-%s", filename, file_uuid)

	//upload to storage
	err = o.storage.UploadFile(c.Context(), newFileKey, contentType, fileData, storage.PrivateBucket)
	if err != nil {
		return apperror.InternalServerError(err, "error uploading file")
	}

	//get key file
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	requestBody := new(dto.UpdateOwnerShipProofRequestBody)
	requestBody.FileKey = newFileKey
	db_err := o.ownershipProofService.UpdateDocument(dormID, requestBody)
	if db_err != nil {
		return db_err
	}

	ownershipProof, get_err := o.ownershipProofService.GetByDormID(dormID)
	if get_err != nil {
		return get_err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)})

}

// Approve godoc
// @Summary Approve ownership proof
// @Description Approve a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 200 {object} dto.SuccessResponse[string] "Ownership proof approved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/approve [post]
func (o *OwnershipProofHandler) Approve(c *fiber.Ctx) error {
	adminID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if body_err := c.BodyParser(dormReqBody); body_err != nil {
		return apperror.BadRequestError(body_err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	updateStatus_err := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Approved"))
	if updateStatus_err != nil {
		return updateStatus_err
	}

	ownershipProof, getOnwership_err := o.ownershipProofService.GetByDormID(dormID)
	if getOnwership_err != nil {
		return getOnwership_err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Dorm's ownership proof": ownershipProof})
}

// Reject godoc
// @Summary Reject ownership proof
// @Description Reject a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 200 {object} dto.SuccessResponse[string] "Ownership proof rejected"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/reject [post]
func (o *OwnershipProofHandler) Reject(c *fiber.Ctx) error {
	adminID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if body_err := c.BodyParser(dormReqBody); body_err != nil {
		return apperror.BadRequestError(body_err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	updateStatus_err := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Rejected"))
	if updateStatus_err != nil {
		return updateStatus_err
	}

	ownershipProof, getOnwership_err := o.ownershipProofService.GetByDormID(dormID)
	if getOnwership_err != nil {
		return getOnwership_err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Dorm's ownership proof": ownershipProof})
}

// GetByDormID godoc
// @Summary Get ownership proof by Dorm ID
// @Description Retrieve ownership proof for a specific dorm
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "Dorm ID (UUID format)"
// @Success 200 {object} dto.SuccessResponse[string] "Ownership proof retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse[] "Invalid UUID format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized request"
// @Failure 404 {object} dto.ErrorResponse "Ownership proof not found"
// @Router /ownership/{id} [get]
func (o *OwnershipProofHandler) GetByDormID(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	ownershipProof, getOnwership_err := o.ownershipProofService.GetByDormID(dormID)
	if getOnwership_err != nil {
		return getOnwership_err
	}

	fileKey := ownershipProof.FileKey
	//get key file
	url, url_err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if url_err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)})

}
