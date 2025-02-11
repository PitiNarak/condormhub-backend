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

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequestBody true "user information"
// @Success 201  {object}  http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "user successfully registered"
// @Failure 400  {object}  http_response.HttpResponse{data=nil} "your request is invalid"
// @Failure 500  {object}  http_response.HttpResponse{data=nil} "system cannot register user"
// @Router /auth/register [post]
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

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     token,
		UserInformation: *gormUser,
	}

	return c.Status(fiber.StatusCreated).JSON(http_response.SuccessResponse("user successfully registered", response))

}

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequestBody true "user information"
// @Success 200 {object} http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "user successfully logged in"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 404 {object} http_response.HttpResponse{data=nil} "user not found"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot login user"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}

	user, token, loginErr := h.userService.Login(req.Email, req.Password)
	if loginErr != nil {
		return loginErr
	}

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     token,
		UserInformation: *user,
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("Login successful", response))
}

// UpdateUserInformation godoc
// @Summary Update user information
// @Description Update user information
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body dto.UserInformationRequestBody true "user information"
// @Success 200 {object} http_response.HttpResponse{data=domain.User} "user successfully updated account information"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot update your account information"
// @Router /user [patch]
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

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.VerifyRequestBody true "token"
// @Success 200 {object} http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "email is verified successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot verify your email"
// @Router /user/verify [post]
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	body := new(dto.VerifyRequestBody)

	if err := c.BodyParser(body); err != nil {
		return error_handler.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return error_handler.BadRequestError(err, "your request body is incorrect")
	}
	accessToken, user, err := h.userService.VerifyUser(body.Token)
	if err != nil {
		return err
	}

	response := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		UserInformation: *user,
	}

	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("email is verified successfully", response))
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resend verification email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordCreateRequestBody true "token"
// @Success 200 {object} http_response.HttpResponse{data=nil} "email is sent to user successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot resend verification email"
// @Router /user/resetpassword [post]
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

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordRequestBody true "token"
// @Success 200 {object} http_response.HttpResponse{data=dto.TokenWithUserInformationResponseBody} "password reset successfully"
// @Failure 400 {object} http_response.HttpResponse{data=nil} "your request is invalid
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot reset password"
// @Router /user/newpassword [post]
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

	user, err := h.userService.ResetPassword(tokenString, body.Password)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("password reset successfully", fiber.Map{
		"userInformation": user,
		"accessToken":     tokenString,
	}))
}

// GetUserInfo godoc
// @Summary Get user information
// @Description Get user information
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {object} http_response.HttpResponse{data=domain.User} "get user information successfully"
// @Failure 401 {object} http_response.HttpResponse{data=nil} "your request is unauthorized"
// @Failure 500 {object} http_response.HttpResponse{data=nil} "system cannot get user information"
// @Router /user/me [get]
func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	return c.Status(fiber.StatusOK).JSON(http_response.SuccessResponse("get user information successfully", user))

}
