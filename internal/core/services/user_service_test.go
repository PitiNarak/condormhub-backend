package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/services"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error {
	args := m.Called(userID, data)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteAccount(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}


// Mock EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerificationEmail(email, name, token string) error {
	args := m.Called(email, name, token)
	return args.Error(0)
}

func (m *MockEmailService) SendResetPasswordEmail(email, name, token string) error {
	args := m.Called(email, name, token)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		user          *domain.User
		setupMocks    func(*MockUserRepository, *MockEmailService)
		expectedError error
	}{
		{
			name: "Successful registration",
			user: &domain.User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
				es.On("SendVerificationEmail", "test@example.com", "testuser", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed registration - database error",
			user: &domain.User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "Failed registration - email service error",
			user: &domain.User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
				es.On("SendVerificationEmail", "test@example.com", "testuser", mock.AnythingOfType("string")).Return(errors.New("email service error"))
			},
			expectedError: errors.New("email service error"),
		},
		// {
		// 	name: "Failed registration - missing information",
		// 	user: &domain.User{
		// 		Email:    "",  // Empty email
		// 		Username: "",  // Empty username
		// 		Password: "password123",
		// 	},
		// 	setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
		// 		// No mock setup needed because the error will occur before the repository or email service is called.
		// 	},
		// 	expectedError: errors.New("missing required information"),  // You can customize the error message as needed
		// },		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)
			
			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo, mockEmailService)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}

			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			_,err := userService.Create(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, "password123", tt.user.Password)
			}

			mockUserRepo.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUserID := uuid.New()

	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:     "Successful login",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByEmail", "test@example.com").Return(&domain.User{
					ID:       testUserID,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedError: false,
		},
		{
			name:     "Invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByEmail", "test@example.com").Return(&domain.User{
					ID:       testUserID,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedError: true,
		},
		{
			name:     "User not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMocks: func(ur *MockUserRepository) {
				var nilUser *domain.User
				ur.On("GetUserByEmail", "nonexistent@example.com").Return(nilUser, errors.New("user not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}

			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			_,token, err := userService.Login(tt.email, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				claims, err := jwtUtils.DecodeJWT(token)
				assert.NoError(t, err)
				assert.NotNil(t, claims)

				assert.NotNil(t, claims.UserID)
				userIDStr := claims.GetUserID()
				assert.Equal(t, testUserID.String(), userIDStr)


				exp := claims.GetExp() // Use the GetExp() method to get expiration time
				expTime := time.Unix(exp, 0)
				assert.True(t, expTime.After(time.Now()))
				assert.True(t, expTime.Before(time.Now().Add(time.Hour*25)))
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

//Add more testt
func TestUserService_VerifyUser(t *testing.T) {
	testUserID := uuid.New()
	tests := []struct {
		name          string
		token         string
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:  "Successful verification",
			token: "valid.jwt.token",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:         testUserID,
					IsVerified: false,
				}, nil)
				ur.On("UpdateUser", mock.MatchedBy(func(user *domain.User) bool {
					return user.ID == testUserID && user.IsVerified == true
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name:  "User not found",
			token: "valid.jwt.token",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByID", testUserID).Return(&domain.User{}, errors.New("user not found"))
			},
			expectedError: true,
		},
		{
			name:  "Update error",
			token: "valid.jwt.token",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:         testUserID,
					IsVerified: false,
				}, nil)
				ur.On("UpdateUser", mock.Anything).Return(errors.New("update error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			// Generate a valid token for testing
			validToken, _ := jwtUtils.GenerateJWT(testUserID)
			if tt.token == "valid.jwt.token" {
				tt.token = validToken
			}

			_, _, err := userService.VerifyUser(tt.token)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateInformation(t *testing.T) {
	testUserID := uuid.New()
	testBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name          string
		userID        uuid.UUID
		updateData    dto.UserInformationRequestBody
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:   "Successful full update",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				Username:        "newusername",
				Password:        "newpassword123",
				Firstname:       "John",
				Lastname:        "Doe",
				NationalID:      "1234567890123",
				Gender:          "Male",
				BirthDate:       testBirthDate,
				StudentEvidence: "evidence123",
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(nil)
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:              testUserID,
					Username:        "newusername",
					Firstname:       "John",
					Lastname:        "Doe",
					NationalID:      "1234567890123",
					Gender:          "Male",
					BirthDate:       testBirthDate,
					StudentEvidence: "evidence123",
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Update personal information only",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				Firstname: "Jane",
				Lastname:  "Smith",
				Gender:    "Female",
				BirthDate: testBirthDate,
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(nil)
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:        testUserID,
					Firstname: "Jane",
					Lastname:  "Smith",
					Gender:    "Female",
					BirthDate: testBirthDate,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Update with password only",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				Password: "newstrongpassword123",
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(nil)
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID: testUserID,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Update student information",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				NationalID:      "9876543210123",
				StudentEvidence: "new_evidence456",
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(nil)
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:              testUserID,
					NationalID:      "9876543210123",
					StudentEvidence: "new_evidence456",
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Update error - database failure",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				Username: "newusername",
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(errors.New("database error"))
			},
			expectedError: true,
		},
		{
			name:   "Get user error after update",
			userID: testUserID,
			updateData: dto.UserInformationRequestBody{
				Username: "newusername",
			},
			setupMocks: func(ur *MockUserRepository) {
				ur.On("UpdateInformation", testUserID, mock.AnythingOfType("dto.UserInformationRequestBody")).Return(nil)
				ur.On("GetUserByID", testUserID).Return(&domain.User{}, errors.New("failed to get user"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			user, err := userService.UpdateInformation(tt.userID, tt.updateData)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				
				if tt.updateData.Password != "" {
					assert.NotEqual(t, tt.updateData.Password, user.Password)
				}
				if tt.updateData.Username != "" {
					assert.Equal(t, tt.updateData.Username, user.Username)
				}
				if tt.updateData.Firstname != "" {
					assert.Equal(t, tt.updateData.Firstname, user.Firstname)
				}
				if tt.updateData.Lastname != "" {
					assert.Equal(t, tt.updateData.Lastname, user.Lastname)
				}
				if tt.updateData.NationalID != "" {
					assert.Equal(t, tt.updateData.NationalID, user.NationalID)
				}
				if tt.updateData.Gender != "" {
					assert.Equal(t, tt.updateData.Gender, user.Gender)
				}
				if !tt.updateData.BirthDate.IsZero() {
					assert.Equal(t, tt.updateData.BirthDate, user.BirthDate)
				}
				if tt.updateData.StudentEvidence != "" {
					assert.Equal(t, tt.updateData.StudentEvidence, user.StudentEvidence)
				}
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	testUserID := uuid.New()
	tests := []struct {
		name          string
		email         string
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:  "Successful get user",
			email: "test@example.com",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByEmail", "test@example.com").Return(&domain.User{
					ID:    testUserID,
					Email: "test@example.com",
				}, nil)
			},
			expectedError: false,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByEmail", "nonexistent@example.com").Return(&domain.User{}, errors.New("user not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			user, err := userService.GetUserByEmail(tt.email)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_ResetPasswordCreate(t *testing.T) {
	testUserID := uuid.New()
	tests := []struct {
		name          string
		email         string
		setupMocks    func(*MockUserRepository, *MockEmailService)
		expectedError bool
	}{
		{
			name:  "Successful reset password request",
			email: "test@example.com",
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("GetUserByEmail", "test@example.com").Return(&domain.User{
					ID:       testUserID,
					Email:    "test@example.com",
					Username: "testuser",
				}, nil)
				es.On("SendResetPasswordEmail", "test@example.com", "testuser", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("GetUserByEmail", "nonexistent@example.com").Return(&domain.User{}, errors.New("user not found"))
			},
			expectedError: true,
		},
		{
			name:  "Email service error",
			email: "test@example.com",
			setupMocks: func(ur *MockUserRepository, es *MockEmailService) {
				ur.On("GetUserByEmail", "test@example.com").Return(&domain.User{
					ID:       testUserID,
					Email:    "test@example.com",
					Username: "testuser",
				}, nil)
				es.On("SendResetPasswordEmail", "test@example.com", "testuser", mock.AnythingOfType("string")).Return(errors.New("email service error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo, mockEmailService)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			err := userService.ResetPasswordCreate(tt.email)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestUserService_ResetPassword(t *testing.T) {
	testUserID := uuid.New()
	tests := []struct {
		name          string
		token         string
		password      string
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:     "Successful password reset",
			token:    "valid.jwt.token",
			password: "newpassword123",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByID", testUserID).Return(&domain.User{
					ID:       testUserID,
					Password: "oldpassword",
				}, nil)
				ur.On("UpdateUser", mock.MatchedBy(func(user *domain.User) bool {
					return user.ID == testUserID && user.Password != "oldpassword"
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "Invalid token",
			token:    "invalid.token",
			password: "newpassword123",
			setupMocks: func(ur *MockUserRepository) {
				// No mocks needed as it will fail at token validation
			},
			expectedError: true,
		},
		{
			name:     "User not found",
			token:    "valid.jwt.token",
			password: "newpassword123",
			setupMocks: func(ur *MockUserRepository) {
				ur.On("GetUserByID", testUserID).Return(&domain.User{}, errors.New("user not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			// Generate a valid token for testing
			if tt.token == "valid.jwt.token" {
				tt.token, _ = jwtUtils.GenerateJWT(testUserID)
			}

			user, err := userService.ResetPassword(tt.token, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteAccount(t *testing.T) {
	testUserID := uuid.New()
	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMocks    func(*MockUserRepository)
		expectedError bool
	}{
		{
			name:   "Successful account deletion",
			userID: testUserID,
			setupMocks: func(ur *MockUserRepository) {
				ur.On("DeleteAccount", testUserID).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "Delete error",
			userID: testUserID,
			setupMocks: func(ur *MockUserRepository) {
				ur.On("DeleteAccount", testUserID).Return(errors.New("delete error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockEmailService := new(MockEmailService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockUserRepo)
			}

			jwtConfig := &utils.JWTConfig{
				JWTSecretKey: "test-secret-key",
				Expiration:   24,
			}
			jwtUtils := utils.NewJWTUtils(jwtConfig)

			userService := services.NewUserService(mockUserRepo, mockEmailService, jwtUtils)

			err := userService.DeleteAccount(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}