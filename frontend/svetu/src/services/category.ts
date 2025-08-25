import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';

export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id: number | null;
  children?: Category[];
  translations?: Record<string, string>;
  icon?: string;
  product_count?: number;
  listing_count?: number;
  count?: number;
  level?: number;
  is_active?: boolean;
  sort_order?: number;
}

export class CategoryService {
  static async getCategories(): Promise<Category[]> {
    const baseUrl = configManager.getApiUrl({ internal: true });
    const url = `${baseUrl}/api/v1/marketplace/categories`;

    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch categories: ${response.status}`);
    }

    const result = await response.json();
    return result.data || [];
  }

  static async getCategoryTree(): Promise<Category[]> {
    const categories = await this.getCategories();
    return this.buildTree(categories);
  }

  private static buildTree(categories: Category[]): Category[] {
    const categoryMap = new Map<number, Category>();
    const rootCategories: Category[] = [];

    // First pass: create map
    categories.forEach((cat) => {
      categoryMap.set(cat.id, { ...cat, children: [] });
    });

    // Second pass: build tree
    categories.forEach((cat) => {
      const category = categoryMap.get(cat.id);
      if (!category) return;

      if (cat.parent_id === null || cat.parent_id === 0) {
        rootCategories.push(category);
      } else {
        const parent = categoryMap.get(cat.parent_id);
        if (parent) {
          if (!parent.children) parent.children = [];
          parent.children.push(category);
        }
      }
    });

    return rootCategories;
  }

  static async getCategoryWithCounts(
    locale: string = 'en'
  ): Promise<Category[]> {
    const baseUrl = configManager.getApiUrl({ internal: true });
    const url = `${baseUrl}/api/v1/marketplace/categories?locale=${locale}`;

    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (!response.ok) {
      // Fallback to regular categories if with-counts endpoint doesn't exist
      return this.getCategories();
    }

    const result = await response.json();
    return result.data || [];
  }
}
