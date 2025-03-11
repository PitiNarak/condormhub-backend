package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
// @Success 200 {object} dto.SuccessResponse[dto.CreateTransactionResponseBody] "Transaction created successfully"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "leasing history not found"
// @Failure 500 {object} dto.ErrorResponse "cannot parse uuid or cannot delete user"
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var reqBody *dto.TransactionRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return apperror.BadRequestError(err, "Failed to parse request body")
	}

	_, url, err := h.tsxService.CreateTransaction(reqBody.OrderID)
	if err != nil {
		return err
	}

	res := dto.Success(dto.CreateTransactionResponseBody{
		CheckoutUrl: *url,
	})

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *TransactionHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()

	event, err := webhook.ConstructEventWithOptions(payload, c.Get("Stripe-Signature"), h.stripeConfig.StripeSignatureKey, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return apperror.BadRequestError(err, "Failed to construct event")
	}

	updateErr := h.tsxService.UpdateTransactionStatus(event)
	if updateErr != nil {
		return updateErr
	}

	return c.SendStatus(fiber.StatusOK)
}
