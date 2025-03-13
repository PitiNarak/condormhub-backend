package test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockDB creates a new mock database for testing
func NewMockDB(t *testing.T) (*database.Database, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	// Wrap the gorm DB with your custom database
	dbWrapper := &database.Database{
		DB: db,
	}

	return dbWrapper, mock
}

// SetupMockDormRepo sets up expected calls for the dorm repository
func SetupMockDormRepo(mock sqlmock.Sqlmock) {
	// You can add common mock expectations here
	// For example, setting up UUID generation functions
	// mock.ExpectExec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).
	// 	WillReturnResult(sqlmock.NewResult(0, 0))
    mock.MatchExpectationsInOrder(false)
    
    // Allow any extension creation statements
    mock.ExpectExec(`CREATE EXTENSION IF NOT EXISTS`).
        WillReturnResult(sqlmock.NewResult(0, 0))
}

// CloseMockDB properly closes the mock database
func CloseMockDB(sqlDB *sql.DB) {
	sqlDB.Close()
}