package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("very good\n"))
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	message := "your id is " + id + "\n"
	w.Write([]byte(message))
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

type Middleware func(http.Handler) http.Handler

func middlewareStack(arr ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(arr); i++ {
			middleware := arr[len(arr)-1-i]
			next = middleware(next)
		}
		return next
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(rw, r)
			log.Println(rw.statusCode, r.Method, r.URL.Path, time.Since(start))
		},
	)
}

func logger2(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(rw, r)
			log.Printf("| %d | %s | %s | %v", rw.statusCode, r.Method, r.URL.Path, time.Since(start))
		},
	)
}

func netHttp() {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, "/health"), healthHandler)
	router.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, "/rsrc/{id}"), resourceHandler)

	stack := middlewareStack(logger, logger2)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))
	server := http.Server{
		Addr:    ":8080",
		Handler: stack(v1),
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()

}
