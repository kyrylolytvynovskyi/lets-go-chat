package middleware

import (
	"log"
	"net/http"
)

func LogHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %v\n", *r)
		next.ServeHTTP(w, r)
		log.Printf("Response: %v\n", w)
	})
}
