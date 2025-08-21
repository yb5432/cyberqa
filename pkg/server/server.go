// Package server sets up the HTTP server for the Cyber Q&A application.
package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"openai-api/pkg/database"
	"openai-api/pkg/handlers"

	"github.com/gorilla/mux"
)

// Start initializes and starts the HTTP server.
func Start() {
	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/submit-user-a", handlers.SubmitUserA).Methods("POST")
	api.HandleFunc("/submit-user-b", handlers.SubmitUserB).Methods("POST")
	api.HandleFunc("/results/{token}", handlers.GetResults).Methods("GET")
	api.HandleFunc("/questions/upload", handlers.UploadQuestions).Methods("POST")
	api.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")

	// Get the path to the dist directory from environment variable or use default
	distPath := os.Getenv("DIST_PATH")
	if distPath == "" {
		// Default to frontend/cyberqa/dist assuming the frontend has been built
		distPath = "frontend/cyberqa/dist"
	}

	// Check if dist directory exists
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		log.Printf("Warning: Dist directory '%s' does not exist. Frontend files will not be served.", distPath)
	} else {
		// Serve static files
		fs := http.FileServer(http.Dir(distPath))
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the requested file exists
			path := filepath.Join(distPath, r.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// If the file doesn't exist, serve index.html for Vue Router
				http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			} else {
				// Otherwise, serve the static file
				fs.ServeHTTP(w, r)
			}
		})
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// Connect to database
	database.Connect()

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at http://localhost:%s/api/", port)
	if _, err := os.Stat(distPath); err == nil {
		log.Printf("Frontend available at http://localhost:%s/", port)
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
