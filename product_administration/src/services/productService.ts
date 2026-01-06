import type { Product, ApiResponse } from "../types/Types";

const API_URL = "./api/product";

export const productService = {
  getAll: async (): Promise<Product[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  create: async (
    name: string,
    type_id: number,
    brand_id: number
  ): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, type_id, brand_id }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Hiba a létrehozáskor");
    return data;
  },

  update: async (
    id: number,
    name: string,
    type_id: number,
    brand_id: number
  ): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id, name, type_id, brand_id }),
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
