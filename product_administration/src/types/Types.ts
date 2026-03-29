export interface Category {
  id: number;
  name: string;
}

export interface SubCategory {
  id: number;
  name: string;
  category_id: number;
  category_name?: string;
}

export interface StoringCondition {
  id: number;
  description: string;
}

export interface ProductType {
  id: number;
  name: string;
  description: string;
  sub_id: number;
  storing_condition_id: number;
}

export interface Brand {
  id: number;
  name: string;
  is_own: number;
  is_temporary: number;
}

export interface Product {
  id: number;
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
  brand_name?: string;
  type_name?: string;
  sub_category_name?: string;
  category_name?: string;
  storing_condition_name?: string;
}

export interface ApiResponse {
  message?: string;
  id?: number;
  name?: string;
}

export interface NetworkStoreResult {
  store_detail: NetworkStoreDetail;
  products: {
    data: Product[];
    total: number;
    page: number;
    limit: number;
  };
}

export interface NetworkStoreDetail {
  address: string;
  storeTypeName?: string;
}
