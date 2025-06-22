import { config } from '@/config';
import type { components } from '@/types/generated/api';

// Типы из сгенерированного API
type Storefront = components['schemas']['models.Storefront'];
type StorefrontCreateDTO = components['schemas']['models.StorefrontCreateDTO'];
type StorefrontUpdateDTO = components['schemas']['models.StorefrontUpdateDTO'];
type StorefrontAnalytics = components['schemas']['models.StorefrontAnalytics'];
type StorefrontRatingSummary =
  components['schemas']['models.StorefrontRatingSummary'];
type StorefrontMapData = components['schemas']['models.StorefrontMapData'];
// TODO: Fix these types after API regeneration - import from storefront.ts for now
import type {
  StorefrontSearchResult,
  StorefrontsListResponse,
} from '@/types/storefront';
type StorefrontHours = components['schemas']['models.StorefrontHours'];
type StorefrontPaymentMethod =
  components['schemas']['models.StorefrontPaymentMethod'];
type StorefrontDeliveryOption =
  components['schemas']['models.StorefrontDeliveryOption'];

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

export interface FileUploadParams {
  storefrontId: number;
  file: File;
  type: 'logo' | 'banner';
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
   * Создает URL с query параметрами
   */
  private createUrl(endpoint: string, params?: Record<string, any>): string {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            // Для массивов добавляем каждый элемент отдельно
            value.forEach((item) => url.searchParams.append(key, String(item)));
          } else {
            url.searchParams.append(key, String(value));
          }
        }
      });
    }

    return url.toString();
  }

  /**
   * Выполняет HTTP запрос с обработкой ошибок
   */
  private async request<T>(url: string, options?: RequestInit): Promise<T> {
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: `HTTP error! status: ${response.status}`,
        }));
        throw new Error(
          errorData.error || errorData.message || 'Network error'
        );
      }

      // Проверяем есть ли контент для парсинга
      const contentType = response.headers.get('content-type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }

      // Возвращаем пустой объект для успешных запросов без контента
      return {} as T;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Unknown error occurred');
    }
  }

  // Основные CRUD операции

  /**
   * Получение списка витрин с фильтрами и пагинацией
   */
  async getStorefronts(
    params?: StorefrontSearchParams
  ): Promise<StorefrontsListResponse> {
    const url = this.createUrl('/api/v1/storefronts', params);
    return this.request<StorefrontsListResponse>(url);
  }

  /**
   * Получение витрины по ID
   */
  async getStorefrontById(id: number): Promise<Storefront> {
    const url = this.createUrl(`/api/v1/storefronts/${id}`);
    return this.request<Storefront>(url);
  }

  /**
   * Получение витрины по slug
   */
  async getStorefrontBySlug(slug: string): Promise<Storefront> {
    const url = this.createUrl(`/api/v1/storefronts/slug/${slug}`);
    return this.request<Storefront>(url);
  }

  /**
   * Получение витрин текущего пользователя
   */
  async getMyStorefronts(): Promise<StorefrontsListResponse> {
    const url = this.createUrl('/api/v1/storefronts/my');
    return this.request<StorefrontsListResponse>(url);
  }

  /**
   * Создание новой витрины
   */
  async createStorefront(data: StorefrontCreateDTO): Promise<Storefront> {
    const url = this.createUrl('/api/v1/storefronts');
    return this.request<Storefront>(url, {
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
    const url = this.createUrl(`/api/v1/storefronts/${id}`);
    return this.request<Storefront>(url, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Удаление витрины
   */
  async deleteStorefront(id: number): Promise<void> {
    const url = this.createUrl(`/api/v1/storefronts/${id}`);
    return this.request<void>(url, {
      method: 'DELETE',
    });
  }

  // Поиск и фильтрация

  /**
   * Поиск витрин через OpenSearch
   */
  async searchStorefronts(
    params: StorefrontSearchParams
  ): Promise<StorefrontSearchResult> {
    const url = this.createUrl('/api/v1/storefronts/search/opensearch');
    return this.request<StorefrontSearchResult>(url, {
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
    const url = this.createUrl('/api/v1/storefronts/map', params);
    return this.request<StorefrontMapData[]>(url);
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
    const url = this.createUrl('/api/v1/storefronts/nearby', params);
    return this.request<Storefront[]>(url);
  }

  // Аналитика и статистика

  /**
   * Получение аналитики витрины
   */
  async getStorefrontAnalytics(
    params: AnalyticsParams
  ): Promise<StorefrontAnalytics> {
    const { storefrontId, ...queryParams } = params;
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/analytics`,
      queryParams
    );
    return this.request<StorefrontAnalytics>(url);
  }

  /**
   * Получение рейтинга витрины
   */
  async getStorefrontRating(id: number): Promise<StorefrontRatingSummary> {
    const url = this.createUrl(`/api/v1/storefronts/${id}/rating`);
    return this.request<StorefrontRatingSummary>(url);
  }

  // Настройки витрины

  /**
   * Получение часов работы витрины
   */
  async getWorkingHours(storefrontId: number): Promise<StorefrontHours[]> {
    const url = this.createUrl(`/api/v1/storefronts/${storefrontId}/hours`);
    return this.request<StorefrontHours[]>(url);
  }

  /**
   * Обновление часов работы витрины
   */
  async updateWorkingHours(
    storefrontId: number,
    hours: StorefrontHours[]
  ): Promise<void> {
    const url = this.createUrl(`/api/v1/storefronts/${storefrontId}/hours`);
    return this.request<void>(url, {
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
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/payment-methods`
    );
    return this.request<StorefrontPaymentMethod[]>(url);
  }

  /**
   * Обновление методов оплаты витрины
   */
  async updatePaymentMethods(
    storefrontId: number,
    methods: StorefrontPaymentMethod[]
  ): Promise<void> {
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/payment-methods`
    );
    return this.request<void>(url, {
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
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/delivery-options`
    );
    return this.request<StorefrontDeliveryOption[]>(url);
  }

  /**
   * Обновление опций доставки витрины
   */
  async updateDeliveryOptions(
    storefrontId: number,
    options: StorefrontDeliveryOption[]
  ): Promise<void> {
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/delivery-options`
    );
    return this.request<void>(url, {
      method: 'PUT',
      body: JSON.stringify({ options }),
    });
  }

  // Загрузка файлов

  /**
   * Загрузка логотипа или баннера витрины
   */
  async uploadFile(params: FileUploadParams): Promise<{ url: string }> {
    const { storefrontId, file, type } = params;

    const formData = new FormData();
    formData.append('file', file);

    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/upload/${type}`
    );

    return this.request<{ url: string }>(url, {
      method: 'POST',
      body: formData,
      headers: {}, // Убираем Content-Type для FormData
    });
  }

  /**
   * Загрузка логотипа витрины
   */
  async uploadLogo(storefrontId: number, file: File): Promise<{ url: string }> {
    return this.uploadFile({ storefrontId, file, type: 'logo' });
  }

  /**
   * Загрузка баннера витрины
   */
  async uploadBanner(
    storefrontId: number,
    file: File
  ): Promise<{ url: string }> {
    return this.uploadFile({ storefrontId, file, type: 'banner' });
  }

  // Управление персоналом

  /**
   * Получение сотрудников витрины
   */
  async getStaff(storefrontId: number): Promise<any[]> {
    const url = this.createUrl(`/api/v1/storefronts/${storefrontId}/staff`);
    return this.request<any[]>(url);
  }

  /**
   * Добавление сотрудника
   */
  async addStaff(
    storefrontId: number,
    userId: number,
    role: string
  ): Promise<void> {
    const url = this.createUrl(`/api/v1/storefronts/${storefrontId}/staff`);
    return this.request<void>(url, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, role }),
    });
  }

  /**
   * Удаление сотрудника
   */
  async removeStaff(storefrontId: number, userId: number): Promise<void> {
    const url = this.createUrl(
      `/api/v1/storefronts/${storefrontId}/staff/${userId}`
    );
    return this.request<void>(url, {
      method: 'DELETE',
    });
  }

  // Запись просмотра

  /**
   * Записать просмотр витрины
   */
  async recordView(storefrontId: number): Promise<void> {
    const url = this.createUrl(`/api/v1/storefronts/${storefrontId}/view`);
    return this.request<void>(url, {
      method: 'POST',
    });
  }
}

// Создаем и экспортируем единственный экземпляр сервиса
export const storefrontApi = new StorefrontApiService();
export default storefrontApi;
