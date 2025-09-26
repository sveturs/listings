import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';
import { logger } from '@/utils/logger';

// Типы для унифицированного поиска
export interface UnifiedSearchParams {
  query: string;
  product_types?: ('marketplace' | 'storefront')[];
  page?: number;
  limit?: number;
  category_id?: string;
  category_ids?: number[]; // Массив ID категорий для множественного выбора
  price_min?: number;
  price_max?: number;
  sort_by?: 'relevance' | 'price' | 'date' | 'popularity';
  sort_order?: 'asc' | 'desc';
  storefront_id?: number;
  city?: string;
  language?: string;
  fuzzy?: boolean; // Параметр для включения нечеткого поиска
  latitude?: number; // Широта для геопоиска
  longitude?: number; // Долгота для геопоиска
  distance?: string; // Радиус поиска (например, "10km", "5000m")
  // Автомобильные параметры
  car_make?: string;
  car_model?: string;
  car_year_from?: number;
  car_year_to?: number;
  car_mileage_max?: number;
  car_fuel_type?: string;
  car_transmission?: string;
  car_body_type?: string[];
  advanced_geo_filters?: {
    travel_time?: {
      center_lat: number;
      center_lng: number;
      max_minutes: number;
      transport_mode: 'walking' | 'driving' | 'cycling' | 'transit';
    };
    poi_filter?: {
      poi_type: string;
      max_distance: number;
      min_count?: number;
    };
    density_filter?: {
      avoid_crowded: boolean;
      max_density?: number;
      min_density?: number;
    };
  };
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
  user?: {
    id: number;
    name: string;
    picture_url?: string;
    is_verified?: boolean;
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
  // Информация об остатках (для товаров витрин)
  stock_quantity?: number;
  stock_status?: 'in_stock' | 'out_of_stock' | 'low_stock';
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
    translations?: Record<string, string>;
  };
  product_id?: number;
  metadata?: {
    source_type?: 'marketplace' | 'storefront';
    storefront_id?: number;
    storefront_slug?: string;
    storefront?: string;
    price?: number;
    image?: string;
    category?: string;
    [key: string]: any;
  };
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

    // Используем правильный эндпоинт в зависимости от наличия геокоординат
    const endpoint =
      params.latitude && params.longitude
        ? '/api/v1/gis/search'
        : '/api/v1/search';
    const fullUrl = baseUrl ? `${baseUrl}${endpoint}` : endpoint; // Для development proxy

    const url = new URL(fullUrl, baseUrl || window.location.origin);

    // Debug logging
    logger.search.debug('UnifiedSearchService - params received:', params);

    // Добавляем параметры в URL
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (key === 'advanced_geo_filters') {
          // Обрабатываем расширенные геофильтры отдельно
          // Они будут отправлены в теле запроса
          return;
        }

        // Backend теперь поддерживает оба параметра: 'q' и 'query'
        // Но используем 'q' для совместимости
        let paramKey = key;
        if (key === 'query') {
          paramKey = 'q';
        } else if (key === 'language') {
          paramKey = 'lang';
        }

        if (Array.isArray(value)) {
          // Для массивов используем запятую в качестве разделителя
          if (key === 'product_types' || key === 'category_ids') {
            logger.search.debug(
              `UnifiedSearchService - adding array param ${key}:`,
              value
            );
            url.searchParams.append(paramKey, value.join(','));
          } else {
            value.forEach((v) => url.searchParams.append(paramKey, v));
          }
        } else {
          url.searchParams.append(paramKey, value.toString());
        }
      }
    });

    logger.search.debug('UnifiedSearchService - final URL:', url.toString());

    // Получаем токен авторизации
    const token = await tokenManager.getAccessToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    // Если есть расширенные геофильтры, используем POST запрос
    const hasAdvancedFilters =
      params.advanced_geo_filters &&
      Object.keys(params.advanced_geo_filters).length > 0;

    const fetchOptions: RequestInit = {
      method: hasAdvancedFilters ? 'POST' : 'GET',
      headers,
      credentials: 'include',
    };

    if (hasAdvancedFilters) {
      fetchOptions.body = JSON.stringify({
        advanced_geo_filters: params.advanced_geo_filters,
      });
    }

    const response = await fetch(url.toString(), fetchOptions);

    if (!response.ok) {
      throw new Error(`Search failed: ${response.status}`);
    }

    const data = await response.json();

    // Адаптируем ответ к унифицированному формату
    let result: UnifiedSearchResult;

    if (data.data && Array.isArray(data.data)) {
      // Старый формат от /marketplace/search
      const items: UnifiedSearchItem[] = data.data.map((item: any) => ({
        id: `marketplace-${item.id}`,
        product_type: 'marketplace' as const,
        product_id: item.id,
        name: item.title || item.name,
        description: item.description,
        price: item.price,
        currency: 'RSD',
        images: (item.images || []).map((img: any) => ({
          url: img.public_url || img.url,
          alt_text: item.title,
          is_main: img.is_main || false,
        })),
        category: {
          id: item.category_id,
          name: item.category_name || '',
          slug: item.category_slug,
        },
        location: {
          city: item.city,
          country: item.country,
          lat: item.latitude,
          lng: item.longitude,
        },
        user: {
          id: item.user_id,
          name: item.user_name || '',
          picture_url: item.user_picture_url,
          is_verified: item.user_is_verified || false,
        },
        score: item.score || 1.0,
        stock_quantity: item.stock_quantity,
        stock_status: item.stock_status,
      }));

      result = {
        items,
        total: data.total || items.length,
        page: params.page || 1,
        limit: params.limit || 20,
        total_pages: Math.ceil(
          (data.total || items.length) / (params.limit || 20)
        ),
        has_more: data.has_more || false,
        took_ms: data.took_ms || 0,
      };
    } else if (data.items) {
      // Новый формат уже готов
      result = data;
    } else {
      // Неизвестный формат
      result = {
        items: [],
        total: 0,
        page: params.page || 1,
        limit: params.limit || 20,
        total_pages: 0,
        has_more: false,
        took_ms: 0,
      };
    }

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
  ): Promise<(SearchSuggestion | EnhancedSuggestion)[]> {
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

    // Обрабатываем ответ от enhanced-suggestions endpoint
    // который возвращает объект с полями: suggestions, categories, popular_items
    if (
      data.data &&
      typeof data.data === 'object' &&
      !Array.isArray(data.data)
    ) {
      const enhancedData = data.data;
      const allSuggestions: SearchSuggestion[] = [];

      // Добавляем текстовые предложения
      if (enhancedData.suggestions && Array.isArray(enhancedData.suggestions)) {
        enhancedData.suggestions.forEach((text: string) => {
          allSuggestions.push({
            text: text,
            type: 'text' as const,
          });
        });
      }

      // Добавляем категории
      if (enhancedData.categories && Array.isArray(enhancedData.categories)) {
        enhancedData.categories.forEach((cat: any) => {
          allSuggestions.push({
            text: cat.name || '',
            type: 'category' as const,
            category: {
              id: cat.id,
              name: cat.name,
              slug: cat.slug,
              translations: cat.translations,
            },
          });
        });
      }

      // Добавляем популярные товары
      if (
        enhancedData.popular_items &&
        Array.isArray(enhancedData.popular_items)
      ) {
        enhancedData.popular_items.forEach((item: any) => {
          allSuggestions.push({
            text: item.title || '',
            type: 'product' as const,
            product_id: item.id,
            metadata: {
              price: item.price,
              currency: item.currency,
              image: item.image,
              location: item.location,
            },
          });
        });
      }

      return allSuggestions;
    }

    // Fallback для старого формата (если backend вернул массив)
    const rawSuggestions = Array.isArray(data.data) ? data.data : [];

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
            translations: item.translations || undefined,
          },
        };
      }

      if (item.type === 'product') {
        const suggestion: SearchSuggestion = {
          text: item.label || item.value,
          type: 'product',
          product_id: item.product_id || item.id,
        };

        // Добавляем metadata если есть
        if (item.metadata) {
          suggestion.metadata = item.metadata;
        }

        return suggestion;
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
