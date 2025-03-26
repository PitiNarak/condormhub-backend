package services

import (
	"context"
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	stripePkg "github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type TransactionService struct {
	tsxRepo            ports.TransactionRepository
	orderRepo          ports.OrderRepository
	leasingHistoryRepo ports.LeasingHistoryRepository
	stripe             *stripePkg.Stripe
	receiptService     ports.ReceiptService
}

func NewTransactionService(tsxRepo ports.TransactionRepository, orderRepo ports.OrderRepository, stripe *stripePkg.Stripe, leasingHistoryRepo ports.LeasingHistoryRepository, receiptService ports.ReceiptService) ports.TransactionService {
	return &TransactionService{
		tsxRepo:            tsxRepo,
		orderRepo:          orderRepo,
		leasingHistoryRepo: leasingHistoryRepo,
		receiptService:     receiptService,
		stripe:             stripe,
	}
}

func (s *TransactionService) CreateTransaction(orderID uuid.UUID) (*domain.Transaction, *string, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, nil, err
	}
	if order.PaidTransactionID != "" {
		return nil, nil, apperror.BadRequestError(fmt.Errorf("order %s is already paid", orderID), "order is already paid")
	}

	productName := order.LeasingHistory.Dorm.Name
	price := order.Price
	customerEmail := order.LeasingHistory.Lessee.Email

	session, sErr := s.stripe.CreateOneTimePaymentSession(productName, int64(price), customerEmail)
	if sErr != nil {
		return nil, nil, apperror.InternalServerError(sErr, "Failed to create payment session")
	}

	tsx := domain.Transaction{
		ID:      session.ID,
		Price:   int64(price),
		OrderID: orderID,
	}
	err = s.tsxRepo.Create(&tsx)
	if err != nil {
		return nil, nil, err
	}

	return &tsx, &session.URL, nil
}

func (s *TransactionService) UpdateTransactionStatus(c context.Context, event stripe.Event) (dto.ReceiptResponseBody, error) {
	var tsx domain.Transaction
	tsx.ID = event.Data.Object["id"].(string)
	switch event.Type {
	case "checkout.session.expired":
		tsx.SessionStatus = domain.StatusExpired
	case "checkout.session.completed":
		tsx.SessionStatus = domain.StatusComplete
	default:
		return dto.ReceiptResponseBody{}, apperror.BadRequestError(fmt.Errorf("event type %s is not supported", event.Type), "Failed to update order status")
	}

	if err := s.tsxRepo.Update(&tsx); err != nil {
		return dto.ReceiptResponseBody{}, err
	}

	tsx, err := s.tsxRepo.GetByID(tsx.ID)
	if err != nil {
		return dto.ReceiptResponseBody{}, err
	}

	if err := s.orderRepo.Update(&domain.Order{
		ID:                tsx.OrderID,
		PaidTransactionID: tsx.ID,
	}); err != nil {
		return dto.ReceiptResponseBody{}, err
	}

	if tsx.SessionStatus == domain.StatusComplete {
		order, err := s.orderRepo.GetByID(tsx.OrderID)
		if err != nil {
			return dto.ReceiptResponseBody{}, err
		}
		history, err := s.leasingHistoryRepo.GetByID(order.LeasingHistoryID)
		if err != nil {
			return dto.ReceiptResponseBody{}, err
		}
		ownerID := history.LesseeID
		receipt, url, receiptErr := s.receiptService.Create(c, ownerID, tsx)
		if receiptErr != nil {
			return dto.ReceiptResponseBody{}, receiptErr
		}

		return receipt.ToDTO(url), nil
	}

	return dto.ReceiptResponseBody{}, nil
}
