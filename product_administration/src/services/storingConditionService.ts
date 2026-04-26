import type { StoringCondition, ApiResponse } from '../types/Types';

const API_URL = './api/storing_condition';

export const storingConditionService = {
  getAll: async (): Promise<StoringCondition[]> => {
    const res = await fetch(API_URL);
    if (!res.ok) return [];
    return res.json();
  },

  checkDeleted: async (description: string): Promise<{ id: number; description: string } | null> => {
    const res = await fetch(`${API_URL}?check_deleted=${encodeURIComponent(description)}`);
    if (!res.ok) return null;
    const data = await res.json();
    return data && data.id ? data : null;
  },

  create: async (description: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ description }),
    });
    const data = await res.json();
    if (!res.ok) {
      const err: any = new Error(data.message || "Hiba a létrehozáskor");
      if (data && data.restorable) {
        err.restorable = true;
        err.restoreId = data.id;
      }
      throw err;
    }
    return data;
  },

  restore: async (id: number, description: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ restore: true, id, description }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Hiba a visszaállításkor");
    return data;
  },

  update: async (id: number, description: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id, description }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Hiba a frissítéskor");
    return data;
  },

  delete: async (id: number): Promise<void> => {
    const res = await fetch(`${API_URL}?id=${id}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      const data = await res.json();
      throw new Error(data.message || "Hiba a törléskor");
    }
  },
};