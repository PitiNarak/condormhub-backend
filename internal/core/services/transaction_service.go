package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	stripePkg "github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type TransactionService struct {
	tsxRepo   ports.TransactionRepository
	orderRepo ports.OrderRepository
	stripe    *stripePkg.Stripe
}

func NewTransactionService(tsxRepo ports.TransactionRepository, orderRepo ports.OrderRepository, stripe *stripePkg.Stripe) ports.TransactionService {
	return &TransactionService{
		tsxRepo:   tsxRepo,
		orderRepo: orderRepo,
		stripe:    stripe,
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

func (s *TransactionService) UpdateTransactionStatus(event stripe.Event) error {
	var tsx domain.Transaction
	tsx.ID = event.Data.Object["id"].(string)
	switch event.Type {
	case "checkout.session.expired":
		tsx.SessionStatus = domain.StatusExpired
	case "checkout.session.completed":
		tsx.SessionStatus = domain.StatusComplete
	default:
		return apperror.BadRequestError(fmt.Errorf("event type %s is not supported", event.Type), "Failed to update order status")
	}

	if err := s.tsxRepo.Update(&tsx); err != nil {
		return err
	}

	tsx, err := s.tsxRepo.GetByID(tsx.ID)
	if err != nil {
		return err
	}

	if err := s.orderRepo.Update(&domain.Order{
		ID:                tsx.OrderID,
		PaidTransactionID: tsx.ID,
	}); err != nil {
		return err
	}

	return nil
}
