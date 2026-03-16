package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dishdice/backend/internal/middleware"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	prefs := ""
	if user.Preferences != nil {
		prefs = *user.Preferences
	}

	response := map[string]string{
		"preferences": prefs,
		"language":    user.Language,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req models.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Default to English if not provided
	if req.Language == "" {
		req.Language = "en"
	}

	// Validate language (only en and de supported)
	if req.Language != "en" && req.Language != "de" {
		http.Error(w, "Language must be 'en' or 'de'", http.StatusBadRequest)
		return
	}

	err := h.userRepo.UpdatePreferences(userID, req.Preferences, req.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"preferences": req.Preferences,
		"language":    req.Language,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
