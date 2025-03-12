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
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"
	"github.com/google/uuid"
)

type OwnershipProofService struct {
	ownershipProofRepo ports.OwnershipProofRepository
	userRepo           ports.UserRepository
	storage            *storage.Storage
}

func NewOwnershipProofService(ownershipProofRepo ports.OwnershipProofRepository, userRepo ports.UserRepository, storage *storage.Storage) ports.OwnershipProofService {
	return &OwnershipProofService{
		ownershipProofRepo: ownershipProofRepo,
		userRepo:           userRepo,
		storage:            storage,
	}
}

func (o *OwnershipProofService) UploadFile(ctx context.Context, dormID uuid.UUID, filename string, contentType string, fileData io.Reader, userID uuid.UUID, isAdmin bool) (string, error) {

	filename = strings.ReplaceAll(filename, " ", "-")
	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("ownership-proof/%s-%s", uuid, filename)

	if err := o.storage.UploadFile(ctx, fileKey, contentType, fileData, storage.PublicBucket); err != nil {
		return "", apperror.InternalServerError(err, "error uploading file")
	}
	url := o.storage.GetPublicUrl(fileKey)

	if _, err := o.ownershipProofRepo.GetByDormID(dormID); err != nil {
		ownershipProof := domain.OwnershipProof{DormID: dormID, FileKey: fileKey}
		if err := o.ownershipProofRepo.Create(&ownershipProof); err != nil {
			return "", err
		}
		return url, nil
	}

	if err := o.ownershipProofRepo.UpdateDocument(dormID, fileKey); err != nil {
		return "", err
	}
	return url, nil

}

func (o *OwnershipProofService) Delete(dormID uuid.UUID) error {
	if err := o.ownershipProofRepo.Delete(dormID); err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) GetByDormID(dormID uuid.UUID) (*domain.OwnershipProof, error) {
	ownershipProof, err := o.ownershipProofRepo.GetByDormID(dormID)
	if err != nil {
		return nil, err
	}
	return ownershipProof, nil
}

func (o *OwnershipProofService) UpdateStatus(dormID uuid.UUID, adminID uuid.UUID, status domain.OwnershipProofStatus) error {
	admin, err := o.userRepo.GetUserByID(adminID)
	if err != nil {
		return err
	}

	if admin == nil || admin.Role == nil {
		return apperror.BadRequestError(errors.New("invalid admin"), "Admin not found or role is missing")
	}

	if *admin.Role != domain.AdminRole {
		return apperror.BadRequestError(errors.New("role mismatch"), "You are not an admin")
	}

	updateStatusRequestBody := new(dto.UpdateOwnerShipProofStatusRequestBody)
	updateStatusRequestBody.Status = status
	updateStatusRequestBody.AdminID = adminID

	if err := o.ownershipProofRepo.UpdateStatus(dormID, updateStatusRequestBody); err != nil {
		return err
	}
	return nil
}

func (o *OwnershipProofService) ConvertToDTO(ownershipProof domain.OwnershipProof, url string, expires time.Time) dto.OwnershipProofResponseBody {
	ownershipProofResponseBody := dto.OwnershipProofResponseBody{
		Url:     url,
		Expires: expires,
		DormID:  ownershipProof.DormID,
		AdminID: ownershipProof.AdminID,
		Status:  ownershipProof.Status,
	}

	return ownershipProofResponseBody
}
