import { api } from './api';
import type { CreateTicketResponse, ValidateTicketResponse } from '../types';

export const ticketService = {
  async createTicket(): Promise<CreateTicketResponse> {
    const response = await api.post<CreateTicketResponse>('/admin/tickets');
    return response.data;
  },

  async validateTicket(token: string): Promise<ValidateTicketResponse> {
    const response = await api.get<ValidateTicketResponse>(`/tickets/${token}/validate`);
    return response.data;
  },
};
