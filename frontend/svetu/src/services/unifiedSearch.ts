import configManager from '@/config';

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

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Search failed: ${response.status}`);
    }

    return response.json();
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
      prefix,
      size: size.toString(),
    });

    const response = await fetch(`${url}?${params}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      console.error('Failed to load suggestions:', response.status);
      return [];
    }

    const data = await response.json();
    return data.data || [];
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
