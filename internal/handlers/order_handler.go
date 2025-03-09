package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=nil} "Order created successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
func (o *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	body := new(dto.OrderRequestBody)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "Your request is invalid")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "Your request body is invalid")
	}

	order, err := o.OrderService.CreateOrder(body.LeasingHistoryID)
	if err != nil {
		return err
	}

	data := dto.OrderResponseBody(*order)
	res := dto.Success(data)

	return c.Status(fiber.StatusCreated).JSON(res)
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
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=nil} "Order retrieved successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
func (o *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequestError(err, "Invalid order ID")
	}

	order, serviceErr := o.OrderService.GetOrderByID(orderID)
	if serviceErr != nil {
		return serviceErr
	}

	data := dto.OrderResponseBody(*order)
	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}

// Get Unpaid Order by ID godoc
// @Summary Get unpaid orders by User ID
// @Description Get unpaid orders by User ID
// @Router /order/unpaid/{userID} [get]
// @Tags order
// @Security Bearer
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param limit query string true "Number of history to be retrieved"
// @Param page query string true "Page to retrieved"
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=dto.PaginationResponseBody} "Order retrieved successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
func (o *OrderHandler) GetUnpaidOrderByUserID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequestError(err, "Invalid user ID")
	}

	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}

	orders, totalPage, totalRows, errHandler := o.OrderService.GetUnpaidOrderByUserID(userID, limit, page)
	if errHandler != nil {
		return errHandler
	}

	responseData := make([]dto.OrderResponseBody, len(orders))
	for i, order := range orders {
		responseData[i] = dto.OrderResponseBody(order)
	}

	pagination := dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	}

	res := dto.SuccessPagination(responseData, pagination)

	return c.Status(fiber.StatusOK).JSON(res)
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
// @Success 200 {object} httpResponse.HttpResponse{data=dto.OrderResponseBody,pagination=dto.PaginationResponseBody} "Unpaid orders retrieved successfully"
// @Failure 400 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is invalid"
// @Failure 401 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "your request is unauthorized"
// @Failure 404 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "leasing history not found"
// @Failure 500 {object} httpResponse.HttpResponse{data=nil,pagination=nil} "cannot parse uuid or cannot delete user"
func (o *OrderHandler) GetMyUnpaidOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	limit := c.QueryInt("limit", 1)
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}

	orders, totalPage, totalRows, errHandler := o.OrderService.GetUnpaidOrderByUserID(userID, limit, page)
	if errHandler != nil {
		return errHandler
	}

	responseData := make([]dto.OrderResponseBody, len(orders))
	for i, order := range orders {
		responseData[i] = dto.OrderResponseBody(order)
	}

	pagination := dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	}

	res := dto.SuccessPagination(responseData, pagination)

	return c.Status(fiber.StatusOK).JSON(res)
}
