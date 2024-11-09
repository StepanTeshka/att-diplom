package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JwtAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			log.Println("Ошибка валидации токена:", err)
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		if _, ok := (*claims)["authorized"].(bool); !ok {
			http.Error(w, "Доступ запрещён", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
