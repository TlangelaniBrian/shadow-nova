package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"strconv"
)

type GitHubHandler struct {
	authService *auth.GitHubAuthService
	db          database.Service
}

func NewGitHubHandler(authService *auth.GitHubAuthService, db database.Service) *GitHubHandler {
	return &GitHubHandler{
		authService: authService,
		db:          db,
	}
}

func (h *GitHubHandler) Login(w http.ResponseWriter, r *http.Request) {
	// State should be random to prevent CSRF
	state := "login" 
	url := h.authService.GetLoginURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GitHubHandler) Connect(w http.ResponseWriter, r *http.Request) {
	// This endpoint is protected, so we know the user ID
	userID := r.Context().Value("user_id").(int)
	
	// Set a cookie to remember who initiated the connection
	// In production, sign this value!
	http.SetCookie(w, &http.Cookie{
		Name:     "github_connect_user_id",
		Value:    fmt.Sprintf("%d", userID),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300, // 5 minutes
	})

	state := "connect"
	url := h.authService.GetLoginURL(state)
	
	// Return the URL as JSON instead of redirecting
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
	})
}

func (h *GitHubHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state") // "login" or "connect"

	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := h.authService.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get GitHub user info
	ghUser, err := h.authService.GetUser(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	var userID int
	var redirectPath string

	if state == "connect" {
		// CONNECT FLOW
		// Retrieve user ID from cookie
		cookie, err := r.Cookie("github_connect_user_id")
		if err != nil {
			http.Error(w, "Session expired or invalid flow", http.StatusBadRequest)
			return
		}
		
		// Clear cookie
		http.SetCookie(w, &http.Cookie{Name: "github_connect_user_id", MaxAge: -1, Path: "/"})
		
		fmt.Sscanf(cookie.Value, "%d", &userID)
		redirectPath = "/profile" // Redirect back to profile page
	} else {
		// LOGIN FLOW
		// Check if user exists by email
		user, err := h.db.GetUserByEmail(r.Context(), ghUser.Email)
		if err != nil {
			// Handle user not found / registration logic here
			// For now, fail if not found
			http.Error(w, "User not found. Please register first.", http.StatusNotFound)
			return
		}
		userID = user.ID
		
		// Generate JWT for login
		jwtToken, err := auth.GenerateJWT(strconv.Itoa(user.ID), user.Username, user.Email)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		redirectPath = fmt.Sprintf("/auth/callback?token=%s", jwtToken)
	}

	// Save integration (this will also update the user's github_username)
	integration := &models.GitHubIntegration{
		UserID:       userID,
		GithubUserID: fmt.Sprintf("%d", ghUser.ID),
		Username:     ghUser.Login, // Add username to the integration model
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}
	
	if err := h.db.SaveGitHubToken(r.Context(), integration); err != nil {
		http.Error(w, "Failed to save integration", http.StatusInternalServerError)
		return
	}

	// Redirect to frontend
	redirectURL := fmt.Sprintf("http://localhost:5173%s", redirectPath)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
