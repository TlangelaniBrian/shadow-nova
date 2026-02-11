package handlers

import (
	"net/http"

	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db database.Service
}

func NewUserHandler(db database.Service) *UserHandler {
	return &UserHandler{db: db}
}

// GetProfile handles GET /api/v1/user/profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	user, err := h.db.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Don't send password hash to client
	user.PasswordHash = ""

	httputil.WriteJSON(w, http.StatusOK, user)
}

// UpdateProfile handles PATCH /api/v1/user/profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req struct {
		Username string `json:"username" validate:"omitempty,min=3,max=100"`
		Email    string `json:"email" validate:"omitempty,email"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	user, err := h.db.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.db.UpdateUser(r.Context(), userID, user); err != nil {
		httputil.HandleError(w, err)
		return
	}

	user.PasswordHash = ""
	httputil.WriteJSON(w, http.StatusOK, user)
}

// UpdatePassword handles PUT /api/v1/user/password
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	// Verify current password
	user, err := h.db.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	if err := h.db.UpdateUserPassword(r.Context(), userID, string(hashedPassword)); err != nil {
		httputil.HandleError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password updated successfully"})
}
