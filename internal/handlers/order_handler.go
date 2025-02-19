package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v81/webhook"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(orderService ports.OrderService) ports.OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var reqBody *dto.OrderBody
	if err := c.BodyParser(&reqBody); err != nil {
		return error_handler.BadRequestError(err, "Failed to parse request body")
	}

	order, err := h.orderService.CreateOrder(domain.InsuranceOrder, reqBody.DormitoryID, reqBody.LessorID, reqBody.LesseeID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("Order created successfully", order))
}

func (h *OrderHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()

	event, err := webhook.ConstructEventWithOptions(payload, c.Get("Stripe-Signature"), "whsec_f178433184d6a3f42bc7da6261096ba5acacd1088d1a8d568113f29f5fa90a2a", webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return error_handler.BadRequestError(err, "Failed to construct event")
	}

	updateErr := h.orderService.UpdateOrderStatus(event)
	if updateErr != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
