package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
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

func (s *TransactionService) CreateTransaction(orderID uuid.UUID) (*domain.Transaction, *string, *errorHandler.ErrorHandler) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, nil, err
	}

	productName := order.LeasingHistory.Dorm.Name
	price := order.Price
	customerEmail := order.LeasingHistory.Lessee.Email

	session, sErr := s.stripe.CreateOneTimePaymentSession(productName, int64(price), customerEmail)
	if err != sErr {
		return nil, nil, errorHandler.InternalServerError(err, "Failed to create payment session")
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

func (s *TransactionService) UpdateTransactionStatus(event stripe.Event) *errorHandler.ErrorHandler {
	var tsx domain.Transaction
	tsx.ID = event.Data.Object["id"].(string)
	switch event.Type {
	case "checkout.session.expired":
		tsx.SessionStatus = stripe.CheckoutSessionStatusExpired
	case "checkout.session.completed":
		tsx.SessionStatus = stripe.CheckoutSessionStatusComplete
	case "checkout.session.async_payment_succeeded":
		tsx.SessionStatus = stripe.CheckoutSessionStatusComplete
	case "checkout.session.async_payment_failed":
		tsx.SessionStatus = stripe.CheckoutSessionStatusExpired
	default:
		return errorHandler.BadRequestError(fmt.Errorf("event type %s is not supported", event.Type), "Failed to update order status")
	}

	err := s.tsxRepo.Update(&tsx)
	if err != nil {
		return err
	}
	return nil
}
