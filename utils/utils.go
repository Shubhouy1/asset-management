package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func RespondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return
	}
}

func RespondError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	resp := Error{
		Error:      errStr,
		StatusCode: statusCode,
		Message:    message,
	}
	encodeErr := json.NewEncoder(w).Encode(resp)
	if encodeErr != nil {
		return
	}
}
func ParseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil

}
func GenerateJWT(userID, sessionID, userRole string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"role":       userRole,
		"exp":        time.Now().Add(time.Minute * 10).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
