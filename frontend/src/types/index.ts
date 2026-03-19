export interface User {
  id: string;
  email: string;
  preferences?: string;
  language: string;
  role: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  ticket?: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterResponse {
  message?: string;
  status: string;
  token?: string;
  user: User;
}

export interface AdminUser {
  id: string;
  email: string;
  role: string;
  status: string;
  language: string;
  created_at: string;
}

export interface MealShoppingItem {
  id: string;
  daily_meal_id: string;
  item_name: string;
  quantity: string;
  unit: string;
  created_at: string;
}

export interface DailyMeal {
  id: string;
  proposal_id: string;
  day_of_week: number;
  menu_name: string;
  recipe: string;
  created_at: string;
  shopping_items: MealShoppingItem[];
}

export interface WeeklyProposal {
  id: string;
  user_id: string;
  week_start_date: string;
  week_preferences?: string;
  current_resources?: string;
  created_at: string;
  daily_meals: DailyMeal[];
}

export interface CreateProposalRequest {
  week_start_date: string;
  week_preferences?: string;
  current_resources?: string;
}

export interface DailyMealOption {
  menu_name: string;
  recipe: string;
  shopping_items: MealShoppingItem[];
}

export interface RegenerateMealResponse {
  options: DailyMealOption[];
}

export interface SelectMealOptionRequest {
  option_index: number;
  menu_name: string;
  recipe: string;
  shopping_items: MealShoppingItem[];
}

export interface ShoppingListItem {
  id: string;
  user_id: string;
  item_name: string;
  quantity: string;
  unit: string;
  is_checked: boolean;
  source: string;
  source_meal_id?: string;
  created_at: string;
  checked_at?: string;
}

export interface AddShoppingItemRequest {
  item_name: string;
  quantity: string;
  unit: string;
}

export interface RegistrationTicket {
  id: string;
  token: string;
  created_by: string;
  used_by?: string;
  created_at: string;
  expires_at: string;
  used_at?: string;
  is_used: boolean;
}

export interface CreateTicketResponse {
  ticket: RegistrationTicket;
  registration_link: string;
}

export interface ValidateTicketResponse {
  valid: boolean;
  message?: string;
  expires_at?: string;
}
