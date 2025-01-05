package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Get connection string from Railway's DATABASE_URL environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// If DATABASE_URL is not set, construct from individual environment variables
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
			os.Getenv("PGUSER"),
			os.Getenv("PGPASSWORD"),
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT"),
			os.Getenv("PGDATABASE"),
		)
	}

	// Open connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(25)  // Maximum number of open connections
	DB.SetMaxIdleConns(5)   // Maximum number of idle connections

	// Ping the database to verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v", err)
	}

	log.Println("Connected to Railway PostgreSQL successfully!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}
}