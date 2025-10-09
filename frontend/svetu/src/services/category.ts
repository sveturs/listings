import { apiClient } from './api-client';

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
    const response = await apiClient.get('/c2c/categories');
    return response.data?.data || [];
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
    try {
      const response = await apiClient.get(`/c2c/categories?locale=${locale}`);
      return response.data?.data || [];
    } catch {
      // Fallback to regular categories if with-counts endpoint doesn't exist
      return this.getCategories();
    }
  }
}
