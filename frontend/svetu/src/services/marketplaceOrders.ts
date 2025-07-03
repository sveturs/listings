import { apiClient } from '@/services/api-client';

// Типы для заказов маркетплейса
export interface MarketplaceOrder {
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

  // Связанные данные
  buyer?: any;
  seller?: any;
  listing?: any;
  payment_transaction?: any;
  status_history?: OrderStatusHistory[];
  messages?: OrderMessage[];
}

export type OrderStatus =
  | 'pending'
  | 'paid'
  | 'shipped'
  | 'delivered'
  | 'completed'
  | 'disputed'
  | 'cancelled'
  | 'refunded';

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
  orders: MarketplaceOrder[];
  total: number;
  page: number;
  limit: number;
}

export const marketplaceOrdersService = {
  // Создать заказ
  async createOrder(data: CreateOrderRequest): Promise<CreateOrderResponse> {
    const response = await apiClient.post(
      '/api/v1/marketplace/orders/create',
      data
    );
    return response.data?.data || response.data;
  },

  // Получить мои покупки
  async getMyPurchases(page = 1, limit = 20): Promise<OrderListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    const response = await apiClient.get(
      `/api/v1/marketplace/orders/my/purchases?${params.toString()}`
    );
    return response.data?.data || response.data;
  },

  // Получить мои продажи
  async getMySales(page = 1, limit = 20): Promise<OrderListResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    const response = await apiClient.get(
      `/api/v1/marketplace/orders/my/sales?${params.toString()}`
    );
    return response.data?.data || response.data;
  },

  // Получить детали заказа
  async getOrderDetails(orderId: number): Promise<MarketplaceOrder> {
    const response = await apiClient.get(
      `/api/v1/marketplace/orders/${orderId}`
    );
    return response.data?.data || response.data;
  },

  // Отметить как отправленный
  async markAsShipped(
    orderId: number,
    shippingMethod: string,
    trackingNumber: string
  ): Promise<void> {
    await apiClient.post(`/api/v1/marketplace/orders/${orderId}/ship`, {
      shipping_method: shippingMethod,
      tracking_number: trackingNumber,
    });
  },

  // Подтвердить доставку
  async confirmDelivery(orderId: number): Promise<void> {
    await apiClient.post(
      `/api/v1/marketplace/orders/${orderId}/confirm-delivery`
    );
  },

  // Открыть спор
  async openDispute(orderId: number, reason: string): Promise<void> {
    await apiClient.post(`/api/v1/marketplace/orders/${orderId}/dispute`, {
      reason,
    });
  },

  // Добавить сообщение
  async addMessage(orderId: number, content: string): Promise<void> {
    await apiClient.post(`/api/v1/marketplace/orders/${orderId}/message`, {
      content,
    });
  },

  // Форматировать статус заказа
  getStatusLabel(status: OrderStatus, locale: string): string {
    const labels: Record<OrderStatus, Record<string, string>> = {
      pending: { ru: 'Ожидает оплаты', en: 'Pending payment' },
      paid: { ru: 'Оплачен', en: 'Paid' },
      shipped: { ru: 'Отправлен', en: 'Shipped' },
      delivered: { ru: 'Доставлен', en: 'Delivered' },
      completed: { ru: 'Завершен', en: 'Completed' },
      disputed: { ru: 'Спор', en: 'Disputed' },
      cancelled: { ru: 'Отменен', en: 'Cancelled' },
      refunded: { ru: 'Возвращен', en: 'Refunded' },
    };

    return labels[status]?.[locale] || status;
  },

  // Получить цвет для статуса
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

  // Проверить можно ли выполнить действие
  canShip(order: MarketplaceOrder, userId: number): boolean {
    return order.status === 'paid' && order.seller_id === userId;
  },

  canConfirmDelivery(order: MarketplaceOrder, userId: number): boolean {
    return order.status === 'shipped' && order.buyer_id === userId;
  },

  canOpenDispute(order: MarketplaceOrder, userId: number): boolean {
    return (
      ['shipped', 'delivered'].includes(order.status) &&
      (order.buyer_id === userId || order.seller_id === userId)
    );
  },

  // Рассчитать оставшееся время защиты
  getProtectionTimeLeft(order: MarketplaceOrder): string | null {
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
