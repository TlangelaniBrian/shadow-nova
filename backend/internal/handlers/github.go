package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/models"
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
	state, err := auth.GenerateState("login")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}
	auth.SetStateCookie(w, state)
	url := h.authService.GetLoginURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GitHubHandler) Connect(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "github_connect_user_id",
		Value:    fmt.Sprintf("%d", userID),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300,
	})

	state, err := auth.GenerateState("connect")
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}
	auth.SetStateCookie(w, state)
	url := h.authService.GetLoginURL(state)

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *GitHubHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		httputil.WriteError(w, http.StatusBadRequest, "Code not found")
		return
	}

	flow, err := auth.ValidateState(r, w, state)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid or expired OAuth state")
		return
	}

	token, err := h.authService.Exchange(r.Context(), code)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	ghUser, err := h.authService.GetUser(r.Context(), token)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to get user info")
		return
	}

	var userID int
	var redirectPath string

	if flow == "connect" {
		cookie, err := r.Cookie("github_connect_user_id")
		if err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "Session expired or invalid flow")
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "github_connect_user_id", MaxAge: -1, Path: "/"})

		userID, err = strconv.Atoi(cookie.Value)
		if err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "Invalid user ID in cookie")
			return
		}
		redirectPath = "/profile"
	} else {
		user, err := h.db.GetUserByEmail(r.Context(), ghUser.Email)
		if err != nil {
			httputil.WriteError(w, http.StatusNotFound, "User not found. Please register first.")
			return
		}
		userID = user.ID

		jwtToken, err := auth.GenerateJWT(strconv.Itoa(user.ID), user.Username, user.Email)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}
		redirectPath = fmt.Sprintf("/auth/callback?token=%s", jwtToken)
	}

	integration := &models.GitHubIntegration{
		UserID:       userID,
		GithubUserID: fmt.Sprintf("%d", ghUser.ID),
		Username:     ghUser.Login,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}

	if err := h.db.SaveGitHubToken(r.Context(), integration); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to save integration")
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	redirectURL := fmt.Sprintf("%s%s", frontendURL, redirectPath)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
