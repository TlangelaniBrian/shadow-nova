package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

type AuthMiddleware struct {
	secret []byte
	db     database.Service
}

func NewAuthMiddleware(secret string, db database.Service) *AuthMiddleware {
	return &AuthMiddleware{
		secret: []byte(secret),
		db:     db,
	}
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}

func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(UserRoleKey).(string)
	return role, ok
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (a *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Try to get token from cookie first
		tokenString, err := extractTokenFromCookie(r)
		if err != nil {
			// Fallback to Authorization header for API clients
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.WriteError(w, http.StatusUnauthorized, "Authorization required")
				return
			}
			tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
		}

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

		// Extract JTI from claims
		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			httputil.WriteError(w, http.StatusUnauthorized, "Invalid token: missing JTI")
			return
		}

		// Check if token is blacklisted
		isBlacklisted, err := a.db.IsTokenBlacklisted(r.Context(), jti)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "Token validation failed")
			return
		}
		if isBlacklisted {
			httputil.WriteError(w, http.StatusUnauthorized, "Token has been revoked")
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

		// Extract role from claims (default to "user" if not present)
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			role = "user"
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
