package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utkarshkhanna2609/blog-api/db"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request recieved, %v",r)
	rows, err := db.DB.Query("SELECT id, email, name FROM Users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var email, name string

		if err := rows.Scan(&id, &email, &name); err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}

		users = append(users, map[string]interface{}{
			"id":    id,
			"email": email,
			"name":  name,
		})
	}

	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if a user with the given email already exists
	var existingUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	err := db.DB.QueryRow("SELECT id, email, name FROM users WHERE email = $1", user.Email).Scan(&existingUser.ID, &existingUser.Email, &existingUser.Name)
	if err == nil {
		// User exists, return their data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingUser)
		return
	}

	// If no user exists, insert the new user
	_, err = db.DB.Exec(
		"INSERT INTO users (email, name) VALUES ($1, $2)",
		user.Email, user.Name,
	)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error executing query: %v\n", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return success response for the new user
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}


