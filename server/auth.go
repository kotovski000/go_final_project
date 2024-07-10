package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Проверка аутентификации пользователя перед выполнением запроса
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPassword := os.Getenv("TODO_PASSWORD")
		if len(expectedPassword) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			valid, err := ValidateJWT(cookie.Value)
			if err != nil || !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Проверка валидность JWT-токена
func ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expectedPassword := os.Getenv("TODO_PASSWORD")
		passwordHash := sha256.Sum256([]byte(expectedPassword))
		return claims["passwordHash"] == hex.EncodeToString(passwordHash[:]), nil
	}

	return false, nil
}
