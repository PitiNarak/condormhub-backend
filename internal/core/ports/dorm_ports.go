package ports

import (
	"context"
	"io"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DormRepository interface {
	Create(dorm *domain.Dorm) error
	GetAll(limit int, page int, search string, min_price int, max_price int, district string, subdistrict string, province string, zipcode string) ([]domain.Dorm, int, int, error)
	GetByID(id uuid.UUID) (*domain.Dorm, error)
	Update(id uuid.UUID, dorm dto.DormUpdateRequestBody) error
	Delete(id uuid.UUID) error
	SaveDormImage(dormImage *domain.DormImage) error
	GetByOwnerID(ownerID uuid.UUID, limit int, page int) ([]domain.Dorm, int, int, error)
	DeleteImageByKey(imageKey string) error
	GetImageByKey(imageKey string) (*domain.DormImage, error)
}

type DormService interface {
	Create(userRole domain.Role, dorm *domain.Dorm) error
	GetAll(limit int, page int, search string, min_price int, max_price int, district string, subdistrict string, province string, zipcode string) ([]dto.DormResponseBody, int, int, error)
	GetByID(id uuid.UUID) (*dto.DormResponseBody, error)
	Update(userID uuid.UUID, isAdmin bool, dormID uuid.UUID, dorm *dto.DormUpdateRequestBody) (*dto.DormResponseBody, error)
	Delete(ctx context.Context, userID uuid.UUID, isAdmin bool, dormID uuid.UUID) error
	UploadDormImage(ctx context.Context, dormID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error)
	GetByOwnerID(ownerID uuid.UUID, limit int, page int) ([]dto.DormResponseBody, int, int, error)
	DeleteImageByURL(ctx context.Context, imageURL string, userID uuid.UUID, isAdmin bool) error
}

type DormHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadDormImage(c *fiber.Ctx) error
	GetByOwnerID(c *fiber.Ctx) error
	DeleteDormImageByURL(c *fiber.Ctx) error
}
