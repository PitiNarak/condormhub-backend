package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TestUploadHandler struct {
	storage *storage.Storage
}

func NewTestUploadHandler(storage *storage.Storage) *TestUploadHandler {
	return &TestUploadHandler{
		storage: storage,
	}
}

func (e *TestUploadHandler) UploadToPrivateBucketHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errorHandler.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return errorHandler.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("test/%s-%s", uuid, filename)

	err = e.storage.UploadFile(c.Context(), fileKey, contentType, fileData, storage.PrivateBucket)
	if err != nil {
		return errorHandler.InternalServerError(err, "error uploading file")
	}

	url, err := e.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return errorHandler.InternalServerError(err, "error getting signed url")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("upload success", http_response.SuccessResponse("upload success", fiber.Map{"url": url, "key": fileKey, "expires": time.Now().Add(time.Minute * 5)})))
}

func (e *TestUploadHandler) UploadToPublicBucketHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errorHandler.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return errorHandler.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("test/%s-%s", uuid, filename)

	err = e.storage.UploadFile(c.Context(), fileKey, contentType, fileData, storage.PublicBucket)
	if err != nil {
		return errorHandler.InternalServerError(err, "error uploading file")
	}

	url := e.storage.GetPublicUrl(fileKey)

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("upload success", fiber.Map{"url": url}))
}

func (e *TestUploadHandler) GetSignedUrlHandler(c *fiber.Ctx) error {
	fileKey := c.Params("*")
	url, err := e.storage.GetSignedUrl(c.Context(), fileKey, time.Minute*5)
	if err != nil {
		return errorHandler.InternalServerError(err, "error getting signed url")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("get signed url success", fiber.Map{"url": url, "key": fileKey, "expires": time.Now().Add(time.Minute * 5)}))
}
