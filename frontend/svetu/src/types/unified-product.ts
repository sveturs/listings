import type { components } from '@/types/generated/api';

// Тип для унифицированного представления изображения
export interface UnifiedProductImage {
  id: number;
  url: string;
  publicUrl: string;
  isMain: boolean;
}

// Тип для информации о витрине
export interface UnifiedStorefrontInfo {
  id: number;
  name: string;
  slug: string;
}

// Тип для информации о продавце/пользователе
export interface UnifiedSellerInfo {
  id: number;
  name: string;
  email?: string;
  pictureUrl?: string;
  rating?: number;
  totalReviews?: number;
  verified?: boolean;
}

// Тип для информации о категории
export interface UnifiedCategoryInfo {
  id: number;
  name: string;
  slug: string;
  translations?: Record<string, string>;
}

// Тип для информации о местоположении
export interface UnifiedLocationInfo {
  address?: string;
  city?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  privacy?: 'exact' | 'street' | 'district' | 'city';
  address_multilingual?: Record<string, string>; // Мультиязычные адреса из геокодирования
  translations?: {
    address?: Record<string, string>;
    city?: Record<string, string>;
    country?: Record<string, string>;
  };
}

// Тип для атрибутов товара
export interface UnifiedProductAttribute {
  id?: number;
  name: string;
  value: string | number | boolean;
  displayValue?: string;
}

// Тип для вариантов товара
export type ProductVariant =
  components['schemas']['backend_internal_domain_models.StorefrontProductVariant'];

// Основной унифицированный тип продукта
export interface UnifiedProduct {
  // Базовая информация
  id: number;
  type: 'marketplace' | 'storefront';
  name: string;
  description?: string;

  // Цена
  price: number;
  oldPrice?: number;
  currency: string;
  hasDiscount?: boolean;
  discountPercentage?: number;

  // Изображения
  images: UnifiedProductImage[];

  // Категория
  category?: UnifiedCategoryInfo;

  // Продавец/Витрина
  seller?: UnifiedSellerInfo;
  storefront?: UnifiedStorefrontInfo;

  // Местоположение
  location?: UnifiedLocationInfo;

  // Состояние и наличие
  condition?: 'new' | 'used' | 'refurbished' | 'like_new';
  stockStatus?: 'in_stock' | 'low_stock' | 'out_of_stock';
  stockQuantity?: number;
  isActive?: boolean;

  // Варианты (для товаров витрин)
  variants?: ProductVariant[];
  hasVariants?: boolean;

  // Атрибуты
  attributes?: UnifiedProductAttribute[];

  // Статистика
  viewsCount?: number;
  favoritesCount?: number;
  soldCount?: number;

  // Метаданные
  createdAt: string;
  updatedAt?: string;

  // SEO и другие метаданные
  metadata?: {
    seo?: {
      keywords?: string;
      slug?: string;
    };
    [key: string]: any;
  };

  // Переводы
  translations?: Record<
    string,
    {
      title?: string;
      description?: string;
      [key: string]: any;
    }
  >;

  // Дополнительные поля для поиска
  score?: number;
  highlights?: Record<string, string[]>;
}

// Типы для адаптеров
export type MarketplaceItemAdapter = (item: any) => UnifiedProduct;
export type StorefrontProductAdapter = (
  product: any,
  storefrontInfo?: UnifiedStorefrontInfo
) => UnifiedProduct;
