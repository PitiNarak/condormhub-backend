package stripe

import (
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type Config struct {
	StripePublishableKey string `env:"PUBLISHABLE_KEY"`
	StripeSecretKey      string `env:"SECRET_KEY"`
	StripeSignatureKey   string `env:"SIGNATURE_KEY"`
	StripeSuccessURL     string `env:"SUCCESS_URL"`
	StripeCancelURL      string `env:"CANCEL_URL"`
}

type Stripe struct {
	config Config
}

func New(config Config) *Stripe {
	return &Stripe{config: config}
}

func (s *Stripe) createSession(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	stripe.Key = s.config.StripeSecretKey
	return session.New(params)
}

func (s *Stripe) CreateOneTimePaymentSession(productName string, price int64, customerEmail string) (*stripe.CheckoutSession, error) {
	stripeParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(price * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(customerEmail),
		SuccessURL:    stripe.String(s.config.StripeSuccessURL),
		CancelURL:     stripe.String(s.config.StripeCancelURL),
	}

	return s.createSession(stripeParams)
}

func (s *Stripe) CreateSubscriptionSession(productName string, price int64, customerEmail string) (*stripe.CheckoutSession, error) {
	stripeParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("thb"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(price * 100),
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(customerEmail),
		SuccessURL:    stripe.String(s.config.StripeSuccessURL),
		CancelURL:     stripe.String(s.config.StripeCancelURL),
	}

	return s.createSession(stripeParams)
}
