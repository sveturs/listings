/**
 * Feature flags для управления функциональностью приложения
 */

export interface FeatureFlags {
  // Унифицированная система атрибутов
  useUnifiedAttributes: boolean;

  // API версии
  useApiV2: boolean;

  // Другие флаги
  enableCarSelector: boolean;
  enableSmartCategoryDetection: boolean;
  enableProductVariants: boolean;
  enableStorefronts: boolean;
}

// Получение feature flags из переменных окружения или localStorage
export function getFeatureFlags(): FeatureFlags {
  // Проверяем localStorage для dev/test окружения
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('featureFlags');
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch {
        console.warn('Failed to parse feature flags from localStorage');
      }
    }
  }

  // Дефолтные значения из переменных окружения
  return {
    useUnifiedAttributes:
      process.env.NEXT_PUBLIC_USE_UNIFIED_ATTRIBUTES === 'true',
    useApiV2: process.env.NEXT_PUBLIC_USE_API_V2 === 'true',
    enableCarSelector: process.env.NEXT_PUBLIC_ENABLE_CAR_SELECTOR !== 'false',
    enableSmartCategoryDetection:
      process.env.NEXT_PUBLIC_ENABLE_SMART_CATEGORY_DETECTION === 'true',
    enableProductVariants:
      process.env.NEXT_PUBLIC_ENABLE_PRODUCT_VARIANTS === 'true',
    enableStorefronts: process.env.NEXT_PUBLIC_ENABLE_STOREFRONTS !== 'false',
  };
}

// Сохранение feature flags в localStorage (для тестирования)
export function setFeatureFlags(flags: Partial<FeatureFlags>): void {
  if (typeof window === 'undefined') return;

  const current = getFeatureFlags();
  const updated = { ...current, ...flags };

  localStorage.setItem('featureFlags', JSON.stringify(updated));

  // Перезагрузка страницы для применения изменений
  if (confirm('Feature flags updated. Reload page to apply changes?')) {
    window.location.reload();
  }
}

// Hook для использования в компонентах React
export function useFeatureFlag(flagName: keyof FeatureFlags): boolean {
  const flags = getFeatureFlags();
  return flags[flagName];
}

// Проверка доступности функции
export function isFeatureEnabled(flagName: keyof FeatureFlags): boolean {
  const flags = getFeatureFlags();
  return flags[flagName];
}

// Экспорт синглтона для использования вне React
export const featureFlags = getFeatureFlags();
