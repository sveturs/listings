import type { UnifiedProduct } from '@/types/unified-product';
import { normalizeImageUrl } from '@/utils/imageUtils';

/**
 * Получает URL продукта в зависимости от его типа
 */
export function getProductUrl(product: UnifiedProduct, locale: string): string {
  if (product.type === 'storefront' && product.storefront?.slug) {
    return `/${locale}/b2c/${product.storefront.slug}/products/${product.id}`;
  }

  // Для обычных объявлений маркетплейса
  return `/${locale}/c2c/${product.id}`;
}

/**
 * Определяет, может ли текущий пользователь добавить товар в корзину
 */
export function canAddToCart(
  product: UnifiedProduct,
  userId?: number
): boolean {
  // Нельзя покупать свой товар (проверяем продавца)
  if (userId && product.seller?.id === userId) {
    return false;
  }

  // Только товары витрин можно добавлять в корзину
  if (product.type !== 'storefront') {
    return false;
  }

  // Товар должен быть активным и в наличии
  if (!product.isActive || product.stockStatus === 'out_of_stock') {
    return false;
  }

  return true;
}

/**
 * Определяет, может ли пользователь начать чат с продавцом
 */
export function canStartChat(
  product: UnifiedProduct,
  userId?: number,
  isAuthenticated?: boolean
): boolean {
  // Нужна авторизация
  if (!isAuthenticated) {
    return false;
  }

  // Нельзя писать самому себе
  if (product.seller?.id === userId) {
    return false;
  }

  return true;
}

/**
 * Получает главное изображение продукта
 */
export function getMainImage(product: UnifiedProduct): string | null {
  const mainImage =
    product.images.find((img) => img.isMain) || product.images[0];
  if (!mainImage?.url) return null;
  return normalizeImageUrl(mainImage.url);
}

/**
 * Форматирует цену с учетом валюты
 */
export function formatPrice(
  price: number,
  currency: string,
  locale: string
): string {
  const formatter = new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: currency || 'RSD',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  });

  return formatter.format(price);
}

/**
 * Вычисляет процент скидки
 */
export function getDiscountPercentage(
  oldPrice?: number,
  currentPrice?: number
): number {
  if (!oldPrice || !currentPrice || oldPrice <= currentPrice) {
    return 0;
  }

  return Math.round(((oldPrice - currentPrice) / oldPrice) * 100);
}

/**
 * Получает бэйдж состояния товара
 */
export function getConditionBadge(
  condition?: string
): { text: string; class: string } | null {
  const badges: Record<string, { text: string; class: string }> = {
    new: { text: 'condition.new', class: 'badge-success' },
    like_new: { text: 'condition.likeNew', class: 'badge-info' },
    used: { text: 'condition.used', class: 'badge-primary' },
    refurbished: { text: 'condition.refurbished', class: 'badge-warning' },
  };

  return condition ? badges[condition] || null : null;
}

/**
 * Получает цвет для статуса наличия товара
 */
export function getStockStatusColor(stockStatus?: string): string {
  const colors: Record<string, string> = {
    in_stock: 'text-success',
    low_stock: 'text-warning',
    out_of_stock: 'text-error',
  };

  return stockStatus
    ? colors[stockStatus] || 'text-base-content'
    : 'text-base-content';
}

/**
 * Вычисляет эко-баллы для б/у товаров
 */
export function getEcoScore(product: UnifiedProduct): number {
  if (product.condition === 'used' || product.condition === 'refurbished') {
    return 8;
  }
  return 0;
}

/**
 * Генерирует временное расстояние на основе ID
 * TODO: Заменить на реальное расстояние из API
 */
export function getMockDistance(productId: number): number {
  const hash = productId
    .toString()
    .split('')
    .reduce((acc, char) => acc + char.charCodeAt(0), 0);
  return (hash % 20) + 0.5; // От 0.5 до 20.5 км
}

/**
 * Определяет, нужно ли показывать бэйдж безопасной сделки
 */
export function showSecureDealBadge(product: UnifiedProduct): boolean {
  // Показываем для всех товаров витрин
  return product.type === 'storefront';
}

/**
 * Получает минимальную цену товара с учетом вариантов
 */
export function getMinPrice(product: UnifiedProduct): number {
  if (
    !product.hasVariants ||
    !product.variants ||
    product.variants.length === 0
  ) {
    return product.price;
  }

  const variantPrices = product.variants
    .filter((v) => v.is_active && v.price !== undefined)
    .map((v) => v.price || 0);

  if (variantPrices.length === 0) {
    return product.price;
  }

  return Math.min(product.price, ...variantPrices);
}

/**
 * Проверяет, есть ли товар в наличии (с учетом вариантов)
 */
export function isInStock(product: UnifiedProduct): boolean {
  // Для товаров без вариантов
  if (!product.hasVariants) {
    // Если stockStatus не определен, считаем в наличии
    if (!product.stockStatus) return true;
    return product.stockStatus !== 'out_of_stock';
  }

  // Для товаров с вариантами - хотя бы один должен быть в наличии
  return (
    product.variants?.some(
      (v) => v.is_active && v.stock_quantity && v.stock_quantity > 0
    ) || false
  );
}
