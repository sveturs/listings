import { config } from '@/config';

interface CategoryHints {
  domain: string;
  productType: string;
  keywords: string[];
}

interface CategoryDetectionRequest {
  keywords?: string[];
  attributes?: Record<string, any>;
  title?: string;
  description?: string;
  language?: string;
  categoryHints?: {
    domain: string;
    productType: string;
    keywords: string[];
  };
}

interface CategoryDetectionResponse {
  categoryId: number;
  categoryName: string;
  categoryPath: string;
  confidenceScore: number;
  algorithm: string;
  processingTimeMs: number;
  keywords: string[];
  feedbackId?: number;
  stats_id?: number;
  alternative_categories?: {
    category_id: number;
    category_name: string;
    category_slug: string;
    confidence_score: number;
  }[];
  alternativeCategories?: {
    categoryId: number;
    categoryName: string;
    confidenceScore: number;
  }[];
  debug_info?: {
    matched_keywords?: string[];
  };
}

export class CategoryDetectorService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = config.api.url || '';
  }

  /**
   * Определяет категорию на основе семантической информации из AI анализа
   * Использует новую улучшенную AI систему с многоуровневым определением
   */
  async detectCategory(
    categoryHints: CategoryHints,
    title?: string,
    description?: string,
    language: string = 'ru'
  ): Promise<CategoryDetectionResponse> {
    try {
      const request: CategoryDetectionRequest = {
        title,
        description,
        keywords: categoryHints.keywords,
        language,
        categoryHints: {
          domain: categoryHints.domain,
          productType: categoryHints.productType,
          keywords: categoryHints.keywords,
        },
      };

      // Получаем токен авторизации из localStorage
      const token = typeof window !== 'undefined' ? localStorage.getItem('svetu_access_token') : null;

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      // Добавляем токен авторизации, если он есть
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      // Используем новый AI эндпоинт с многоуровневой системой определения
      const response = await fetch(
        `${this.baseUrl}/api/v1/marketplace/ai/detect-category`,
        {
          method: 'POST',
          headers,
          body: JSON.stringify(request),
        }
      );

      if (!response.ok) {
        const errorData = await response.text();
        console.error('Category detection error:', response.status, errorData);
        throw new Error(
          `Category detection failed: ${response.status} - ${errorData}`
        );
      }

      const result = await response.json();
      return result.data || result;
    } catch (error) {
      console.error('Category detection error:', error);
      throw error;
    }
  }

  /**
   * Подтверждает выбор категории пользователем для улучшения алгоритма
   * Используется для обучения AI системы на основе обратной связи
   */
  async confirmCategorySelection(
    statsId: number,
    confirmed: boolean,
    selectedCategoryId?: number
  ): Promise<boolean> {
    try {
      // Получаем токен авторизации из localStorage
      const token = typeof window !== 'undefined' ? localStorage.getItem('svetu_access_token') : null;

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      // Добавляем токен авторизации, если он есть
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${this.baseUrl}/api/v1/marketplace/ai/detect-category/${statsId}/confirm`,
        {
          method: 'PUT',
          headers,
          body: JSON.stringify({
            confirmed,
            selected_category_id: selectedCategoryId,
          }),
        }
      );

      if (!response.ok) {
        console.error(
          'Category confirmation error:',
          response.status,
          await response.text()
        );
        return false;
      }

      const result = await response.json();
      return result.data?.updated || false;
    } catch (error) {
      console.error('Category confirmation error:', error);
      return false;
    }
  }

  /**
   * Получает ключевые слова для указанной категории
   */
  async getCategoryKeywords(
    categoryId: number,
    language?: string
  ): Promise<
    Array<{
      id: number;
      keyword: string;
      language: string;
      weight: number;
      keyword_type: string;
      is_negative?: boolean;
      usage_count: number;
      success_rate: number;
    }>
  > {
    try {
      let url = `${this.baseUrl}/api/v1/marketplace/categories/${categoryId}/keywords`;
      if (language) {
        url += `?language=${language}`;
      }

      const response = await fetch(url);

      if (!response.ok) {
        console.error(
          'Get keywords error:',
          response.status,
          await response.text()
        );
        return [];
      }

      const result = await response.json();
      return result.data || [];
    } catch (error) {
      console.error('Get keywords error:', error);
      return [];
    }
  }
}

export const categoryDetector = new CategoryDetectorService();
