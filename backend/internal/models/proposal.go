package models

import "time"

type WeeklyProposal struct {
	ID               string      `json:"id"`
	UserID           string      `json:"user_id"`
	WeekStartDate    time.Time   `json:"week_start_date"`
	WeekPreferences  *string     `json:"week_preferences"`
	CurrentResources *string     `json:"current_resources"`
	CreatedAt        time.Time   `json:"created_at"`
	DailyMeals       []DailyMeal `json:"daily_meals,omitempty"`
}

type DailyMeal struct {
	ID           string              `json:"id"`
	ProposalID   string              `json:"proposal_id"`
	DayOfWeek    int                 `json:"day_of_week"`
	MenuName     string              `json:"menu_name"`
	Recipe       string              `json:"recipe"`
	CreatedAt    time.Time           `json:"created_at"`
	ShoppingItems []MealShoppingItem `json:"shopping_items,omitempty"`
}

type MealShoppingItem struct {
	ID          string    `json:"id"`
	DailyMealID string    `json:"daily_meal_id"`
	ItemName    string    `json:"item_name"`
	Quantity    string    `json:"quantity"`
	Unit        string    `json:"unit"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProposalRequest struct {
	WeekStartDate    string  `json:"week_start_date"`
	WeekPreferences  *string `json:"week_preferences"`
	CurrentResources *string `json:"current_resources"`
}

type RegenerateMealRequest struct {
	// Empty for now, might add options later
}

type RegenerateMealResponse struct {
	Options []DailyMealOption `json:"options"`
}

type DailyMealOption struct {
	MenuName      string              `json:"menu_name"`
	Recipe        string              `json:"recipe"`
	ShoppingItems []MealShoppingItem  `json:"shopping_items"`
}

type SelectMealOptionRequest struct {
	OptionIndex int              `json:"option_index"`
	MenuName    string           `json:"menu_name"`
	Recipe      string           `json:"recipe"`
	ShoppingItems []MealShoppingItem `json:"shopping_items"`
}

type MealHistory struct {
	MealName    string    `json:"meal_name"`
	GeneratedAt time.Time `json:"generated_at"`
}
