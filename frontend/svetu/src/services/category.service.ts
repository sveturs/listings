import { apiClient } from '@/lib/api-client';

export interface CategoryTranslations {
  en: string;
  ru: string;
  sr: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  icon: string;
  parent_id?: number;
  created_at: string;
  level: number;
  path: string;
  listing_count: number;
  children_count: number;
  translations: CategoryTranslations;
  children?: Category[];
}

export interface CategoryTreeResponse {
  data: Category[];
  success: boolean;
}

class CategoryService {
  async getCategoryTree(): Promise<Category[]> {
    try {
      const response = await apiClient.get<CategoryTreeResponse>('/api/v1/marketplace/category-tree');
      return response.data.data;
    } catch (error) {
      console.error('Error fetching category tree:', error);
      return [];
    }
  }

  async getTopLevelCategories(): Promise<Category[]> {
    const allCategories = await this.getCategoryTree();
    return allCategories.filter(category => category.level === 1);
  }

  async getCategoryById(id: number): Promise<Category | null> {
    const allCategories = await this.getCategoryTree();
    return this.findCategoryInTree(allCategories, id);
  }

  private findCategoryInTree(categories: Category[], id: number): Category | null {
    for (const category of categories) {
      if (category.id === id) {
        return category;
      }
      if (category.children) {
        const found = this.findCategoryInTree(category.children, id);
        if (found) return found;
      }
    }
    return null;
  }

  getCategoryName(category: Category, locale: string): string {
    if (locale === 'en') return category.translations.en;
    if (locale === 'ru') return category.translations.ru;
    if (locale === 'sr') return category.translations.sr;
    return category.name;
  }

  flattenCategories(categories: Category[]): Category[] {
    const flattened: Category[] = [];
    const seenIds = new Set<number>();
    
    const addCategory = (category: Category) => {
      if (!seenIds.has(category.id)) {
        seenIds.add(category.id);
        flattened.push(category);
        if (category.children) {
          category.children.forEach(addCategory);
        }
      }
    };
    
    categories.forEach(addCategory);
    
    return flattened;
  }
}

export const categoryService = new CategoryService();