package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dishdice/backend/internal/middleware"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type ProposalHandler struct {
	proposalService *services.ProposalService
	shoppingService *services.ShoppingService
}

func NewProposalHandler(proposalService *services.ProposalService, shoppingService *services.ShoppingService) *ProposalHandler {
	return &ProposalHandler{
		proposalService: proposalService,
		shoppingService: shoppingService,
	}
}

func (h *ProposalHandler) ListProposals(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	proposals, err := h.proposalService.GetProposals(userID, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proposals)
}

func (h *ProposalHandler) GetProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	proposalID := chi.URLParam(r, "id")
	if proposalID == "" {
		http.Error(w, "Proposal ID is required", http.StatusBadRequest)
		return
	}

	proposal, err := h.proposalService.GetProposal(proposalID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err.Error() == "proposal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proposal)
}

func (h *ProposalHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.WeekStartDate == "" {
		http.Error(w, "Week start date is required", http.StatusBadRequest)
		return
	}

	// Parse date
	weekStartDate, err := time.Parse("2006-01-02", req.WeekStartDate)
	if err != nil {
		http.Error(w, "Invalid date format (use YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	proposal, err := h.proposalService.CreateWeeklyProposal(r.Context(), userID, weekStartDate, req.WeekPreferences, req.CurrentResources)
	if err != nil {
		if err.Error() == "proposal already exists for this week" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proposal)
}

func (h *ProposalHandler) DeleteProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	proposalID := chi.URLParam(r, "id")
	if proposalID == "" {
		http.Error(w, "Proposal ID is required", http.StatusBadRequest)
		return
	}

	err := h.proposalService.DeleteProposal(proposalID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err.Error() == "proposal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProposalHandler) AddToShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	proposalID := chi.URLParam(r, "id")
	if proposalID == "" {
		http.Error(w, "Proposal ID is required", http.StatusBadRequest)
		return
	}

	err := h.shoppingService.AddProposalToShoppingList(proposalID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err.Error() == "proposal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
