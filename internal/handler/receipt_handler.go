package handler

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	receiptService ports.ReceiptService
}

func NewReceiptHandler(service ports.ReceiptService) ports.ReceiptHandler {
	return &ReceiptHandler{receiptService: service}
}

// Create godoc
// @Summary create receipt
// @Description create receipt
// @Tags receipt
// @Security Bearer
// @Param transactionId path string true "Transaction ID"
// @Produce json
// @Success 200 {object}  dto.SuccessResponse[dto.ReceiptResponseBody]
// @Failure 400 {object} dto.ErrorResponse "Incorrect UUID format"
// @Failure 404 {object} dto.ErrorResponse "Receipt not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to create file"
// @Router /ownership/{transactionId} [post]
func (r *ReceiptHandler) Create(c *fiber.Ctx) error {
	ownerID := c.Locals("userID").(uuid.UUID)
	id := c.Params("transactionId")

	if err := uuid.Validate(id); err != nil {
		return apperror.BadRequestError(err, "Incorrect UUID format")
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		return apperror.InternalServerError(err, "Can not parse UUID")
	}

	receipt, url, err := r.receiptService.Create(c.Context(), ownerID, transactionID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
	}
	fmt.Println(receipt)

	return c.Status(fiber.StatusOK).JSON(dto.Success(receipt.ToDTO(url)))
}

func (r *ReceiptHandler) GetByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	limit := min(50, c.QueryInt("limit", 10))
	if limit <= 0 {
		limit = 10
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}
	receipts, totalPage, totalRows, err := r.receiptService.GetByUserID(userID, limit, page)
	if err != nil {
		return err
	}

	resData := make([]dto.ReceiptResponseBody, len(receipts))
	for i, v := range receipts {
		url, err := r.receiptService.GetUrl(c.Context(), v)
		if err != nil {
			return err
		}
		resData[i] = v.ToDTO(url)
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}
