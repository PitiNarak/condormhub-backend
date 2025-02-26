package handlers

import (
	"errors"

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

	responseData := dto.OrderResponseBody(*order)

	return c.Status(fiber.StatusCreated).JSON(httpResponse.SuccessResponse("Order successfully created", responseData))
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

func (o *OrderHandler) GetUnpaidOrderByUserID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errorHandler.BadRequestError(err, "Invalid user ID")
	}

	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return errorHandler.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return errorHandler.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}

	orders, totalPage, totalRows, errHandler := o.OrderService.GetUnpaidOrderByUserID(userID, limit, page)
	if errHandler != nil {
		return errHandler
	}

	responseData := make([]dto.OrderResponseBody, len(orders))
	for i, order := range orders {
		responseData[i] = dto.OrderResponseBody(order)
	}

	pageination := dto.PaginationResponseBody{
		Currentpage: page,
		Lastpage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessPageResponse("Orders successfully retrieved", responseData, pageination))
}

func (o *OrderHandler) GetMyUnpaidOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return errorHandler.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return errorHandler.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}

	orders, totalPage, totalRows, errHandler := o.OrderService.GetUnpaidOrderByUserID(userID, limit, page)
	if errHandler != nil {
		return errHandler
	}

	responseData := make([]dto.OrderResponseBody, len(orders))
	for i, order := range orders {
		responseData[i] = dto.OrderResponseBody(order)
	}

	pageination := dto.PaginationResponseBody{
		Currentpage: page,
		Lastpage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	}

	return c.Status(fiber.StatusOK).JSON(httpResponse.SuccessPageResponse("Orders successfully retrieved", responseData, pageination))
}
