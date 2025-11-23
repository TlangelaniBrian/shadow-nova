package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthService struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	jwtSecret []byte
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func NewGoogleAuthService() (*GoogleAuthService, error) {
	ctx := context.Background()
	
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID not set")
	}
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/callback"
	}
	
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	
	return &GoogleAuthService{
		config:   config,
		verifier: verifier,
		jwtSecret: []byte(jwtSecret),
	}, nil
}

// GetAuthURL returns the Google OAuth URL for login
func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// VerifyGoogleToken verifies a Google ID token and returns user info
func (s *GoogleAuthService) VerifyGoogleToken(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	token, err := s.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	
	var userInfo GoogleUserInfo
	if err := token.Claims(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}
	
	return &userInfo, nil
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (s *GoogleAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

// GenerateJWT creates a JWT token for an authenticated user
func (s *GoogleAuthService) GenerateJWT(userID, email, name string) (string, error) {
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
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT validates and parses a JWT token
func (s *GoogleAuthService) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

// HandleGoogleCallback handles the OAuth callback from Google
func (s *GoogleAuthService) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}
	
	// Exchange code for token
	oauth2Token, err := s.ExchangeCodeForToken(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	
	// Extract ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token in response", http.StatusInternalServerError)
		return
	}
	
	// Verify ID token
	userInfo, err := s.VerifyGoogleToken(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify token", http.StatusUnauthorized)
		return
	}
	
	// Generate our JWT
	jwtToken, err := s.GenerateJWT(userInfo.Sub, userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// Return JWT to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": jwtToken,
		"user": map[string]string{
			"id":      userInfo.Sub,
			"email":   userInfo.Email,
			"name":    userInfo.Name,
			"picture": userInfo.Picture,
		},
	})
}
