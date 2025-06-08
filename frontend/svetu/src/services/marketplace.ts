import {
  MarketplaceSearchParams,
  MarketplaceSearchResponse,
} from '@/types/marketplace';
import configManager from '@/config';

export interface CategoryResponse {
  data: Array<{
    id: number;
    name: string;
    slug?: string;
    parent_id?: number | null;
    icon?: string;
    translations?: Record<string, string>;
    level?: number;
    count?: number;
  }>;
  success: boolean;
}

export interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: 'text' | 'number' | 'select' | 'multiselect' | 'boolean';
  icon?: string;
  options?: {
    values?: string[];
    min?: number;
    max?: number;
    step?: number;
    multiselect?: boolean;
  };
  validation_rules?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_required: boolean;
  sort_order: number;
  created_at: string;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
  custom_component?: string;
}

export interface CategoryAttributeMapping {
  category_id: number;
  attribute_id: number;
  is_enabled: boolean;
  is_required: boolean;
  sort_order: number;
  custom_component?: string;
  attribute?: CategoryAttribute;
}

export interface CategoryAttributesResponse {
  success: boolean;
  data: CategoryAttribute[];
}

export class MarketplaceService {
  static async search(
    params: MarketplaceSearchParams = {}
  ): Promise<MarketplaceSearchResponse> {
    const searchParams = new URLSearchParams();

    if (params.sort_by) searchParams.append('sort_by', params.sort_by);
    if (params.page !== undefined)
      searchParams.append('page', params.page.toString());
    if (params.size !== undefined)
      searchParams.append('size', params.size.toString());
    if (params.category) searchParams.append('category', params.category);
    if (params.minPrice !== undefined)
      searchParams.append('minPrice', params.minPrice.toString());
    if (params.maxPrice !== undefined)
      searchParams.append('maxPrice', params.maxPrice.toString());
    if (params.search) searchParams.append('search', params.search);

    const url = `${configManager.getApiUrl()}/api/v1/marketplace/search?${searchParams.toString()}`;

    try {
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(
          `Failed to fetch marketplace items: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching marketplace items:', error);
      throw error;
    }
  }

  static async getCategories(): Promise<CategoryResponse> {
    const url = `${configManager.getApiUrl()}/api/v1/marketplace/categories`;

    try {
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`Failed to fetch categories: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching categories:', error);
      throw error;
    }
  }

  static async getCategoryAttributes(
    categoryId: number
  ): Promise<CategoryAttributesResponse> {
    const url = `${configManager.getApiUrl()}/api/v1/marketplace/categories/${categoryId}/attributes`;

    try {
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(
          `Failed to fetch category attributes: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching category attributes:', error);
      throw error;
    }
  }
}
