package handlers

import (
	"net/http"

	"github.com/ARKTEEK/shorty/internal/middleware"
	"github.com/ARKTEEK/shorty/internal/models"
	"github.com/ARKTEEK/shorty/internal/services"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := DecodeBody(r, &req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response, err := h.auth.Register(r.Context(), req)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := DecodeBody(r, &req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response, err := h.auth.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	var req models.DeactivateRequest
	if err := DecodeBody(r, &req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	request := &models.DeactivateRequest{
		UserID:   userID,
		Password: req.Password,
	}

	success, err := h.auth.Deactivate(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	WriteJSON(w, http.StatusAccepted, success)
}
