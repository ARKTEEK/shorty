package handlers

import (
	"net/http"

	"github.com/ARKTEEK/shorty/internal/models"
	"github.com/ARKTEEK/shorty/internal/services"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	ls *services.LinkService
}

func NewLinkHandler(ls *services.LinkService) *LinkHandler {
	return &LinkHandler{ls: ls}
}

func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLinkRequest

	if err := DecodeBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OriginalUrl == "" {
		http.Error(w, "original_url is required", http.StatusBadRequest)
		return
	}

	link, err := h.ls.CreateShortLink(r.Context(), &req)
	if err != nil {
		http.Error(w, "Could not create short link", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusCreated, link)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	if shortCode == "" {
		http.Error(w, "Invalid short code", http.StatusBadRequest)
		return
	}

	originalUrl, err := h.ls.GetOriginalUrl(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	err = h.ls.IncrementVisits(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "Failed to increment visit count", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusFound)
}
