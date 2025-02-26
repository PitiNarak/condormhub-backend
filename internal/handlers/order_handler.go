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

// Create Order godoc
// @Summary Create an order
// @Description Create an order
// @Router /order [post]
// @Tags order
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body dto.OrderRequestBody true "Order request body"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=nil} "account successfully deleted"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
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

// Get Order by ID godoc
// @Summary Get an order by ID
// @Description Get an order by ID
// @Router /order/{id} [get]
// @Tags order
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=nil} "account successfully deleted"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
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

// Get Unpaid Order by ID godoc
// @Summary Get unpaid orders by ID
// @Description Get unpaid orders by ID
// @Router /order/unpaid/{userID} [get]
// @Tags order
// @Security Bearer
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param limit query string true "Number of history to be retrieved"
// @Param page query string true "Page to retrieved"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=dto.PaginationResponseBody} "account successfully deleted"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
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

// Get MT Unpaid Order by ID godoc
// @Summary Get my unpaid orders by ID
// @Description Get my unpaid orders by ID
// @Router /order/unpaid/me [get]
// @Tags order
// @Security Bearer
// @Accept json
// @Produce json
// @Param limit query string true "Number of history to be retrieved"
// @Param page query string true "Page to retrieved"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=dto.PaginationResponseBody} "account successfully deleted"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
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
