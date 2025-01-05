package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/utkarshkhanna2609/blog-api/db"
)

func BlogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBlogs(w, r)
	case http.MethodPost:
		createBlog(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getBlogs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request recieved, %v",r)
	rows, err := db.DB.Query("SELECT id, userid, title, images, description FROM blogs")
	if err != nil {
		http.Error(w, "Failed to fetch blogs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var blogs []map[string]interface{}
	for rows.Next() {
		var id int
		var userid int
		var title, description string
		var images []byte

		if err := rows.Scan(&id, &userid, &title, &images, &description); err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}

		blogs = append(blogs, map[string]interface{}{
			"id":          id,
			"userid":      userid,
			"title":       title,
			"images":      string(images),
			"description": description,
		})
	}

	json.NewEncoder(w).Encode(blogs)
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	var blog struct {
		UserID      int      `json:"userid"`
		Title       string   `json:"title"`
		Images      []string `json:"images"`
		Description string   `json:"description"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Convert the images array to the PostgreSQL array format
	imagesArray := "{" + strings.Join(blog.Images, ",") + "}"

	// Execute the query with the formatted images array
	_, err := db.DB.Exec(
		"INSERT INTO blogs (userid, title, images, description) VALUES ($1, $2, $3, $4)",
		blog.UserID, blog.Title, imagesArray, blog.Description,
	)
	if err != nil {
		fmt.Println("Error executing query:", err)
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Blog created successfully"))
}
