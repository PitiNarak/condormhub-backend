package database

import (
	"fmt"
	"log"
	"math"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"NAME"`
	SSLMode  string `env:"SSLMODE"`
}

type Database struct {
	*gorm.DB
}

func New(config Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database")
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) Paginate(value any, tx *gorm.DB, limit int, page int, order string) (int, int, error) {
	var totalRows int64

	offset := (page - 1) * limit
	if err := db.Model(value).Count(&totalRows).Offset(offset).Limit(limit).Order(order).Find(value).Error; err != nil {
		return 0, 0, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return totalPages, int(totalRows), nil

}
