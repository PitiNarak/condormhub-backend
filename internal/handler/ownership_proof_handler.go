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

	return c.Status(fiber.StatusOK).JSON(dto.Success(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)}))
}

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

	return c.Status(fiber.StatusOK).JSON(dto.Success(fiber.Map{"url": url, "user's ownership proof": ownershipProof, "expires": time.Now().Add(time.Minute * 5)}))

}

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

	return c.Status(fiber.StatusOK).JSON(dto.Success(fiber.Map{"Dorm's ownership proof": ownershipProof}))
}
