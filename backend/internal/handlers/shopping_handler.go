package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dishdice/backend/internal/middleware"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type ShoppingHandler struct {
	shoppingService *services.ShoppingService
}

func NewShoppingHandler(shoppingService *services.ShoppingService) *ShoppingHandler {
	return &ShoppingHandler{
		shoppingService: shoppingService,
	}
}

func (h *ShoppingHandler) GetShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	showChecked := r.URL.Query().Get("show_checked") == "true"

	items, err := h.shoppingService.GetList(userID, showChecked)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ShoppingHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req models.AddShoppingItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemName == "" || req.Quantity == "" {
		http.Error(w, "Item name and quantity are required", http.StatusBadRequest)
		return
	}

	item, err := h.shoppingService.AddItem(userID, req.ItemName, req.Quantity, req.Unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *ShoppingHandler) ToggleItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	itemID := chi.URLParam(r, "id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	err := h.shoppingService.ToggleItem(itemID, userID)
	if err != nil {
		if err.Error() == "item not found or unauthorized" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShoppingHandler) DeleteChecked(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	err := h.shoppingService.DeleteCheckedItems(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShoppingHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	itemID := chi.URLParam(r, "id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	err := h.shoppingService.DeleteItem(itemID, userID)
	if err != nil {
		if err.Error() == "item not found or unauthorized" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
