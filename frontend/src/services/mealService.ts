import { api } from './api';
import type { DailyMeal, RegenerateMealResponse, SelectMealOptionRequest, ShoppingListItem } from '../types';

export const mealService = {
  async regenerateMeal(id: string): Promise<RegenerateMealResponse> {
    const response = await api.post<RegenerateMealResponse>(`/meals/${id}/regenerate`);
    return response.data;
  },

  async selectOption(id: string, data: SelectMealOptionRequest): Promise<DailyMeal> {
    const response = await api.put<DailyMeal>(`/meals/${id}/select`, data);
    return response.data;
  },

  async addToShoppingList(id: string): Promise<ShoppingListItem[]> {
    const response = await api.post<ShoppingListItem[]>(`/meals/${id}/save-to-shopping`);
    return response.data;
  },
};
