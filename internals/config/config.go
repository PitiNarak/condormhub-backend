package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	Email    string `env:"EMAIL,required"`
	Password string `env:"PASSWORD,required"`
}

type JWTConfig struct {
	JWTSecretKey string `env:"SECRET,required"`
	Expiration   int    `env:"EXPIRATION_HOURS,required"`
}

type AppConfig struct {
	SMTP SMTPConfig `envPrefix:"SMTP_"`
	JWT  JWTConfig  `envPrefix:"JWT_"`
}

// Load email config from .env file
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
