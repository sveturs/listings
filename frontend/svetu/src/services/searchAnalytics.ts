import { apiClient } from './api-client';

// Типы для аналитических данных
export interface SearchMetrics {
  total_searches: number;
  unique_searches: number;
  average_results_count: number;
  average_search_duration_ms: number;
  top_queries: TopQuery[];
  search_trends: SearchTrend[];
  click_metrics: ClickMetrics;
}

export interface TopQuery {
  query: string;
  count: number;
  ctr: number;
  avg_position: number;
  avg_results: number;
}

export interface SearchTrend {
  date: string;
  searches_count: number;
  clicks_count: number;
  ctr: number;
}

export interface ClickMetrics {
  total_clicks: number;
  average_click_position: number;
  ctr: number;
  conversion_rate: number;
}

export interface ItemPerformance {
  item_id: string;
  item_title: string;
  impressions: number;
  clicks: number;
  ctr: number;
  average_position: number;
  conversions: number;
  revenue?: number;
}

export interface AnalyticsFilters {
  period?: 'day' | 'week' | 'month' | 'custom';
  date_from?: string;
  date_to?: string;
  query?: string;
  limit?: number;
  offset?: number;
}

class SearchAnalyticsService {
  private async fetchWithAuth<T>(url: string): Promise<T> {
    const response = await apiClient.get(url);

    if (!response.data) {
      throw new Error('API request failed');
    }

    return response.data.data || response.data;
  }

  /**
   * Получить общие метрики поиска
   */
  async getSearchMetrics(filters?: AnalyticsFilters): Promise<SearchMetrics> {
    const params = new URLSearchParams();

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const queryString = params.toString();
    const url = `/analytics/metrics/search${queryString ? `?${queryString}` : ''}`;

    return this.fetchWithAuth<SearchMetrics>(url);
  }

  /**
   * Получить производительность товаров
   */
  async getItemsPerformance(
    filters?: AnalyticsFilters
  ): Promise<ItemPerformance[]> {
    const params = new URLSearchParams();

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const queryString = params.toString();
    const url = `/analytics/metrics/items${queryString ? `?${queryString}` : ''}`;

    return this.fetchWithAuth<ItemPerformance[]>(url);
  }

  /**
   * Получить CTR по времени
   */
  async getCTRTrends(filters?: AnalyticsFilters): Promise<SearchTrend[]> {
    const searchMetrics = await this.getSearchMetrics(filters);
    return searchMetrics.search_trends || [];
  }

  /**
   * Получить топ запросы по CTR
   */
  async getTopQueriesByCTR(filters?: AnalyticsFilters): Promise<TopQuery[]> {
    const searchMetrics = await this.getSearchMetrics(filters);
    return searchMetrics.top_queries?.sort((a, b) => b.ctr - a.ctr) || [];
  }

  /**
   * Получить статистику конверсий
   */
  async getConversionMetrics(
    filters?: AnalyticsFilters
  ): Promise<ClickMetrics> {
    const searchMetrics = await this.getSearchMetrics(filters);
    return (
      searchMetrics.click_metrics || {
        total_clicks: 0,
        average_click_position: 0,
        ctr: 0,
        conversion_rate: 0,
      }
    );
  }
}

export const searchAnalyticsService = new SearchAnalyticsService();
