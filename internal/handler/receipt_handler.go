package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	receiptService ports.ReceiptService
}

func NewReceiptHandler(service ports.ReceiptService) ports.ReceiptHandler {
	return &ReceiptHandler{receiptService: service}
}

// GetByUserID godoc
// @Summary Retrieve receipts by user ID
// @Description Get all receipts associated with a specific user
// @Tags receipt
// @Param limit query int false "Number of receipts to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.ReceiptResponseBody] "Receipts retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve receipts"
// @Router /receipt [get]
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
