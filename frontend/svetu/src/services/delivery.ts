/**
 * Delivery Service API Client
 *
 * Централизованный API wrapper для работы с delivery endpoints.
 *
 * ✅ Использует apiClient (BFF proxy /api/v2)
 * ❌ НЕ делает прямые fetch к backend
 *
 * @module services/delivery
 */

import { apiClient, ApiResponse } from './api-client';

// ========================================
// TYPE DEFINITIONS
// ========================================

/**
 * Адрес для доставки
 */
export interface Address {
  /** Имя получателя/отправителя */
  name?: string;
  /** Телефон */
  phone?: string;
  /** Email */
  email?: string;
  /** Улица и номер дома */
  street: string;
  /** Город */
  city: string;
  /** Почтовый индекс */
  postalCode: string;
  /** Код страны (ISO 3166-1 alpha-2) */
  country: string;
}

/**
 * Посылка/упаковка
 */
export interface Package {
  /** Вес в килограммах */
  weight: number;
  /** Длина в сантиметрах */
  length?: number;
  /** Ширина в сантиметрах */
  width?: number;
  /** Высота в сантиметрах */
  height?: number;
  /** Описание содержимого */
  description?: string;
  /** Объявленная ценность */
  declaredValue?: number;
  /** Наложенный платеж (COD) */
  cashOnDelivery?: boolean;
  /** Сумма наложенного платежа */
  codAmount?: number;
}

/**
 * Запрос расчета стоимости доставки
 */
export interface CalculateRateRequest {
  /** Код провайдера: post_express, bex, aks, d_express, city_express */
  provider: string;
  /** Адрес отправителя (город) */
  from_city: string;
  /** Адрес получателя (город) */
  to_city: string;
  /** Вес посылки в килограммах */
  weight: number;
  /** Длина в сантиметрах */
  length?: number;
  /** Ширина в сантиметрах */
  width?: number;
  /** Высота в сантиметрах */
  height?: number;
  /** Наложенный платеж */
  cash_on_delivery?: boolean;
  /** Сумма наложенного платежа */
  cod_amount?: number;
  /** Страхование */
  insurance?: boolean;
}

/**
 * Ответ с расчетом стоимости доставки
 */
export interface CalculateRateResponse {
  /** Базовая стоимость */
  base_price: number;
  /** Стоимость страхования */
  insurance: number;
  /** Комиссия за наложенный платеж */
  cod_fee: number;
  /** Доплата за вес */
  weight_fee: number;
  /** Доплата за расстояние */
  distance_fee: number;
  /** Итоговая стоимость */
  total_cost: number;
  /** Расчетная дата доставки (ISO 8601) */
  estimated_delivery: string;
  /** Валюта */
  currency?: string;
}

/**
 * Запрос создания отправления
 */
export interface CreateShipmentRequest {
  /** ID заказа */
  order_id: number;
  /** Код провайдера */
  provider_code: string;
  /** Адрес отправителя (полный) */
  from_address: Address;
  /** Адрес получателя (полный) */
  to_address: Address;
  /** Список посылок */
  packages: Package[];
  /** Дополнительные услуги */
  services?: ShipmentServices;
}

/**
 * Дополнительные услуги доставки
 */
export interface ShipmentServices {
  /** Страхование */
  insurance: boolean;
  /** SMS уведомление */
  sms_notification: boolean;
  /** Email уведомление */
  email_notification: boolean;
}

/**
 * Ответ при создании отправления
 */
export interface CreateShipmentResponse {
  /** ID отправления */
  shipment_id: number;
  /** Трекинг номер */
  tracking_number: string;
  /** Код провайдера */
  provider_code: string;
  /** Статус */
  status: string;
  /** Ссылка на печатную этикетку */
  label_url?: string;
}

/**
 * Информация о трекинге отправления
 */
export interface TrackingInfo {
  /** ID отправления */
  shipment_id: number;
  /** Трекинг номер */
  tracking_number: string;
  /** Текущий статус */
  status: string;
  /** Текущее местоположение */
  current_location?: string;
  /** Расчетная дата доставки (ISO 8601) */
  estimated_delivery?: string;
  /** Фактическая дата доставки (ISO 8601) */
  actual_delivery?: string;
  /** История событий */
  events: TrackingEvent[];
}

/**
 * Событие в истории трекинга
 */
export interface TrackingEvent {
  /** Время события (ISO 8601) */
  timestamp: string;
  /** Местоположение */
  location?: string;
  /** Статус */
  status: string;
  /** Описание события */
  description: string;
}

/**
 * Провайдер доставки
 */
export interface DeliveryProvider {
  /** Код провайдера */
  code: string;
  /** Название */
  name: string;
  /** Описание */
  description?: string;
  /** Активен ли провайдер */
  enabled: boolean;
  /** Логотип */
  logo_url?: string;
}

/**
 * Населенный пункт (для Post Express)
 */
export interface Settlement {
  /** ID населенного пункта */
  id: string;
  /** Название */
  name: string;
  /** Почтовый индекс */
  postal_code: string;
  /** Регион */
  region?: string;
}

/**
 * Улица в населенном пункте (для Post Express)
 */
export interface Street {
  /** ID улицы */
  id: string;
  /** Название улицы */
  name: string;
  /** ID населенного пункта */
  settlement_id: string;
}

/**
 * Парцел локер / постомат (для Post Express)
 */
export interface ParcelLocker {
  /** ID постомата */
  id: string;
  /** Название/номер постомата */
  name: string;
  /** Адрес */
  address: string;
  /** Город */
  city: string;
  /** Координаты (широта) */
  latitude?: number;
  /** Координаты (долгота) */
  longitude?: number;
  /** Часы работы */
  working_hours?: string;
  /** Доступен ли постомат */
  available: boolean;
}

// ========================================
// API SERVICE
// ========================================

/**
 * Delivery Service - API клиент для работы с доставкой
 *
 * @example
 * ```typescript
 * import { deliveryService } from '@/services/delivery';
 *
 * // Рассчитать стоимость доставки
 * const rate = await deliveryService.calculateRate({
 *   provider: 'post_express',
 *   from_city: 'Belgrade',
 *   to_city: 'Novi Sad',
 *   weight: 2.5,
 * });
 *
 * // Создать отправление
 * const shipment = await deliveryService.createShipment({
 *   order_id: 123,
 *   provider_code: 'post_express',
 *   from_address: { ... },
 *   to_address: { ... },
 *   packages: [{ weight: 2.5 }],
 * });
 *
 * // Отследить отправление
 * const tracking = await deliveryService.trackShipment('PE123456789');
 * ```
 */
export const deliveryService = {
  /**
   * Рассчитать стоимость доставки
   *
   * @param request - Параметры расчета
   * @returns Стоимость доставки и детализация
   *
   * @example
   * ```typescript
   * const response = await deliveryService.calculateRate({
   *   provider: 'post_express',
   *   from_city: 'Belgrade',
   *   to_city: 'Novi Sad',
   *   weight: 2.5,
   *   cash_on_delivery: true,
   *   cod_amount: 5000,
   * });
   *
   * if (response.data) {
   *   console.log('Total cost:', response.data.total_cost);
   *   console.log('Estimated delivery:', response.data.estimated_delivery);
   * }
   * ```
   */
  async calculateRate(
    request: CalculateRateRequest
  ): Promise<ApiResponse<CalculateRateResponse>> {
    return apiClient.post<CalculateRateResponse>('/delivery/calculate-rate', request);
  },

  /**
   * Рассчитать доставку для корзины
   *
   * @param storefrontId - ID витрины
   * @returns Стоимость доставки корзины
   *
   * @example
   * ```typescript
   * const response = await deliveryService.calculateCart(123);
   *
   * if (response.data) {
   *   console.log('Delivery cost for cart:', response.data.total_cost);
   * }
   * ```
   */
  async calculateCart(storefrontId: number): Promise<ApiResponse<CalculateRateResponse>> {
    return apiClient.post<CalculateRateResponse>(
      `/delivery/calculate-cart/${storefrontId}`
    );
  },

  /**
   * Получить список доступных провайдеров доставки
   *
   * @returns Список провайдеров
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getProviders();
   *
   * if (response.data) {
   *   response.data.forEach(provider => {
   *     console.log(provider.name, provider.code);
   *   });
   * }
   * ```
   */
  async getProviders(): Promise<ApiResponse<DeliveryProvider[]>> {
    return apiClient.get<DeliveryProvider[]>('/delivery/providers');
  },

  /**
   * Создать отправление
   *
   * @param request - Данные отправления
   * @returns Информация о созданном отправлении
   *
   * @example
   * ```typescript
   * const response = await deliveryService.createShipment({
   *   order_id: 123,
   *   provider_code: 'post_express',
   *   from_address: {
   *     name: 'Store Name',
   *     phone: '+381601234567',
   *     street: 'Main Street 1',
   *     city: 'Belgrade',
   *     postalCode: '11000',
   *     country: 'RS',
   *   },
   *   to_address: {
   *     name: 'Customer Name',
   *     phone: '+381607654321',
   *     street: 'Customer Street 5',
   *     city: 'Novi Sad',
   *     postalCode: '21000',
   *     country: 'RS',
   *   },
   *   packages: [
   *     {
   *       weight: 2.5,
   *       length: 30,
   *       width: 20,
   *       height: 10,
   *       description: 'Order #123',
   *     },
   *   ],
   * });
   *
   * if (response.data) {
   *   console.log('Tracking number:', response.data.tracking_number);
   * }
   * ```
   */
  async createShipment(
    request: CreateShipmentRequest
  ): Promise<ApiResponse<CreateShipmentResponse>> {
    return apiClient.post<CreateShipmentResponse>('/delivery/shipments', request);
  },

  /**
   * Отследить отправление по трекинг токену
   *
   * @param trackingToken - Трекинг токен или номер
   * @returns Информация о трекинге
   *
   * @example
   * ```typescript
   * const response = await deliveryService.trackShipment('PE123456789');
   *
   * if (response.data) {
   *   console.log('Status:', response.data.status);
   *   console.log('Current location:', response.data.current_location);
   *   console.log('Events:', response.data.events);
   * }
   * ```
   */
  async trackShipment(trackingToken: string): Promise<ApiResponse<TrackingInfo>> {
    return apiClient.get<TrackingInfo>(`/delivery/track/${trackingToken}`);
  },

  /**
   * Отменить отправление
   *
   * @param shipmentId - ID отправления
   * @returns Результат отмены
   *
   * @example
   * ```typescript
   * const response = await deliveryService.cancelShipment(123);
   *
   * if (!response.error) {
   *   console.log('Shipment cancelled successfully');
   * }
   * ```
   */
  async cancelShipment(shipmentId: number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/delivery/shipments/${shipmentId}`);
  },

  /**
   * Получить delivery атрибуты продукта
   *
   * @param productId - ID продукта
   * @param type - Тип продукта (c2c/b2c)
   * @returns Атрибуты доставки
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getProductAttributes('123', 'c2c');
   *
   * if (response.data) {
   *   console.log('Weight:', response.data.weight);
   *   console.log('Dimensions:', response.data.dimensions);
   * }
   * ```
   */
  async getProductAttributes(productId: string, type: string): Promise<ApiResponse<any>> {
    return apiClient.get<any>(`/products/${productId}/delivery-attributes?type=${type}`);
  },

  /**
   * Получить дефолтные настройки доставки категории
   *
   * @param categoryId - ID категории
   * @returns Настройки доставки по умолчанию
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getCategoryDefaults('123');
   *
   * if (response.data) {
   *   console.log('Default weight:', response.data.default_weight);
   * }
   * ```
   */
  async getCategoryDefaults(categoryId: string): Promise<ApiResponse<any>> {
    return apiClient.get<any>(`/categories/${categoryId}/delivery-defaults`);
  },

  // ========================================
  // POST EXPRESS SPECIFIC
  // ========================================

  /**
   * Получить список населенных пунктов (для Post Express)
   *
   * @param searchTerm - Поисковый запрос (опционально)
   * @returns Список населенных пунктов
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getSettlements('Beograd');
   *
   * if (response.data) {
   *   response.data.forEach(settlement => {
   *     console.log(settlement.name, settlement.postal_code);
   *   });
   * }
   * ```
   */
  async getSettlements(searchTerm?: string): Promise<ApiResponse<Settlement[]>> {
    const url = searchTerm
      ? `/delivery/settlements?search=${encodeURIComponent(searchTerm)}`
      : '/delivery/settlements';
    return apiClient.get<Settlement[]>(url);
  },

  /**
   * Получить список улиц в населенном пункте (для Post Express)
   *
   * @param settlementId - ID населенного пункта
   * @param searchTerm - Поисковый запрос (опционально)
   * @returns Список улиц
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getStreets('123', 'Main');
   *
   * if (response.data) {
   *   response.data.forEach(street => {
   *     console.log(street.name);
   *   });
   * }
   * ```
   */
  async getStreets(
    settlementId: string,
    searchTerm?: string
  ): Promise<ApiResponse<Street[]>> {
    const url = searchTerm
      ? `/delivery/settlements/${settlementId}/streets?search=${encodeURIComponent(searchTerm)}`
      : `/delivery/settlements/${settlementId}/streets`;
    return apiClient.get<Street[]>(url);
  },

  /**
   * Получить список постоматов/парцел локеров (для Post Express)
   *
   * @param cityId - ID города (опционально)
   * @returns Список постоматов
   *
   * @example
   * ```typescript
   * const response = await deliveryService.getParcelLockers('123');
   *
   * if (response.data) {
   *   response.data.forEach(locker => {
   *     console.log(locker.name, locker.address);
   *   });
   * }
   * ```
   */
  async getParcelLockers(cityId?: string): Promise<ApiResponse<ParcelLocker[]>> {
    const url = cityId
      ? `/delivery/parcel-lockers?city_id=${encodeURIComponent(cityId)}`
      : '/delivery/parcel-lockers';
    return apiClient.get<ParcelLocker[]>(url);
  },
};

/**
 * Экспорт для обратной совместимости и тестирования
 */
export default deliveryService;
