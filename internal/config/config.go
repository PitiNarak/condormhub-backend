package config

import (
	"log"

	"github.com/PitiNarak/condormhub-backend/internal/databases"
	"github.com/PitiNarak/condormhub-backend/internal/server"
	"github.com/PitiNarak/condormhub-backend/internal/storage"
	"github.com/PitiNarak/condormhub-backend/pkg/email"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/PitiNarak/condormhub-backend/pkg/redis"
	"github.com/PitiNarak/condormhub-backend/pkg/stripe"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SMTP         email.SMTPConfig `envPrefix:"SMTP_"`
	JWT          jwt.JWTConfig    `envPrefix:"JWT_"`
	Server       server.Config    `envPrefix:"SERVER_"`
	Database     databases.Config `envPrefix:"DB_"`
	Storage      storage.Config   `envPrefix:"STORAGE_"`
	StripeConfig stripe.Config    `envPrefix:"STRIPE_"`
	Redis        redis.Config     `envPrefix:"REDIS_"`
}

// Load configs from .env file
func Load() *AppConfig {
	config := &AppConfig{}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to parse env vars: %s", err)
	}

	return config
}
