package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v81/webhook"
)

type TransactionHandler struct {
	tsxService   ports.TransactionService
	stripeConfig *stripe.Config
}

func NewTransactionHandler(orderService ports.TransactionService, stripeConfig *stripe.Config) ports.TransactionHandler {
	return &TransactionHandler{tsxService: orderService, stripeConfig: stripeConfig}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var reqBody *dto.TransactionRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return error_handler.BadRequestError(err, "Failed to parse request body")
	}

	_, url, err := h.tsxService.CreateTransaction(reqBody.OrderID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("Order created successfully", dto.CreateTransactionResponseBody{
		CheckoutUrl: *url,
	}))
}

func (h *TransactionHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()

	event, err := webhook.ConstructEventWithOptions(payload, c.Get("Stripe-Signature"), h.stripeConfig.StripeSignatureKey, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return error_handler.BadRequestError(err, "Failed to construct event")
	}

	updateErr := h.tsxService.UpdateTransactionStatus(event)
	if updateErr != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
