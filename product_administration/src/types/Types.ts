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

export interface ApiResponse {
  message?: string;
  id?: number;
  name?: string;
}