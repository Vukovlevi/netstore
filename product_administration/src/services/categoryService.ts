import type { Category, ApiResponse } from '../types/Types';

const API_URL = 'http://localhost/netstore_api/crud_operations.php/category';

export const categoryService = {
  getAll: async (): Promise<Category[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  getOne: async (id: number): Promise<Category> => {
    const res = await fetch(`${API_URL}?id=${id}`);
    if (!res.ok) throw new Error('Category not found');
    return res.json();
  },

  create: async (name: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Error creating category');
    return data;
  },

  update: async (id: number, name: string): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, name }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Error updating category');
    return data;
  },

  delete: async (id: number): Promise<void> => {
    const res = await fetch(`${API_URL}?id=${id}`, {
      method: 'DELETE',
    });
    
    if (!res.ok) {
      const data = await res.json();
      throw new Error(data.message || 'Error deleting category');
    }
  }
};