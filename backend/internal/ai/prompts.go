package ai

import (
	"fmt"
	"strings"
)

func BuildWeeklyPrompt(req WeeklyPlanRequest) string {
	var sb strings.Builder

	// Determine language instruction
	languageInstruction := "Respond in English."
	if req.Language == "de" {
		languageInstruction = "Antworte auf Deutsch. Alle Rezepte, Zutaten und Anweisungen müssen auf Deutsch sein."
	}

	sb.WriteString("You are a creative meal planning assistant. Generate a weekly meal plan (Monday through Sunday) with diverse, delicious meals.\n\n")
	sb.WriteString(fmt.Sprintf("IMPORTANT: %s\n\n", languageInstruction))

	sb.WriteString("USER PREFERENCES:\n")
	if req.UserPreferences != "" {
		sb.WriteString(req.UserPreferences)
	} else {
		sb.WriteString("No specific preferences provided.")
	}
	sb.WriteString("\n\n")

	if req.WeekPreferences != "" {
		sb.WriteString("THIS WEEK'S SPECIFIC REQUESTS:\n")
		sb.WriteString(req.WeekPreferences)
		sb.WriteString("\n\n")
	}

	if req.CurrentResources != "" {
		sb.WriteString("AVAILABLE INGREDIENTS:\n")
		sb.WriteString(req.CurrentResources)
		sb.WriteString("\n(Try to incorporate these ingredients where appropriate)\n\n")
	}

	if len(req.RecentMeals) > 0 {
		sb.WriteString("RECENTLY GENERATED MEALS (AVOID REPEATING):\n")
		for _, meal := range req.RecentMeals {
			sb.WriteString(fmt.Sprintf("- %s\n", meal))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("REQUIREMENTS:\n")
	sb.WriteString("1. Generate 7 different meals, one for each day of the week\n")
	sb.WriteString("2. Ensure variety - no repeated main proteins or cuisines\n")
	sb.WriteString("3. Keep recipes practical and not overly complex\n")
	sb.WriteString("4. Include complete shopping lists with quantities\n")
	sb.WriteString("5. Avoid meals from the recent history list\n")
	sb.WriteString("6. Make meals family-friendly and delicious\n")
	sb.WriteString("7. CRITICAL: For shopping_items, quantity must be a NUMBER only (e.g., '2', '250', '0.5'), and unit must be separate\n")
	sb.WriteString("8. CRITICAL: NO fuzzy quantities! Use specific amounts:\n")
	sb.WriteString("   - Instead of 'nach Geschmack' (to taste) → use '1' with unit 'Prise'\n")
	sb.WriteString("   - Instead of 'einige' (some) → use '3' or '4' with appropriate unit\n")
	sb.WriteString("   - Instead of '1 Packung' → use '400' with unit 'g' or actual weight\n\n")

	sb.WriteString("Respond with ONLY a valid JSON object in this exact format:\n")
	sb.WriteString(`{
  "monday": {
    "menu_name": "Meal name",
    "recipe": "Detailed cooking instructions with steps",
    "shopping_items": [
      {"item_name": "Tomatoes", "quantity": "4", "unit": "Stück"},
      {"item_name": "Olive Oil", "quantity": "2", "unit": "EL"},
      {"item_name": "Salt", "quantity": "1", "unit": "Prise"},
      {"item_name": "Ground Beef", "quantity": "500", "unit": "g"}
    ]
  },
  "tuesday": {...},
  "wednesday": {...},
  "thursday": {...},
  "friday": {...},
  "saturday": {...},
  "sunday": {...}
}`)

	return sb.String()
}

func BuildDayOptionsPrompt(req DayOptionsRequest) string {
	var sb strings.Builder

	// Determine language instruction
	languageInstruction := "Respond in English."
	dayName := req.DayName
	if req.Language == "de" {
		languageInstruction = "Antworte auf Deutsch. Alle Rezepte, Zutaten und Anweisungen müssen auf Deutsch sein."
		// Translate day names
		dayNames := map[string]string{
			"Monday": "Montag", "Tuesday": "Dienstag", "Wednesday": "Mittwoch",
			"Thursday": "Donnerstag", "Friday": "Freitag", "Saturday": "Samstag", "Sunday": "Sonntag",
		}
		if translated, ok := dayNames[req.DayName]; ok {
			dayName = translated
		}
	}

	sb.WriteString(fmt.Sprintf("You are a creative meal planning assistant. Generate 3 diverse meal options for %s.\n\n", dayName))
	sb.WriteString(fmt.Sprintf("IMPORTANT: %s\n\n", languageInstruction))

	sb.WriteString("USER PREFERENCES:\n")
	if req.UserPreferences != "" {
		sb.WriteString(req.UserPreferences)
	} else {
		sb.WriteString("No specific preferences provided.")
	}
	sb.WriteString("\n\n")

	if req.WeekPreferences != "" {
		sb.WriteString("THIS WEEK'S SPECIFIC REQUESTS:\n")
		sb.WriteString(req.WeekPreferences)
		sb.WriteString("\n\n")
	}

	if req.CurrentResources != "" {
		sb.WriteString("AVAILABLE INGREDIENTS:\n")
		sb.WriteString(req.CurrentResources)
		sb.WriteString("\n\n")
	}

	if len(req.OtherDaysInWeek) > 0 {
		sb.WriteString("OTHER MEALS THIS WEEK:\n")
		for i, meal := range req.OtherDaysInWeek {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, meal.MenuName))
		}
		sb.WriteString("\n")
	}

	if len(req.RecentMeals) > 0 {
		sb.WriteString("RECENTLY GENERATED MEALS (AVOID REPEATING):\n")
		for _, meal := range req.RecentMeals {
			sb.WriteString(fmt.Sprintf("- %s\n", meal))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("REQUIREMENTS:\n")
	sb.WriteString("1. Generate 3 completely different meal options\n")
	sb.WriteString("2. Each option should have a different cuisine or style\n")
	sb.WriteString("3. Don't repeat proteins/styles from other days this week\n")
	sb.WriteString("4. Keep recipes practical and family-friendly\n")
	sb.WriteString("5. Include complete shopping lists with quantities\n")
	sb.WriteString("6. CRITICAL: For shopping_items, quantity must be a NUMBER only (e.g., '2', '250', '0.5'), and unit must be separate\n")
	sb.WriteString("7. CRITICAL: NO fuzzy quantities! Use specific amounts:\n")
	sb.WriteString("   - Instead of 'nach Geschmack' (to taste) → use '1' with unit 'Prise'\n")
	sb.WriteString("   - Instead of 'einige' (some) → use '3' or '4' with appropriate unit\n")
	sb.WriteString("   - Instead of '1 Packung' → use '400' with unit 'g' or actual weight\n\n")

	sb.WriteString("Respond with ONLY a valid JSON object in this exact format:\n")
	sb.WriteString(`{
  "options": [
    {
      "menu_name": "Meal name",
      "recipe": "Detailed cooking instructions",
      "shopping_items": [
        {"item_name": "Tomatoes", "quantity": "4", "unit": "Stück"},
        {"item_name": "Olive Oil", "quantity": "2", "unit": "EL"},
        {"item_name": "Salt", "quantity": "1", "unit": "Prise"}
      ]
    },
    {
      "menu_name": "Second option",
      "recipe": "...",
      "shopping_items": [...]
    },
    {
      "menu_name": "Third option",
      "recipe": "...",
      "shopping_items": [...]
    }
  ]
}`)

	return sb.String()
}
