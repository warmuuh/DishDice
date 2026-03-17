package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dishdice/backend/internal/models"
	"github.com/google/uuid"
)

type ProposalRepository struct {
	db *sql.DB
}

func NewProposalRepository(db *sql.DB) *ProposalRepository {
	return &ProposalRepository{db: db}
}

func (r *ProposalRepository) CreateProposal(userID string, weekStartDate time.Time, weekPreferences, currentResources *string) (*models.WeeklyProposal, error) {
	proposal := &models.WeeklyProposal{
		ID:               uuid.New().String(),
		UserID:           userID,
		WeekStartDate:    weekStartDate,
		WeekPreferences:  weekPreferences,
		CurrentResources: currentResources,
	}

	query := `
		INSERT INTO weekly_proposals (id, user_id, week_start_date, week_preferences, current_resources)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err := r.db.QueryRow(query, proposal.ID, proposal.UserID, proposal.WeekStartDate, proposal.WeekPreferences, proposal.CurrentResources).
		Scan(&proposal.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create proposal: %w", err)
	}

	return proposal, nil
}

func (r *ProposalRepository) CreateDailyMeal(tx *sql.Tx, proposalID string, dayOfWeek int, menuName, recipe string) (*models.DailyMeal, error) {
	meal := &models.DailyMeal{
		ID:         uuid.New().String(),
		ProposalID: proposalID,
		DayOfWeek:  dayOfWeek,
		MenuName:   menuName,
		Recipe:     recipe,
	}

	query := `
		INSERT INTO daily_meals (id, proposal_id, day_of_week, menu_name, recipe)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err := tx.QueryRow(query, meal.ID, meal.ProposalID, meal.DayOfWeek, meal.MenuName, meal.Recipe).
		Scan(&meal.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create daily meal: %w", err)
	}

	return meal, nil
}

func (r *ProposalRepository) CreateMealShoppingItems(tx *sql.Tx, mealID string, items []models.MealShoppingItem) error {
	if len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO meal_shopping_items (id, daily_meal_id, item_name, quantity, unit)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range items {
		_, err := tx.Exec(query, uuid.New().String(), mealID, item.ItemName, item.Quantity, item.Unit)
		if err != nil {
			return fmt.Errorf("failed to create shopping item: %w", err)
		}
	}

	return nil
}

func (r *ProposalRepository) GetProposalsByUser(userID string, limit, offset int) ([]models.WeeklyProposal, error) {
	query := `
		SELECT id, user_id, week_start_date, week_preferences, current_resources, created_at
		FROM weekly_proposals
		WHERE user_id = $1
		ORDER BY week_start_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposals: %w", err)
	}
	defer rows.Close()

	var proposals []models.WeeklyProposal
	for rows.Next() {
		var p models.WeeklyProposal
		err := rows.Scan(&p.ID, &p.UserID, &p.WeekStartDate, &p.WeekPreferences, &p.CurrentResources, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proposal: %w", err)
		}
		proposals = append(proposals, p)
	}

	return proposals, nil
}

func (r *ProposalRepository) GetProposalByID(proposalID string) (*models.WeeklyProposal, error) {
	proposal := &models.WeeklyProposal{}
	query := `
		SELECT id, user_id, week_start_date, week_preferences, current_resources, created_at
		FROM weekly_proposals
		WHERE id = $1
	`

	err := r.db.QueryRow(query, proposalID).Scan(
		&proposal.ID,
		&proposal.UserID,
		&proposal.WeekStartDate,
		&proposal.WeekPreferences,
		&proposal.CurrentResources,
		&proposal.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Get daily meals
	meals, err := r.GetDailyMealsByProposal(proposalID)
	if err != nil {
		return nil, err
	}
	proposal.DailyMeals = meals

	return proposal, nil
}

func (r *ProposalRepository) GetDailyMealsByProposal(proposalID string) ([]models.DailyMeal, error) {
	query := `
		SELECT id, proposal_id, day_of_week, menu_name, recipe, created_at
		FROM daily_meals
		WHERE proposal_id = $1
		ORDER BY day_of_week
	`

	rows, err := r.db.Query(query, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily meals: %w", err)
	}
	defer rows.Close()

	var meals []models.DailyMeal
	for rows.Next() {
		var m models.DailyMeal
		err := rows.Scan(&m.ID, &m.ProposalID, &m.DayOfWeek, &m.MenuName, &m.Recipe, &m.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meal: %w", err)
		}

		// Get shopping items for this meal
		items, err := r.GetMealShoppingItems(m.ID)
		if err != nil {
			return nil, err
		}
		m.ShoppingItems = items

		meals = append(meals, m)
	}

	return meals, nil
}

func (r *ProposalRepository) GetMealShoppingItems(mealID string) ([]models.MealShoppingItem, error) {
	query := `
		SELECT id, daily_meal_id, item_name, quantity, unit, created_at
		FROM meal_shopping_items
		WHERE daily_meal_id = $1
		ORDER BY item_name
	`

	rows, err := r.db.Query(query, mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shopping items: %w", err)
	}
	defer rows.Close()

	var items []models.MealShoppingItem
	for rows.Next() {
		var item models.MealShoppingItem
		err := rows.Scan(&item.ID, &item.DailyMealID, &item.ItemName, &item.Quantity, &item.Unit, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shopping item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ProposalRepository) GetDailyMealByID(mealID string) (*models.DailyMeal, error) {
	meal := &models.DailyMeal{}
	query := `
		SELECT id, proposal_id, day_of_week, menu_name, recipe, created_at
		FROM daily_meals
		WHERE id = $1
	`

	err := r.db.QueryRow(query, mealID).Scan(
		&meal.ID,
		&meal.ProposalID,
		&meal.DayOfWeek,
		&meal.MenuName,
		&meal.Recipe,
		&meal.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get daily meal: %w", err)
	}

	// Get shopping items
	items, err := r.GetMealShoppingItems(meal.ID)
	if err != nil {
		return nil, err
	}
	meal.ShoppingItems = items

	return meal, nil
}

func (r *ProposalRepository) UpdateDailyMeal(tx *sql.Tx, mealID, menuName, recipe string) error {
	query := `
		UPDATE daily_meals
		SET menu_name = $1, recipe = $2
		WHERE id = $3
	`

	_, err := tx.Exec(query, menuName, recipe, mealID)
	if err != nil {
		return fmt.Errorf("failed to update daily meal: %w", err)
	}

	return nil
}

func (r *ProposalRepository) DeleteMealShoppingItems(tx *sql.Tx, mealID string) error {
	query := `DELETE FROM meal_shopping_items WHERE daily_meal_id = $1`
	_, err := tx.Exec(query, mealID)
	if err != nil {
		return fmt.Errorf("failed to delete shopping items: %w", err)
	}
	return nil
}

func (r *ProposalRepository) DeleteProposal(proposalID string) error {
	query := `DELETE FROM weekly_proposals WHERE id = $1`
	result, err := r.db.Exec(query, proposalID)
	if err != nil {
		return fmt.Errorf("failed to delete proposal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("proposal not found")
	}

	return nil
}

func (r *ProposalRepository) GetMealHistory(userID string, limit int) ([]models.MealHistory, error) {
	query := `
		SELECT meal_name, generated_at
		FROM meal_generation_history
		WHERE user_id = $1
		ORDER BY generated_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal history: %w", err)
	}
	defer rows.Close()

	var history []models.MealHistory
	for rows.Next() {
		var h models.MealHistory
		err := rows.Scan(&h.MealName, &h.GeneratedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

func (r *ProposalRepository) AddToHistory(tx *sql.Tx, userID, mealName string) error {
	query := `
		INSERT INTO meal_generation_history (id, user_id, meal_name)
		VALUES ($1, $2, $3)
	`

	_, err := tx.Exec(query, uuid.New().String(), userID, mealName)
	if err != nil {
		return fmt.Errorf("failed to add to history: %w", err)
	}

	return nil
}

func (r *ProposalRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *ProposalRepository) GetProposalByUserAndDate(userID string, weekStartDate time.Time) (*models.WeeklyProposal, error) {
	proposal := &models.WeeklyProposal{}
	query := `
		SELECT id, user_id, week_start_date, week_preferences, current_resources, created_at
		FROM weekly_proposals
		WHERE user_id = $1 AND week_start_date = $2
	`

	err := r.db.QueryRow(query, userID, weekStartDate).Scan(
		&proposal.ID,
		&proposal.UserID,
		&proposal.WeekStartDate,
		&proposal.WeekPreferences,
		&proposal.CurrentResources,
		&proposal.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal by user and date: %w", err)
	}

	return proposal, nil
}
