package middleware

import (
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	if allowedOrigins == "" {
		log.Println("необходимо установить ALLOWED_ORIGINS")
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			if slices.Contains(strings.Split(allowedOrigins, ","), origin) {
				log.Println(origin)
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
