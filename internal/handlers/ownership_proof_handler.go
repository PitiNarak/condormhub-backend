package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/internal/storage"

	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
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

func (o *OwnershipProofHandler) Create(c *fiber.Ctx) error {
	// userIDstr := c.Locals("userID").(string)
	// userID, err := uuid.Parse(userIDstr)
	// if err != nil {
	// 	return errorHandler.InternalServerError(err, "Can not parse UUID")
	// }
	userID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return errorHandler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	//extract file from http
	file, err := c.FormFile("file")
	if err != nil {
		return errorHandler.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return errorHandler.InternalServerError(err, "error opening file")
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
		return errorHandler.InternalServerError(err, "error uploading file")
	}

	//get key file
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return errorHandler.InternalServerError(err, "error getting signed url")
	}

	ownershipProof := &domain.OwnershipProof{
		LessorID: userID,
		FileKey:  fileKey,
	}

	db_err := o.ownershipProofService.Create(ownershipProof)

	if db_err != nil {
		return db_err
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("upload success", httpResponse.SuccessResponse("upload success", fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)})))
}

func (o *OwnershipProofHandler) Delete(c *fiber.Ctx) error {
	// LessorIDstr := c.Locals("userID").(string)
	// LessorID, err := uuid.Parse(LessorIDstr)
	// if err != nil {
	// 	return errorHandler.InternalServerError(err, "cannot parse uuid")
	// }
	LessorID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return errorHandler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	ownershipProof, err := o.ownershipProofService.GetByLessorID(LessorID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	storage_err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if storage_err != nil {
		return errorHandler.InternalServerError(err, "error deleting file")
	}

	err = o.ownershipProofService.Delete(LessorID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Ownership proof successfully deleted", nil))
}

func (o *OwnershipProofHandler) Update(c *fiber.Ctx) error {
	LessorID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return errorHandler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	//extract file from http
	file, err := c.FormFile("file")
	if err != nil {
		return errorHandler.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return errorHandler.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	ownershipProof, err := o.ownershipProofService.GetByLessorID(LessorID)
	if err != nil {
		return err
	}
	fileKey := ownershipProof.FileKey
	storage_err := o.storage.DeleteFile(c.Context(), fileKey, storage.PrivateBucket)

	if storage_err != nil {
		return errorHandler.InternalServerError(err, "error deleting file")
	}

	//prepare element for upload to storage
	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	file_uuid := uuid.New().String()
	newFileKey := fmt.Sprintf("ownership-proof/%s-%s", filename, file_uuid)

	//upload to storage
	err = o.storage.UploadFile(c.Context(), newFileKey, contentType, fileData, storage.PrivateBucket)
	if err != nil {
		return errorHandler.InternalServerError(err, "error uploading file")
	}

	//get key file
	url, err := o.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return errorHandler.InternalServerError(err, "error getting signed url")
	}

	requestBody := new(dto.UpdateOwnerShipProofRequestBody)
	requestBody.FileKey = newFileKey
	db_err := o.ownershipProofService.UpdateDocument(LessorID, requestBody)
	if db_err != nil {
		return db_err
	}

	ownershipProof, get_err := o.ownershipProofService.GetByLessorID(LessorID)
	if get_err != nil {
		return get_err
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("upload success", httpResponse.SuccessResponse("upload success", fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)})))

}

func (o *OwnershipProofHandler) Verify(c *fiber.Ctx) error {
	adminID, locals_err := c.Locals("userID").(uuid.UUID)
	if !locals_err {
		return errorHandler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := o.ownershipProofService.UpdateStatus()

	return nil
}
