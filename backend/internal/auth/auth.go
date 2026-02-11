package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const oauthStateCookie = "oauth_state"

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for an authenticated user
func GenerateJWT(userID, name, email string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shadow-nova",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

// GenerateState creates a CSRF-safe state token with a flow prefix
func GenerateState(flowType string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return flowType + ":" + base64.URLEncoding.EncodeToString(b), nil
}

// SetStateCookie stores the OAuth state in an HttpOnly cookie for CSRF verification.
func SetStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})
}

// ValidateState checks that the state from the OAuth callback matches the cookie,
// clears the cookie, and returns the flow type.
func ValidateState(r *http.Request, w http.ResponseWriter, queryState string) (string, error) {
	cookie, err := r.Cookie(oauthStateCookie)
	if err != nil {
		return "", fmt.Errorf("missing state cookie")
	}

	// Clear the one-time-use cookie
	http.SetCookie(w, &http.Cookie{
		Name:   oauthStateCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if cookie.Value == "" || cookie.Value != queryState {
		return "", fmt.Errorf("state mismatch")
	}

	return parseState(queryState), nil
}

// parseState extracts the flow type from a state token.
func parseState(state string) string {
	parts := strings.SplitN(state, ":", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
