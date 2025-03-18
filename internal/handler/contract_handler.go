package handler

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	contractService ports.ContractService
}

func NewContractHandler(contractService ports.ContractService) ports.ContractHandler {
	return &ContractHandler{contractService: contractService}
}

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

	res, err := ct.contractService.GetByDormID(reqBody.DormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get contract error")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res))

}
