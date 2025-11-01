/**
 * Unified Marketplace API Service
 *
 * Объединяет функциональность C2C и B2C marketplace в единый API client.
 * Все запросы идут через BFF proxy /api/v2 → backend /api/v1
 *
 * ВАЖНО: НЕ используйте прямые обращения к backend!
 * Всегда используйте apiClient, который автоматически:
 * - Добавляет JWT токены из httpOnly cookies
 * - Обрабатывает CORS
 * - Управляет ошибками
 */

import { apiClient } from './api-client';
import type {
  MarketplaceSearchParams,
  MarketplaceSearchResponse,
} from '@/types/c2c';

// ============================================================================
// TYPES - Categories & Attributes
// ============================================================================

export interface CategoryResponse {
  data: Array<{
    id: number;
    name: string;
    slug?: string;
    parent_id?: number | null;
    icon?: string;
    translations?: Record<string, string>;
    level?: number;
    count?: number;
  }>;
  success: boolean;
}

export interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: 'text' | 'number' | 'select' | 'multiselect' | 'boolean';
  icon?: string;
  options?: {
    values?: string[];
    min?: number;
    max?: number;
    step?: number;
    multiselect?: boolean;
  };
  validation_rules?: any;
  is_searchable: boolean;
  is_filterable: boolean;
  is_required: boolean;
  sort_order: number;
  created_at: string;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
  custom_component?: string;
}

export interface CategoryAttributesResponse {
  success: boolean;
  data: CategoryAttribute[];
}

// ============================================================================
// TYPES - C2C Orders
// ============================================================================

export type OrderStatus =
  | 'pending'
  | 'paid'
  | 'shipped'
  | 'delivered'
  | 'completed'
  | 'disputed'
  | 'cancelled'
  | 'refunded';

export interface C2COrder {
  id: number;
  buyer_id: number;
  seller_id: number;
  listing_id: number;
  item_price: number;
  platform_fee_rate: number;
  platform_fee_amount: number;
  seller_payout_amount: number;
  payment_transaction_id?: number;
  status: OrderStatus;
  protection_period_days: number;
  protection_expires_at?: string;
  shipping_method?: string;
  tracking_number?: string;
  shipped_at?: string;
  delivered_at?: string;
  created_at: string;
  updated_at: string;
  buyer?: any;
  seller?: any;
  listing?: any;
  payment_transaction?: any;
  status_history?: OrderStatusHistory[];
  messages?: OrderMessage[];
}

export interface OrderStatusHistory {
  id: number;
  order_id: number;
  old_status?: string;
  new_status: string;
  reason?: string;
  created_by?: number;
  created_at: string;
}

export interface OrderMessage {
  id: number;
  order_id: number;
  sender_id: number;
  message_type:
    | 'text'
    | 'shipping_update'
    | 'dispute_opened'
    | 'dispute_message'
    | 'system';
  content: string;
  metadata?: Record<string, any>;
  created_at: string;
  sender?: any;
}

export interface CreateOrderRequest {
  listing_id: number;
  message?: string;
  payment_method: string;
}

export interface CreateOrderResponse {
  order_id: number;
  payment_url: string;
}

export interface OrderListResponse {
  orders: C2COrder[];
  total: number;
  page: number;
  limit: number;
}

// ============================================================================
// MARKETPLACE API SERVICE
// ============================================================================

export const marketplaceApi = {
  // --------------------------------------------------------------------------
  // SEARCH & LISTINGS
  // --------------------------------------------------------------------------

  /**
   * Поиск объявлений с фильтрами
   * @param params - параметры поиска (категория, цена, локация и т.д.)
   */
  async search(
    params: MarketplaceSearchParams = {}
  ): Promise<MarketplaceSearchResponse> {
    const searchParams = new URLSearchParams();

    if (params.sort_by) searchParams.append('sort_by', params.sort_by);
    if (params.page !== undefined)
      searchParams.append('page', params.page.toString());
    if (params.size !== undefined)
      searchParams.append('size', params.size.toString());
    if (params.category) searchParams.append('category', params.category);
    if (params.minPrice !== undefined)
      searchParams.append('minPrice', params.minPrice.toString());
    if (params.maxPrice !== undefined)
      searchParams.append('maxPrice', params.maxPrice.toString());
    if (params.search) searchParams.append('search', params.search);
    if (params.locale) searchParams.append('lang', params.locale);

    const queryString = searchParams.toString();
    const endpoint = `/marketplace/search${queryString ? `?${queryString}` : ''}`;

    try {
      // Добавляем таймаут для SSR запросов
      const response = await apiClient.get<MarketplaceSearchResponse>(
        endpoint,
        {
          timeout: 5000, // 5 секунд
        }
      );

      if (response.error) {
        console.error('Marketplace search error:', response.error);
        // Возвращаем пустой результат при ошибке
        return {
          data: [],
          meta: {
            total: 0,
            page: 1,
            size: 20,
            has_more: false,
            total_pages: 0,
          },
        };
      }

      return response.data!;
    } catch (error) {
      console.error('Error fetching marketplace items:', error);
      // Возвращаем пустой результат при таймауте/ошибке
      return {
        data: [],
        meta: {
          total: 0,
          page: 1,
          size: 20,
          has_more: false,
          total_pages: 0,
        },
      };
    }
  },

  // --------------------------------------------------------------------------
  // CATEGORIES
  // --------------------------------------------------------------------------

  /**
   * Получить все категории
   * @param locale - язык для переводов (sr, en, ru)
   */
  async getCategories(locale?: string): Promise<CategoryResponse> {
    const params = new URLSearchParams();
    if (locale) {
      params.append('lang', locale);
    }

    const queryString = params.toString();
    const endpoint = `/marketplace/categories${queryString ? `?${queryString}` : ''}`;

    const response = await apiClient.get<CategoryResponse>(endpoint);

    if (response.error) {
      throw new Error(`Failed to fetch categories: ${response.error.message}`);
    }

    return response.data!;
  },

  /**
   * Получить атрибуты категории
   * @param categoryId - ID категории
   */
  async getCategoryAttributes(
    categoryId: number
  ): Promise<CategoryAttributesResponse> {
    const response = await apiClient.get<CategoryAttributesResponse>(
      `/marketplace/categories/${categoryId}/attributes`
    );

    if (response.error) {
      throw new Error(
        `Failed to fetch category attributes: ${response.error.message}`
      );
    }

    return response.data!;
  },

  // --------------------------------------------------------------------------
  // C2C ORDERS
  // --------------------------------------------------------------------------

  /**
   * Создать заказ C2C
   */
  async createOrder(data: CreateOrderRequest): Promise<CreateOrderResponse> {
    const response = await apiClient.post<CreateOrderResponse>(
      '/marketplace/orders/create',
      data
    );

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!;
  },

  /**
   * Получить мои покупки
   */
  async getMyPurchases(page = 1, limit = 20): Promise<OrderListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });

    const response = await apiClient.get<OrderListResponse>(
      `/marketplace/orders/my/purchases?${params.toString()}`
    );

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!;
  },

  /**
   * Получить мои продажи
   */
  async getMySales(page = 1, limit = 20): Promise<OrderListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });

    const response = await apiClient.get<OrderListResponse>(
      `/marketplace/orders/my/sales?${params.toString()}`
    );

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!;
  },

  /**
   * Получить детали заказа
   */
  async getOrderDetails(orderId: number): Promise<C2COrder> {
    const response = await apiClient.get<C2COrder>(
      `/marketplace/orders/${orderId}`
    );

    if (response.error) {
      throw new Error(response.error.message);
    }

    return response.data!;
  },

  /**
   * Отметить заказ как отправленный
   */
  async markAsShipped(orderId: number, shippingMethod: string): Promise<void> {
    const response = await apiClient.post(
      `/marketplace/orders/${orderId}/ship`,
      {
        shipping_method: shippingMethod,
      }
    );

    if (response.error) {
      throw new Error(response.error.message);
    }
  },

  /**
   * Подтвердить доставку заказа
   */
  async confirmDelivery(orderId: number): Promise<void> {
    const response = await apiClient.post(
      `/marketplace/orders/${orderId}/confirm-delivery`
    );

    if (response.error) {
      throw new Error(response.error.message);
    }
  },

  /**
   * Открыть спор по заказу
   */
  async openDispute(orderId: number, reason: string): Promise<void> {
    const response = await apiClient.post(
      `/marketplace/orders/${orderId}/dispute`,
      {
        reason,
      }
    );

    if (response.error) {
      throw new Error(response.error.message);
    }
  },

  /**
   * Добавить сообщение к заказу
   */
  async addMessage(orderId: number, content: string): Promise<void> {
    const response = await apiClient.post(
      `/marketplace/orders/${orderId}/message`,
      {
        content,
      }
    );

    if (response.error) {
      throw new Error(response.error.message);
    }
  },

  // --------------------------------------------------------------------------
  // ORDER HELPERS
  // --------------------------------------------------------------------------

  /**
   * Получить локализованный статус заказа
   */
  getStatusLabel(status: OrderStatus, locale: string): string {
    const labels: Record<OrderStatus, Record<string, string>> = {
      pending: {
        ru: 'Ожидает оплаты',
        en: 'Pending payment',
        sr: 'Чека на уплату',
      },
      paid: { ru: 'Оплачен', en: 'Paid', sr: 'Плаћено' },
      shipped: { ru: 'Отправлен', en: 'Shipped', sr: 'Послато' },
      delivered: { ru: 'Доставлен', en: 'Delivered', sr: 'Достављено' },
      completed: { ru: 'Завершен', en: 'Completed', sr: 'Завршено' },
      disputed: { ru: 'Спор', en: 'Disputed', sr: 'Спор' },
      cancelled: { ru: 'Отменен', en: 'Cancelled', sr: 'Отказано' },
      refunded: { ru: 'Возвращен', en: 'Refunded', sr: 'Враћено' },
    };

    return labels[status]?.[locale] || status;
  },

  /**
   * Получить CSS класс для статуса заказа
   */
  getStatusColor(status: OrderStatus): string {
    const colors: Record<OrderStatus, string> = {
      pending: 'badge-warning',
      paid: 'badge-info',
      shipped: 'badge-info',
      delivered: 'badge-primary',
      completed: 'badge-success',
      disputed: 'badge-error',
      cancelled: 'badge-ghost',
      refunded: 'badge-ghost',
    };

    return colors[status] || 'badge-ghost';
  },

  /**
   * Проверить можно ли отправить заказ
   */
  canShip(order: C2COrder, userId: number): boolean {
    return order.status === 'paid' && order.seller_id === userId;
  },

  /**
   * Проверить можно ли подтвердить доставку
   */
  canConfirmDelivery(order: C2COrder, userId: number): boolean {
    return order.status === 'shipped' && order.buyer_id === userId;
  },

  /**
   * Проверить можно ли открыть спор
   */
  canOpenDispute(order: C2COrder, userId: number): boolean {
    return (
      ['shipped', 'delivered'].includes(order.status) &&
      (order.buyer_id === userId || order.seller_id === userId)
    );
  },

  /**
   * Рассчитать оставшееся время защиты покупателя
   */
  getProtectionTimeLeft(order: C2COrder): string | null {
    if (!order.protection_expires_at) return null;

    const expiresAt = new Date(order.protection_expires_at);
    const now = new Date();
    const diff = expiresAt.getTime() - now.getTime();

    if (diff <= 0) return null;

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

    if (days > 0) {
      return `${days} дн. ${hours} ч.`;
    } else {
      return `${hours} ч.`;
    }
  },
};

/**
 * Backward compatibility - default export
 */
export default marketplaceApi;

/**
 * Legacy service class for backward compatibility
 * @deprecated Use marketplaceApi instead
 */
export class MarketplaceService {
  static async search(
    params: MarketplaceSearchParams = {}
  ): Promise<MarketplaceSearchResponse> {
    return marketplaceApi.search(params);
  }

  static async getCategories(locale?: string): Promise<CategoryResponse> {
    return marketplaceApi.getCategories(locale);
  }

  static async getCategoryAttributes(
    categoryId: number
  ): Promise<CategoryAttributesResponse> {
    return marketplaceApi.getCategoryAttributes(categoryId);
  }
}

/**
 * Legacy orders service for backward compatibility
 * @deprecated Use marketplaceApi instead
 */
export const marketplaceOrdersService = {
  createOrder: marketplaceApi.createOrder,
  getMyPurchases: marketplaceApi.getMyPurchases,
  getMySales: marketplaceApi.getMySales,
  getOrderDetails: marketplaceApi.getOrderDetails,
  markAsShipped: marketplaceApi.markAsShipped,
  confirmDelivery: marketplaceApi.confirmDelivery,
  openDispute: marketplaceApi.openDispute,
  addMessage: marketplaceApi.addMessage,
  getStatusLabel: marketplaceApi.getStatusLabel,
  getStatusColor: marketplaceApi.getStatusColor,
  canShip: marketplaceApi.canShip,
  canConfirmDelivery: marketplaceApi.canConfirmDelivery,
  canOpenDispute: marketplaceApi.canOpenDispute,
  getProtectionTimeLeft: marketplaceApi.getProtectionTimeLeft,
};
