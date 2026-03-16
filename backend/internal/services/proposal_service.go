package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dishdice/backend/internal/ai"
	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type ProposalService struct {
	proposalRepo *repository.ProposalRepository
	userRepo     *repository.UserRepository
	aiClient     *ai.Client
}

func NewProposalService(proposalRepo *repository.ProposalRepository, userRepo *repository.UserRepository, aiClient *ai.Client) *ProposalService {
	return &ProposalService{
		proposalRepo: proposalRepo,
		userRepo:     userRepo,
		aiClient:     aiClient,
	}
}

func (s *ProposalService) CreateWeeklyProposal(ctx context.Context, userID string, weekStartDate time.Time, weekPreferences, currentResources *string) (*models.WeeklyProposal, error) {
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
	if weekPreferences != nil {
		weekPrefs = *weekPreferences
	}

	resources := ""
	if currentResources != nil {
		resources = *currentResources
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

	// Call AI to generate weekly plan
	aiReq := ai.WeeklyPlanRequest{
		UserPreferences:  userPrefs,
		WeekPreferences:  weekPrefs,
		CurrentResources: resources,
		RecentMeals:      recentMeals,
		Language:         user.Language,
	}

	plan, err := s.aiClient.GenerateWeeklyPlan(ctx, aiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate weekly plan: %w", err)
	}

	// Begin transaction
	tx, err := s.proposalRepo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create proposal
	proposal, err := s.proposalRepo.CreateProposal(userID, weekStartDate, weekPreferences, currentResources)
	if err != nil {
		return nil, fmt.Errorf("failed to create proposal: %w", err)
	}

	// Helper to create a daily meal
	createDayMeal := func(dayOfWeek int, dayMeal ai.DayMeal) error {
		meal, err := s.proposalRepo.CreateDailyMeal(tx, proposal.ID, dayOfWeek, dayMeal.MenuName, dayMeal.Recipe)
		if err != nil {
			return err
		}

		// Convert AI shopping items to model items
		var items []models.MealShoppingItem
		for _, item := range dayMeal.ShoppingItems {
			items = append(items, models.MealShoppingItem{
				ItemName: item.ItemName,
				Quantity: item.Quantity,
				Unit:     item.Unit,
			})
		}

		if err := s.proposalRepo.CreateMealShoppingItems(tx, meal.ID, items); err != nil {
			return err
		}

		// Add to history
		if err := s.proposalRepo.AddToHistory(tx, userID, dayMeal.MenuName); err != nil {
			return err
		}

		return nil
	}

	// Create all 7 daily meals
	dayMeals := []ai.DayMeal{
		plan.Monday,
		plan.Tuesday,
		plan.Wednesday,
		plan.Thursday,
		plan.Friday,
		plan.Saturday,
		plan.Sunday,
	}

	for i, dayMeal := range dayMeals {
		if err := createDayMeal(i, dayMeal); err != nil {
			return nil, fmt.Errorf("failed to create day %d meal: %w", i, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return complete proposal with meals
	fullProposal, err := s.proposalRepo.GetProposalByID(proposal.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get complete proposal: %w", err)
	}

	return fullProposal, nil
}

func (s *ProposalService) GetProposals(userID string, page, limit int) ([]models.WeeklyProposal, error) {
	offset := (page - 1) * limit
	return s.proposalRepo.GetProposalsByUser(userID, limit, offset)
}

func (s *ProposalService) GetProposal(proposalID, userID string) (*models.WeeklyProposal, error) {
	proposal, err := s.proposalRepo.GetProposalByID(proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}
	if proposal == nil {
		return nil, fmt.Errorf("proposal not found")
	}

	// Check authorization
	if proposal.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	return proposal, nil
}

func (s *ProposalService) DeleteProposal(proposalID, userID string) error {
	// Check authorization first
	proposal, err := s.proposalRepo.GetProposalByID(proposalID)
	if err != nil {
		return fmt.Errorf("failed to get proposal: %w", err)
	}
	if proposal == nil {
		return fmt.Errorf("proposal not found")
	}
	if proposal.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.proposalRepo.DeleteProposal(proposalID)
}
