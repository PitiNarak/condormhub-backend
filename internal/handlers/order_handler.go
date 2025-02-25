package handlers

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/httpResponse"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderService ports.OrderService
}

func NewOrderHandler(service ports.OrderService) ports.OrderHandler {
	return &OrderHandler{OrderService: service}
}

func (o *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	body := new(dto.OrderRequestBody)
	if err := c.BodyParser(body); err != nil {
		return errorHandler.BadRequestError(err, "Your request is invalid")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return errorHandler.BadRequestError(err, "Your request body is invalid")
	}

	order, err := o.OrderService.CreateOrder(body.LeasingHistoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(httpResponse.SuccessResponse("Order successfully created", order))
}

func (o *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errorHandler.BadRequestError(err, "Invalid order ID")
	}

	order, errHandler := o.OrderService.GetOrderByID(orderID)
	if errHandler != nil {
		return err
	}

	responseData := dto.OrderResponseBody(*order)

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessResponse("Order successfully retrieved", responseData))
}
