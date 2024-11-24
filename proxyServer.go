package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	PORT               = "8989"
	WIKIPEDIA_BASE_URL = "https://wikipedia.org"
	MODIFIED_BASE_URL  = "https://m-wikipedia.org"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// Helper function to fetch resources with proper headers
func fetchResource(url string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 WikiProxy/1.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching resource: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response: %v", err)
	}

	return body, resp.Header.Get("Content-Type"), nil
}

func fetchWikipediaContent(path string) ([]byte, string, error) {
	wikipediaURL := WIKIPEDIA_BASE_URL + path
	return fetchResource(wikipediaURL)
}

func modifyWikipediaContent(content []byte, contentType string) []byte {
	// Only modify HTML content
	if !strings.Contains(contentType, "text/html") {
		return content
	}

	modified := content

	// Update resource URLs to absolute paths
	patterns := []string{
		`href="https?://[^/]*wikipedia\.org([^"]*)"`,
		`href="//[^/]*wikipedia\.org([^"]*)"`,
		`href="/wiki/([^"]*)"`,
		`href="/w/([^"]*)"`, // Add support for Wikipedia resource URLs
		`src="//([^"]*)"`,
		`src="/static/([^"]*)"`,
		`url\(['"]?//([^'"]*?)['"]?\)`,
		`@import "//([^"]*)"`,
		`@import url\("//([^"]*?)"\)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		modified = re.ReplaceAllFunc(modified, func(match []byte) []byte {
			link := string(match)
			switch {
			case strings.Contains(link, "//"):
				if strings.Contains(link, "wikipedia.org") {
					return []byte(strings.Replace(link, "wikipedia.org", "m-wikipedia.org", 1))
				}
				// Handle protocol-relative URLs
				if strings.HasPrefix(link, `src="//`) {
					return []byte(strings.Replace(link, `src="//`, `src="https://`, 1))
				}
				if strings.Contains(link, "url(//") {
					return []byte(strings.Replace(link, "url(//", "url(https://", 1))
				}
				if strings.Contains(link, "@import") {
					return []byte(strings.Replace(link, "//", "https://", 1))
				}
			case strings.HasPrefix(link, `href="/`):
				if strings.Contains(link, "/wiki/") {
					return []byte(fmt.Sprintf(`href="https://m-wikipedia.org%s"`, link[6:]))
				}
				if strings.Contains(link, "/w/") {
					return []byte(fmt.Sprintf(`href="https://wikipedia.org%s"`, link[6:]))
				}
			case strings.HasPrefix(link, `src="/`):
				return []byte(fmt.Sprintf(`src="https://wikipedia.org%s"`, link[5:]))
			}
			return match
		})
	}

	// Add custom CSS and JS before </head>
	headEndTag := []byte("</head>")
	if bytes.Contains(modified, headEndTag) {
		customStyles := fmt.Sprintf(`
<link rel="stylesheet" href="https://wikipedia.org/w/load.php?debug=false&lang=en&modules=site.styles&only=styles&skin=vector">
<link rel="stylesheet" href="/static/custom.css">
<script src="/static/custom.js"></script>
</head>`)
		modified = bytes.Replace(modified, headEndTag, []byte(customStyles), 1)
	}

	return modified
}

func staticHandler(w http.ResponseWriter, r *http.Request) {

	filename := strings.TrimPrefix(r.URL.Path, "/static/")

	var content []byte
	var contentType string

	filepath := fmt.Sprintf("./static/%s", filename) // Path to CSS/JS files at project root

	switch filename {
	case "custom.css":
		contentType = "text/css"
	case "custom.js":
		contentType = "application/javascript"
	default:
		http.NotFound(w, r)
		return
	}

	// Read file content
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading %s: %v", filename, err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	io.WriteString(w, string(content))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle static files
	if strings.HasPrefix(r.URL.Path, "/static/") {
		staticHandler(w, r)
		return
	}

	// Handle resource files from Wikipedia
	if strings.HasPrefix(r.URL.Path, "/w/") || strings.HasPrefix(r.URL.Path, "/static/") {
		content, contentType, err := fetchWikipediaContent(r.URL.Path)
		if err != nil {
			log.Printf("Error fetching resource: %v", err)
			http.Error(w, "Error fetching resource", http.StatusInternalServerError)
			return
		}

		// Set appropriate content type based on file extension
		if contentType == "" {
			ext := filepath.Ext(r.URL.Path)
			contentType = mime.TypeByExtension(ext)
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(content)
		return
	}

	content, contentType, err := fetchWikipediaContent(r.URL.Path)
	if err != nil {
		log.Printf("Error fetching content: %v", err)
		http.Error(w, "Error fetching content", http.StatusInternalServerError)
		return
	}

	modifiedContent := modifyWikipediaContent(content, contentType)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprint(len(modifiedContent)))
	w.Header().Set("Server", "WikiProxy/1.0")

	if _, err := io.Copy(w, bytes.NewReader(modifiedContent)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func proxyServer() {
	http.HandleFunc("/", proxyHandler)

	serverAddr := ":" + PORT
	log.Printf("Starting enhanced proxy server on http://localhost%s", serverAddr)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
