package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func server() {
	mux := http.NewServeMux()

	// Method 1: Using http.FileServer to serve an entire directory
	// This will serve files from the "static" directory
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Method 2: Using http.FileServer with embedded files
	// Uncomment and modify the next lines if you want to use Go 1.16+ embed
	/*
	   //go:embed static
	   var staticFiles embed.FS
	   fsys := http.FS(staticFiles)
	   mux.Handle("/embedded/", http.StripPrefix("/embedded/", http.FileServer(fsys)))
	*/

	// Method 3: Custom file handler with more control
	mux.HandleFunc("GET /files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		// Validate filename to prevent directory traversal
		if filepath.Ext(filename) == "" {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		// Set content type based on file extension
		switch filepath.Ext(filename) {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		default:
			w.Header().Set("Content-Type", "text/plain")
		}

		// Serve the file from the static directory
		http.ServeFile(w, r, filepath.Join("static", filename))
	})

	// Method 4: Serving a specific file directly
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	// Example middleware for logging static file requests
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Static file requested: %s", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware to the entire mux
	wrappedMux := loggingMiddleware(mux)

	// Create the static directory if it doesn't exist
	createDirStructure()

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}

func createDirStructure() {
	// Create example directory structure and files
	fmt.Println("Creating example directory structure...")
	fmt.Println(`
    static/
    ├── css/
    │   └── styles.css
    ├── js/
    │   └── script.js
    ├── images/
    │   ├── photo.jpg
    │   └── icon.png
    └── favicon.ico
    
    Access files at:
    - http://localhost:8080/static/css/styles.css
    - http://localhost:8080/static/js/script.js
    - http://localhost:8080/static/images/photo.jpg
    - http://localhost:8080/files/photo.jpg
    - http://localhost:8080/favicon.ico
    `)
}
