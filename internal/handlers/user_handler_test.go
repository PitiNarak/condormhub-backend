package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/handlers"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/http_response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of ports.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) Login(email, password string) (*domain.User, string, error) {
	args := m.Called(email, password)
	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

func (m *MockUserService) VerifyUser(token string) (string, *domain.User, error) {
	args := m.Called(token)
	return args.String(0), args.Get(1).(*domain.User), args.Error(2)
}

func (m *MockUserService) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error) {
	args := m.Called(userID, data)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) ResetPasswordCreate(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(token string, password string) (*domain.User, error) {
	args := m.Called(token, password)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) DeleteAccount(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func setupTest() (*fiber.App, *MockUserService) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			// var e *error_handler.ErrorHandler
			if err != nil {
				if e, ok := err.(*error_handler.ErrorHandler); ok {
					code = e.Code
					message = e.Message
				} else {
					message = err.Error()
				}
			}

			return c.Status(code).JSON(http_response.HttpResponse{
				Success: false,
				Message: message,
				Data:    nil,
			})
		},
	})
	
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	// Setup routes
	app.Post("/auth/register", handler.Register)
	app.Post("/auth/login", handler.Login)
	app.Post("/user/verify", handler.VerifyEmail)
	app.Delete("/user", handler.DeleteAccount)

	return app, mockService
}

func TestRegister(t *testing.T) {
	app, mockService := setupTest()

	t.Run("successful registration", func(t *testing.T) {
		registerData := dto.RegisterRequestBody{
			Email:    "test@example.com",
			UserName: "testuser",
			Password: "password123",
		}

		expectedUser := &domain.User{
			Email:    registerData.Email,
			Username: registerData.UserName,
			Password: registerData.Password,
		}

		expectedToken := "dummy-token"
		mockService.On("Create", mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == expectedUser.Email && u.Username == expectedUser.Username
		})).Return(expectedToken, nil).Once()

		body, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response http_response.HttpResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "successfully registered")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Send empty request body
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response http_response.HttpResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "request body is incorrect")
	})
}

func TestLogin(t *testing.T) {
	app, mockService := setupTest()

	t.Run("successful login", func(t *testing.T) {
		loginData := dto.LoginRequestBody{
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := &domain.User{
			Email:    loginData.Email,
			Password: loginData.Password,
		}

		expectedToken := "dummy-token"
		mockService.On("Login", loginData.Email, loginData.Password).
			Return(expectedUser, expectedToken, nil).Once()

		body, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response http_response.HttpResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "Login successful")
	})
}

func TestVerifyEmail(t *testing.T) {
	app, mockService := setupTest()

	t.Run("successful email verification", func(t *testing.T) {
		verifyData := dto.VerifyRequestBody{
			Token: "valid-token",
		}

		expectedUser := &domain.User{
			Email:      "test@example.com",
			IsVerified: true,
		}

		mockService.On("VerifyUser", verifyData.Token).
			Return(verifyData.Token, expectedUser, nil).Once()

		body, _ := json.Marshal(verifyData)
		req, _ := http.NewRequest("POST", "/user/verify", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response http_response.HttpResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "email is verified successfully")
	})
}

func TestDeleteAccount(t *testing.T) {
	app, mockService := setupTest()

	t.Run("successful account deletion", func(t *testing.T) {
		userID := uuid.New()

		// Setup middleware to inject user ID into context
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user", &domain.User{ID: userID})
			c.Locals("userID", userID.String())
			return c.Next()
		})

		mockService.On("DeleteAccount", userID).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/user", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response http_response.HttpResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "account successfully deleted")
	})

	t.Run("missing user ID in context", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/user", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}