package models

import "time"

type ShoppingListItem struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	ItemName     string     `json:"item_name"`
	Quantity     string     `json:"quantity"`
	Unit         string     `json:"unit"`
	IsChecked    bool       `json:"is_checked"`
	Source       string     `json:"source"`
	SourceMealID *string    `json:"source_meal_id"`
	CreatedAt    time.Time  `json:"created_at"`
	CheckedAt    *time.Time `json:"checked_at"`
}

type AddShoppingItemRequest struct {
	ItemName string `json:"item_name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}
