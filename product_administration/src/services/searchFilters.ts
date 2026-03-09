import type { Product } from "../types/Types";

const API_URL = "./api/search_product";
const NET_API_URL = "http://localhost:8000/api/network-search"

export interface SearchFilters {
  name?: string;
  description?: string;
  category_id?: number;
  sub_category_id?: number;
  type_id?: number;
  brand_id?: number;
  storing_condition_id?: number;

  amount_min?: number;
  amount_max?: number;
  show_expired: boolean;

  size_val?: number;
  size_type?: string;

  price_min?: number;
  price_max?: number;

  has_warranty: boolean;
  is_discounted: boolean;

  page: number;
}

export interface SearchResponse {
  data: Product[];
  total: number;
  page: number;
  limit: number;
}

export const searchService = {
  search: async (filters: SearchFilters, NetworkSearch: boolean = false): Promise<SearchResponse> => {
    const params: any = {};

    if (filters.name && filters.name.trim() !== "") {
      params.name = filters.name.trim();
    }

    if (filters.description && filters.description.trim() !== "") {
      params.description = filters.description.trim();
    }

    if (filters.category_id !== undefined && filters.category_id !== 0) {
      params.category_id = filters.category_id.toString();
    }

    if (filters.sub_category_id !== undefined && filters.sub_category_id !== 0) {
      params.sub_category_id = filters.sub_category_id.toString();
    }

    if (filters.type_id !== undefined && filters.type_id !== 0) {
      params.type_id = filters.type_id.toString();
    }

    if (filters.brand_id !== undefined && filters.brand_id !== 0) {
      params.brand_id = filters.brand_id.toString();
    }

    if (filters.storing_condition_id !== undefined && filters.storing_condition_id !== 0) {
      params.storing_condition_id = filters.storing_condition_id.toString();
    }

    if (filters.amount_min !== undefined && filters.amount_min !== 0) {
      params.amount_min = filters.amount_min.toString();
    }

    if (filters.amount_max !== undefined && filters.amount_max !== 0) {
      params.amount_max = filters.amount_max.toString();
    }

    params.show_expired = filters.show_expired ? "true" : "false";

    if (filters.size_val !== undefined && filters.size_val !== 0) {
      params.size_val = filters.size_val.toString();
    }

    if (filters.size_type && filters.size_type.trim() !== "") {
      params.size_type = filters.size_type;
    }

    if (filters.price_min !== undefined && filters.price_min !== 0) {
      params.price_min = filters.price_min.toString();
    }

    if (filters.price_max !== undefined && filters.price_max !== 0) {
      params.price_max = filters.price_max.toString();
    }

    if (filters.has_warranty) {
      params.has_warranty = "true";
    }

    if (filters.is_discounted) {
      params.is_discounted = "true";
    }

    params.page = filters.page.toString();

    const res = await fetch(`${NetworkSearch ? NET_API_URL : API_URL}`, {
      method: 'POST',
      headers: {
        'Content-Type':'application/json'
      },
      body: JSON.stringify(params)
    });
    if (!res.ok) {
      throw new Error("Hiba a keresés során");
    }
    return res.json();
  },
};