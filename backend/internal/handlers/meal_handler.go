package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dishdice/backend/internal/middleware"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type MealHandler struct {
	mealService     *services.MealService
	shoppingService *services.ShoppingService
}

func NewMealHandler(mealService *services.MealService, shoppingService *services.ShoppingService) *MealHandler {
	return &MealHandler{
		mealService:     mealService,
		shoppingService: shoppingService,
	}
}

func (h *MealHandler) RegenerateMeal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	mealID := chi.URLParam(r, "id")
	if mealID == "" {
		http.Error(w, "Meal ID is required", http.StatusBadRequest)
		return
	}

	options, err := h.mealService.RegenerateDayOptions(r.Context(), mealID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err.Error() == "meal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := models.RegenerateMealResponse{
		Options: options,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MealHandler) SelectMealOption(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	mealID := chi.URLParam(r, "id")
	if mealID == "" {
		http.Error(w, "Meal ID is required", http.StatusBadRequest)
		return
	}

	var req models.SelectMealOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MenuName == "" || req.Recipe == "" {
		http.Error(w, "Menu name and recipe are required", http.StatusBadRequest)
		return
	}

	// Create the selected option from the request data
	selectedOption := models.DailyMealOption{
		MenuName:      req.MenuName,
		Recipe:        req.Recipe,
		ShoppingItems: req.ShoppingItems,
	}

	meal, err := h.mealService.SelectDayOption(r.Context(), mealID, userID, req.OptionIndex, selectedOption)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

func (h *MealHandler) AddToShoppingList(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	mealID := chi.URLParam(r, "id")
	if mealID == "" {
		http.Error(w, "Meal ID is required", http.StatusBadRequest)
		return
	}

	err := h.shoppingService.AddMealToShoppingList(mealID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if err.Error() == "meal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
