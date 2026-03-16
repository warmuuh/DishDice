package ai

type ShoppingItem struct {
	ItemName string `json:"item_name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

type DayMeal struct {
	MenuName      string         `json:"menu_name"`
	Recipe        string         `json:"recipe"`
	ShoppingItems []ShoppingItem `json:"shopping_items"`
}

type WeeklyPlan struct {
	Monday    DayMeal `json:"monday"`
	Tuesday   DayMeal `json:"tuesday"`
	Wednesday DayMeal `json:"wednesday"`
	Thursday  DayMeal `json:"thursday"`
	Friday    DayMeal `json:"friday"`
	Saturday  DayMeal `json:"saturday"`
	Sunday    DayMeal `json:"sunday"`
}

type WeeklyPlanRequest struct {
	UserPreferences  string
	WeekPreferences  string
	CurrentResources string
	RecentMeals      []string
	Language         string
}

type DayOptionsRequest struct {
	UserPreferences  string
	WeekPreferences  string
	CurrentResources string
	RecentMeals      []string
	OtherDaysInWeek  []DayMeal
	DayName          string
	Language         string
}

type DayOptions struct {
	Options []DayMeal `json:"options"`
}
