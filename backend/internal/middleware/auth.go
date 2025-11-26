package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secret: []byte(secret),
	}
}

func (a *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and user_id
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get user_id from claims and convert to int
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			fmt.Printf("[Auth] Failed to extract user_id from claims. Claims: %+v\n", claims)
			http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
			return
		}

		// Convert string user_id to int
		userID := 0
		fmt.Sscanf(userIDStr, "%d", &userID)
		fmt.Printf("[Auth] Extracted user_id: '%s' -> %d\n", userIDStr, userID)
		if userID == 0 {
			http.Error(w, "Invalid user_id format", http.StatusUnauthorized)
			return
		}

		// Add user_id to context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
