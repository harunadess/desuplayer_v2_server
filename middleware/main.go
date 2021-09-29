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
			origin := r.Header.Get("Origin")

			if err != nil {
				log.Println("failed to unescape request URI")
				log.Printf("%v - origin: %v, request URI: %v", r.Method, origin, r.RequestURI)
			} else {
				log.Printf("%v - origin: %v, request URI: %v", r.Method, origin, path)
			}
			next.ServeHTTP(w, r)
		},
	)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			allowed := false
			origin := r.Header.Get("Origin")

			for _, v := range allowedOrigins {
				if strings.EqualFold(origin, v) {
					allowed = true
					break
				}
			}

			if !allowed {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("origin not allowed"))
				return
			}

			allowedOriginsForHeader := strings.Join(allowedOrigins, ",")
			w.Header().Add("Access-Control-Allow-Origin", allowedOriginsForHeader)
			next.ServeHTTP(w, r)
		},
	)
}
