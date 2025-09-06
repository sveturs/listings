'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { X } from 'lucide-react';

interface ActiveFiltersChipsProps {
  filters: Record<string, any>;
  onRemoveFilter: (key: string) => void;
  className?: string;
}

export const ActiveFiltersChips: React.FC<ActiveFiltersChipsProps> = ({
  filters,
  onRemoveFilter,
  className = '',
}) => {
  const t = useTranslations('search');

  const formatFilterDisplay = (key: string, value: any): string => {
    const filterLabels: Record<string, string> = {
      price_min: t('priceMin'),
      price_max: t('priceMax'),
      condition: t('condition'),
      city: t('location'),
      distance: t('radius'),
      make: t('make'),
      model: t('model'),
      yearFrom: t('yearFrom'),
      yearTo: t('yearTo'),
      mileage: t('mileage'),
      fuelType: t('fuelType'),
      transmission: t('transmission'),
      propertyType: t('propertyType'),
      rooms: t('rooms'),
      area: t('area'),
      floor: t('floor'),
      brand: t('brand'),
    };

    const label = filterLabels[key] || key;

    if (key === 'condition') {
      return `${label}: ${t(`conditions.${value}`)}`;
    }

    if (key === 'distance') {
      return `${label}: ${value} km`;
    }

    if (key === 'price_min' || key === 'price_max') {
      return `${label}: €${value}`;
    }

    return `${label}: ${value}`;
  };

  const filterEntries = Object.entries(filters).filter(
    ([_key, value]) => value !== undefined && value !== '' && value !== null
  );

  if (filterEntries.length === 0) return null;

  return (
    <div className={`${className}`}>
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-semibold text-base-content/70">
          {t('activeFilters')}
        </span>
        <button
          onClick={() => filterEntries.forEach(([key]) => onRemoveFilter(key))}
          className="btn btn-ghost btn-xs text-error"
        >
          {t('clearAll')}
        </button>
      </div>
      <div className="flex flex-wrap gap-2">
        {filterEntries.map(([key, value]) => (
          <div key={key} className="badge badge-primary badge-lg gap-2 py-3">
            <span className="text-xs">{formatFilterDisplay(key, value)}</span>
            <button
              onClick={() => onRemoveFilter(key)}
              className="hover:bg-primary-focus rounded-full p-0.5"
            >
              <X className="w-3 h-3" />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
