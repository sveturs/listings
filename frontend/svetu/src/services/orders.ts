import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

type StorefrontOrder = components['schemas']['models.StorefrontOrder'];
type CreateOrderRequest = components['schemas']['models.CreateOrderRequest'];
type CancelOrderRequest = components['schemas']['models.CancelOrderRequest'];

export const ordersService = {
  // Создать заказ
  async createOrder(orderData: CreateOrderRequest): Promise<StorefrontOrder> {
    console.log('[ordersService] Creating order with data:', orderData);
    const response = await apiClient.post('/api/v1/orders', orderData);
    if (response.error) {
      console.error('[ordersService] Order creation failed:', response.error);
      throw new Error(response.error.message || 'Failed to create order');
    }
    console.log('[ordersService] Order created successfully:', response.data);
    return response.data.data;
  },

  // Получить заказ по ID
  async getOrder(orderId: number): Promise<StorefrontOrder> {
    const response = await apiClient.get(`/api/v1/orders/${orderId}`);
    if (!response.error && response.data?.data) {
      return response.data.data;
    }
    throw new Error(response.error?.message || 'Order not found');
  },

  // Получить список заказов пользователя (покупки)
  async getUserOrders(params?: {
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ orders: StorefrontOrder[]; total: number }> {
    // Используем marketplace endpoint для получения заказов где пользователь - покупатель
    const page = params?.offset
      ? Math.floor(params.offset / (params.limit || 20)) + 1
      : 1;
    const limit = params?.limit || 20;

    const searchParams = new URLSearchParams();
    searchParams.append('page', page.toString());
    searchParams.append('limit', limit.toString());
    if (params?.status) searchParams.append('status', params.status);

    const response = await apiClient.get(`/api/v1/orders?${searchParams}`);
    console.log('[ordersService] getUserOrders response:', response);

    if (!response.error && response.data?.data) {
      const { orders = [], total = 0 } = response.data.data;
      return {
        orders: Array.isArray(orders) ? orders : [],
        total: total || 0,
      };
    }

    return {
      orders: [],
      total: 0,
    };
  },

  // Отменить заказ
  async cancelOrder(
    orderId: number,
    reason?: string
  ): Promise<StorefrontOrder> {
    const cancelData: CancelOrderRequest = { reason };
    const response = await apiClient.put(
      `/api/v1/orders/${orderId}/cancel`,
      cancelData
    );
    if (response.error) {
      throw new Error(response.error.message || 'Failed to cancel order');
    }
    return response.data.data;
  },
};
