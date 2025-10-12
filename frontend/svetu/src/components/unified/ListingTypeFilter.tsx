/**
 * Фильтр по типу источника unified listings (All / C2C / B2C)
 *
 * Компонент позволяет пользователю переключаться между:
 * - Все (C2C + B2C вместе)
 * - Только C2C объявления
 * - Только B2C товары из витрин
 */

'use client';

import { useTranslations } from 'next-intl';
import type { ListingSourceType } from '@/types/unified-listing';

export interface ListingTypeFilterProps {
  /** Текущее выбранное значение */
  value: ListingSourceType;
  /** Коллбэк при изменении значения */
  onChange: (value: ListingSourceType) => void;
  /** Дополнительные CSS классы */
  className?: string;
  /** Показывать количество товаров рядом с кнопкой */
  showCounts?: boolean;
  /** Объект с количеством для каждого типа */
  counts?: {
    all?: number;
    c2c?: number;
    b2c?: number;
  };
  /** Размер кнопок */
  size?: 'sm' | 'md' | 'lg';
  /** Вариант отображения */
  variant?: 'buttons' | 'tabs' | 'pills';
}

const SIZES = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg',
} as const;

const VARIANTS = {
  buttons: {
    base: 'rounded-lg font-medium transition-all duration-200',
    active: 'bg-blue-600 text-white shadow-md scale-105',
    inactive: 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:scale-102',
  },
  tabs: {
    base: 'border-b-2 pb-2 font-medium transition-colors duration-200',
    active: 'border-blue-600 text-blue-600',
    inactive:
      'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300',
  },
  pills: {
    base: 'rounded-full font-medium transition-all duration-200',
    active: 'bg-blue-600 text-white shadow-md',
    inactive: 'bg-gray-100 text-gray-700 hover:bg-gray-200',
  },
} as const;

export function ListingTypeFilter({
  value,
  onChange,
  className = '',
  showCounts = false,
  counts,
  size = 'md',
  variant = 'buttons',
}: ListingTypeFilterProps) {
  const t = useTranslations('unified');

  const options: Array<{
    value: ListingSourceType;
    label: string;
    icon?: string;
    description?: string;
  }> = [
    {
      value: 'all',
      label: t('filter.all'),
      icon: '🌐',
      description: t('filter.all_description'),
    },
    {
      value: 'c2c',
      label: t('filter.c2c'),
      icon: '👤',
      description: t('filter.c2c_description'),
    },
    {
      value: 'b2c',
      label: t('filter.b2c'),
      icon: '🏪',
      description: t('filter.b2c_description'),
    },
  ];

  const variantStyles = VARIANTS[variant];
  const sizeClasses = SIZES[size];

  return (
    <div
      className={`flex gap-2 ${variant === 'tabs' ? 'border-b border-gray-200' : ''} ${className}`}
      role="group"
      aria-label={t('filter.aria_label')}
    >
      {options.map((option) => {
        const isActive = value === option.value;
        const count = counts?.[option.value];

        return (
          <button
            key={option.value}
            onClick={() => onChange(option.value)}
            className={`
              ${variantStyles.base}
              ${isActive ? variantStyles.active : variantStyles.inactive}
              ${sizeClasses}
              ${variant === 'tabs' ? 'flex-1' : ''}
              focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
              disabled:opacity-50 disabled:cursor-not-allowed
            `}
            title={option.description}
            aria-pressed={isActive}
            aria-label={`${option.label}${showCounts && count !== undefined ? ` (${count})` : ''}`}
          >
            <span className="flex items-center gap-2">
              {/* Icon */}
              {option.icon && <span className="text-lg">{option.icon}</span>}

              {/* Label */}
              <span>{option.label}</span>

              {/* Count badge */}
              {showCounts && count !== undefined && (
                <span
                  className={`
                    ml-1 px-2 py-0.5 rounded-full text-xs font-semibold
                    ${isActive ? 'bg-blue-500 text-white' : 'bg-gray-300 text-gray-700'}
                  `}
                >
                  {count}
                </span>
              )}
            </span>
          </button>
        );
      })}
    </div>
  );
}

/**
 * Компактная версия фильтра для мобильных устройств
 */
export function ListingTypeFilterCompact({
  value,
  onChange,
  className = '',
  counts,
}: Omit<ListingTypeFilterProps, 'size' | 'variant'>) {
  const t = useTranslations('unified');

  const options: Array<{ value: ListingSourceType; label: string }> = [
    { value: 'all', label: t('filter.all_short') },
    { value: 'c2c', label: t('filter.c2c_short') },
    { value: 'b2c', label: t('filter.b2c_short') },
  ];

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value as ListingSourceType)}
      className={`
        w-full px-4 py-2 rounded-lg border border-gray-300
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
        bg-white text-gray-900
        ${className}
      `}
      aria-label={t('filter.aria_label')}
    >
      {options.map((option) => {
        const count = counts?.[option.value];
        return (
          <option key={option.value} value={option.value}>
            {option.label}
            {count !== undefined ? ` (${count})` : ''}
          </option>
        );
      })}
    </select>
  );
}

/**
 * Responsive версия - автоматически переключается между полной и компактной версией
 */
export function ListingTypeFilterResponsive(props: ListingTypeFilterProps) {
  return (
    <>
      {/* Desktop: кнопки */}
      <div className="hidden md:block">
        <ListingTypeFilter {...props} />
      </div>

      {/* Mobile: select */}
      <div className="block md:hidden">
        <ListingTypeFilterCompact {...props} />
      </div>
    </>
  );
}
