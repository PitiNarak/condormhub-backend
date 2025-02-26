package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
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

// Create Transaction godoc
// @Summary Create a transaction
// @Description Create a transaction
// @Router /transaction [post]
// @Tags transaction
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.TransactionRequestBody true "Transaction request body"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.CreateTransactionResponseBody,pagination=nil} "account successfully deleted"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var reqBody *dto.TransactionRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return errorHandler.BadRequestError(err, "Failed to parse request body")
	}

	_, url, err := h.tsxService.CreateTransaction(reqBody.OrderID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(httpResponse.SuccessResponse("Order created successfully", dto.CreateTransactionResponseBody{
		CheckoutUrl: *url,
	}))
}

func (h *TransactionHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()

	event, err := webhook.ConstructEventWithOptions(payload, c.Get("Stripe-Signature"), h.stripeConfig.StripeSignatureKey, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return errorHandler.BadRequestError(err, "Failed to construct event")
	}

	updateErr := h.tsxService.UpdateTransactionStatus(event)
	if updateErr != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
