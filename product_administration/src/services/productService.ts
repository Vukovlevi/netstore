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

  create: async (data: ProductPayload): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    const responseData = await res.json();
    if (!res.ok)
      throw new Error(responseData.message || "Hiba a létrehozáskor");
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