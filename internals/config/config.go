package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     int
	Email    string
	Password string
}

// Load email config from .env file
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("Invalid SMTP_PORT value in .env file")
	}

	return &Config{Host: os.Getenv("SMTP_HOST"), Port: port, Email: os.Getenv("SMTP_EMAIL"), Password: os.Getenv("SMTP_PASSWORD")}
}
