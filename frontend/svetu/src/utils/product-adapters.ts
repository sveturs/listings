import type { MarketplaceItem } from '@/types/marketplace';
import type { components } from '@/types/generated/api';
import type {
  UnifiedProduct,
  UnifiedProductImage,
  UnifiedStorefrontInfo,
} from '@/types/unified-product';
import config from '@/config';

type StorefrontProduct =
  components['schemas']['backend_internal_domain_models.StorefrontProduct'];

/**
 * Конвертирует изображения маркетплейса в унифицированный формат
 */
function adaptMarketplaceImages(
  images?: MarketplaceItem['images']
): UnifiedProductImage[] {
  if (!images || images.length === 0) return [];

  return images.map((img) => ({
    id: img.id,
    url: config.buildImageUrl(img.public_url),
    publicUrl: img.public_url,
    isMain: img.is_main,
  }));
}

/**
 * Конвертирует изображения витрины в унифицированный формат
 */
function adaptStorefrontImages(
  images?: StorefrontProduct['images']
): UnifiedProductImage[] {
  if (!images || images.length === 0) return [];

  return images.map((img, index) => ({
    id: img.id || index,
    url: img.image_url || '',
    publicUrl: img.image_url || '',
    isMain: img.is_default || false,
  }));
}

/**
 * Конвертирует атрибуты в унифицированный формат
 */
function adaptAttributes(
  attributes?: MarketplaceItem['attributes']
): UnifiedProduct['attributes'] {
  if (!attributes || attributes.length === 0) return [];

  return attributes.map((attr) => ({
    id: attr.attribute_id,
    name: attr.attribute_name || attr.name || '',
    value:
      attr.value ||
      attr.text_value ||
      attr.numeric_value ||
      attr.boolean_value ||
      '',
    displayValue: String(
      attr.value ||
        attr.text_value ||
        attr.numeric_value ||
        attr.boolean_value ||
        ''
    ),
  }));
}

/**
 * Адаптер для конвертации MarketplaceItem в UnifiedProduct
 */
export function adaptMarketplaceItem(item: MarketplaceItem): UnifiedProduct {
  // Попытка найти condition в атрибутах, если его нет в основном объекте
  let condition = item.condition;
  if (!condition && item.attributes) {
    const conditionAttr = item.attributes.find(
      (attr) => attr.attribute_name === 'condition' || attr.name === 'condition'
    );
    if (conditionAttr) {
      condition = String(conditionAttr.value || conditionAttr.text_value || '');
    }
  }

  return {
    // Базовая информация
    id: item.id,
    type: item.product_type || 'marketplace',
    name: item.title,
    description: item.description,

    // Цена
    price: item.price || 0,
    oldPrice: item.old_price,
    currency: 'RSD', // TODO: Получать из API
    hasDiscount: item.has_discount || false,

    // Изображения
    images: adaptMarketplaceImages(item.images),

    // Категория
    category: item.category
      ? {
          id: item.category.id,
          name: item.category.name,
          slug: item.category.slug,
          translations: item.category.translations,
        }
      : undefined,

    // Продавец/Витрина
    seller: item.user
      ? {
          id: item.user.id,
          name: item.user.name,
          email: item.user.email,
          pictureUrl: item.user.picture_url,
          rating: 4.8, // TODO: Получать из API
          totalReviews: 25, // TODO: Получать из API
          verified: false, // TODO: Получать из API
        }
      : item.storefront
        ? {
            id: item.storefront.id,
            name: item.storefront.name,
            email: '',
            pictureUrl: '',
            rating: 4.8, // TODO: Получать из API
            totalReviews: 25, // TODO: Получать из API
            verified: true, // Витрины считаем верифицированными
          }
        : undefined,
    storefront: item.storefront
      ? {
          id: item.storefront.id,
          name: item.storefront.name,
          slug: item.storefront.slug,
        }
      : undefined,

    // Местоположение
    location: {
      address: item.location,
      city: item.city,
      country: item.country,
      privacy: item.location_privacy,
      translations: item.translations,
    },

    // Состояние и наличие
    condition: condition as UnifiedProduct['condition'],
    isActive: item.status === 'active' || item.status === undefined, // По умолчанию считаем активным
    stockStatus:
      item.stock_status ||
      (item.product_type === 'storefront' ? 'in_stock' : undefined),
    stockQuantity:
      item.stock_quantity !== undefined ? item.stock_quantity : undefined, // Используем реальные остатки из API

    // Варианты - для товаров маркетплейса их нет
    variants: [],
    hasVariants: false,

    // Атрибуты
    attributes: adaptAttributes(item.attributes),

    // Статистика
    viewsCount: item.views_count,

    // Метаданные
    createdAt: item.created_at,
    updatedAt: item.updated_at,
    metadata: item.metadata,
  };
}

/**
 * Адаптер для конвертации StorefrontProduct в UnifiedProduct
 */
export function adaptStorefrontProduct(
  product: StorefrontProduct,
  storefrontInfo?: UnifiedStorefrontInfo
): UnifiedProduct {
  // Определяем статус наличия на основе количества
  const stockStatus =
    product.stock_status ||
    (product.stock_quantity === 0
      ? 'out_of_stock'
      : product.stock_quantity && product.stock_quantity < 10
        ? 'low_stock'
        : 'in_stock');

  return {
    // Базовая информация
    id: product.id || 0,
    type: 'storefront',
    name: product.name || '',
    description: product.description,

    // Цена
    price: product.price || 0,
    currency: product.currency || 'RSD',
    hasDiscount: false, // TODO: Добавить поддержку скидок для витрин

    // Изображения
    images: adaptStorefrontImages(product.images),

    // Категория
    category: product.category
      ? {
          id: product.category.id || 0,
          name: product.category.name || '',
          slug: product.category.slug || '',
          translations: product.category.translations,
        }
      : undefined,

    // Продавец (витрина)
    seller: storefrontInfo
      ? {
          id: storefrontInfo.id,
          name: storefrontInfo.name,
          email: '',
          pictureUrl: '',
          rating: 4.8, // TODO: Получать из API
          totalReviews: 25, // TODO: Получать из API
          verified: true, // Витрины считаем верифицированными
        }
      : {
          id: product.storefront_id || 0,
          name: 'Store', // TODO: Получать из API
          email: '',
          pictureUrl: '',
          rating: 4.8, // TODO: Получать из API
          totalReviews: 25, // TODO: Получать из API
          verified: true,
        },

    // Витрина
    storefront: storefrontInfo || {
      id: product.storefront_id || 0,
      name: 'Store', // TODO: Получать из API
      slug: 'store', // TODO: Получать из API
    },

    // Местоположение (для товаров с индивидуальным адресом)
    location: product.has_individual_location
      ? {
          address: product.individual_address,
          latitude: product.individual_latitude,
          longitude: product.individual_longitude,
          privacy: product.location_privacy as
            | 'exact'
            | 'street'
            | 'district'
            | 'city'
            | undefined,
        }
      : undefined,

    // Состояние и наличие
    condition: 'new', // Товары витрин обычно новые
    stockStatus: stockStatus as UnifiedProduct['stockStatus'],
    stockQuantity: product.stock_quantity,
    isActive: product.is_active,

    // Варианты
    variants: product.variants,
    hasVariants: product.variants && product.variants.length > 0,

    // Атрибуты
    attributes: product.attributes
      ? Object.entries(product.attributes).map(([key, value]) => ({
          name: key,
          value: value as string | number | boolean,
          displayValue: String(value),
        }))
      : [],

    // Статистика
    viewsCount: product.view_count,
    soldCount: product.sold_count,

    // Метаданные
    createdAt: product.created_at || new Date().toISOString(),
    updatedAt: product.updated_at,
  };
}

/**
 * Универсальный адаптер, который определяет тип продукта и применяет соответствующий адаптер
 */
export function adaptProduct(data: any): UnifiedProduct {
  // Определяем тип по наличию характерных полей
  if (data.title && data.user_id !== undefined) {
    // Это MarketplaceItem
    return adaptMarketplaceItem(data as MarketplaceItem);
  } else if (data.name && data.stock_quantity !== undefined) {
    // Это StorefrontProduct
    return adaptStorefrontProduct(data as StorefrontProduct);
  } else if (data.product_type) {
    // Используем явно указанный тип
    if (data.product_type === 'storefront') {
      return adaptStorefrontProduct(data);
    } else {
      return adaptMarketplaceItem(data);
    }
  }

  // По умолчанию пытаемся как MarketplaceItem
  return adaptMarketplaceItem(data);
}
