'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { MarketplaceService } from '@/services/marketplace';
import { SmartFilters } from './SmartFilters';

export interface FilterValues {
  // Постоянные фильтры
  priceMin?: number;
  priceMax?: number;
  condition?: string;
  sellerType?: 'private' | 'company'; // для определения storefrontID
  hasDelivery?: boolean;

  // Динамические фильтры атрибутов
  attributeFilters?: Record<string, any>;
}

interface ListingFiltersProps {
  selectedCategoryId?: number | null;
  filters: FilterValues;
  onFiltersChange: (filters: FilterValues) => void;
  className?: string;
}

export default function ListingFilters({
  selectedCategoryId,
  filters,
  onFiltersChange,
  className = '',
}: ListingFiltersProps) {
  const t = useTranslations('marketplace');
  const tRoot = useTranslations();
  const locale = useLocale();
  const [isExpanded, setIsExpanded] = useState(true);
  const [filterableAttributes, setFilterableAttributes] = useState<any[]>([]);
  const [loadingAttributes, setLoadingAttributes] = useState(false);

  // Загрузка атрибутов для выбранной категории
  useEffect(() => {
    const loadAttributes = async () => {
      if (!selectedCategoryId) {
        setFilterableAttributes([]);
        return;
      }

      setLoadingAttributes(true);
      try {
        const response =
          await MarketplaceService.getCategoryAttributes(selectedCategoryId);

        if (response.data) {
          // Фильтруем только те атрибуты, которые можно использовать для фильтрации
          const filterable = response.data.filter(
            (attr: any) => attr.is_filterable === true
          );
          setFilterableAttributes(filterable);
        }
      } catch (error) {
        console.error('Failed to load category attributes:', error);
        setFilterableAttributes([]);
      } finally {
        setLoadingAttributes(false);
      }
    };

    loadAttributes();
  }, [selectedCategoryId, locale]);

  // Обработчик изменения постоянных фильтров
  const handlePermanentFilterChange = useCallback(
    (field: keyof FilterValues, value: any) => {
      onFiltersChange({
        ...filters,
        [field]: value || undefined, // убираем пустые значения
      });
    },
    [filters, onFiltersChange]
  );

  // Обработчик изменения атрибутов
  const handleAttributeFiltersChange = useCallback(
    (attributeFilters: Record<string, any>) => {
      onFiltersChange({
        ...filters,
        attributeFilters,
      });
    },
    [filters, onFiltersChange]
  );

  // Очистка всех фильтров
  const handleClearFilters = useCallback(() => {
    onFiltersChange({});
  }, [onFiltersChange]);

  // Проверка наличия активных фильтров
  const hasActiveFilters = Object.values(filters).some((value) => {
    if (value === null || value === undefined || value === '') return false;
    if (typeof value === 'object' && Object.keys(value).length === 0)
      return false;
    return true;
  });

  return (
    <div
      className={`card bg-base-100 shadow-sm border border-base-200 ${className}`}
    >
      <div className="card-body p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-base">
            {t('filters.title')}
            {hasActiveFilters && (
              <span className="badge badge-primary badge-sm ml-2">
                {Object.keys(filters).length}
              </span>
            )}
          </h3>
          <div className="flex items-center gap-2">
            {hasActiveFilters && (
              <button
                onClick={handleClearFilters}
                className="btn btn-ghost btn-xs text-error"
              >
                {t('filters.clearAll')}
              </button>
            )}
            <button
              onClick={() => setIsExpanded(!isExpanded)}
              className="btn btn-ghost btn-xs btn-square"
            >
              <svg
                className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </button>
          </div>
        </div>

        {isExpanded && (
          <div className="space-y-6">
            {/* Постоянные фильтры */}
            <div className="space-y-4">
              {/* Диапазон цен */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    {t('filters.priceRange')}
                  </span>
                </label>
                <div className="flex gap-2 items-center">
                  <input
                    type="number"
                    placeholder={t('filters.priceFrom')}
                    className="input input-bordered input-sm flex-1"
                    value={filters.priceMin || ''}
                    onChange={(e) =>
                      handlePermanentFilterChange(
                        'priceMin',
                        parseFloat(e.target.value) || null
                      )
                    }
                  />
                  <span className="text-base-content/60">—</span>
                  <input
                    type="number"
                    placeholder={t('filters.priceTo')}
                    className="input input-bordered input-sm flex-1"
                    value={filters.priceMax || ''}
                    onChange={(e) =>
                      handlePermanentFilterChange(
                        'priceMax',
                        parseFloat(e.target.value) || null
                      )
                    }
                  />
                </div>
              </div>

              {/* Состояние товара */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    {t('filters.condition')}
                  </span>
                </label>
                <select
                  className="select select-bordered select-sm"
                  value={filters.condition || ''}
                  onChange={(e) =>
                    handlePermanentFilterChange('condition', e.target.value)
                  }
                >
                  <option value="">{t('filters.anyCondition')}</option>
                  <option value="new">{tRoot('condition.new')}</option>
                  <option value="used">{tRoot('condition.used')}</option>
                  <option value="refurbished">
                    {tRoot('condition.refurbished')}
                  </option>
                  <option value="damaged">{tRoot('condition.damaged')}</option>
                </select>
              </div>

              {/* Тип продавца */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    {t('filters.sellerType')}
                  </span>
                </label>
                <select
                  className="select select-bordered select-sm"
                  value={filters.sellerType || ''}
                  onChange={(e) =>
                    handlePermanentFilterChange('sellerType', e.target.value)
                  }
                >
                  <option value="">{t('filters.anySeller')}</option>
                  <option value="private">{t('filters.privateSeller')}</option>
                  <option value="company">{t('filters.companySeller')}</option>
                </select>
              </div>

              {/* С доставкой */}
              <div className="form-control">
                <label className="label cursor-pointer justify-start gap-3">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={filters.hasDelivery || false}
                    onChange={(e) =>
                      handlePermanentFilterChange(
                        'hasDelivery',
                        e.target.checked ? true : null
                      )
                    }
                  />
                  <span className="label-text">
                    {t('filters.withDelivery')}
                  </span>
                </label>
              </div>
            </div>

            {/* Динамические фильтры атрибутов */}
            {selectedCategoryId && (
              <div className="border-t border-base-200 pt-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('filters.categoryFilters')}
                    </span>
                  </label>

                  {loadingAttributes ? (
                    <div className="flex justify-center py-4">
                      <span className="loading loading-spinner loading-sm"></span>
                    </div>
                  ) : filterableAttributes.length > 0 ? (
                    <SmartFilters
                      categoryId={selectedCategoryId}
                      onChange={handleAttributeFiltersChange}
                      lang={locale}
                      className="mt-2"
                    />
                  ) : (
                    <p className="text-sm text-base-content/60 mt-2">
                      {t('filters.noAttributeFilters')}
                    </p>
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
