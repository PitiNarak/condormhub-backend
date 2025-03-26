package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/storage"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type ReceiptService struct {
	receiptRepo        ports.ReceiptRepository
	userRepo           ports.UserRepository
	transactionRepo    ports.TransactionRepository
	orderRepo          ports.OrderRepository
	leasingHistoryRepo ports.LeasingHistoryRepository
	dormRepo           ports.DormRepository
	storage            *storage.Storage
}

func NewReceiptService(receiptRepo ports.ReceiptRepository, userRepo ports.UserRepository, transactionRepo ports.TransactionRepository, orderRepo ports.OrderRepository, leasingHistoryRepo ports.LeasingHistoryRepository, dormRepo ports.DormRepository, storage *storage.Storage) ports.ReceiptService {
	return &ReceiptService{
		receiptRepo:        receiptRepo,
		userRepo:           userRepo,
		transactionRepo:    transactionRepo,
		orderRepo:          orderRepo,
		leasingHistoryRepo: leasingHistoryRepo,
		dormRepo:           dormRepo,
		storage:            storage,
	}
}

func (r *ReceiptService) Create(c context.Context, ownerID uuid.UUID, transaction domain.Transaction) error {

	if err := r.validateTransaction(transaction); err != nil {
		return err
	}

	buff, buffErr := r.generatePDF(ownerID, transaction)
	if buffErr != nil {
		return buffErr
	}

	fileKey, saveErr := r.saveFile(c, buff, transaction.ID)
	if saveErr != nil {
		return saveErr
	}

	receipt := &domain.Receipt{
		OwnerID:       ownerID,
		TransactionID: transaction.ID,
		FileKey:       fileKey,
	}

	if err := r.receiptRepo.Create(receipt); err != nil {
		return err
	}

	return nil

}

func (r *ReceiptService) GetUrl(c context.Context, receipt domain.Receipt) (string, error) {
	url, err := r.storage.GetSignedUrl(c, receipt.FileKey, time.Minute*60)
	if err != nil {
		return "", apperror.InternalServerError(err, "Fail to upload file to storage")
	}
	return url, nil
}

func (r *ReceiptService) validateTransaction(transaction domain.Transaction) error {
	if transaction.SessionStatus != domain.StatusComplete {
		return apperror.BadRequestError(errors.New("invalid transaction"), "Receive unpaid transcation")
	}
	return nil
}

func (r *ReceiptService) validateOwner(ownerID uuid.UUID, leasingHistory domain.LeasingHistory) error {
	if ownerID != leasingHistory.LesseeID {
		return apperror.BadRequestError(errors.New("user mismatch"), "user is not an owner of this transcation")
	}
	return nil
}

func (r *ReceiptService) generatePDF(ownerID uuid.UUID, transaction domain.Transaction) (*bytes.Buffer, error) {
	order, err := r.orderRepo.GetByID(transaction.OrderID)
	if err != nil {
		return nil, err
	}
	history, err := r.leasingHistoryRepo.GetByID(order.LeasingHistoryID)
	if err != nil {
		return nil, err
	}

	ownerErr := r.validateOwner(ownerID, *history)
	if ownerErr != nil {
		return nil, ownerErr
	}

	dorm, err := r.dormRepo.GetByID(history.DormID)
	if err != nil {
		return nil, err
	}
	lessor, err := r.userRepo.GetUserByID(dorm.OwnerID)
	if err != nil {
		return nil, err
	}
	lessee, err := r.userRepo.GetUserByID(ownerID)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()

	// Header
	pdf.SetXY(10, 10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Payment Receipt")
	pdf.Ln(12)

	// Transaction Details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", transaction.ID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Lessee: %s %s", lessee.Firstname, lessee.Lastname))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Lessor: %s %s", lessor.Firstname, lessor.Lastname))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Dorm: %s", dorm.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Amount Paid: %.2f", float64(transaction.Price)))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Issued At: %s", time.Now()))
	pdf.Ln(12)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(40, 10, "Thank you for your payment!")

	// Save to a buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, apperror.InternalServerError(err, "Fail to generate PDF file")
	}

	return &buf, nil
}

func (r *ReceiptService) saveFile(c context.Context, pdfBuffer *bytes.Buffer, transactionID string) (string, error) {
	filename := fmt.Sprintf("receipt-%s.pdf", transactionID)
	filename = strings.ReplaceAll(filename, " ", "-")

	uuid := uuid.New().String()
	fileKey := fmt.Sprintf("receipts/%s-%s", uuid, filename)
	contentType := "application/pdf"
	err := r.storage.UploadFile(c, fileKey, contentType, pdfBuffer, storage.PrivateBucket)
	if err != nil {
		return "", err
	}

	return fileKey, nil

}

func (r *ReceiptService) GetByUserID(userID uuid.UUID, limit, page int) ([]domain.Receipt, int, int, error) {
	user, userErr := r.userRepo.GetUserByID(userID)
	if userErr != nil {
		return nil, 0, 0, userErr
	}
	if user == nil || user.Role == "" {
		return nil, 0, 0, apperror.BadRequestError(errors.New("invalid user"), "user not found or role is missing")
	}
	if user.Role != domain.LesseeRole {
		return nil, 0, 0, apperror.BadRequestError(errors.New("invalid user"), "role mismatch")
	}
	return r.receiptRepo.GetByUserID(userID, limit, page)
}
