package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetCorsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			next.ServeHTTP(w, r)
		})
	}
}