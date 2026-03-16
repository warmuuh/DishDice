package services

import (
	"fmt"

	"github.com/dishdice/backend/internal/models"
	"github.com/dishdice/backend/internal/repository"
)

type ShoppingService struct {
	shoppingRepo *repository.ShoppingRepository
	proposalRepo *repository.ProposalRepository
}

func NewShoppingService(shoppingRepo *repository.ShoppingRepository, proposalRepo *repository.ProposalRepository) *ShoppingService {
	return &ShoppingService{
		shoppingRepo: shoppingRepo,
		proposalRepo: proposalRepo,
	}
}

func (s *ShoppingService) AddItem(userID, itemName, quantity, unit string) (*models.ShoppingListItem, error) {
	return s.shoppingRepo.Create(userID, itemName, quantity, unit, "manual", nil)
}

func (s *ShoppingService) AddMealToShoppingList(mealID, userID string) error {
	// Get the meal
	meal, err := s.proposalRepo.GetDailyMealByID(mealID)
	if err != nil {
		return fmt.Errorf("failed to get meal: %w", err)
	}
	if meal == nil {
		return fmt.Errorf("meal not found")
	}

	// Get proposal and check authorization
	proposal, err := s.proposalRepo.GetProposalByID(meal.ProposalID)
	if err != nil {
		return fmt.Errorf("failed to get proposal: %w", err)
	}
	if proposal.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Add each shopping item to the shopping list (with merging)
	for _, item := range meal.ShoppingItems {
		err := s.shoppingRepo.CreateOrMerge(userID, item.ItemName, item.Quantity, item.Unit, "meal", &mealID)
		if err != nil {
			return fmt.Errorf("failed to add item to shopping list: %w", err)
		}
	}

	return nil
}

func (s *ShoppingService) AddProposalToShoppingList(proposalID, userID string) error {
	// Get proposal and check authorization
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

	// Add all meals from the proposal
	for _, meal := range proposal.DailyMeals {
		for _, item := range meal.ShoppingItems {
			err := s.shoppingRepo.CreateOrMerge(userID, item.ItemName, item.Quantity, item.Unit, "meal", &meal.ID)
			if err != nil {
				return fmt.Errorf("failed to add item to shopping list: %w", err)
			}
		}
	}

	return nil
}

func (s *ShoppingService) GetList(userID string, showChecked bool) ([]models.ShoppingListItem, error) {
	return s.shoppingRepo.GetByUser(userID, showChecked)
}

func (s *ShoppingService) ToggleItem(itemID, userID string) error {
	return s.shoppingRepo.ToggleChecked(itemID, userID)
}

func (s *ShoppingService) DeleteCheckedItems(userID string) error {
	return s.shoppingRepo.BulkDelete(userID)
}

func (s *ShoppingService) DeleteItem(itemID, userID string) error {
	return s.shoppingRepo.Delete(itemID, userID)
}
