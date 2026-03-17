import { api } from './api';
import type { AdminUser } from '../types';

export const adminService = {
  async getAllUsers(): Promise<AdminUser[]> {
    const response = await api.get<AdminUser[]>('/admin/users');
    return response.data;
  },

  async getPendingUsers(): Promise<AdminUser[]> {
    const response = await api.get<AdminUser[]>('/admin/users/pending');
    return response.data;
  },

  async approveUser(userId: string): Promise<void> {
    await api.put(`/admin/users/${userId}/approve`);
  },

  async rejectUser(userId: string): Promise<void> {
    await api.put(`/admin/users/${userId}/reject`);
  },
};
