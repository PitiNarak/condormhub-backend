package main

import (
	"github.com/PitiNarak/condormhub-backend/internals/databases"
	"github.com/PitiNarak/condormhub-backend/internals/server"
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

	s := server.NewServer(db)
	s.Start("3000")
}
