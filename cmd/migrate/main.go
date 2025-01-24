package main

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/databases"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warnf("Warning: No .env file found")
	}

	db, err := databases.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	db.AutoMigrate(&domain.SampleLog{})
	fmt.Println("Migration completed")
}
