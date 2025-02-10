package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(UserService ports.UserService) ports.UserHandler {
	return &UserHandler{userService: UserService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequestBody)
	err := c.BodyParser(&user)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	gormUser := &domain.User{
		Email:    user.Email,
		Username: user.UserName,
		Password: user.Password,
	}
	token, err := h.userService.Create(gormUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("user successfully registered", fiber.Map{"token": token, "user": gormUser}))

}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	err := c.BodyParser(&req)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	token, loginErr := h.userService.Login(req.Email, req.Password)
	if loginErr != nil {
		error_handler.InternalServerError(err, "system cannot login to your account")
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("user successfully registered", fiber.Map{"token": token}))
}

func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var requestBody *dto.UserInformationRequestBody

	user := c.Locals("user").(*domain.User)
	if user == nil {
		return error_handler.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := c.BodyParser(&requestBody)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(requestBody); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.userService.UpdateInformation(user.ID, *requestBody)

	if err != nil {
		return error_handler.InternalServerError(err, "system cannot update your account information")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("user successfully updated account information", userInfo))

}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Get("token")
	if tokenString == "" {
		return error_handler.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	if err := h.userService.VerifyUser(tokenString); err != nil {
		return error_handler.InternalServerError(err, "cannot verify your account")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("email is sent to user successfully", nil))
}

func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordCreateRequestBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	err := h.userService.ResetPasswordCreate(body.Email)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("email is sent to user successfully", nil))
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordRequestBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	tokenString := body.Token
	if tokenString == "" {
		return error_handler.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	user, err := h.userService.ResetPasswordResponse(tokenString, body.Password)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("password reset successfully", user))
}

func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("get user information successfully", user))

}
