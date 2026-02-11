package handlers

import (
	"net/http"
	"strconv"

	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
	"shadow-nova/backend/internal/middleware"
	"shadow-nova/backend/internal/models"
	"shadow-nova/backend/internal/validator"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsersHandler struct {
	db database.Service
}

func NewAdminUsersHandler(db database.Service) *AdminUsersHandler {
	return &AdminUsersHandler{db: db}
}

// GET /api/v1/admin/users - List all users with pagination
func (h *AdminUsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	pagination := models.ParsePagination(r)
	offset := (pagination.Page - 1) * pagination.Limit

	users, err := h.db.GetUsers(r.Context(), pagination.Limit, offset)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	total, err := h.db.GetUsersCount(r.Context())
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	// Remove password hashes (should already be empty due to json:"-" tag)
	for i := range users {
		users[i].PasswordHash = ""
	}

	response := models.NewPaginatedResponse(users, pagination.Page, pagination.Limit, total)
	httputil.WriteJSON(w, http.StatusOK, response)
}

// GET /api/v1/admin/users/{id} - Get specific user
func (h *AdminUsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.db.GetUserByID(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	user.PasswordHash = ""
	httputil.WriteJSON(w, http.StatusOK, user)
}

// POST /api/v1/admin/users - Create new user (admin only)
func (h *AdminUsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=8"`
		Role     string `json:"role" validate:"required,oneof=user admin"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	if err := h.db.CreateUser(r.Context(), user); err != nil {
		httputil.HandleError(w, err)
		return
	}

	user.PasswordHash = ""
	httputil.WriteJSON(w, http.StatusCreated, user)
}

// PUT /api/v1/admin/users/{id} - Update user (admin only)
func (h *AdminUsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Prevent admins from modifying their own account via this endpoint
	currentUserID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if currentUserID == id {
		httputil.WriteError(w, http.StatusBadRequest, "Cannot modify your own account via admin endpoint")
		return
	}

	var req struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Username string `json:"username" validate:"omitempty,min=3,max=100"`
		Role     string `json:"role" validate:"omitempty,oneof=user admin"`
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.WriteValidationError(w, err)
		return
	}

	user, err := h.db.GetUserByID(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := h.db.UpdateUser(r.Context(), id, user); err != nil {
		httputil.HandleError(w, err)
		return
	}

	user.PasswordHash = ""
	httputil.WriteJSON(w, http.StatusOK, user)
}

// DELETE /api/v1/admin/users/{id} - Soft delete user (admin only)
func (h *AdminUsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Prevent admins from deleting their own account
	currentUserID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if currentUserID == id {
		httputil.WriteError(w, http.StatusBadRequest, "Cannot delete your own account")
		return
	}

	if err := h.db.DeleteUser(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
