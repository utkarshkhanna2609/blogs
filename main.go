package main

import (
	"log"
	"net/http"
	"os"

	"github.com/utkarshkhanna2609/blog-api/db"
	"github.com/utkarshkhanna2609/blog-api/handlers"
)

func main() {
	// Initialize the database
	db.InitDB()
	defer db.CloseDB()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Could not read portfolio page", http.StatusInternalServerError)
			log.Printf("Error reading file: %v\n", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(file)
	})

	// Routes
	http.HandleFunc("/api/blogs", handlers.BlogsHandler)
	http.HandleFunc("/api/users", handlers.UsersHandler)

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
