import { api } from './api';
import type { ShoppingListItem, AddShoppingItemRequest } from '../types';

export const shoppingService = {
  async getShoppingList(showChecked = false): Promise<ShoppingListItem[]> {
    const response = await api.get<ShoppingListItem[]>('/shopping-list', {
      params: { show_checked: showChecked },
    });
    return response.data;
  },

  async addItem(data: AddShoppingItemRequest): Promise<ShoppingListItem> {
    const response = await api.post<ShoppingListItem>('/shopping-list', data);
    return response.data;
  },

  async toggleItem(id: string): Promise<void> {
    await api.put(`/shopping-list/${id}/toggle`);
  },

  async deleteChecked(): Promise<void> {
    await api.delete('/shopping-list/checked');
  },

  async deleteItem(id: string): Promise<void> {
    await api.delete(`/shopping-list/${id}`);
  },
};
