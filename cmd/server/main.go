package main

import (
	"github.com/PitiNarak/condormhub-backend/internals/databases"
	"github.com/PitiNarak/condormhub-backend/internals/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db, err := databases.NewDatabaseConnection()
	if err != nil {
		panic("Error connecting to database")
	}

	s := server.NewServer(db)
	s.Start("3000")
}
