import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';

// Типы для унифицированного поиска
export interface UnifiedSearchParams {
  query: string;
  product_types?: ('marketplace' | 'storefront')[];
  page?: number;
  limit?: number;
  category_id?: string;
  price_min?: number;
  price_max?: number;
  sort_by?: 'relevance' | 'price' | 'date' | 'popularity';
  sort_order?: 'asc' | 'desc';
  storefront_id?: number;
  city?: string;
  language?: string;
}

export interface UnifiedSearchItem {
  id: string;
  product_type: 'marketplace' | 'storefront';
  product_id: number;
  name: string;
  description: string;
  price: number;
  currency: string;
  images: Array<{
    url: string;
    alt_text?: string;
    is_main: boolean;
  }>;
  category: {
    id: number;
    name: string;
    slug?: string;
  };
  location?: {
    city?: string;
    country?: string;
    lat?: number;
    lng?: number;
  };
  storefront?: {
    id: number;
    name: string;
    slug?: string;
    rating?: number;
    is_verified: boolean;
  };
  score: number;
  highlights?: Record<string, string[]>;
}

export interface UnifiedSearchResult {
  items: UnifiedSearchItem[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_more: boolean;
  took_ms: number;
  facets?: Record<string, any>;
}

export interface SearchSuggestion {
  text: string;
  type: 'text' | 'category' | 'product';
  category?: {
    id: number;
    name: string;
    slug: string;
  };
  product_id?: number;
}

export interface EnhancedSuggestion {
  id?: string;
  text: string;
  type: 'query' | 'category' | 'product' | 'brand' | 'location';
  category?: {
    id: number;
    name: string;
    slug: string;
    count?: number;
  };
  product?: {
    id: number;
    name: string;
    price: number;
    currency: string;
    image_url?: string;
  };
  metadata?: {
    count?: number;
    popularity?: number;
    brand?: string;
    location?: string;
  };
  highlight?: string;
}

export class UnifiedSearchService {
  /**
   * Выполняет унифицированный поиск по всем типам товаров
   */
  static async search(
    params: UnifiedSearchParams
  ): Promise<UnifiedSearchResult> {
    const baseUrl = configManager.getApiUrl({ internal: true });
    const fullUrl = baseUrl ? `${baseUrl}/api/v1/search` : '/api/v1/search'; // Для development proxy

    const url = new URL(fullUrl, baseUrl || window.location.origin);

    // Добавляем параметры в URL
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          // Для массивов используем запятую в качестве разделителя
          if (key === 'product_types') {
            url.searchParams.append(key, value.join(','));
          } else {
            value.forEach((v) => url.searchParams.append(key, v));
          }
        } else {
          url.searchParams.append(key, value.toString());
        }
      }
    });

    // Получаем токен авторизации
    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Search failed: ${response.status}`);
    }

    const result = await response.json();

    // Сохраняем успешный поиск в историю
    if (params.query && result.items && result.items.length > 0) {
      this.saveToHistory(params.query);
    }

    return result;
  }

  /**
   * Получает предложения для автодополнения
   */
  static async getSuggestions(
    prefix: string,
    size: number = 10
  ): Promise<SearchSuggestion[]> {
    const url = `${configManager.getApiUrl({ internal: true })}/api/v1/marketplace/enhanced-suggestions`;
    const params = new URLSearchParams({
      query: prefix,
      limit: size.toString(),
    });

    // Получаем токен авторизации
    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${url}?${params}`, {
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (!response.ok) {
      console.error('Failed to load suggestions:', response.status);
      return [];
    }

    const data = await response.json();
    const rawSuggestions = data.data || [];

    // Преобразуем данные от backend в нужный формат для frontend
    return rawSuggestions.map((item: any): SearchSuggestion => {
      if (item.type === 'category') {
        return {
          text: item.label || item.value,
          type: 'category',
          category: {
            id: item.category_id,
            name: item.label || item.value,
            slug: item.value || item.category_id.toString(),
          },
        };
      }

      if (item.type === 'product') {
        return {
          text: item.label || item.value,
          type: 'product',
          product_id: item.product_id || item.id,
        };
      }

      // Для text/query типов
      return {
        text: item.label || item.value || item.text,
        type: 'text',
      };
    });
  }

  /**
   * Получает расширенные предложения для автодополнения с дополнительной информацией
   */
  static async getEnhancedSuggestions(
    prefix: string,
    options?: {
      size?: number;
      include_categories?: boolean;
      include_products?: boolean;
      include_brands?: boolean;
      include_locations?: boolean;
      language?: string;
    }
  ): Promise<EnhancedSuggestion[]> {
    const baseUrl = configManager.getApiUrl({ internal: true });
    const url = new URL(
      `${baseUrl}/api/v1/marketplace/enhanced-suggestions`,
      baseUrl || window.location.origin
    );

    // Добавляем параметры
    url.searchParams.append('query', prefix);
    url.searchParams.append('limit', (options?.size || 10).toString());

    if (options?.include_categories !== undefined) {
      url.searchParams.append(
        'include_categories',
        options.include_categories.toString()
      );
    }
    if (options?.include_products !== undefined) {
      url.searchParams.append(
        'include_products',
        options.include_products.toString()
      );
    }
    if (options?.include_brands !== undefined) {
      url.searchParams.append(
        'include_brands',
        options.include_brands.toString()
      );
    }
    if (options?.include_locations !== undefined) {
      url.searchParams.append(
        'include_locations',
        options.include_locations.toString()
      );
    }
    if (options?.language) {
      url.searchParams.append('language', options.language);
    }

    // Получаем токен авторизации
    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url.toString(), {
        method: 'GET',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        console.error('Failed to load enhanced suggestions:', response.status);
        return [];
      }

      const data = await response.json();
      const rawSuggestions = data.data || [];

      // Преобразуем данные от backend в нужный формат для frontend
      return rawSuggestions.map((item: any): EnhancedSuggestion => {
        if (item.type === 'category') {
          return {
            text: item.label || item.value,
            type: 'category',
            category: {
              id: item.category_id,
              name: item.label || item.value,
              slug: item.value || item.category_id.toString(),
              count: item.metadata?.count,
            },
          };
        }

        if (item.type === 'product') {
          return {
            text: item.label || item.value,
            type: 'product',
            product: {
              id: item.product_id || item.id,
              name: item.label || item.value,
              price: item.price || 0,
              currency: item.currency || 'RSD',
              image_url: item.image_url,
            },
          };
        }

        // Для query/text типов
        return {
          text: item.label || item.value || item.text,
          type: item.type === 'query' ? 'query' : 'query',
        };
      });
    } catch (error) {
      console.error('Error fetching enhanced suggestions:', error);
      return [];
    }
  }

  /**
   * Получает популярные поисковые запросы
   */
  static async getTrendingSearches(): Promise<string[]> {
    // TODO: Реализовать endpoint для трендовых запросов
    return [
      'iPhone',
      'Samsung',
      'Laptop',
      'Квартира',
      'Автомобиль',
      'Велосипед',
    ];
  }

  /**
   * Получает историю поиска пользователя
   */
  static async getSearchHistory(): Promise<string[]> {
    // Получаем из localStorage
    const history = localStorage.getItem('search_history');
    if (history) {
      try {
        return JSON.parse(history);
      } catch {
        return [];
      }
    }
    return [];
  }

  /**
   * Сохраняет поисковый запрос в историю
   */
  static saveToHistory(query: string): void {
    // Получаем историю синхронно
    const historyString = localStorage.getItem('search_history');
    let history: string[] = [];
    if (historyString) {
      try {
        history = JSON.parse(historyString);
      } catch {
        history = [];
      }
    }

    const filtered = history.filter((q) => q !== query);
    const updated = [query, ...filtered].slice(0, 10); // Храним только 10 последних
    localStorage.setItem('search_history', JSON.stringify(updated));
  }

  /**
   * Очищает историю поиска
   */
  static clearHistory(): void {
    localStorage.removeItem('search_history');
  }
}
