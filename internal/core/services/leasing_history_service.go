package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/google/uuid"
)

type LeasingHistoryService struct {
	historyRepo ports.LeasingHistoryRepository
	dormRepo    ports.DormRepository
	storage     *storage.Storage
}

func NewLeasingHistoryService(historyRepo ports.LeasingHistoryRepository, dormRepo ports.DormRepository, storage *storage.Storage) ports.LeasingHistoryService {
	return &LeasingHistoryService{historyRepo: historyRepo, dormRepo: dormRepo, storage: storage}
}

func (s *LeasingHistoryService) GetImageUrl(reviewImage []domain.ReviewImage) []string {
	urls := make([]string, len(reviewImage))
	for i, v := range reviewImage {
		urls[i] = s.storage.GetPublicUrl(v.ImageKey)
	}
	return urls
}

func (s *LeasingHistoryService) Create(userID uuid.UUID, dormID uuid.UUID) (*domain.LeasingHistory, error) {
	_, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	createTime := time.Now()
	dorm, err := s.dormRepo.GetByID(dormID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	leasingHistory := &domain.LeasingHistory{DormID: dormID, LesseeID: userID, Start: createTime, Price: dorm.Price}
	err = s.historyRepo.Create(leasingHistory)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	leasingHistory, err = s.historyRepo.GetByID(leasingHistory.ID)
	if err != nil {
		return &domain.LeasingHistory{}, err
	}
	return leasingHistory, nil
}

func (s *LeasingHistoryService) Delete(id uuid.UUID) error {
	err := s.historyRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingHistoryService) GetByID(id uuid.UUID) (*domain.LeasingHistory, error) {
	leasingHistory, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return leasingHistory, nil
}

func (s *LeasingHistoryService) GetByUserID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	leasingHistory, totalPage, totalRows, err := s.historyRepo.GetByUserID(id, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingHistory, totalPage, totalRows, nil
}
func (s *LeasingHistoryService) GetByDormID(id uuid.UUID, limit, page int) ([]domain.LeasingHistory, int, int, error) {
	leasingHistory, totalPage, totalRows, err := s.historyRepo.GetByDormID(id, limit, page)
	if err != nil {
		return nil, totalPage, totalRows, err
	}
	return leasingHistory, totalPage, totalRows, nil
}
func (s *LeasingHistoryService) SetEndTimestamp(id uuid.UUID) error {
	leasingHistory, err := s.historyRepo.GetByID(id)
	if err != nil {
		return err
	}
	leasingHistory.End = time.Now()
	err = s.historyRepo.Update(leasingHistory)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingHistoryService) CreateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error) {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateUserForReview(user, history, true)
	if err != nil {
		return nil, err
	}
	review := domain.Review{
		Message: Message,
		Rate:    Rate,
	}
	history.Review = review
	history.ReviewFlag = true
	err = s.historyRepo.Update(history)
	if err != nil {
		return nil, err
	}
	history, err = s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &history.Review, nil
}

func (s *LeasingHistoryService) UpdateReview(user *domain.User, id uuid.UUID, Message string, Rate int) (*domain.Review, error) {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateUserForReview(user, history, false)
	if err != nil {
		return nil, err
	}
	review := domain.Review{
		Message: Message,
		Rate:    Rate,
	}
	history.Review = review
	err = s.historyRepo.Update(history)
	if err != nil {
		return nil, err
	}
	history, err = s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &history.Review, nil
}

func (s *LeasingHistoryService) DeleteReview(user *domain.User, id uuid.UUID) error {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return err
	}
	err = utils.ValidateUserForReview(user, history, false)
	if err != nil {
		return err
	}
	history.ReviewFlag = false
	err = s.historyRepo.DeleteReview(history)
	if err != nil {
		return err
	}
	return nil
}

func (s *LeasingHistoryService) UploadReviewImage(ctx context.Context, historyID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error) {
	history, err := s.historyRepo.GetByID(historyID)
	if err != nil {
		return "", err
	}

	if err = checkPermission(history.Lessee.ID, userID, isAdmin); err != nil {
		return "", apperror.ForbiddenError(err, "You do not have permission to upload image to this dorm")
	}

	filename = strings.ReplaceAll(filename, " ", "-")
	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("reviews/%s-%s", uuid, filename)

	err = s.storage.UploadFile(ctx, fileKey, contentType, fileData, storage.PublicBucket)
	if err != nil {
		return "", apperror.InternalServerError(err, "error uploading file")
	}
	reviewImage := &domain.ReviewImage{HistoryID: historyID, ImageKey: fileKey}
	if err = s.historyRepo.SaveReviewImage(reviewImage); err != nil {
		return "", err
	}

	url := s.storage.GetPublicUrl(fileKey)

	return url, nil
}

func (s *LeasingHistoryService) DeleteImageByURL(ctx context.Context, imageURL string, userID uuid.UUID, isAdmin bool) error {
	imageKey, err := s.storage.GetFileKeyFromPublicUrl(imageURL)
	if err != nil {
		return apperror.InternalServerError(err, "Failed to parse URL")
	}

	reviewImage, err := s.historyRepo.GetImageByKey(imageKey)
	if err != nil {
		return err
	}

	history, err := s.historyRepo.GetByID(reviewImage.HistoryID)
	if err != nil {
		return err
	}

	if err := checkPermission(history.Lessee.ID, userID, isAdmin); err != nil {
		return apperror.ForbiddenError(err, "You do not have permission to delete this review image")
	}

	err = s.storage.DeleteFile(ctx, imageKey, storage.PublicBucket)
	if err != nil {
		return apperror.InternalServerError(err, "Failed to delete images")
	}

	return s.historyRepo.DeleteImageByKey(imageKey)
}

func (s *LeasingHistoryService) GetReportedReviews(limit int, page int) ([]domain.LeasingHistory, int, int, error) {
	return s.historyRepo.GetReportedReviews(limit, page)
}

func (s *LeasingHistoryService) ReportReview(id uuid.UUID) (*domain.LeasingHistory, error) {
	history, err := s.historyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !history.ReviewFlag {
		return nil, apperror.NotFoundError(errors.New("review not found"), "review not found")
	}
	if history.Review.ReportFlag {
		return nil, apperror.ConflictError(errors.New("review already reported"), "review already reported")
	}
	history.Review.ReportFlag = true
	err = s.historyRepo.Update(history)
	if err != nil {
		return nil, err
	}
	return history, nil
}
