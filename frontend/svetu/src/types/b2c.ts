import type { components } from '@/types/generated/api';

// Базовые типы из API
export type StorefrontBase = components['schemas']['models.Storefront'];
export type B2CStoreCreateDTO =
  components['schemas']['models.B2CStoreCreateDTO'];
export type B2CStoreUpdateDTO =
  components['schemas']['models.B2CStoreUpdateDTO'];
export type B2CProduct =
  components['schemas']['models.B2CProduct'];

// Расширенный тип Storefront с дополнительными полями
export interface B2CStore extends StorefrontBase {
  // Дополнительные поля для UI
  stats?: {
    average_rating?: number;
    total_products?: number;
    total_reviews?: number;
    total_sales?: number;
    view_count?: number;
  };
  hours?: StorefrontHours[];
  payment_methods?: StorefrontPaymentMethod[];
  delivery_options?: StorefrontDeliveryOption[];
  location?: Location;
  business_type?: string;
  registration_number?: string;
  tax_number?: string;
  vat_number?: string;
  banner_image_url?: string;
  logo_url?: string;
  website_url?: string;
  social_links?: {
    facebook?: string;
    instagram?: string;
    telegram?: string;
  };
}
export type StorefrontAnalytics =
  components['schemas']['models.StorefrontAnalytics'];
export type StorefrontRatingSummary =
  components['schemas']['models.StorefrontRatingSummary'];
export type B2CStoreMapData =
  components['schemas']['models.B2CStoreMapData'];
export type StorefrontSearchResult = any; // Type removed from swagger schema
export type StorefrontSearchItem = any; // Type removed from swagger schema
export type StorefrontsListResponse =
  components['schemas']['handler.StorefrontsListResponse'];

// Связанные типы
export type StorefrontStaff = components['schemas']['models.StorefrontStaff'];
export type StorefrontHours = components['schemas']['models.StorefrontHours'];
export type StorefrontPaymentMethod =
  components['schemas']['models.StorefrontPaymentMethod'];
export type StorefrontDeliveryOption =
  components['schemas']['models.StorefrontDeliveryOption'];
export type Location = components['schemas']['models.Location'];
export type JSONB = any; // Type removed from swagger schema

// Enum типы
export type SubscriptionPlan =
  components['schemas']['models.SubscriptionPlanType'];
export type StaffRole = components['schemas']['models.StaffRole'];
export type PaymentMethodType =
  components['schemas']['models.PaymentMethodType'];
export type DeliveryProvider = components['schemas']['models.DeliveryProvider'];

// Дополнительные интерфейсы для UI

/**
 * Состояние формы создания витрины
 */
export interface B2CStoreFormData
  extends Omit<B2CStoreCreateDTO, 'location'> {
  // Координаты отдельно для удобства
  latitude?: number;
  longitude?: number;

  // Файлы для загрузки
  logoFile?: File;
  bannerFile?: File;

  // Дополнительные настройки
  theme?: {
    primaryColor?: string;
    secondaryColor?: string;
    style?: string;
  };

  seoMeta?: {
    title?: string;
    description?: string;
    keywords?: string;
  };
}

/**
 * Фильтры для поиска витрин
 */
export interface StorefrontFilters {
  search?: string;
  city?: string;
  isActive?: boolean;
  isVerified?: boolean;
  minRating?: number;
  paymentMethods?: PaymentMethodType[];
  hasDelivery?: boolean;
  hasSelfPickup?: boolean;
  isOpenNow?: boolean;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
  sortBy?: 'rating' | 'distance' | 'created_at' | 'products_count';
  sortOrder?: 'asc' | 'desc';
}

/**
 * Параметры пагинации
 */
export interface PaginationParams {
  limit: number;
  offset: number;
}

/**
 * Настройки витрины для UI
 */
export interface StorefrontSettings {
  // Общие настройки
  timezone?: string;
  currency?: string;
  language?: string;

  // Настройки уведомлений
  notifications?: {
    email?: boolean;
    telegram?: boolean;
    newOrders?: boolean;
    newMessages?: boolean;
    newReviews?: boolean;
  };

  // Настройки продаж
  sales?: {
    autoAcceptOrders?: boolean;
    requireOrderConfirmation?: boolean;
    allowBulkOrders?: boolean;
    minimumOrderAmount?: number;
  };

  // Настройки доставки
  delivery?: {
    autoCalculateShipping?: boolean;
    defaultShippingCost?: number;
    freeShippingThreshold?: number;
    allowSelfPickup?: boolean;
  };

  // SEO настройки
  seo?: {
    metaTitle?: string;
    metaDescription?: string;
    metaKeywords?: string[];
    openGraphImage?: string;
  };

  // Интеграции
  integrations?: {
    googleAnalytics?: string;
    facebookPixel?: string;
    telegramBot?: string;
  };
}

/**
 * Статистика витрины для дашборда
 */
export interface StorefrontDashboardStats {
  // Основные метрики
  totalViews: number;
  totalOrders: number;
  totalRevenue: number;
  averageRating: number;

  // Изменения за период
  viewsChange: number;
  ordersChange: number;
  revenueChange: number;
  ratingChange: number;

  // Дополнительные метрики
  conversionRate: number;
  averageOrderValue: number;
  repeatCustomerRate: number;
  topProducts: Array<{
    id: number;
    name: string;
    sales: number;
    revenue: number;
  }>;

  // География продаж
  topCities: Array<{
    name: string;
    orders: number;
    revenue: number;
  }>;

  // Методы оплаты
  paymentMethodsStats: Array<{
    method: PaymentMethodType;
    count: number;
    percentage: number;
  }>;
}

/**
 * Рабочие часы витрины в удобном формате
 */
export interface WorkingHoursData {
  [key: number]: {
    dayOfWeek: number;
    isOpen: boolean;
    openTime?: string;
    closeTime?: string;
    specialNote?: string;
  };
}

/**
 * Состояние загрузки для различных операций
 */
export interface StorefrontLoadingState {
  isLoading: boolean;
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;
  isUploadingLogo: boolean;
  isUploadingBanner: boolean;
  isLoadingAnalytics: boolean;
  isLoadingSettings: boolean;
}

/**
 * Ошибки валидации формы
 */
export interface StorefrontFormErrors {
  name?: string;
  description?: string;
  email?: string;
  phone?: string;
  website?: string;
  address?: string;
  city?: string;
  latitude?: string;
  longitude?: string;
  logo?: string;
  banner?: string;
  general?: string;
}

/**
 * Опции для создания витрины
 */
export interface CreateStorefrontOptions {
  redirect?: boolean;
  showSuccess?: boolean;
  autoSave?: boolean;
}

/**
 * Конфигурация карты витрин
 */
export interface StorefrontMapConfig {
  center: {
    lat: number;
    lng: number;
  };
  zoom: number;
  bounds?: {
    north: number;
    south: number;
    east: number;
    west: number;
  };
  clustering?: boolean;
  showFilters?: boolean;
  allowSearch?: boolean;
}

/**
 * Элемент для автокомплита
 */
export interface AutocompleteItem {
  id: string;
  label: string;
  description?: string;
  coordinates?: {
    lat: number;
    lng: number;
  };
}

/**
 * Результат геокодирования
 */
export interface GeocodeResult {
  address: string;
  latitude: number;
  longitude: number;
  city?: string;
  country?: string;
  postalCode?: string;
}

/**
 * Константы для витрин
 */
export const STOREFRONT_CONSTANTS = {
  // Лимиты
  MAX_NAME_LENGTH: 255,
  MAX_DESCRIPTION_LENGTH: 2000,
  MAX_LOGO_SIZE: 5 * 1024 * 1024, // 5MB
  MAX_BANNER_SIZE: 10 * 1024 * 1024, // 10MB

  // Форматы файлов
  ALLOWED_IMAGE_FORMATS: ['image/jpeg', 'image/png', 'image/webp'],

  // Пагинация
  DEFAULT_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,

  // Поиск
  MIN_SEARCH_LENGTH: 2,
  SEARCH_DEBOUNCE_MS: 300,

  // Карта
  DEFAULT_MAP_ZOOM: 12,
  MAX_MAP_ZOOM: 18,
  MIN_MAP_ZOOM: 3,

  // Фильтры
  MAX_RADIUS_KM: 100,
  DEFAULT_RADIUS_KM: 10,
} as const;

/**
 * Статусы витрины
 */
export enum StorefrontStatus {
  DRAFT = 'draft',
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
  DELETED = 'deleted',
}

/**
 * Типы событий витрины
 */
export enum StorefrontEventType {
  CREATED = 'created',
  UPDATED = 'updated',
  DELETED = 'deleted',
  VIEW_RECORDED = 'view_recorded',
  ORDER_RECEIVED = 'order_received',
  REVIEW_RECEIVED = 'review_received',
  STAFF_ADDED = 'staff_added',
  STAFF_REMOVED = 'staff_removed',
}

/**
 * Приоритеты уведомлений
 */
export enum NotificationPriority {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  URGENT = 'urgent',
}
