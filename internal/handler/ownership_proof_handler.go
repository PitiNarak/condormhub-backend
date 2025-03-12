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
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofWithFileResponseBody] "Ownership proof created"
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
	if err = o.storage.UploadFile(c.Context(), fileKey, contentType, fileData, storage.PrivateBucket); err != nil {
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

	if err := o.ownershipProofService.Create(ownershipProof); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create ownership proof error")
	}

	ownershipProofResponse := o.ownershipProofService.ConvertToDTOWithFile(*ownershipProof, url, time.Now().Add(time.Minute*60))

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponse))
}

// Delete godoc
// @Summary Delete ownership proof
// @Description Delete an ownership proof file
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 204 "Ownership proof successfully deleted"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ownership [delete]
func (o *OwnershipProofHandler) Delete(c *fiber.Ctx) error {

	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	ownershipProof, err := o.ownershipProofService.GetByDormID(dormID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	if err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error deleting file")
	}

	if err = o.ownershipProofService.Delete(dormID); err != nil {
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
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofWithFileResponseBody] "Ownership proof updated"
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
	if err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error deleting file")
	}
	//prepare element for upload to storage
	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	file_uuid := uuid.New().String()
	newFileKey := fmt.Sprintf("ownership-proof/%s-%s", filename, file_uuid)

	//upload to storage
	if err = o.storage.UploadFile(c.Context(), newFileKey, contentType, fileData, storage.PrivateBucket); err != nil {
		return apperror.InternalServerError(err, "error uploading file")
	}

	//get key file
	url, err := o.storage.GetSignedUrl(c.Context(), newFileKey, time.Minute*60)
	if err != nil {
		return apperror.InternalServerError(err, "error getting signed url")
	}

	if err := o.ownershipProofService.UpdateDocument(dormID, newFileKey); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
	}
	ownershipProofResponse := o.ownershipProofService.ConvertToDTOWithFile(*ownershipProof, url, time.Now().Add(time.Minute*60))

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponse))

}

// Approve godoc
// @Summary Approve ownership proof
// @Description Approve a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 200 {object} dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof approved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/approve [post]
func (o *OwnershipProofHandler) Approve(c *fiber.Ctx) error {
	adminID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	dormID, dormIdRrr := uuid.Parse(id)
	if dormIdRrr != nil {
		return apperror.InternalServerError(dormIdRrr, "Can not parse UUID")
	}

	if err := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Approved")); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
	}
	ownershipProofResponse := o.ownershipProofService.ConvertToDTO(*ownershipProof)

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponse))

}

// Reject godoc
// @Summary Reject ownership proof
// @Description Reject a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Accept json
// @Produce json
// @Param ownership body dto.DormIDForOwnershipProofRequestBody true "Dorm ID"
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof rejected"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/reject [post]
func (o *OwnershipProofHandler) Reject(c *fiber.Ctx) error {
	adminID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	dormID, dormIdRrr := uuid.Parse(id)
	if dormIdRrr != nil {
		return apperror.InternalServerError(dormIdRrr, "Can not parse UUID")
	}

	if err := o.ownershipProofService.UpdateStatus(dormID, adminID, domain.OwnershipProofStatus("Rejected")); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error updating file")
	}

	ownershipProof, getErr := o.ownershipProofService.GetByDormID(dormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "error getting new file")
	}
	ownershipProofResponse := o.ownershipProofService.ConvertToDTO(*ownershipProof)

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponse))
}

// GetByDormID godoc
// @Summary Get ownership proof by Dorm ID
// @Description Retrieve ownership proof for a specific dorm
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "Dorm ID (UUID format)"
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofWithFileResponseBody] "Ownership proof retrieved successfully"
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

	ownershipProofResponse := o.ownershipProofService.ConvertToDTOWithFile(*ownershipProof, url, time.Now().Add(time.Minute*60))

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponse))

}
