package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	stripePkg "github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type OrderService struct {
	orderRepo ports.OrderRepository
	stripe    *stripePkg.Stripe
}

func NewOrderService(orderRepo ports.OrderRepository, stripe *stripePkg.Stripe) ports.OrderService {
	return &OrderService{orderRepo: orderRepo, stripe: stripe}
}

func (s *OrderService) CreateOrder(orderType domain.OrderType, dormitoryID uuid.UUID, lessorID uuid.UUID, lessee *domain.User) (*domain.Order, *string, *error_handler.ErrorHandler) {

	// implement logic to get dorm name and price from dormitoryID
	dormName := "TEMPORARY_DORM_NAME"
	var price int64 = 4000
	session, err := s.stripe.CreateOneTimePaymentSession(dormName, price, lessee.Email)
	if err != nil {
		return nil, nil, error_handler.InternalServerError(err, "Failed to create Stripe session")
	}

	order := domain.Order{
		Type:        orderType,
		Price:       price,
		SessionID:   session.ID,
		DormitoryID: dormitoryID,
		LessorID:    lessorID,
		LesseeID:    lessee.ID,
	}

	repoError := s.orderRepo.Create(&order)
	if repoError != nil {
		return nil, nil, repoError
	}

	return &order, &session.URL, nil
}

func (s *OrderService) UpdateOrderStatus(event stripe.Event) *error_handler.ErrorHandler {
	var order domain.Order
	order.SessionID = event.Data.Object["id"].(string)
	switch event.Type {
	case "checkout.session.expired":
		order.SessionStatus = stripe.CheckoutSessionStatusExpired
	case "checkout.session.completed":
		order.SessionStatus = stripe.CheckoutSessionStatusComplete
	case "checkout.session.async_payment_succeeded":
		order.SessionStatus = stripe.CheckoutSessionStatusComplete
	case "checkout.session.async_payment_failed":
		order.SessionStatus = stripe.CheckoutSessionStatusExpired
	default:
		return error_handler.BadRequestError(fmt.Errorf("event type %s is not supported", event.Type), "Failed to update order status")
	}

	err := s.orderRepo.Update(&order)
	if err != nil {
		return err
	}
	return nil
}
