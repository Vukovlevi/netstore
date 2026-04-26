import type { SubCategory, ApiResponse } from '../types/Types';

const API_URL = './api/sub_category';

export const subCategoryService = {
  getAll: async (): Promise<SubCategory[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  checkDeleted: async (name: string): Promise<{ id: number; name: string; category_id: number } | null> => {
    const res = await fetch(`${API_URL}?check_deleted=${encodeURIComponent(name)}`);
    if (!res.ok) return null;
    const data = await res.json();
    return data && data.id ? data : null;
  },

  create: async (name: string, category_id: number): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, category_id }),
    });
    const data = await res.json();
    if (!res.ok) {
      const err: any = new Error(data.message || 'Hiba a létrehozáskor');
      if (data && data.restorable) {
        err.restorable = true;
        err.restoreId = data.id;
      }
      throw err;
    }
    return data;
  },

  restore: async (id: number, name: string, category_id: number): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ restore: true, id, name, category_id }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Hiba a visszaállításkor');
    return data;
  },

  update: async (id: number, name: string, category_id: number): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, name, category_id }),
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