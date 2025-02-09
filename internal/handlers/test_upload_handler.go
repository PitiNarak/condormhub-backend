package handlers

import (
	"fmt"
	"strings"

	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
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

func (e *TestUploadHandler) UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return error_handler.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return error_handler.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	uuid := uuid.New().String()

	url, err := e.storage.UploadFile(c.Context(), fmt.Sprintf("test/%s-%s", uuid, filename), contentType, fileData)
	if err != nil {
		return error_handler.InternalServerError(err, "error uploading file")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("upload success", fiber.Map{"url": url}))
	// return error_handler.InternalServerError(errors.New("error from system"), "your error message")
}
