package handler

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ContractHandler struct {
	contractService ports.ContractService
}

func NewContractHandler(contractService ports.ContractService) ports.ContractHandler {
	return &ContractHandler{contractService: contractService}
}

// CreateContract godoc
// @Summary Create a new contract
// @Description Create a contract between a lessor and lessee for a dorm
// @Tags contracts
// @Accept json
// @Produce json
// @Param contract body dto.ContractRequestBody true "Contract details"
// @Success 201 {object} dto.SuccessResponse[dto.ContractResponseBody] "Contract created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 500 {object} dto.ErrorResponse "Failed to create contract"
// @Router /contract [post]
func (ct *ContractHandler) Create(c *fiber.Ctx) error {
	var reqBody *dto.ContractRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return apperror.BadRequestError(err, "Failed to parse request body")
	}

	contract := &domain.Contract{
		LessorID: reqBody.LessorID,
		LesseeID: reqBody.LesseeID,
		DormID:   reqBody.DormID,
	}

	if err := ct.contractService.Create(contract); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "create contract error")
	}

	res, err := ct.contractService.GetContractByContractID(contract.ContractID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get contract error")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res.ToDTO()))

}

// SignContract godoc
// @Summary Sign an existing contract
// @Description Sign a contract by providing contract ID
// @Tags contracts
// @Param contractId path string true "Contract ID"
// @Produce json
// @Success 200 {object} dto.SuccessResponse[dto.ContractResponseBody] "Contract signed successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid contract ID format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to sign contract"
// @Router /contract/{contractId}/sign [patch]
func (ct *ContractHandler) SignContract(c *fiber.Ctx) error {
	userID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	contractID, parseErr := uuid.Parse(c.Params("contractId"))
	if parseErr != nil {
		return apperror.BadRequestError(parseErr, "Invalid contract ID format")
	}

	if err := ct.contractService.UpdateStatus(contractID, domain.Signed, userID); err != nil {
		return err
	}

	res, getErr := ct.contractService.GetContractByContractID(contractID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "get contract error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(res.ToDTO()))

}

// CancelContract godoc
// @Summary Cancel an existing contract
// @Description Cancel a contract if it is not signed
// @Tags contracts
// @Param contractId path string true "Contract ID"
// @Produce json
// @Success 200 {object} dto.SuccessResponse[dto.ContractResponseBody] "Contract cancelled successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid contract ID format or contract cannot be cancelled"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to cancel contract"
// @Router /contract/{contractId}/cancel [patch]
func (ct *ContractHandler) CancelContract(c *fiber.Ctx) error {
	userID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	contractID, parseErr := uuid.Parse(c.Params("contractId"))
	if parseErr != nil {
		return apperror.BadRequestError(parseErr, "Invalid contract ID format")
	}
	contract, getErr := ct.contractService.GetContractByContractID(contractID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "get contract error")
	}

	if contract.Status == domain.Signed {
		return apperror.BadRequestError(errors.New("contract is already signed"), "You cannot cancel signed contract")
	}

	if contract.Status == domain.Cancelled {
		return apperror.BadRequestError(errors.New("contract is already cancelld"), "You cannot cancel cancelled contract")
	}

	if err := ct.contractService.UpdateStatus(contractID, domain.Cancelled, userID); err != nil {
		return err
	}

	res, getErr := ct.contractService.GetContractByContractID(contractID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "get contract error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(res.ToDTO()))

}

// DeleteContract godoc
// @Summary Delete a contract
// @Description Permanently delete a contract by ID
// @Tags contracts
// @Param contractId path string true "Contract ID"
// @Success 204 "Contract deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid contract ID format"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete contract"
// @Router /contract/{contractId} [delete]
func (ct *ContractHandler) Delete(c *fiber.Ctx) error {
	contractID, parseErr := uuid.Parse(c.Params("contractId"))
	if parseErr != nil {
		return apperror.BadRequestError(parseErr, "Invalid contract ID format")
	}

	if err := ct.contractService.DeleteContract(contractID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)

}

// GetContractByContractID godoc
// @Summary Retrieve a contract by contract ID
// @Description Get details of a specific contract
// @Tags contracts
// @Param contractId path string true "Contract ID"
// @Produce json
// @Success 200 {object} dto.SuccessResponse[dto.ContractResponseBody] "Contract retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid contract ID format"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve contract"
// @Router /contract/{contractId} [get]
func (ct *ContractHandler) GetContractByContractID(c *fiber.Ctx) error {
	contractID, parseErr := uuid.Parse(c.Params("contractId"))
	if parseErr != nil {
		return apperror.BadRequestError(parseErr, "Invalid contract ID format")
	}

	contract, err := ct.contractService.GetContractByContractID(contractID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(contract.ToDTO()))

}

// GetContractByUserID godoc
// @Summary Retrieve contracts by user ID
// @Description Get all contracts associated with a specific user
// @Tags contracts
// @Param limit query int false "Number of contracts to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.ContractResponseBody] "Contracts retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve contracts"
// @Router /contract [get]
func (ct *ContractHandler) GetContractByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	limit := min(50, c.QueryInt("limit", 10))
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	contracts, totalPage, totalRows, err := ct.contractService.GetByUserID(userID, limit, page)
	if err != nil {
		return err
	}

	resData := make([]dto.ContractResponseBody, len(*contracts))
	for i, v := range *contracts {
		resData[i] = v.ToDTO()
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetContractByDormID godoc
// @Summary Retrieve contracts by dorm ID
// @Description Get all contracts associated with a specific dorm
// @Tags contracts
// @Param dormId path string true "Dorm ID"
// @Param limit query int false "Number of contracts to retrieve (default 10, max 50)"
// @Param page query int false "Page number to retrieve (default 1)"
// @Produce json
// @Success 200 {object} dto.PaginationResponse[dto.ContractResponseBody] "Contracts retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid dorm ID format or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve contracts"
// @Router /contract/{dormId} [get]
func (ct *ContractHandler) GetContractByDormID(c *fiber.Ctx) error {
	id := c.Params("dormId")
	dormID, err := uuid.Parse(id)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "Can not parse UUID")
	}
	limit := min(50, c.QueryInt("limit", 10))
	if limit <= 0 {
		return apperror.BadRequestError(errors.New("limit parameter is incorrect"), "limit parameter is incorrect")
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return apperror.BadRequestError(errors.New("page parameter is incorrect"), "page parameter is incorrect")
	}
	contracts, totalPage, totalRows, err := ct.contractService.GetByDormID(dormID, limit, page)
	if err != nil {
		return err
	}
	resData := make([]dto.ContractResponseBody, len(*contracts))
	for i, v := range *contracts {
		resData[i] = v.ToDTO()
	}

	res := dto.SuccessPagination(resData, dto.Pagination{
		CurrentPage: page,
		LastPage:    totalPage,
		Limit:       limit,
		Total:       totalRows,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}
