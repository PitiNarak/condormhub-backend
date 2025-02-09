package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SampleLogHandler struct {
	sampleLogService ports.SampleLogService
}

func NewSampleLogHandler(sampleLogService ports.SampleLogService) *SampleLogHandler {
	return &SampleLogHandler{
		sampleLogService: sampleLogService,
	}
}

func (s *SampleLogHandler) Save(c *fiber.Ctx) error {
	payload := new(domain.SampleLog)

	if err := c.BodyParser(payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := s.sampleLogService.Save(payload.Message); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Saved"})
}

func (s *SampleLogHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := uuid.Validate(id); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	uuid, _ := uuid.Parse(id)
	if err := s.sampleLogService.Delete(uuid); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}

func (s *SampleLogHandler) GetAll(c *fiber.Ctx) error {
	messages, err := s.sampleLogService.GetAll()
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(messages)
}

func (s *SampleLogHandler) EditMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	payload := new(domain.SampleLog)

	if err := c.BodyParser(payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := uuid.Validate(id); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	uuid, _ := uuid.Parse(id)
	if err := s.sampleLogService.EditMessage(uuid, payload.Message); err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Edited"})
}
