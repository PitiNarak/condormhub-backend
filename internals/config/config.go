package config

import (
	"log"

	"github.com/PitiNarak/condormhub-backend/internals/core/services"
	"github.com/PitiNarak/condormhub-backend/internals/core/utils"
	"github.com/PitiNarak/condormhub-backend/internals/databases"
	"github.com/PitiNarak/condormhub-backend/internals/server"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SMTP     services.SMTPConfig `envPrefix:"SMTP_"`
	JWT      utils.JWTConfig     `envPrefix:"JWT_"`
	Server   server.Config       `envPrefix:"SERVER_"`
	Database databases.Config    `envPrefix:"DB_"`
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
