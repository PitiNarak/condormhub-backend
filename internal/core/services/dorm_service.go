package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"
	"github.com/google/uuid"
)

type DormService struct {
	dormRepo ports.DormRepository
	storage  *storage.Storage
}

func NewDormService(repo ports.DormRepository, storage *storage.Storage) ports.DormService {
	return &DormService{dormRepo: repo, storage: storage}
}

func checkPermission(ownerID uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	if ownerID != userID && !isAdmin {
		return errors.New("unauthorized action")
	}
	return nil
}

func (s *DormService) getImageUrl(dormImage []domain.DormImage) []string {
	urls := make([]string, len(dormImage))
	for i, v := range dormImage {
		urls[i] = s.storage.GetPublicUrl(v.ImageKey)
	}
	return urls
}

func (s *DormService) Create(userRole domain.Role, dorm *domain.Dorm) error {
	if userRole != domain.AdminRole && userRole != domain.LessorRole {
		return apperror.ForbiddenError(errors.New("unauthorized action"), "You do not have permission to create a dorm")
	}
	return s.dormRepo.Create(dorm)
}

func (s *DormService) GetAll(
	limit int, page int,
	search string,
	minPrice int, maxPrice int,
	district string,
	subdistrict string,
	province string,
	zipcode string,
) ([]dto.DormResponseBody, int, int, error) {
	dorms, totalPages, totalRows, err := s.dormRepo.GetAll(limit, page, search, minPrice, maxPrice, district, subdistrict, province, zipcode)
	if err != nil {
		return nil, totalPages, totalRows, err
	}
	resData := make([]dto.DormResponseBody, len(dorms))
	for i, v := range dorms {
		resData[i] = v.ToDTO()
		resData[i].Images = s.getImageUrl(v.Images)
	}
	return resData, totalPages, totalRows, nil
}

func (s *DormService) GetByID(id uuid.UUID) (*dto.DormResponseBody, error) {
	dorm, err := s.dormRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	resData := dorm.ToDTO()
	resData.Images = s.getImageUrl(dorm.Images)
	return &resData, nil
}

func (s *DormService) Update(userID uuid.UUID, isAdmin bool, dormID uuid.UUID, updateData *dto.DormUpdateRequestBody) (*dto.DormResponseBody, error) {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return nil, err
	}

	if err = checkPermission(dorm.OwnerID, userID, isAdmin); err != nil {
		return nil, apperror.ForbiddenError(err, "You do not have permission to update this dorm")
	}

	if err := s.dormRepo.Update(dormID, *updateData); err != nil {
		return nil, err
	}

	return s.GetByID(dormID)
}

func (s *DormService) Delete(ctx context.Context, userID uuid.UUID, isAdmin bool, dormID uuid.UUID) error {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return err
	}

	if err := checkPermission(dorm.OwnerID, userID, isAdmin); err != nil {
		return apperror.ForbiddenError(err, "You do not have permission to delete this dorm")
	}

	if len(dorm.Images) > 0 {
		for _, image := range dorm.Images {
			err = s.storage.DeleteFile(ctx, image.ImageKey, storage.PublicBucket)
			if err != nil {
				return apperror.InternalServerError(err, "Failed to delete images")
			}
		}
	}

	return s.dormRepo.Delete(dormID)
}

func (s *DormService) UploadDormImage(ctx context.Context, dormID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error) {
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return "", err
	}

	if err = checkPermission(dorm.OwnerID, userID, isAdmin); err != nil {
		return "", apperror.ForbiddenError(err, "You do not have permission to upload image to this dorm")
	}

	filename = strings.ReplaceAll(filename, " ", "-")
	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("dorms/%s-%s", uuid, filename)

	err = s.storage.UploadFile(ctx, fileKey, contentType, fileData, storage.PublicBucket)
	if err != nil {
		return "", apperror.InternalServerError(err, "error uploading file")
	}

	dormImage := &domain.DormImage{DormID: dormID, ImageKey: fileKey}
	if err = s.dormRepo.SaveDormImage(dormImage); err != nil {
		return "", err
	}

	url := s.storage.GetPublicUrl(fileKey)

	return url, nil
}

func (s *DormService) GetByOwnerID(ownerID uuid.UUID, limit int, page int) ([]dto.DormResponseBody, int, int, error) {
	dorms, totalPages, totalRows, err := s.dormRepo.GetByOwnerID(ownerID, limit, page)
	if err != nil {
		return nil, totalPages, totalRows, err
	}
	resData := make([]dto.DormResponseBody, len(dorms))
	for i, v := range dorms {
		resData[i] = v.ToDTO()
		resData[i].Images = s.getImageUrl(v.Images)
	}
	return resData, totalPages, totalRows, nil
}
