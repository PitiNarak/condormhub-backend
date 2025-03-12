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
	dormID, err := uuid.Parse(dormIDStr)

	if err != nil {
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
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*60)
	if err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	ownershipProof := &domain.OwnershipProof{
		DormID:  dormID,
		FileKey: fileKey,
	}

	createErr := o.ownershipProofService.Create(ownershipProof)

	if createErr != nil {
		if apperror.IsAppError(createErr) {
			return createErr
		}
		return apperror.InternalServerError(createErr, "create ownership proof error")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 60)})
}

// Delete godoc
// @Summary Delete ownership proof
// @Description Delete an ownership proof file
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 204 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ownership [delete]
func (o *OwnershipProofHandler) Delete(c *fiber.Ctx) error {

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if err := c.BodyParser(dormReqBody); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	ownershipProof, err := o.ownershipProofService.GetByDormID(dormID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	deleteErr := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if deleteErr != nil {
		if apperror.IsAppError(deleteErr) {
			return deleteErr
		}
		return apperror.InternalServerError(deleteErr, "error deleting file")
	}

	err = o.ownershipProofService.Delete(dormID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)

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
	dormID, err := uuid.Parse(dormIDStr)

	if err != nil {
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
	deleteErr := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if deleteErr != nil {
		if apperror.IsAppError(deleteErr) {
			return deleteErr
		}
		return apperror.InternalServerError(deleteErr, "error deleting file")
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
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*60)
	if err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	updateErr := o.ownershipProofService.UpdateDocument(dormID, newFileKey)
	if updateErr != nil {
		if apperror.IsAppError(updateErr) {
			return updateErr
		}
		return apperror.InternalServerError(updateErr, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 60)})

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
	adminID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if err := c.BodyParser(dormReqBody); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	updateErr := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Approved"))
	if updateErr != nil {
		if apperror.IsAppError(updateErr) {
			return updateErr
		}
		return apperror.InternalServerError(updateErr, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
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
	adminID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	dormReqBody := new(dto.DormIDForOwnershipProofRequestBody)
	if err := c.BodyParser(dormReqBody); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	dormID := dormReqBody.DormID

	updateErr := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Rejected"))
	if updateErr != nil {
		if apperror.IsAppError(updateErr) {
			return updateErr
		}
		return apperror.InternalServerError(updateErr, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
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
// @Failure 400 {object} dto.ErrorResponse "Invalid UUID format"
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

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
	}

	fileKey := ownershipProof.FileKey
	//get key file
	url, urlErr := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*60)
	if urlErr != nil {
		return apperror.InternalServerError(urlErr, "error getting signed url")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 60)})

}
