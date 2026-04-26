import type { Product, ApiResponse } from "../types/Types";

const API_URL = "./api/product";

export interface ProductPayload {
  name: string;
  description: string;
  amount: number;
  size: string;
  size_type: string;
  expires_at: string | null;
  price: number;
  discount: number;
  warranty: string | null;
  type_id: number;
  brand_id: number;
}

export const productService = {
  getAll: async (): Promise<Product[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  getSizeTypes: async (): Promise<string[]> => {
    const res = await fetch(`${API_URL}?distinct_size_types=1`);
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },

  checkDeleted: async (name: string, brand_id: number): Promise<{ id: number; name: string; brand_id: number } | null> => {
    const res = await fetch(`${API_URL}?check_deleted=${encodeURIComponent(name)}&brand_id=${brand_id}`);
    if (!res.ok) return null;
    const data = await res.json();
    return data && data.id ? data : null;
  },

  create: async (data: ProductPayload): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    const responseData = await res.json();
    if (!res.ok) {
      const err: any = new Error(responseData.message || "Hiba a létrehozáskor");
      if (responseData && responseData.restorable) {
        err.restorable = true;
        err.restoreId = responseData.id;
      }
      throw err;
    }
    return responseData;
  },

  restore: async (id: number, data: ProductPayload): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ restore: true, id, ...data }),
    });
    const responseData = await res.json();
    if (!res.ok)
      throw new Error(responseData.message || "Hiba a visszaállításkor");
    return responseData;
  },

  update: async (id: number, data: ProductPayload): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id, ...data }),
    });
    const responseData = await res.json();
    if (!res.ok) throw new Error(responseData.message || "Hiba a frissítéskor");
    return responseData;
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