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

	res, err := ct.contractService.GetContract(reqBody.LessorID, reqBody.LesseeID, reqBody.DormID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get contract error")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res))

}

func (ct *ContractHandler) SignContract(c *fiber.Ctx) error {
	userID, err := c.Locals("userID").(uuid.UUID)
	if !err {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	var reqBody *dto.ContractRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return apperror.BadRequestError(err, "Failed to parse request body")
	}

	if userID != reqBody.LessorID && userID != reqBody.LesseeID {
		return apperror.BadRequestError(errors.New("user ID does not match anyone in the contract"), "You are not authorized to update this contract")
	}

	if err := ct.contractService.UpdateStatus(*reqBody, domain.Signed, userID); err != nil {
		return err
	}

	res, getErr := ct.contractService.GetContract(reqBody.LessorID, reqBody.LesseeID, reqBody.DormID)
	if getErr != nil {
		if apperror.IsAppError(getErr) {
			return getErr
		}
		return apperror.InternalServerError(getErr, "get contract error")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res))

}
