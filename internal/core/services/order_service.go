package services

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	stripePkg "github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type OrderService struct {
	orderRepo    ports.OrderRepository
	stripeConfig *stripePkg.Config
}

func NewOrderService(orderRepo ports.OrderRepository, config *stripePkg.Config) ports.OrderService {
	return &OrderService{orderRepo: orderRepo, stripeConfig: config}
}

func (s *OrderService) CreateOrder(orderType domain.OrderType, dormitoryID uuid.UUID, lessorID uuid.UUID, lesseeID uuid.UUID) (*domain.Order, *error_handler.ErrorHandler) {
	stripe.Key = s.stripeConfig.StripeSecretKey

	// implement logic to get dorm name and price from dormitoryID
	dormName := "TEMPORARY_DORM_NAME"
	var price int64 = 4000

	stripeParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(dormName),
					},
					UnitAmount: stripe.Int64(price * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String("sern.dev@gmail.com"),
		SuccessURL:    stripe.String("https://example.com/success"),
		CancelURL:     stripe.String("https://example.com/cancel"),
	}

	stripeSession, stripeErr := session.New(stripeParams)
	if stripeErr != nil {
		return nil, error_handler.InternalServerError(stripeErr, "Failed to create Stripe session")
	}

	order := domain.Order{
		Type:        orderType,
		Price:       price,
		SessionID:   stripeSession.ID,
		DormitoryID: dormitoryID,
		LessorID:    lessorID,
		LesseeID:    lesseeID,
		CheckoutUrl: stripeSession.URL,
	}

	err := s.orderRepo.Create(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
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
