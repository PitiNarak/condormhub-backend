package ports

import (
	"context"
	"io"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeasingHistoryRepository interface {
	Create(LeasingHistory *domain.LeasingHistory) error
	Update(LeasingHistory *domain.LeasingHistory) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.LeasingHistory, error)
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	DeleteReview(leasingHistory *domain.LeasingHistory) error
	SaveReviewImage(reviewImage *domain.ReviewImage) error
	DeleteImageByKey(imageKey string) error
	GetImageByKey(imageKey string) (*domain.ReviewImage, error)
	GetReportedReviews(limit int, page int) ([]domain.LeasingHistory, int, int, error)
}

type LeasingHistoryService interface {
	Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error)
	CreateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error)
	UpdateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error)
	DeleteReview(user *domain.User, id uuid.UUID) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.LeasingHistory, error)
	GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error)
	SetEndTimestamp(id uuid.UUID) error
	UploadReviewImage(ctx context.Context, historyID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error)
	DeleteImageByURL(ctx context.Context, imageURL string, userID uuid.UUID, isAdmin bool) error
	GetImageUrl(reviewImage []domain.ReviewImage) []string
	GetReportedReviews(limit int, page int) ([]domain.LeasingHistory, int, int, error)
	ReportReview(id uuid.UUID) (*domain.LeasingHistory, error)
}

type LeasingHistoryHandler interface {
	Delete(c *fiber.Ctx) error
	CreateReview(c *fiber.Ctx) error
	UpdateReview(c *fiber.Ctx) error
	DeleteReview(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	GetByDormID(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	SetEndTimestamp(c *fiber.Ctx) error
	UploadReviewImage(c *fiber.Ctx) error
	DeleteReviewImageByURL(c *fiber.Ctx) error
	GetReportedReviews(c *fiber.Ctx) error
	ReportReview(c *fiber.Ctx) error
}
