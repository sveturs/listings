import { config } from '@/config';

interface CategoryHints {
  domain: string;
  productType: string;
  keywords: string[];
}

interface CategoryDetectionRequest {
  keywords?: string[];
  attributes?: Record<string, any>;
  domain?: string;
  productType?: string;
  language?: string;
  title?: string;
  description?: string;
  include_debug_info?: boolean;
}

interface CategoryDetectionResponse {
  category_id: number;
  category_name: string;
  category_slug: string;
  confidence_score: number;
  method: string;
  stats_id?: number;
  alternative_categories?: {
    category_id: number;
    category_name: string;
    category_slug: string;
    confidence_score: number;
  }[];
  processing_time_ms: number;
  debug_info?: {
    keyword_score: number;
    similarity_score: number;
    matched_keywords: string[];
  };
}

export class CategoryDetectorService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = config.api.url || '';
  }

  /**
   * Определяет категорию на основе семантической информации из AI анализа
   */
  async detectCategory(
    categoryHints: CategoryHints,
    title?: string,
    description?: string,
    language: string = 'ru'
  ): Promise<CategoryDetectionResponse> {
    try {
      const request: CategoryDetectionRequest = {
        keywords: categoryHints.keywords,
        domain: categoryHints.domain,
        productType: categoryHints.productType,
        title,
        description,
        language,
        include_debug_info: true,
      };

      const response = await fetch(
        `${this.baseUrl}/api/v1/marketplace/categories/detect`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
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
   */
  async confirmCategorySelection(
    statsId: number,
    confirmed: boolean,
    selectedCategoryId?: number
  ): Promise<boolean> {
    try {
      const response = await fetch(
        `${this.baseUrl}/api/v1/marketplace/categories/detect/${statsId}/confirm`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
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
