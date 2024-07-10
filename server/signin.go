package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Структура запроса на аутентификацию
type SigninRequest struct {
	Password string `json:"password"`
}

// Структура ответа на аутентификацию
type SigninResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

// Обработчик запросов на аутентификацию
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != expectedPassword {
		resp := SigninResponse{
			Error: "Неверный пароль",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	token, err := GenerateJWT(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SigninResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Генерация JWT-токен на основе пароля
func GenerateJWT(password string) (string, error) {
	// Вычисляем хэш пароля
	passwordHash := sha256.Sum256([]byte(password))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passwordHash": hex.EncodeToString(passwordHash[:]),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	})

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	return token.SignedString(secretKey)
}
