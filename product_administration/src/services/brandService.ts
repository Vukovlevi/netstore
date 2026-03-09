import type { Brand, ApiResponse } from "../types/Types";

const API_URL = "./api/brand";

export const brandService = {
  getAll: async (): Promise<Brand[]> => {
    const res = await fetch(API_URL);
    return res.json();
  },

  create: async (
    name: string,
    is_own: boolean,
    is_temporary: boolean
  ): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name,
        is_own: is_own ? 1 : 0,
        is_temporary: is_temporary ? 1 : 0,
      }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Hiba a létrehozáskor");
    return data;
  },

  update: async (
    id: number,
    name: string,
    is_own: boolean,
    is_temporary: boolean
  ): Promise<ApiResponse> => {
    const res = await fetch(API_URL, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id,
        name,
        is_own: is_own ? 1 : 0,
        is_temporary: is_temporary ? 1 : 0,
      }),
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
