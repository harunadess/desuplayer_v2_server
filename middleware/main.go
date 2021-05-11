package middleware

import (
	"log"
	"net/http"
	"strings"
)

var allowedOrigins = []string{"http://localhost:8080"}

func CustomMiddleware(next http.Handler) http.Handler {
	handler := LoggerMiddleware(next)
	return CorsMiddleware(handler)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%v - %v", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		},
	)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			allowedOriginsForHeader := strings.Join(allowedOrigins, ",")
			w.Header().Add("Access-Control-Allow-Origin", allowedOriginsForHeader)
			next.ServeHTTP(w, r)
		},
	)
}
