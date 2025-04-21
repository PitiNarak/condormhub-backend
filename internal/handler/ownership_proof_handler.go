package handler

import (
	"errors"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/go-pkg/apperror"

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
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Param id path string true "DormID"
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof created"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Ownershop proof not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to upload file"
// @Router /ownership/{id}/upload [post]
func (o *OwnershipProofHandler) UploadFile(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	contentType := file.Header.Get("Content-Type")
	url, err := o.ownershipProofService.UploadFile(c.Context(), dormID, file.Filename, contentType, fileData)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
	}

	ownershipProof, err := o.ownershipProofService.GetByDormID(dormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error getting ownership proof")
	}
	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		Url:     url,
		DormID:  dormID,
		AdminID: ownershipProof.AdminID,
		Status:  dto.OwnershipProofStatus(ownershipProof.Status),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponseBody))
}

// Delete godoc
// @Summary Delete ownership proof
// @Description Delete an ownership proof file
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 204 "Ownership proof successfully deleted"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Ownership file not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete Ownership file"
// @Router /ownership/{id} [delete]
func (o *OwnershipProofHandler) Delete(c *fiber.Ctx) error {

	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}
	dormID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	if err = o.ownershipProofService.Delete(c.Context(), dormID); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "error deleting file")
	}
	return c.SendStatus(fiber.StatusNoContent)

}

// Approve godoc
// @Summary Approve ownership proof
// @Description Approve a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object} dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof approved"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Ownership file not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/{id}/approve [post]
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

	url, urlErr := o.ownershipProofService.GetUrl(c.Context(), dormID)

	if urlErr != nil {
		if apperror.IsAppError(urlErr) {
			return urlErr
		}
		return apperror.InternalServerError(urlErr, "error getting url")
	}
	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		Url:     url,
		DormID:  dormID,
		AdminID: ownershipProof.AdminID,
		Status:  dto.OwnershipProofStatus(ownershipProof.Status),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponseBody))

}

// Reject godoc
// @Summary Reject ownership proof
// @Description Reject a submitted ownership proof for a dorm
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof rejected"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Ownership file not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /ownership/{id}/reject [post]
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
	url, urlErr := o.ownershipProofService.GetUrl(c.Context(), dormID)
	if urlErr != nil {
		if apperror.IsAppError(urlErr) {
			return urlErr
		}
		return apperror.InternalServerError(urlErr, "error getting url")
	}
	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		Url:     url,
		DormID:  dormID,
		AdminID: ownershipProof.AdminID,
		Status:  dto.OwnershipProofStatus(ownershipProof.Status),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponseBody))
}

// GetByDormID godoc
// @Summary Get ownership proof by Dorm ID
// @Description Retrieve ownership proof for a specific dorm
// @Tags ownership
// @Security Bearer
// @Produce json
// @Param id path string true "DormID"
// @Success 200 {object}  dto.SuccessResponse[dto.OwnershipProofResponseBody] "Ownership proof retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Ownership file not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
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
		return apperror.InternalServerError(getErr, "error getting ownership proof")
	}

	fileKey := ownershipProof.FileKey
	//get key file
	url, urlErr := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*60)
	if urlErr != nil {
		return apperror.InternalServerError(urlErr, "error getting url")
	}

	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		Url:     url,
		DormID:  dormID,
		AdminID: ownershipProof.AdminID,
		Status:  dto.OwnershipProofStatus(ownershipProof.Status),
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(ownershipProofResponseBody))

}
