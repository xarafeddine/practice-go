package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("very good\n"))
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	message := "your id is " + id + "\n"
	w.Write([]byte(message))
}

func netHttp() {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, "/health"), healthHandler)
	router.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, "/rsrc/{id}"), resourceHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()

}
