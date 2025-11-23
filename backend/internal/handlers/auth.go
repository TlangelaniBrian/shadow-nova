package handlers

import (
	"encoding/json"
	"net/http"

	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"

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

// GoogleLogin initiates Google OAuth flow
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := "random-state-token" // TODO: Generate and store secure state token
	url := h.googleAuth.GetAuthURL(state)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"auth_url": url,
	})
}

// GoogleCallback handles the OAuth callback from Google
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	h.googleAuth.HandleGoogleCallback(w, r)
}

// VerifyGoogleToken verifies a Google ID token from the frontend
func (h *AuthHandler) VerifyGoogleToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}
	
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}
	
	// Verify the Google token
	userInfo, err := h.googleAuth.VerifyGoogleToken(r.Context(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid Google token", http.StatusUnauthorized)
		return
	}
	
	// Generate our JWT
	jwtToken, err := h.googleAuth.GenerateJWT(userInfo.Sub, userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// TODO: Save or update user in database
	
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

// Register handles user registration (traditional email/password)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	
	// Validate request
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	// Create user in database
	if err := h.db.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "User registered successfully",
		Data: map[string]string{
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

// Login handles user authentication (traditional email/password)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	
	// Validate request
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}
	
	// Get user from database
	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	// Note: Using GoogleAuthService's GenerateJWT for consistency, but ideally should be a separate TokenService
	// Using User ID as sub
	jwtToken, err := h.googleAuth.GenerateJWT(string(rune(user.ID)), user.Email, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Login successful",
		Data: map[string]string{
			"token": jwtToken,
			"username": user.Username,
			"email": user.Email,
		},
	})
}

// UpdateProgress handles learning path progress updates
func (h *AuthHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProgressRequest
	
	// Validate request
	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}
	
	// TODO: Update progress in database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SuccessResponse{
		Message: "Progress updated successfully",
	})
}
