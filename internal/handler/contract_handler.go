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

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res.ToDTO()))

}

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

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res.ToDTO()))

}

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
