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
  type_id: number;
  brand_id: number;
  is_active: number;
  type_name?: string;
  brand_name?: string;
}

export interface ApiResponse {
  message?: string;
  id?: number;
  name?: string;
}