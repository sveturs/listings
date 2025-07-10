'use client';

import { useTranslations } from 'next-intl';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';

interface SearchFiltersProps {
  filters: {
    category_id?: string;
    price_min?: number;
    price_max?: number;
    product_types?: string[];
    sort_by?: string;
    sort_order?: string;
    city?: string;
  };
  searchQuery: string;
  resultsCount: number;
  onFilterChange: (filters: Partial<SearchFiltersProps['filters']>) => void;
}

/**
 * Компонент фильтров поиска с интегрированным трекингом
 * Автоматически отслеживает изменения фильтров
 */
export default function SearchFilters({
  filters,
  searchQuery,
  resultsCount,
  onFilterChange,
}: SearchFiltersProps) {
  const t = useTranslations('search');
  const { trackSearchFilterApplied } = useBehaviorTracking();

  const handleFilterChange = async (
    filterType: string,
    filterValue: any,
    newFilters: Partial<SearchFiltersProps['filters']>
  ) => {
    // Трекинг изменения фильтра
    if (searchQuery) {
      try {
        await trackSearchFilterApplied({
          search_query: searchQuery,
          filter_type: filterType,
          filter_value: JSON.stringify(filterValue),
          results_count_before: resultsCount,
          results_count_after: 0, // Будет обновлено после применения фильтра
        });
      } catch (error) {
        console.error('Failed to track filter change:', error);
      }
    }

    // Применяем изменения фильтра
    onFilterChange(newFilters);
  };

  const handleProductTypeToggle = (productType: string) => {
    const currentTypes = filters.product_types || [];
    let newTypes: string[];

    if (currentTypes.includes(productType) && currentTypes.length > 1) {
      newTypes = currentTypes.filter((t) => t !== productType);
    } else if (!currentTypes.includes(productType)) {
      newTypes = [...currentTypes, productType];
    } else {
      return; // Не разрешаем убрать последний тип
    }

    handleFilterChange('product_types', newTypes, { product_types: newTypes });
  };

  const handlePriceChange = (type: 'min' | 'max', value: string) => {
    const numValue = value ? Number(value) : undefined;
    const filterType = type === 'min' ? 'price_min' : 'price_max';
    const newFilters = { [filterType]: numValue };

    handleFilterChange(filterType, numValue, newFilters);
  };

  const handleCityChange = (value: string) => {
    const newFilters = { city: value || undefined };
    handleFilterChange('city', value, newFilters);
  };

  const handleCategoryChange = (value: string) => {
    const newFilters = { category_id: value || undefined };
    handleFilterChange('category_id', value, newFilters);
  };

  return (
    <div className="card bg-base-100 shadow-md">
      <div className="card-body">
        <h3 className="card-title text-lg flex items-center gap-2">
          <svg
            className="w-5 h-5 text-primary"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
            />
          </svg>
          {t('filters')}
        </h3>

        <div className="space-y-6">
          {/* Тип товаров */}
          <div>
            <label className="label">
              <span className="label-text font-medium">
                {t('productTypes')}
              </span>
            </label>
            <div className="grid grid-cols-2 gap-3">
              <div
                className={`card card-compact cursor-pointer transition-all ${
                  filters.product_types?.includes('marketplace')
                    ? 'ring-2 ring-primary bg-primary/5'
                    : 'bg-base-200 hover:bg-base-300'
                }`}
                onClick={() => handleProductTypeToggle('marketplace')}
              >
                <div className="card-body items-center text-center">
                  <svg
                    className="w-6 h-6 text-primary mb-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                    />
                  </svg>
                  <span className="text-xs font-medium">{t('private')}</span>
                </div>
              </div>
              <div
                className={`card card-compact cursor-pointer transition-all ${
                  filters.product_types?.includes('storefront')
                    ? 'ring-2 ring-primary bg-primary/5'
                    : 'bg-base-200 hover:bg-base-300'
                }`}
                onClick={() => handleProductTypeToggle('storefront')}
              >
                <div className="card-body items-center text-center">
                  <svg
                    className="w-6 h-6 text-secondary mb-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                    />
                  </svg>
                  <span className="text-xs font-medium">{t('stores')}</span>
                </div>
              </div>
            </div>
          </div>

          <div className="divider my-4"></div>

          {/* Диапазон цен */}
          <div>
            <label className="label">
              <span className="label-text font-medium flex items-center gap-2">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                {t('priceRange')}
              </span>
            </label>
            <div className="flex gap-2">
              <div className="form-control flex-1">
                <label className="input-group">
                  <span className="bg-base-200">{t('from')}</span>
                  <input
                    type="number"
                    className="input input-bordered w-full"
                    value={filters.price_min || ''}
                    onChange={(e) => handlePriceChange('min', e.target.value)}
                  />
                </label>
              </div>
              <div className="form-control flex-1">
                <label className="input-group">
                  <span className="bg-base-200">{t('to')}</span>
                  <input
                    type="number"
                    className="input input-bordered w-full"
                    value={filters.price_max || ''}
                    onChange={(e) => handlePriceChange('max', e.target.value)}
                  />
                </label>
              </div>
            </div>
          </div>

          <div className="divider my-4"></div>

          {/* Город */}
          <div>
            <label className="label">
              <span className="label-text font-medium flex items-center gap-2">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                {t('city')}
              </span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              placeholder={t('enterCity')}
              value={filters.city || ''}
              onChange={(e) => handleCityChange(e.target.value)}
            />
          </div>

          <div className="divider my-4"></div>

          {/* Категория */}
          <div>
            <label className="label">
              <span className="label-text font-medium flex items-center gap-2">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                  />
                </svg>
                {t('category')}
              </span>
            </label>
            <select
              className="select select-bordered w-full"
              value={filters.category_id || ''}
              onChange={(e) => handleCategoryChange(e.target.value)}
            >
              <option value="">{t('allCategories')}</option>
              {/* TODO: Добавить динамическую загрузку категорий */}
            </select>
          </div>
        </div>
      </div>
    </div>
  );
}
