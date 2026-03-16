import { api } from './api';
import type { WeeklyProposal, CreateProposalRequest } from '../types';

export const proposalService = {
  async getProposals(page = 1, limit = 20): Promise<WeeklyProposal[]> {
    const response = await api.get<WeeklyProposal[]>('/proposals', {
      params: { page, limit },
    });
    return response.data;
  },

  async getProposal(id: string): Promise<WeeklyProposal> {
    const response = await api.get<WeeklyProposal>(`/proposals/${id}`);
    return response.data;
  },

  async createProposal(data: CreateProposalRequest): Promise<WeeklyProposal> {
    const response = await api.post<WeeklyProposal>('/proposals', data);
    return response.data;
  },

  async deleteProposal(id: string): Promise<void> {
    await api.delete(`/proposals/${id}`);
  },

  async addProposalToShoppingList(id: string): Promise<void> {
    await api.post(`/proposals/${id}/save-to-shopping`);
  },
};
