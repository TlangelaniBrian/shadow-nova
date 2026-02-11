package handlers

import (
	"net/http"
	"strconv"

	"shadow-nova/backend/internal/auth"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/httputil"
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

	jwtToken, err := auth.GenerateJWT(userInfo.Sub, userInfo.Name, userInfo.Email)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": jwtToken,
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

	jwtToken, err := auth.GenerateJWT(userInfo.Sub, userInfo.Name, userInfo.Email)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": jwtToken,
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
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	jwtToken, err := auth.GenerateJWT(strconv.Itoa(user.ID), user.Username, user.Email)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	httputil.WriteSuccess(w, "Login successful", map[string]string{
		"token":    jwtToken,
		"username": user.Username,
		"email":    user.Email,
	})
}
