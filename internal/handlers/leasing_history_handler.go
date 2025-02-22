package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type LeasingHistoryHandler struct {
	service ports.LeasingHistoryService
}

func NewLeasingHistoryHandler(service ports.LeasingHistoryService) ports.LeasingHistoryHandler {
	return &LeasingHistoryHandler{service: service}
}

func (h *LeasingHistoryHandler) Create(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) Update(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) Delete(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) GetByUserID(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) GetByDormID(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) AddNewOrder(c *fiber.Ctx) error {
	return nil
}
func (h *LeasingHistoryHandler) PatchEndTimestamp(c *fiber.Ctx) error {
	return nil
}
