package middleware

import (
	"log"
	"net/http"
	"net/url"
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
			path, err := url.PathUnescape(r.RequestURI)
			if err != nil {
				log.Println("failed to unescape request URI")
				log.Printf("%v - %v", r.Method, r.RequestURI)
			} else {
				log.Printf("%v - %v", r.Method, path)
			}
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
