package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/PitiNarak/condormhub-backend/internal/config"
	"github.com/PitiNarak/condormhub-backend/internal/databases"
	"github.com/PitiNarak/condormhub-backend/internal/server"
	"github.com/PitiNarak/condormhub-backend/pkg/redis"
	"github.com/gofiber/fiber/v2/log"
	// _ "github.com/PitiNarak/condormhub-backend/docs"
)

// @title Condormhub API
// @version 1.0
// @description This is the API for the Condormhub project.

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	config := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db, err := databases.New(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	redis, err := redis.New(config.Redis)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	s := server.NewServer(config.Server, config.SMTP, config.JWT, config.Storage, config.StripeConfig, redis, db)
	s.Start(ctx, stop)
}
