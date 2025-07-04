import { apiClient } from '@/services/api-client';
import type { components } from '@/types/generated/api';

type UserBalance = components['schemas']['models.UserBalance'];
type TransactionHistory = any;
type DepositResponse = components['schemas']['utils.SuccessResponseSwag'] & {
  data?: components['schemas']['models.PaymentSession'];
};

export const balanceService = {
  // Получить баланс пользователя
  async getBalance(): Promise<UserBalance> {
    const response = await apiClient.get('/api/v1/balance');
    // Backend возвращает данные в формате { success: true, data: {...} }
    return response.data?.data || response.data;
  },

  // Получить историю транзакций
  async getTransactionHistory(params?: {
    limit?: number;
    offset?: number;
    type?: string;
  }): Promise<TransactionHistory> {
    const queryParams = new URLSearchParams();
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    if (params?.type) queryParams.append('type', params.type);

    const response = await apiClient.get(
      `/api/v1/balance/transactions?${queryParams}`
    );
    // Backend возвращает данные в формате { success: true, data: {...} }
    return response.data?.data || response.data;
  },

  // Создать депозит (пополнение баланса)
  async createDeposit(data: {
    amount: number;
    currency: string;
    payment_method: string;
    return_url?: string;
  }): Promise<DepositResponse> {
    const response = await apiClient.post('/api/v1/balance/deposit', {
      ...data,
      return_url: data.return_url,
    });
    // Backend возвращает данные в формате { success: true, data: {...} }
    return response.data;
  },

  // Форматировать сумму с валютой
  formatAmount(amount: number | string, currency: string): string {
    const numAmount = typeof amount === 'string' ? parseFloat(amount) : amount;

    // Определяем формат для разных валют
    const formatOptions: Intl.NumberFormatOptions = {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    };

    const formatted = numAmount.toLocaleString('sr-RS', formatOptions);

    // Позиция символа валюты
    switch (currency) {
      case 'EUR':
        return `${formatted} €`;
      case 'USD':
        return `$${formatted}`;
      case 'RSD':
      default:
        return `${formatted} RSD`;
    }
  },

  // Получить символ валюты
  getCurrencySymbol(currency: string): string {
    switch (currency) {
      case 'EUR':
        return '€';
      case 'USD':
        return '$';
      case 'RSD':
      default:
        return 'RSD';
    }
  },
};
