package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shubhouy1/asset-management/database/dbhelpers"
	"github.com/Shubhouy1/asset-management/utils"
	"github.com/golang-jwt/jwt/v5"

	"net/http"
	"os"
)

type AuthContext struct {
	UserID    string
	SessionID string
	Role      string
}

type contextKey string

const authContextKey contextKey = "authContext"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "authorization header required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "invalid authorization header", nil)
			return
		}

		tokenStr := parts[1]
		if tokenStr == "" {
			utils.RespondError(w, http.StatusUnauthorized, "authorization header required", nil)
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			utils.RespondError(w, http.StatusInternalServerError, "server configuration error", nil)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token", nil)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token claims", nil)
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token data", nil)
			return
		}
		sessionID, ok := claims["session_id"].(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token data", nil)
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token data", nil)
			return
		}
		dbUserID, err := dbhelpers.GetUserIDFromSession(sessionID)
		if err != nil || dbUserID != userID {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token", nil)
			return
		}
		authContext := AuthContext{
			UserID:    userID,
			SessionID: sessionID,
			Role:      role,
		}
		ctx := context.WithValue(r.Context(), authContextKey, authContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func RequiredRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok := GetAuthContext(r)
			if !ok {
				utils.RespondError(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}
			for _, role := range allowedRoles {
				if auth.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.RespondError(w, http.StatusForbidden, "forbidden", nil)
		})
	}
}
func GetAuthContext(r *http.Request) (AuthContext, bool) {
	auth, ok := r.Context().Value(authContextKey).(AuthContext)
	return auth, ok
}
