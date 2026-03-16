package services

import (
	"context"
	"fmt"

	"github.com/dishdice/backend/internal/ai"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type MealService struct {
	proposalRepo *repository.ProposalRepository
	userRepo     *repository.UserRepository
	aiClient     *ai.Client
}

func NewMealService(proposalRepo *repository.ProposalRepository, userRepo *repository.UserRepository, aiClient *ai.Client) *MealService {
	return &MealService{
		proposalRepo: proposalRepo,
		userRepo:     userRepo,
		aiClient:     aiClient,
	}
}

func (s *MealService) RegenerateDayOptions(ctx context.Context, mealID, userID string) ([]models.DailyMealOption, error) {
	// Get the meal
	meal, err := s.proposalRepo.GetDailyMealByID(mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}
	if meal == nil {
		return nil, fmt.Errorf("meal not found")
	}

	// Get the proposal
	proposal, err := s.proposalRepo.GetProposalByID(meal.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check authorization
	if proposal.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get user preferences
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userPrefs := ""
	if user.Preferences != nil {
		userPrefs = *user.Preferences
	}

	weekPrefs := ""
	if proposal.WeekPreferences != nil {
		weekPrefs = *proposal.WeekPreferences
	}

	resources := ""
	if proposal.CurrentResources != nil {
		resources = *proposal.CurrentResources
	}

	// Get other meals in the week
	var otherMeals []ai.DayMeal
	for _, m := range proposal.DailyMeals {
		if m.ID != mealID {
			var items []ai.ShoppingItem
			for _, item := range m.ShoppingItems {
				items = append(items, ai.ShoppingItem{
					ItemName: item.ItemName,
					Quantity: item.Quantity,
				})
			}
			otherMeals = append(otherMeals, ai.DayMeal{
				MenuName:      m.MenuName,
				Recipe:        m.Recipe,
				ShoppingItems: items,
			})
		}
	}

	// Get recent meal history
	history, err := s.proposalRepo.GetMealHistory(userID, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal history: %w", err)
	}

	recentMeals := make([]string, len(history))
	for i, h := range history {
		recentMeals[i] = h.MealName
	}

	// Get day name
	dayNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	dayName := dayNames[meal.DayOfWeek]

	// Call AI to generate options
	aiReq := ai.DayOptionsRequest{
		UserPreferences:  userPrefs,
		WeekPreferences:  weekPrefs,
		CurrentResources: resources,
		RecentMeals:      recentMeals,
		OtherDaysInWeek:  otherMeals,
		DayName:          dayName,
		Language:         user.Language,
	}

	aiOptions, err := s.aiClient.GenerateDayOptions(ctx, aiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate day options: %w", err)
	}

	// Convert AI options to model options
	var options []models.DailyMealOption
	for _, opt := range aiOptions.Options {
		var items []models.MealShoppingItem
		for _, item := range opt.ShoppingItems {
			items = append(items, models.MealShoppingItem{
				ItemName: item.ItemName,
				Quantity: item.Quantity,
				Unit:     item.Unit,
			})
		}
		options = append(options, models.DailyMealOption{
			MenuName:      opt.MenuName,
			Recipe:        opt.Recipe,
			ShoppingItems: items,
		})
	}

	return options, nil
}

func (s *MealService) SelectDayOption(ctx context.Context, mealID, userID string, optionIndex int, option models.DailyMealOption) (*models.DailyMeal, error) {
	// Get the meal
	meal, err := s.proposalRepo.GetDailyMealByID(mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}
	if meal == nil {
		return nil, fmt.Errorf("meal not found")
	}

	// Get the proposal and check authorization
	proposal, err := s.proposalRepo.GetProposalByID(meal.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}
	if proposal.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Begin transaction
	tx, err := s.proposalRepo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update meal
	err = s.proposalRepo.UpdateDailyMeal(tx, mealID, option.MenuName, option.Recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	// Delete old shopping items
	err = s.proposalRepo.DeleteMealShoppingItems(tx, mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old items: %w", err)
	}

	// Insert new shopping items
	err = s.proposalRepo.CreateMealShoppingItems(tx, mealID, option.ShoppingItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create new items: %w", err)
	}

	// Add to history
	err = s.proposalRepo.AddToHistory(tx, userID, option.MenuName)
	if err != nil {
		return nil, fmt.Errorf("failed to add to history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated meal
	updatedMeal, err := s.proposalRepo.GetDailyMealByID(mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated meal: %w", err)
	}

	return updatedMeal, nil
}
