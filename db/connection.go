package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Connection string
	connStr := "postgres://blog_user:blog_password@localhost:5432/blog_db?sslmode=disable"

	// Open connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// Ping the database to verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v", err)
	}

	log.Println("Connected to the database successfully!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}
}
