package middleware

import (
	"log"
	"net/http"
)

func CustomMiddleware(next http.Handler) http.Handler {
	return LoggerMiddleWare(next)
}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%v - %v", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		},
	)
}
