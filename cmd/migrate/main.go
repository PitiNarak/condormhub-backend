package main

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/config"
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/databases"

	"github.com/gofiber/fiber/v2/log"
)

func main() {
	config := config.Load()

	db, err := databases.NewDatabaseConnection(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.AutoMigrate(
		&domain.SampleLog{},
		&domain.User{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed")
}
