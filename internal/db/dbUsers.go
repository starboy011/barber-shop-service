package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes and returns a database connection
func InitUsersDB() (*sql.DB, error) {
	// Retrieve connection string components from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbName, dbSSLMode)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return db, nil
}
