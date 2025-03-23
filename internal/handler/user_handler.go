package handler

import (
	"errors"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(UserService ports.UserService) *UserHandler {
	return &UserHandler{userService: UserService}
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.VerifyRequestBody true "token"
// @Success 200 {object} dto.SuccessResponse[dto.TokenWithUserInformationResponseBody] "email is verified successfully"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "system cannot verify your email"
// @Router /user/verify [post]
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	body := new(dto.VerifyRequestBody)

	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	user, accessToken, refreshToken, err := h.userService.VerifyUser(c.Context(), body.Token)
	if err != nil {
		return err
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: h.userService.ConvertToDTO(*user),
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateUserInformation godoc
// @Summary Update user information
// @Description Update user information
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body dto.UserInformationRequestBody true "user information"
// @Success 200 {object} dto.SuccessResponse[dto.UserResponse] "user successfully updated account information"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "system cannot update your account information"
// @Router /user [patch]
func (h *UserHandler) UpdateUserInformation(c *fiber.Ctx) error {
	var requestBody *dto.UserInformationRequestBody

	user := c.Locals("user").(*domain.User)
	if user == nil {
		return apperror.UnauthorizedError(errors.New("no user in context"), "your request is unauthorized")
	}

	err := c.BodyParser(&requestBody)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()
	lifeStyleErr := validate.RegisterValidation("lifestyle", utils.ValidateLifestyles)
	if lifeStyleErr != nil {
		return apperror.BadRequestError(lifeStyleErr, "your lifestyle-tag is incorrect format")
	}

	phoneNumberErr := validate.RegisterValidation("phoneNumber", utils.ValidatePhone)
	if phoneNumberErr != nil {
		return apperror.BadRequestError(phoneNumberErr, "your phone number is incorrect format")
	}

	if err := validate.Struct(requestBody); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	userInfo, err := h.userService.UpdateInformation(user.ID, *requestBody)

	if err != nil {
		return apperror.InternalServerError(err, "system cannot update your account information")
	}

	res := dto.Success(h.userService.ConvertToDTO(*userInfo))

	return c.Status(fiber.StatusOK).JSON(res)

}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordRequestBody true "token"
// @Success 200 {object} dto.SuccessResponse[dto.TokenWithUserInformationResponseBody] "password reset successfully"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid
// @Failure 500 {object} dto.ErrorResponse "system cannot reset password"
// @Router /user/newpassword [post]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordRequestBody)

	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	tokenString := body.Token
	if tokenString == "" {
		return apperror.BadRequestError(errors.New("no token in header"), "your request header is incorrect")
	}

	user, accessToken, refreshToken, err := h.userService.ResetPassword(c.Context(), tokenString, body.Password)
	if err != nil {
		return err
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: user.ToDTO(),
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resend verification email
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.ResetPasswordCreateRequestBody true "token"
// @Success 204 "email is sent to user successfully"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid
// @Failure 500 {object} dto.ErrorResponse "system cannot resend verification email"
// @Router /user/resetpassword [post]
func (h *UserHandler) ResetPasswordCreate(c *fiber.Ctx) error {
	body := new(dto.ResetPasswordCreateRequestBody)

	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	err := h.userService.ResetPasswordCreate(c.Context(), body.Email)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequestBody true "user information"
// @Success 201 {object} dto.SuccessResponse[dto.TokenWithUserInformationResponseBody] "user successfully registered"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid"
// @Failure 500 {object} dto.ErrorResponse "system cannot register user"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(dto.RegisterRequestBody)
	err := c.BodyParser(&user)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}
	gormUser := &domain.User{
		Email:    user.Email,
		Username: user.UserName,
		Password: user.Password,
	}

	accessToken, refreshToken, err := h.userService.Create(c.Context(), gormUser)
	if err != nil {
		return err
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: h.userService.ConvertToDTO(*gormUser),
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusCreated).JSON(res)

}

// Login godoc
// @Summary Refresh user
// @Description Refresh user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RefreshTokenRequestBody true "user information"
// @Success 200 {object} dto.SuccessResponse[dto.TokenResponseBody] "user successfully Refresh in"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "user not found"
// @Failure 500 {object} dto.ErrorResponse "system cannot refresh user"
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	accessToken, refreshToken, loginErr := h.userService.RefreshToken(c.Context(), req.RefreshToken)
	if loginErr != nil {
		return loginErr
	}

	data := dto.TokenResponseBody{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequestBody true "user information"
// @Success 200 {object} dto.SuccessResponse[dto.TokenWithUserInformationResponseBody] "user successfully logged in"
// @Failure 400 {object} dto.ErrorResponse "your request is invalid"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "user not found"
// @Failure 500 {object} dto.ErrorResponse "system cannot login user"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestBody
	err := c.BodyParser(&req)
	if err != nil {
		return apperror.BadRequestError(err, "your request is invalid")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperror.BadRequestError(err, "your request body is incorrect")
	}

	user, accessToken, refreshToken, loginErr := h.userService.Login(c.Context(), req.Email, req.Password)
	if loginErr != nil {
		return loginErr
	}

	data := dto.TokenWithUserInformationResponseBody{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserInformation: h.userService.ConvertToDTO(*user),
	}

	res := dto.Success(data)

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetUserInfo godoc
// @Summary Get user information
// @Description Get user information
// @Tags user
// @Security Bearer
// @Produce json
// @Success 200 {object} dto.SuccessResponse[dto.UserResponse] "get user information successfully"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "system cannot get user information"
// @Router /user/me [get]
func (h *UserHandler) GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	res := dto.Success(h.userService.ConvertToDTO(*user))
	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteAccount godoc
// @Summary Delete a user account
// @Description Delete a user account
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Success 204 "account successfully deleted"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "cannot parse uuid or cannot delete user"
// @Router /user/ [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	err := h.userService.DeleteAccount(userID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetUserByID godoc
// @Summary GetUserByID
// @Description Get User By ID
// @Tags user
// @Security Bearer
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} dto.SuccessResponse[dto.UserResponse] "get user information successfully"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 500 {object} dto.ErrorResponse "system cannot get user information"
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequestError(err, "invalid user id")
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get user by id failed")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(h.userService.ConvertToDTO(*user)))
}

// UploadStudentEvidence godoc
// @Summary Upload an evidence for student verification
// @Description Upload an image of a student ID card, by attaching the image as a value for the key field name "image", as a multipart form-data
// @Tags user
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Student ID image"
// @Success 200 {object} dto.SuccessResponse[dto.StudentEvidenceUploadResponseBody] "Evidence uploaded successfully"
// @Failure 400 {object} dto.ErrorResponse "File is required"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to upload file"
// @Router /user/studentEvidence [post]
func (h *UserHandler) UploadStudentEvidence(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	file, err := c.FormFile("image")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	contentType := file.Header.Get("Content-Type")
	url, err := h.userService.UploadStudentEvidence(c.Context(), file.Filename, contentType, fileData, userID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(dto.StudentEvidenceUploadResponseBody{ImageUrl: url, Expired: time.Now().Add(time.Hour)}))
}

// GetStudentEvidenceByID godoc
// @Summary Get student evidence by user id
// @Description Get student evidence by user id
// @Tags user
// @Security Bearer
// @Produce json
// @Param id path string true "userID"
// @Success 200 {object} dto.SuccessResponse[dto.StudentEvidenceUploadResponseBody] "Get student evidence successfully"
// @Failure 400 {object} dto.ErrorResponse "invalid user id"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 403 {object} dto.ErrorResponse "unauthorized to view this evidence"
// @Failure 404 {object} dto.ErrorResponse "User or evidence not found"
// @Failure 500 {object} dto.ErrorResponse "system cannot get user's student evidence"
// @Router /user/{id}/studentEvidence [get]
func (h *UserHandler) GetStudentEvidenceByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequestError(err, "invalid user id")
	}

	localUser := c.Locals("user").(*domain.User)
	if localUser.Role == "" {
		return apperror.UnauthorizedError(errors.New("unauthorized"), "user role is missing")
	}
	isSelf := localUser.ID == userID
	isAdmin := localUser.Role == domain.AdminRole

	evidence, err := h.userService.GetStudentEvidenceByID(c.Context(), userID, isSelf, isAdmin)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(err, "get student evidence by id failed")
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(evidence))
}

// SendConfirmationEmailAgain godoc
// @Summary SendConfirmationEmailAgain
// @Description Re send the confirmation email
// @Tags user
// @Security Bearer
// @Produce json
// @Success 204 "resend verification email successfully"
// @Failure 500 {object} dto.ErrorResponse "system cannot verification email"
// @Router /user/resend [post]
func (h *UserHandler) ResendVerificationEmailHandler(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("can't get user form context"), "get user information error")
	}
	if err := h.userService.ResendVerificationEmailService(c.Context(), user.Email); err != nil {
		if apperror.IsAppError(err) {
			return err
		}
		return apperror.InternalServerError(errors.New("can't send confirmation email"), "resend email error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// UploadProfilePicture godoc
// @Summary Upload user profile picture
// @Description Upload an profile picture for the current user, by attaching the image as a value for the key field name "image", as a multipart form-data
// @Tags user
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile picture image"
// @Success 200 {object} dto.SuccessResponse[dto.ProfilePictureUploadResponseBody] "Profile picture updated"
// @Failure 400 {object} dto.ErrorResponse "File is required"
// @Failure 401 {object} dto.ErrorResponse "your request is unauthorized"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Server failed to upload file"
// @Router /user/profilePic [post]
func (h *UserHandler) UploadProfilePicture(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	file, err := c.FormFile("image")
	if err != nil {
		return apperror.BadRequestError(err, "file is required")
	}

	fileData, err := file.Open()
	if err != nil {
		return apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	contentType := file.Header.Get("Content-Type")
	url, err := h.userService.UploadProfilePicture(c.Context(), file.Filename, contentType, fileData, userID)
	if err != nil {
		if apperror.IsAppError(err) {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.Success(dto.ProfilePictureUploadResponseBody{ImageURL: url}))
}
