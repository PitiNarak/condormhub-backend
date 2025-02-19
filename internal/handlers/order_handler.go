package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v81/webhook"
)

type OrderHandler struct {
	orderService ports.OrderService
	stripeConfig *stripe.Config
}

func NewOrderHandler(orderService ports.OrderService, stripeConfig *stripe.Config) ports.OrderHandler {
	return &OrderHandler{orderService: orderService, stripeConfig: stripeConfig}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var reqBody *dto.OrderRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return error_handler.BadRequestError(err, "Failed to parse request body")
	}

	user := c.Locals("user").(*domain.User)
	order, url, err := h.orderService.CreateOrder(domain.InsuranceOrder, reqBody.DormitoryID, reqBody.LessorID, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("Order created successfully", dto.CreateOrderResponseBody{
		Order:       *order,
		CheckoutUrl: *url,
	}))
}

func (h *OrderHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()

	event, err := webhook.ConstructEventWithOptions(payload, c.Get("Stripe-Signature"), h.stripeConfig.StripeSignatureKey, webhook.ConstructEventOptions{
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
