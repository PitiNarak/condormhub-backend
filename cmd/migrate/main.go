package main

import (
	"fmt"

	"github.com/PitiNarak/condormhub-backend/internal/config"
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/gofiber/fiber/v2/log"
)

func CreateEnum(db *database.Database) {
	query := `
    DO $$ BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'lifestyle_tag') THEN
            CREATE TYPE lifestyle_tag AS ENUM (
                'Active',
                'Creative',
                'Social',
                'Relaxed',

                'Football',
                'Basketball',
                'Tennis',
                'Swimming',
                'Running',
                'Cycling',
                'Badminton',
                'Yoga',
                'Gym & Fitness',

                'Music',
                'Dancing',
                'Photography',
                'Painting',
                'Gaming',
                'Reading',
                'Writing',
                'DIY & Crafting',
                'Cooking',

                'Extrovert',
                'Introvert',
                'Night Owl',
                'Early Bird',

                'Traveler',
                'Backpacker',
                'Nature Lover',
                'Camping',
                'Beach Lover',

                'Dog Lover',
                'Cat Lover',

                'Freelancer',
                'Entrepreneur',
                'Office Worker',
                'Remote Worker',
                'Student',
                'Self-Employed'
            );
        END IF;
    END $$;
    `
	db.Exec(query)
}

func main() {
	config := config.Load()

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	CreateEnum(db)
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.AutoMigrate(
		&domain.SampleLog{},
		&domain.User{},
		&domain.Transaction{},
		&domain.Dorm{},
		&domain.LeasingHistory{},
		&domain.Order{},
<<<<<<< HEAD
		&domain.LeasingRequest{},
=======
		&domain.OwnershipProof{},
>>>>>>> dev
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed")
}
