package handlers

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService ports.UserService
}

func NewUserHandler(UserService ports.UserService) ports.UserHandler {
	return &UserHandler{UserService: UserService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(domain.UserBody)
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
		UserName: user.UserName,
		Password: user.Password,
	}
	token, err := h.UserService.Create(gormUser)
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

	token, loginErr := h.UserService.Login(req.Email, req.Password)
	if loginErr != nil {
		error_handler.InternalServerError(err, "system cannot login to your account")
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("user successfully registered", fiber.Map{"token": token}))
}

func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var user domain.User
	err := c.BodyParser(&user)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	updateInfo := domain.UpdateInfo{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		NationalID:      user.NationalID,
		Gender:          user.Gender,
		BirthDate:       user.BirthDate,
		StudentEvidence: user.StudentEvidence,
	}

	err = h.UserService.Update(user, updateInfo)

	if err != nil {
		return error_handler.InternalServerError(err, "system cannot update your account information")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("user successfully updated account information", nil))

}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	tokenString := c.Get("token")
	if tokenString == "" {
		return error_handler.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	if err := h.UserService.VerifyUser(tokenString); err != nil {
		return error_handler.InternalServerError(err, "cannot verify your account")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("email is sent to user successfully", nil))
}

func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(domain.ResetPasswordBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	err := h.UserService.ResetPasswordCreate(body.Email)
	if err != nil {
		return error_handler.InternalServerError(err, "cannot sent email to reset password")
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("email is sent to user successfully", nil))
}

func (h *UserHandler) ResetPasswordResponse(c *fiber.Ctx) error {
	body := new(domain.ResponseResetPasswordBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	tokenString := c.Get("token")
	if tokenString == "" {
		return error_handler.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	err := h.UserService.ResetPasswordResponse(tokenString, body.Password)
	if err != nil {
		return error_handler.InternalServerError(err, "cannot reset user password")
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("password reset successfully", nil))
}

func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	var getInfoRequest domain.GetInfoRequest
	err := c.BodyParser(&getInfoRequest)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(getInfoRequest); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.UserService.GetUserByEmail(getInfoRequest.Email)

	if err != nil {
		return error_handler.InternalServerError(err, "cannot get user information")
	}

	publicUserInfo := domain.UserInfo{
		UserName:          userInfo.UserName,
		Email:             userInfo.Email,
		FirstName:         userInfo.FirstName,
		LastName:          userInfo.LastName,
		NationalID:        userInfo.NationalID,
		Gender:            userInfo.Gender,
		BirthDate:         userInfo.BirthDate,
		IsVerified:        userInfo.IsVerified,
		Role:              userInfo.Role,
		StudentEvidence:   userInfo.StudentEvidence,
		IsStudentVerified: userInfo.IsStudentVerified,
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("get user information successfully", publicUserInfo))

}
