package repository

import (
	"database/sql"
	"fmt"

	"github.com/dishdice/backend/internal/models"
	"github.com/google/uuid"
)

type ShoppingRepository struct {
	db *sql.DB
}

func NewShoppingRepository(db *sql.DB) *ShoppingRepository {
	return &ShoppingRepository{db: db}
}

func (r *ShoppingRepository) Create(userID, itemName, quantity, unit, source string, sourceMealID *string) (*models.ShoppingListItem, error) {
	item := &models.ShoppingListItem{
		ID:           uuid.New().String(),
		UserID:       userID,
		ItemName:     itemName,
		Quantity:     quantity,
		Unit:         unit,
		Source:       source,
		SourceMealID: sourceMealID,
		IsChecked:    false,
	}

	query := `
		INSERT INTO shopping_list (id, user_id, item_name, quantity, unit, source, source_meal_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.db.QueryRow(query, item.ID, item.UserID, item.ItemName, item.Quantity, item.Unit, item.Source, item.SourceMealID).
		Scan(&item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create shopping item: %w", err)
	}

	return item, nil
}

// CreateOrMerge adds an item to the shopping list, merging with existing items if same name and unit
func (r *ShoppingRepository) CreateOrMerge(userID, itemName, quantity, unit, source string, sourceMealID *string) error {
	// Check if item with same name and unit exists
	checkQuery := `
		SELECT id, quantity
		FROM shopping_list
		WHERE user_id = $1 AND LOWER(item_name) = LOWER($2) AND unit = $3 AND is_checked = false
		LIMIT 1
	`

	var existingID, existingQty string
	err := r.db.QueryRow(checkQuery, userID, itemName, unit).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		// No existing item, create new one
		_, err := r.Create(userID, itemName, quantity, unit, source, sourceMealID)
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to check existing item: %w", err)
	}

	// Item exists, merge quantities
	var existingQtyFloat, newQtyFloat float64
	fmt.Sscanf(existingQty, "%f", &existingQtyFloat)
	fmt.Sscanf(quantity, "%f", &newQtyFloat)
	mergedQty := existingQtyFloat + newQtyFloat

	updateQuery := `
		UPDATE shopping_list
		SET quantity = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(updateQuery, fmt.Sprintf("%.2f", mergedQty), existingID)
	if err != nil {
		return fmt.Errorf("failed to merge quantities: %w", err)
	}

	return nil
}

func (r *ShoppingRepository) GetByUser(userID string, showChecked bool) ([]models.ShoppingListItem, error) {
	query := `
		SELECT id, user_id, item_name, quantity, unit, is_checked, source, source_meal_id, created_at, checked_at
		FROM shopping_list
		WHERE user_id = $1
	`

	if !showChecked {
		query += " AND is_checked = false"
	}

	query += " ORDER BY is_checked, created_at DESC"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shopping list: %w", err)
	}
	defer rows.Close()

	var items []models.ShoppingListItem
	for rows.Next() {
		var item models.ShoppingListItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ItemName,
			&item.Quantity,
			&item.Unit,
			&item.IsChecked,
			&item.Source,
			&item.SourceMealID,
			&item.CreatedAt,
			&item.CheckedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shopping item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ShoppingRepository) ToggleChecked(itemID, userID string) error {
	query := `
		UPDATE shopping_list
		SET is_checked = NOT is_checked,
		    checked_at = CASE WHEN is_checked THEN NULL ELSE CURRENT_TIMESTAMP END
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(query, itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to toggle item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("item not found or unauthorized")
	}

	return nil
}

func (r *ShoppingRepository) BulkDelete(userID string) error {
	query := `DELETE FROM shopping_list WHERE user_id = $1 AND is_checked = true`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to bulk delete: %w", err)
	}
	return nil
}

func (r *ShoppingRepository) Delete(itemID, userID string) error {
	query := `DELETE FROM shopping_list WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("item not found or unauthorized")
	}

	return nil
}
