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

func (s *Stripe) CreateSession(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	stripe.Key = s.config.StripeSecretKey
	return session.New(params)
}
