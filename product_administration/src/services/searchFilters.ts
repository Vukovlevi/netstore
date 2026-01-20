import type { Product } from "../types/Types";

const API_URL = "./api/search_product";

export interface SearchFilters {
  name?: string;
  category_id?: number;
  sub_category_id?: number;
  type_id?: number;
  brand_id?: number;
  storing_condition_id?: number;

  amount_min?: number;
  amount_max?: number;
  show_expired?: boolean;

  size_val?: number;
  size_type?: string;

  price_min?: number;
  price_max?: number;

  has_warranty?: boolean;
  is_discounted?: boolean;
  other_properties?: string;

  page: number;
}

export interface SearchResponse {
  data: Product[];
  total: number;
  page: number;
  limit: number;
}

export const searchService = {
  search: async (filters: SearchFilters): Promise<SearchResponse> => {
    const params = new URLSearchParams();

    if (filters.name) params.append("name", filters.name);
    if (filters.category_id)
      params.append("category_id", filters.category_id.toString());
    if (filters.sub_category_id)
      params.append("sub_category_id", filters.sub_category_id.toString());
    if (filters.type_id) params.append("type_id", filters.type_id.toString());
    if (filters.brand_id)
      params.append("brand_id", filters.brand_id.toString());
    if (filters.storing_condition_id)
      params.append(
        "storing_condition_id",
        filters.storing_condition_id.toString(),
      );

    if (filters.amount_min)
      params.append("amount_min", filters.amount_min.toString());
    if (filters.amount_max)
      params.append("amount_max", filters.amount_max.toString());

    if (filters.show_expired) params.append("show_expired", "true");
    else params.append("show_expired", "false");

    if (filters.size_val)
      params.append("size_val", filters.size_val.toString());
    if (filters.size_type) params.append("size_type", filters.size_type);

    if (filters.price_min)
      params.append("price_min", filters.price_min.toString());
    if (filters.price_max)
      params.append("price_max", filters.price_max.toString());

    if (filters.has_warranty) params.append("has_warranty", "true");
    if (filters.is_discounted) params.append("is_discounted", "true");
    if (filters.other_properties)
      params.append("other_properties", filters.other_properties);

    params.append("page", filters.page.toString());

    const res = await fetch(`${API_URL}?${params.toString()}`);
    if (!res.ok) throw new Error("Hiba a keresés során");
    return res.json();
  },
};