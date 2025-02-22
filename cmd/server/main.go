package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/PitiNarak/condormhub-backend/docs"

	"github.com/PitiNarak/condormhub-backend/internal/config"
	"github.com/PitiNarak/condormhub-backend/internal/databases"
	"github.com/PitiNarak/condormhub-backend/internal/server"
	"github.com/gofiber/fiber/v2/log"
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

	db, err := databases.NewDatabaseConnection(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	s := server.NewServer(config.Server, config.SMTP, config.JWT, config.Storage, db)
	s.Start(ctx, stop)
}
