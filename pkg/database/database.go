// Package database provides database connection functionality for the Cyber Q&A application.
package database

import (
	"log"
	"os"
	"strings"

	"openai-api/pkg/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global database instance.
var DB *gorm.DB

// Connect initializes the database connection using environment variables.
func Connect() {
	// Get database DSN from environment variable, default to "sqlite:cyberqa.db"
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "sqlite:cyberqa.db"
	}

	// Parse the DSL to determine database type and connection string
	parts := strings.SplitN(dsn, ":", 2)
	if len(parts) != 2 {
		log.Fatal("Invalid database DSL format. Expected format: type:connection_string")
	}

	dbType := strings.ToLower(parts[0])
	connectionString := parts[1]

	var err error
	switch dbType {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	case "mysql":
		DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	case "postgres", "postgresql":
		DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	default:
		log.Fatal("Unsupported database type:", dbType)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = DB.AutoMigrate(&models.UserA{}, &models.UserB{}, &models.Session{}, &models.Question{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
