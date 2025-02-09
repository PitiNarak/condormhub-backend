package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/PitiNarak/condormhub-backend/internals/databases"
	"github.com/PitiNarak/condormhub-backend/internals/server"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	config := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db, err := databases.NewDatabaseConnection(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	s := server.NewServer(config.Server, config.SMTP, config.JWT, db)
	s.Start(ctx, stop, config.JWT)
}
