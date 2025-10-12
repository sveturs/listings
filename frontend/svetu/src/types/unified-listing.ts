/**
 * Типы для Unified Listings API (C2C + B2C объединенные)
 *
 * Этот файл содержит типы для работы с unified API, который объединяет
 * C2C listings (marketplace) и B2C products (storefronts) без дублирования данных.
 */

export type ListingSourceType = 'c2c' | 'b2c' | 'all';

export interface UnifiedImage {
  id: number;
  url: string;
  thumbnail_url?: string;
  is_main: boolean;
  display_order: number;
}

export interface UnifiedListing {
  // Базовая информация
  id: number;
  source_type: 'c2c' | 'b2c';
  user_id: number;
  category_id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  status: string;

  // Местоположение
  location: string;
  latitude?: number;
  longitude?: number;
  city: string;
  country: string;
  show_on_map: boolean;

  // Метаданные
  views_count: number;
  created_at: string;
  updated_at: string;
  original_language: string;
  storefront_id?: number; // только для B2C
  metadata?: Record<string, any>;

  // Связанные данные
  images: UnifiedImage[];
  user?: {
    id: number;
    name: string;
    email?: string;
    picture_url?: string;
    rating?: number;
    total_reviews?: number;
    verified?: boolean;
  };
  category?: {
    id: number;
    name: string;
    slug: string;
    translations?: Record<string, string>;
  };
  storefront?: {
    id: number;
    name: string;
    slug: string;
    logo_url?: string;
    rating?: number;
  };
  translations?: Record<string, any>;

  // UI флаги
  is_favorite?: boolean;
  has_discount?: boolean;
}

export interface UnifiedListingsFilters {
  source_type?: ListingSourceType;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  condition?: string;
  query?: string;
  storefront_id?: number;
  limit?: number;
  offset?: number;
}

export interface UnifiedListingsResponse {
  success: boolean;
  data: UnifiedListing[];
  total: number;
  limit: number;
  offset: number;
}

/**
 * Параметры для получения unified listing по ID
 */
export interface GetUnifiedListingByIdParams {
  id: number;
  source_type: 'c2c' | 'b2c';
}

/**
 * Константы для unified listings
 */
export const UNIFIED_LISTING_CONSTANTS = {
  // Значения по умолчанию
  DEFAULT_SOURCE_TYPE: 'all' as ListingSourceType,
  DEFAULT_LIMIT: 20,
  DEFAULT_OFFSET: 0,

  // Лимиты
  MAX_LIMIT: 100,
  MIN_LIMIT: 1,

  // Фильтры
  CONDITION_VALUES: ['new', 'used', 'refurbished', 'like_new'] as const,
} as const;

/**
 * Type guards
 */
export function isC2CListing(listing: UnifiedListing): boolean {
  return listing.source_type === 'c2c';
}

export function isB2CListing(listing: UnifiedListing): boolean {
  return listing.source_type === 'b2c';
}

/**
 * Получить URL детальной страницы для unified listing
 */
export function getUnifiedListingDetailUrl(
  listing: UnifiedListing,
  locale: string
): string {
  if (listing.source_type === 'c2c') {
    return `/${locale}/c2c/listings/${listing.id}`;
  } else {
    return `/${locale}/b2c/products/${listing.id}`;
  }
}

/**
 * Получить главное изображение из массива
 */
export function getMainImage(images: UnifiedImage[]): UnifiedImage | undefined {
  return images.find((img) => img.is_main) || images[0];
}

/**
 * Отсортировать изображения (главное первым, затем по display_order)
 */
export function sortImages(images: UnifiedImage[]): UnifiedImage[] {
  return [...images].sort((a, b) => {
    if (a.is_main && !b.is_main) return -1;
    if (!a.is_main && b.is_main) return 1;
    return a.display_order - b.display_order;
  });
}
