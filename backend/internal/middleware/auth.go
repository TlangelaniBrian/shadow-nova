package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"shadow-nova/backend/internal/httputil"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secret: []byte(secret),
	}
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}

func (a *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.secret, nil
		})

		if err != nil || !token.Valid {
			httputil.WriteError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httputil.WriteError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			httputil.WriteError(w, http.StatusUnauthorized, "Invalid user_id in token")
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID == 0 {
			httputil.WriteError(w, http.StatusUnauthorized, "Invalid user_id format")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
