package main

import (
	"fmt"
	"log"

	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/databases"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	db, err := databases.NewDatabaseConnection()
	if err != nil {
		log.Fatal("Error connecting to database")
		panic(err)
	}

	db.AutoMigrate(&domain.SampleLog{})
	fmt.Println("Migration completed")
}
