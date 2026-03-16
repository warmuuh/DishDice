import { api } from './api';

export const userService = {
  async getPreferences(): Promise<{ preferences: string; language: string }> {
    const response = await api.get<{ preferences: string; language: string }>('/user/preferences');
    return response.data;
  },

  async updatePreferences(preferences: string, language: string): Promise<{ preferences: string; language: string }> {
    const response = await api.put<{ preferences: string; language: string }>('/user/preferences', { preferences, language });
    return response.data;
  },
};
