import type { StoringCondition, ApiResponse } from '../types/Types';

const API_URL = './api/storing_condition';

export const storingConditionService = {
  getAll: async (): Promise<StoringCondition[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  create: async (description: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ description }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Hiba a létrehozáskor');
    return data;
  },

  update: async (id: number, description: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, description }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Hiba a frissítéskor');
    return data;
  },

  delete: async (id: number): Promise<void> => {
    const res = await fetch(`${API_URL}?id=${id}`, {
      method: 'DELETE',
    });
    
    if (!res.ok) {
      const data = await res.json();
      throw new Error(data.message || 'Hiba a törléskor');
    }
  }
};