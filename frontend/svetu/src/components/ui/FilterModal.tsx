'use client';

import React, { useState, useEffect } from 'react';
import { X, Filter } from 'lucide-react';
import { useTranslations } from 'next-intl';

interface FilterModalProps {
  isOpen: boolean;
  onClose: () => void;
  filters: Record<string, any>;
  onFiltersChange: (filters: Record<string, any>) => void;
  selectedCategoryId?: number | null;
}

export const FilterModal: React.FC<FilterModalProps> = ({
  isOpen,
  onClose,
  filters,
  onFiltersChange,
  selectedCategoryId,
}) => {
  const t = useTranslations('marketplace');
  const [localFilters, setLocalFilters] = useState(filters);

  useEffect(() => {
    setLocalFilters(filters);
  }, [filters]);

  if (!isOpen) return null;

  const handleApply = () => {
    onFiltersChange(localFilters);
    onClose();
  };

  const handleReset = () => {
    setLocalFilters({});
  };

  const handlePriceChange = (field: 'priceMin' | 'priceMax', value: string) => {
    const numValue = value ? parseInt(value) : undefined;
    setLocalFilters((prev) => ({
      ...prev,
      [field]: numValue,
    }));
  };

  const handleConditionChange = (condition: string) => {
    setLocalFilters((prev) => ({
      ...prev,
      condition: prev.condition === condition ? undefined : condition,
    }));
  };

  const handleSellerTypeChange = (type: string) => {
    setLocalFilters((prev) => ({
      ...prev,
      sellerType: prev.sellerType === type ? undefined : type,
    }));
  };

  const activeFiltersCount = Object.values(localFilters).filter(
    (v) => v !== undefined && v !== ''
  ).length;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 z-50 lg:hidden"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="fixed inset-x-0 bottom-0 z-50 lg:hidden">
        <div className="bg-base-100 rounded-t-3xl shadow-xl max-h-[85vh] flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-base-200">
            <div className="flex items-center gap-2">
              <Filter className="w-5 h-5 text-primary" />
              <h3 className="text-lg font-bold">{t('filters.title')}</h3>
              {activeFiltersCount > 0 && (
                <span className="badge badge-primary badge-sm">
                  {activeFiltersCount}
                </span>
              )}
            </div>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-4 space-y-6">
            {/* Price Range */}
            <div>
              <h4 className="font-medium mb-3">{t('filters.priceRange')}</h4>
              <div className="flex gap-2 items-center">
                <input
                  type="number"
                  placeholder={t('filters.priceMin')}
                  className="input input-bordered input-sm flex-1"
                  value={localFilters.priceMin || ''}
                  onChange={(e) =>
                    handlePriceChange('priceMin', e.target.value)
                  }
                />
                <span className="text-base-content/60">—</span>
                <input
                  type="number"
                  placeholder={t('filters.priceMax')}
                  className="input input-bordered input-sm flex-1"
                  value={localFilters.priceMax || ''}
                  onChange={(e) =>
                    handlePriceChange('priceMax', e.target.value)
                  }
                />
              </div>
            </div>

            {/* Condition */}
            <div>
              <h4 className="font-medium mb-3">{t('filters.condition')}</h4>
              <div className="flex flex-wrap gap-2">
                {['new', 'like_new', 'used', 'refurbished'].map((condition) => (
                  <button
                    key={condition}
                    onClick={() => handleConditionChange(condition)}
                    className={`btn btn-sm ${
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

            {/* Seller Type */}
            <div>
              <h4 className="font-medium mb-3">{t('filters.sellerType')}</h4>
              <div className="flex gap-2">
                <button
                  onClick={() => handleSellerTypeChange('private')}
                  className={`btn btn-sm flex-1 ${
                    localFilters.sellerType === 'private'
                      ? 'btn-primary'
                      : 'btn-outline'
                  }`}
                >
                  {t('filters.privateSeller')}
                </button>
                <button
                  onClick={() => handleSellerTypeChange('company')}
                  className={`btn btn-sm flex-1 ${
                    localFilters.sellerType === 'company'
                      ? 'btn-primary'
                      : 'btn-outline'
                  }`}
                >
                  {t('filters.companySeller')}
                </button>
              </div>
            </div>

            {/* Location Range */}
            <div>
              <h4 className="font-medium mb-3">{t('filters.locationRange')}</h4>
              <select
                className="select select-bordered select-sm w-full"
                value={localFilters.distanceKm || ''}
                onChange={(e) =>
                  setLocalFilters((prev) => ({
                    ...prev,
                    distanceKm: e.target.value
                      ? parseInt(e.target.value)
                      : undefined,
                  }))
                }
              >
                <option value="">{t('filters.anyDistance')}</option>
                <option value="5">5 км</option>
                <option value="10">10 км</option>
                <option value="25">25 км</option>
                <option value="50">50 км</option>
                <option value="100">100 км</option>
              </select>
            </div>
          </div>

          {/* Footer */}
          <div className="p-4 border-t border-base-200 flex gap-2">
            <button
              onClick={handleReset}
              className="btn btn-outline flex-1"
              disabled={activeFiltersCount === 0}
            >
              {t('filters.reset')}
            </button>
            <button onClick={handleApply} className="btn btn-primary flex-1">
              {t('filters.apply')}{' '}
              {activeFiltersCount > 0 && `(${activeFiltersCount})`}
            </button>
          </div>
        </div>
      </div>
    </>
  );
};
