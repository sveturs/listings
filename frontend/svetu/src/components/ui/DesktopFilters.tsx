'use client';

import React, { useState, useEffect } from 'react';
import { Filter } from 'lucide-react';
import { useTranslations } from 'next-intl';

interface DesktopFiltersProps {
  filters: Record<string, any>;
  onFiltersChange: (filters: Record<string, any>) => void;
  selectedCategoryId?: number | null;
}

export const DesktopFilters: React.FC<DesktopFiltersProps> = ({
  filters,
  onFiltersChange,
}) => {
  const t = useTranslations('marketplace');
  const [localFilters, setLocalFilters] = useState(filters);
  const [isExpanded, setIsExpanded] = useState(false);

  useEffect(() => {
    setLocalFilters(filters);
  }, [filters]);

  const handlePriceChange = (field: 'priceMin' | 'priceMax', value: string) => {
    const numValue = value ? parseInt(value) : undefined;
    const newFilters = {
      ...localFilters,
      [field]: numValue,
    };
    setLocalFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const handleConditionChange = (condition: string) => {
    const newFilters = {
      ...localFilters,
      condition: localFilters.condition === condition ? undefined : condition,
    };
    setLocalFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const handleSellerTypeChange = (type: string) => {
    const newFilters = {
      ...localFilters,
      sellerType: localFilters.sellerType === type ? undefined : type,
    };
    setLocalFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const handleReset = () => {
    setLocalFilters({});
    onFiltersChange({});
  };

  const activeFiltersCount = Object.values(localFilters).filter(
    (v) => v !== undefined && v !== ''
  ).length;

  return (
    <div className="hidden lg:block col-span-1 row-span-1 bg-base-100 rounded-2xl shadow-xl overflow-hidden">
      {/* Заголовок */}
      <div
        className="p-6 cursor-pointer hover:bg-base-200/50 transition-colors"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-secondary/10 rounded-xl">
              <Filter className="w-6 h-6 text-secondary" />
            </div>
            <div>
              <h3 className="text-lg font-bold">{t('filters.title')}</h3>
              <p className="text-sm text-base-content/60">
                {activeFiltersCount > 0
                  ? `${activeFiltersCount} активных`
                  : 'Настройте параметры'}
              </p>
            </div>
          </div>
          <button className="btn btn-primary btn-sm">
            {isExpanded ? 'Скрыть' : 'Показать'}
          </button>
        </div>
      </div>

      {/* Развернутый контент */}
      {isExpanded && (
        <div className="px-6 pb-6 space-y-4 border-t border-base-200">
          {/* Ценовой диапазон */}
          <div className="pt-4">
            <h4 className="font-medium mb-2 text-sm">
              {t('filters.priceRange')}
            </h4>
            <div className="flex gap-2 items-center">
              <input
                type="number"
                placeholder="От"
                className="input input-bordered input-sm flex-1"
                value={localFilters.priceMin || ''}
                onChange={(e) => handlePriceChange('priceMin', e.target.value)}
              />
              <span className="text-base-content/60">—</span>
              <input
                type="number"
                placeholder="До"
                className="input input-bordered input-sm flex-1"
                value={localFilters.priceMax || ''}
                onChange={(e) => handlePriceChange('priceMax', e.target.value)}
              />
            </div>
          </div>

          {/* Состояние */}
          <div>
            <h4 className="font-medium mb-2 text-sm">
              {t('filters.condition')}
            </h4>
            <div className="grid grid-cols-2 gap-2">
              {['new', 'like_new', 'used', 'refurbished'].map((condition) => (
                <button
                  key={condition}
                  onClick={() => handleConditionChange(condition)}
                  className={`btn btn-xs ${
                    localFilters.condition === condition
                      ? 'btn-primary'
                      : 'btn-outline'
                  }`}
                >
                  {t(`condition.${condition}`)}
                </button>
              ))}
            </div>
          </div>

          {/* Тип продавца */}
          <div>
            <h4 className="font-medium mb-2 text-sm">
              {t('filters.sellerType')}
            </h4>
            <div className="grid grid-cols-2 gap-2">
              <button
                onClick={() => handleSellerTypeChange('private')}
                className={`btn btn-xs ${
                  localFilters.sellerType === 'private'
                    ? 'btn-primary'
                    : 'btn-outline'
                }`}
              >
                {t('filters.privateSeller')}
              </button>
              <button
                onClick={() => handleSellerTypeChange('company')}
                className={`btn btn-xs ${
                  localFilters.sellerType === 'company'
                    ? 'btn-primary'
                    : 'btn-outline'
                }`}
              >
                {t('filters.companySeller')}
              </button>
            </div>
          </div>

          {/* Кнопка сброса */}
          {activeFiltersCount > 0 && (
            <div className="pt-2">
              <button
                onClick={handleReset}
                className="btn btn-outline btn-sm btn-block"
              >
                {t('filters.reset')}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};
