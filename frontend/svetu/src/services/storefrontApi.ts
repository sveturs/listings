import { config } from '@/config';
import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

// Типы из сгенерированного API
type Storefront =
  components['schemas']['backend_internal_domain_models.Storefront'];
type StorefrontCreateDTO =
  components['schemas']['backend_internal_domain_models.StorefrontCreateDTO'];
type StorefrontUpdateDTO =
  components['schemas']['backend_internal_domain_models.StorefrontUpdateDTO'];
type StorefrontAnalytics =
  components['schemas']['backend_internal_domain_models.StorefrontAnalytics'];
type StorefrontRatingSummary =
  components['schemas']['backend_internal_domain_models.StorefrontRatingSummary'];
type StorefrontMapData =
  components['schemas']['backend_internal_domain_models.StorefrontMapData'];
type StorefrontSearchResult =
  components['schemas']['backend_internal_proj_storefronts_storage_opensearch.StorefrontSearchResult'];
type StorefrontsListResponse =
  components['schemas']['internal_proj_storefronts_handler.StorefrontsListResponse'];
type StorefrontHours =
  components['schemas']['backend_internal_domain_models.StorefrontHours'];
type StorefrontPaymentMethod =
  components['schemas']['backend_internal_domain_models.StorefrontPaymentMethod'];
type StorefrontDeliveryOption =
  components['schemas']['backend_internal_domain_models.StorefrontDeliveryOption'];

// Интерфейсы для параметров запросов
export interface StorefrontSearchParams {
  search?: string;
  city?: string;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
  minRating?: number;
  isActive?: boolean;
  isVerified?: boolean;
  isOpenNow?: boolean;
  paymentMethods?: string[];
  hasDelivery?: boolean;
  hasSelfPickup?: boolean;
  sortBy?: string;
  sortOrder?: string;
  limit?: number;
  offset?: number;
}

export interface GeoBounds {
  north: number;
  south: number;
  east: number;
  west: number;
}

export interface AnalyticsParams {
  storefrontId: number;
  from?: string;
  to?: string;
}


/**
 * API сервис для работы с витринами
 */
class StorefrontApiService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = config.api.url;
  }

  /**
   * Создает endpoint с query параметрами
   */
  private createEndpoint(
    endpoint: string,
    params?: Record<string, any>
  ): string {
    if (!params) {
      return endpoint;
    }

    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          // Для массивов добавляем каждый элемент отдельно
          value.forEach((item) => searchParams.append(key, String(item)));
        } else {
          searchParams.append(key, String(value));
        }
      }
    });

    const queryString = searchParams.toString();
    return queryString ? `${endpoint}?${queryString}` : endpoint;
  }

  /**
   * Выполняет HTTP запрос с обработкой ошибок используя apiClient
   */
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const response = await apiClient.request<T>(endpoint, {
      ...options,
      includeAuth: true, // Всегда включаем авторизацию для storefrontApi
    });

    if (response.error) {
      throw new Error(response.error.message || 'Request failed');
    }

    if (response.data !== undefined) {
      return response.data;
    }

    // Возвращаем пустой объект для успешных запросов без контента
    return {} as T;
  }

  // Основные CRUD операции

  /**
   * Получение списка витрин с фильтрами и пагинацией
   */
  async getStorefronts(
    params?: StorefrontSearchParams
  ): Promise<StorefrontsListResponse> {
    const endpoint = this.createEndpoint('/api/v1/storefronts', params);
    return this.request<StorefrontsListResponse>(endpoint);
  }

  /**
   * Получение витрины по ID
   */
  async getStorefrontById(id: number): Promise<Storefront> {
    const endpoint = `/api/v1/storefronts/${id}`;
    return this.request<Storefront>(endpoint);
  }

  /**
   * Получение витрины по slug
   */
  async getStorefrontBySlug(slug: string): Promise<Storefront> {
    const endpoint = `/api/v1/storefronts/slug/${slug}`;
    return this.request<Storefront>(endpoint);
  }

  /**
   * Получение витрин текущего пользователя
   */
  async getMyStorefronts(): Promise<StorefrontsListResponse> {
    const endpoint = '/api/v1/storefronts/my';
    return this.request<StorefrontsListResponse>(endpoint);
  }

  /**
   * Создание новой витрины
   */
  async createStorefront(data: StorefrontCreateDTO): Promise<Storefront> {
    const endpoint = '/api/v1/storefronts';
    return this.request<Storefront>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * Обновление витрины
   */
  async updateStorefront(
    id: number,
    data: StorefrontUpdateDTO
  ): Promise<Storefront> {
    const endpoint = `/api/v1/storefronts/${id}`;
    return this.request<Storefront>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Удаление витрины
   */
  async deleteStorefront(id: number): Promise<void> {
    const endpoint = `/api/v1/storefronts/${id}`;
    return this.request<void>(endpoint, {
      method: 'DELETE',
    });
  }

  /**
   * Загрузка логотипа витрины
   */
  async uploadLogo(storefrontId: number, file: File): Promise<{ url: string }> {
    const formData = new FormData();
    formData.append('logo', file);

    const response = await fetch(
      `${config.api.baseUrl}/api/v1/storefronts/${storefrontId}/logo`,
      {
        method: 'POST',
        body: formData,
        credentials: 'include',
      }
    );

    if (!response.ok) {
      throw new Error('Failed to upload logo');
    }

    const data = await response.json();
    return { url: data.logo_url || data.url };
  }

  /**
   * Загрузка баннера витрины
   */
  async uploadBanner(
    storefrontId: number,
    file: File
  ): Promise<{ url: string }> {
    const formData = new FormData();
    formData.append('banner', file);

    const response = await fetch(
      `${config.api.baseUrl}/api/v1/storefronts/${storefrontId}/banner`,
      {
        method: 'POST',
        body: formData,
        credentials: 'include',
      }
    );

    if (!response.ok) {
      throw new Error('Failed to upload banner');
    }

    const data = await response.json();
    return { url: data.banner_url || data.url };
  }

  // Поиск и фильтрация

  /**
   * Поиск витрин через OpenSearch
   */
  async searchStorefronts(
    params: StorefrontSearchParams
  ): Promise<StorefrontSearchResult> {
    const endpoint = '/api/v1/storefronts/search/opensearch';
    return this.request<StorefrontSearchResult>(endpoint, {
      method: 'POST',
      body: JSON.stringify(params),
    });
  }

  /**
   * Получение витрин для карты в заданных границах
   */
  async getStorefrontsForMap(bounds: GeoBounds): Promise<StorefrontMapData[]> {
    const params = {
      north: bounds.north,
      south: bounds.south,
      east: bounds.east,
      west: bounds.west,
    };
    const endpoint = this.createEndpoint('/api/v1/storefronts/map', params);
    return this.request<StorefrontMapData[]>(endpoint);
  }

  /**
   * Получение витрин поблизости
   */
  async getNearbyStorefronts(
    latitude: number,
    longitude: number,
    radiusKm: number,
    limit?: number
  ): Promise<Storefront[]> {
    const params = { latitude, longitude, radius_km: radiusKm, limit };
    const endpoint = this.createEndpoint('/api/v1/storefronts/nearby', params);
    return this.request<Storefront[]>(endpoint);
  }

  // Аналитика и статистика

  /**
   * Получение аналитики витрины
   */
  async getStorefrontAnalytics(
    params: AnalyticsParams
  ): Promise<StorefrontAnalytics> {
    const { storefrontId, ...queryParams } = params;
    const endpoint = this.createEndpoint(
      `/api/v1/storefronts/${storefrontId}/analytics`,
      queryParams
    );
    return this.request<StorefrontAnalytics>(endpoint);
  }

  /**
   * Получение рейтинга витрины
   */
  async getStorefrontRating(id: number): Promise<StorefrontRatingSummary> {
    const endpoint = `/api/v1/storefronts/${id}/rating`;
    return this.request<StorefrontRatingSummary>(endpoint);
  }

  // Настройки витрины

  /**
   * Получение часов работы витрины
   */
  async getWorkingHours(storefrontId: number): Promise<StorefrontHours[]> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/hours`;
    return this.request<StorefrontHours[]>(endpoint);
  }

  /**
   * Обновление часов работы витрины
   */
  async updateWorkingHours(
    storefrontId: number,
    hours: StorefrontHours[]
  ): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/hours`;
    return this.request<void>(endpoint, {
      method: 'PUT',
      body: JSON.stringify({ hours }),
    });
  }

  /**
   * Получение методов оплаты витрины
   */
  async getPaymentMethods(
    storefrontId: number
  ): Promise<StorefrontPaymentMethod[]> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/payment-methods`;
    return this.request<StorefrontPaymentMethod[]>(endpoint);
  }

  /**
   * Обновление методов оплаты витрины
   */
  async updatePaymentMethods(
    storefrontId: number,
    methods: StorefrontPaymentMethod[]
  ): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/payment-methods`;
    return this.request<void>(endpoint, {
      method: 'PUT',
      body: JSON.stringify({ methods }),
    });
  }

  /**
   * Получение опций доставки витрины
   */
  async getDeliveryOptions(
    storefrontId: number
  ): Promise<StorefrontDeliveryOption[]> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/delivery-options`;
    return this.request<StorefrontDeliveryOption[]>(endpoint);
  }

  /**
   * Обновление опций доставки витрины
   */
  async updateDeliveryOptions(
    storefrontId: number,
    options: StorefrontDeliveryOption[]
  ): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/delivery-options`;
    return this.request<void>(endpoint, {
      method: 'PUT',
      body: JSON.stringify({ options }),
    });
  }


  // Управление персоналом

  /**
   * Получение сотрудников витрины
   */
  async getStaff(storefrontId: number): Promise<any[]> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/staff`;
    return this.request<any[]>(endpoint);
  }

  /**
   * Добавление сотрудника
   */
  async addStaff(
    storefrontId: number,
    userId: number,
    role: string
  ): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/staff`;
    return this.request<void>(endpoint, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, role }),
    });
  }

  /**
   * Удаление сотрудника
   */
  async removeStaff(storefrontId: number, userId: number): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/staff/${userId}`;
    return this.request<void>(endpoint, {
      method: 'DELETE',
    });
  }

  // Запись просмотра

  /**
   * Записать просмотр витрины
   */
  async recordView(storefrontId: number): Promise<void> {
    const endpoint = `/api/v1/storefronts/${storefrontId}/view`;
    return this.request<void>(endpoint, {
      method: 'POST',
    });
  }
}

// Создаем и экспортируем единственный экземпляр сервиса
export const storefrontApi = new StorefrontApiService();
export default storefrontApi;
