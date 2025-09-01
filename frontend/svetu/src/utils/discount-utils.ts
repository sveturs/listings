// Utility functions for discount visualization

export interface DiscountVisualConfig {
  badgeClass: string;
  borderClass: string;
  glowClass: string;
  pulseAnimation: boolean;
  fireEmoji: boolean;
  isHot: boolean;
}

/**
 * Определяет визуальную конфигурацию для отображения скидки
 * @param discountPercentage процент скидки
 * @returns конфигурация для визуального отображения
 */
export function getDiscountVisualConfig(
  discountPercentage: number | null | undefined
): DiscountVisualConfig | null {
  if (!discountPercentage || discountPercentage <= 0) {
    return null;
  }

  // Скидка более 30% - "горячее" предложение
  if (discountPercentage >= 30) {
    return {
      badgeClass: 'badge-error text-white font-bold',
      borderClass: 'border-error ring-2 ring-error/50 ring-offset-2',
      glowClass: 'shadow-lg shadow-error/50',
      pulseAnimation: true,
      fireEmoji: true,
      isHot: true,
    };
  }

  // Скидка от 20% до 29% - привлекательное предложение
  if (discountPercentage >= 20) {
    return {
      badgeClass: 'badge-warning text-warning-content font-semibold',
      borderClass: 'border-warning ring-1 ring-warning/30',
      glowClass: 'shadow-md shadow-warning/30',
      pulseAnimation: false,
      fireEmoji: false,
      isHot: true,
    };
  }

  // Скидка от 10% до 19% - стандартное выделение
  if (discountPercentage >= 10) {
    return {
      badgeClass: 'badge-success text-success-content',
      borderClass: 'border-success/50',
      glowClass: '',
      pulseAnimation: false,
      fireEmoji: false,
      isHot: false,
    };
  }

  // Скидка менее 10% - минимальное выделение
  return {
    badgeClass: 'badge-ghost',
    borderClass: '',
    glowClass: '',
    pulseAnimation: false,
    fireEmoji: false,
    isHot: false,
  };
}

/**
 * Форматирует бейдж скидки с эмодзи
 * @param discountPercentage процент скидки
 * @param showEmoji показывать ли эмодзи
 * @returns отформатированная строка для бейджа
 */
export function formatDiscountBadge(
  discountPercentage: number,
  showEmoji: boolean = true
): string {
  const config = getDiscountVisualConfig(discountPercentage);
  if (!config) return '';

  let badge = `-${discountPercentage}%`;

  if (showEmoji && config.fireEmoji) {
    badge = `🔥 ${badge}`;
  } else if (showEmoji && config.isHot) {
    badge = `⚡ ${badge}`;
  }

  return badge;
}

/**
 * Вычисляет процент скидки
 * @param oldPrice старая цена
 * @param newPrice новая цена
 * @returns процент скидки
 */
export function calculateDiscountPercentage(
  oldPrice: number,
  newPrice: number
): number {
  if (oldPrice <= 0 || newPrice <= 0 || oldPrice <= newPrice) {
    return 0;
  }
  return Math.round(((oldPrice - newPrice) / oldPrice) * 100);
}
