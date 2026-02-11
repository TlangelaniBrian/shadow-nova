package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	googleAuth *auth.GoogleAuthService
	db         database.Service
}

func NewAuthHandler(googleAuth *auth.GoogleAuthService, db database.Service) *AuthHandler {
	return &AuthHandler{
		googleAuth: googleAuth,
		db:         db,
	}
}

// Helper function to determine if cookies should be secure (HTTPS only)
func isSecureCookie() bool {
	// Allow insecure cookies in development (when ENV=development)
	return os.Getenv("ENV") != "development"
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := auth.GenerateState("google")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}
	auth.SetStateCookie(w, state)
	url := h.googleAuth.GetAuthURL(state)

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"auth_url": url})
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	state := r.URL.Query().Get("state")
	if _, err := auth.ValidateState(r, w, state); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid or expired OAuth state")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		httputil.WriteError(w, http.StatusBadRequest, "Missing code parameter")
		return
	}

	oauth2Token, err := h.googleAuth.ExchangeCodeForToken(ctx, code)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		httputil.WriteError(w, http.StatusInternalServerError, "No id_token in response")
		return
	}

	userInfo, err := h.googleAuth.VerifyGoogleToken(ctx, rawIDToken)
	if err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Failed to verify token")
		return
	}

	jwtToken, err := auth.GenerateJWT(userInfo.Sub, userInfo.Name, userInfo.Email, "user")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set HttpOnly cookie instead of returning token in JSON
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   86400, // 24 hours
	})

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]string{
			"id":      userInfo.Sub,
			"email":   userInfo.Email,
			"name":    userInfo.Name,
			"picture": userInfo.Picture,
		},
	})
}

func (h *AuthHandler) VerifyGoogleToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	userInfo, err := h.googleAuth.VerifyGoogleToken(r.Context(), req.IDToken)
	if err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid Google token")
		return
	}

	jwtToken, err := auth.GenerateJWT(userInfo.Sub, userInfo.Name, userInfo.Email, "user")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set HttpOnly cookie instead of returning token in JSON
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   86400, // 24 hours
	})

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]string{
			"id":      userInfo.Sub,
			"email":   userInfo.Email,
			"name":    userInfo.Name,
			"picture": userInfo.Picture,
		},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := h.db.CreateUser(r.Context(), user); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	httputil.WriteCreated(w, "User registered successfully", map[string]string{
		"username": req.Username,
		"email":    req.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Don't leak information about whether the user exists
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	jwtToken, err := auth.GenerateJWT(strconv.Itoa(user.ID), user.Username, user.Email, user.Role)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Set HttpOnly cookie instead of returning token in JSON
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   86400, // 24 hours
	})

	httputil.WriteSuccess(w, "Login successful", map[string]string{
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract token from request (cookie or header)
	token := extractTokenFromRequest(r)
	if token == "" {
		httputil.WriteError(w, http.StatusBadRequest, "No token provided")
		return
	}

	// Validate and extract claims
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Get user ID from context (set by middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Blacklist the token
	expiresAt := claims.ExpiresAt.Time
	err = h.db.BlacklistToken(r.Context(), claims.ID, userID, expiresAt, "user_logout")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1, // Delete cookie
	})

	httputil.WriteSuccess(w, "Logout successful", nil)
}

func extractTokenFromRequest(r *http.Request) string {
	// Try cookie first
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value
	}

	// Fallback to Authorization header
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

func (h *AuthHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	w.Header().Set("X-CSRF-Token", token)
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"csrf_token": token,
	})
}
