package ports

import (
	"context"
	"io"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetUserByID(userID uuid.UUID) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data domain.User) error
	UpdateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	DeleteAccount(userID uuid.UUID) error
	GetLessorIncome(lessorID uuid.UUID) (float64, error)
	GetPending(limit int, page int) ([]domain.User, int, int, error)
}

type UserService interface {
	ConvertToDTO(user domain.User) dto.UserResponse
	GetStudentEvidenceDTO(c context.Context, studentEvidence string) (*dto.StudentEvidenceUploadResponseBody, error)
	Create(ctx context.Context, user *domain.User) (string, string, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	FirstFillInformation(userID uuid.UUID, data dto.UserFirstFillRequestBody) (*domain.User, error)
	UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error)
	Login(context.Context, string, string) (*domain.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	VerifyUser(ctx context.Context, token string) (*domain.User, string, string, error)
	ResetPasswordCreate(context.Context, string) error
	ResetPassword(context.Context, string, string) (*domain.User, string, string, error)
	DeleteAccount(userID uuid.UUID) error
	UploadStudentEvidence(ctx context.Context, filename string, contentType string, fileData io.Reader, userID uuid.UUID) (string, error)
	GetStudentEvidenceByID(ctx context.Context, id uuid.UUID, isSelf bool, isAdmin bool) (*dto.StudentEvidenceUploadResponseBody, error)
	ResendVerificationEmailService(ctx context.Context, email string) error
	UploadProfilePicture(ctx context.Context, filename string, contentType string, fileData io.Reader, userID uuid.UUID) (string, error)
	GetLessorIncome(lessorID uuid.UUID, userRole domain.Role) (float64, error)
	UpdateUserBanStatus(id uuid.UUID, ban bool) (*domain.User, error)
	GetPending(limit int, page int) ([]domain.User, int, int, error)
	UpdateVerificationStatus(lesseeID uuid.UUID, status domain.VerificationStatus) (*domain.User, error)
}

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	FirstFillInformation(c *fiber.Ctx) error
	UpdateUserInformation(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResetPasswordCreate(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	DeleteAccount(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	UploadStudentEvidence(c *fiber.Ctx) error
	GetStudentEvidenceByID(c *fiber.Ctx) error
	ResendVerificationEmailHandler(c *fiber.Ctx) error
	UploadProfilePicture(c *fiber.Ctx) error
	GetLessorIncome(c *fiber.Ctx) error
	BanUser(c *fiber.Ctx) error
	UnbanUser(c *fiber.Ctx) error
	GetPending(c *fiber.Ctx) error
	VerifyStudentVerification(c *fiber.Ctx) error
	RejectStudentVerification(c *fiber.Ctx) error
}
