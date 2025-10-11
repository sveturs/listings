/**
 * API клиент для Unified Listings (C2C + B2C объединенные)
 *
 * Этот сервис предоставляет методы для работы с unified API,
 * который объединяет C2C listings и B2C products без дублирования данных.
 *
 * ВАЖНО: Все запросы идут через BFF proxy /api/v2
 */

import { apiClient } from './api-client';
import type {
  UnifiedListing,
  UnifiedListingsFilters,
  UnifiedListingsResponse,
  GetUnifiedListingByIdParams,
} from '@/types/unified-listing';

export class UnifiedListingsService {
  /**
   * Получить unified listings (C2C + B2C) с фильтрами
   *
   * @example
   * // Получить все listings
   * const all = await unifiedListingsApi.getListings();
   *
   * // Только C2C
   * const c2c = await unifiedListingsApi.getListings({ source_type: 'c2c' });
   *
   * // Только B2C
   * const b2c = await unifiedListingsApi.getListings({ source_type: 'b2c' });
   *
   * // С фильтрами
   * const filtered = await unifiedListingsApi.getListings({
   *   source_type: 'all',
   *   category_id: 10,
   *   min_price: 100,
   *   max_price: 1000,
   *   limit: 20,
   *   offset: 0
   * });
   */
  async getListings(
    filters?: UnifiedListingsFilters
  ): Promise<UnifiedListingsResponse> {
    const params = new URLSearchParams();

    if (filters?.source_type) {
      params.append('source_type', filters.source_type);
    }
    if (filters?.category_id) {
      params.append('category_id', String(filters.category_id));
    }
    if (filters?.min_price !== undefined) {
      params.append('min_price', String(filters.min_price));
    }
    if (filters?.max_price !== undefined) {
      params.append('max_price', String(filters.max_price));
    }
    if (filters?.condition) {
      params.append('condition', filters.condition);
    }
    if (filters?.query) {
      params.append('query', filters.query);
    }
    if (filters?.storefront_id) {
      params.append('storefront_id', String(filters.storefront_id));
    }
    if (filters?.limit !== undefined) {
      params.append('limit', String(filters.limit));
    }
    if (filters?.offset !== undefined) {
      params.append('offset', String(filters.offset));
    }

    const queryString = params.toString();
    const endpoint = `/unified/listings${queryString ? `?${queryString}` : ''}`;

    const response = await apiClient.get<UnifiedListingsResponse>(endpoint);

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!;
  }

  /**
   * Получить конкретный unified listing по ID
   *
   * @param id - ID листинга/продукта
   * @param sourceType - Тип источника ('c2c' или 'b2c')
   *
   * @example
   * const c2cListing = await unifiedListingsApi.getListingById({ id: 1067, source_type: 'c2c' });
   * const b2cProduct = await unifiedListingsApi.getListingById({ id: 1061, source_type: 'b2c' });
   */
  async getListingById(
    params: GetUnifiedListingByIdParams
  ): Promise<UnifiedListing> {
    const { id, source_type } = params;
    const endpoint = `/unified/listings/${id}?source_type=${source_type}`;

    const response = await apiClient.get<{ data: UnifiedListing }>(endpoint);

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!.data;
  }

  /**
   * Получить несколько unified listings по массиву ID
   *
   * @param ids - Массив пар {id, source_type}
   *
   * @example
   * const listings = await unifiedListingsApi.getListingsByIds([
   *   { id: 1067, source_type: 'c2c' },
   *   { id: 1061, source_type: 'b2c' }
   * ]);
   */
  async getListingsByIds(
    ids: GetUnifiedListingByIdParams[]
  ): Promise<UnifiedListing[]> {
    const promises = ids.map((params) => this.getListingById(params));
    return Promise.all(promises);
  }

  /**
   * Получить unified listings для конкретной категории
   *
   * @example
   * const electronicsListings = await unifiedListingsApi.getListingsByCategory(10, 'all');
   * const c2cElectronics = await unifiedListingsApi.getListingsByCategory(10, 'c2c');
   */
  async getListingsByCategory(
    categoryId: number,
    sourceType: 'all' | 'c2c' | 'b2c' = 'all',
    limit: number = 20,
    offset: number = 0
  ): Promise<UnifiedListingsResponse> {
    return this.getListings({
      category_id: categoryId,
      source_type: sourceType,
      limit,
      offset,
    });
  }

  /**
   * Получить unified listings для конкретной витрины (только B2C)
   *
   * @example
   * const storefrontProducts = await unifiedListingsApi.getListingsByStorefront(43);
   */
  async getListingsByStorefront(
    storefrontId: number,
    limit: number = 20,
    offset: number = 0
  ): Promise<UnifiedListingsResponse> {
    return this.getListings({
      source_type: 'b2c',
      storefront_id: storefrontId,
      limit,
      offset,
    });
  }

  /**
   * Поиск unified listings по ключевому слову
   *
   * @example
   * const results = await unifiedListingsApi.search('батарейка Nokia', 'all');
   */
  async search(
    query: string,
    sourceType: 'all' | 'c2c' | 'b2c' = 'all',
    limit: number = 20,
    offset: number = 0
  ): Promise<UnifiedListingsResponse> {
    return this.getListings({
      query,
      source_type: sourceType,
      limit,
      offset,
    });
  }

  /**
   * Получить unified listings в ценовом диапазоне
   *
   * @example
   * const affordable = await unifiedListingsApi.getListingsByPriceRange(100, 1000);
   */
  async getListingsByPriceRange(
    minPrice: number,
    maxPrice: number,
    sourceType: 'all' | 'c2c' | 'b2c' = 'all',
    limit: number = 20,
    offset: number = 0
  ): Promise<UnifiedListingsResponse> {
    return this.getListings({
      min_price: minPrice,
      max_price: maxPrice,
      source_type: sourceType,
      limit,
      offset,
    });
  }

  /**
   * Получить unified listings по состоянию
   *
   * @example
   * const newItems = await unifiedListingsApi.getListingsByCondition('new');
   * const usedItems = await unifiedListingsApi.getListingsByCondition('used');
   */
  async getListingsByCondition(
    condition: string,
    sourceType: 'all' | 'c2c' | 'b2c' = 'all',
    limit: number = 20,
    offset: number = 0
  ): Promise<UnifiedListingsResponse> {
    return this.getListings({
      condition,
      source_type: sourceType,
      limit,
      offset,
    });
  }
}

// Экспортируем singleton instance
export const unifiedListingsApi = new UnifiedListingsService();

// For named import compatibility
export default unifiedListingsApi;
