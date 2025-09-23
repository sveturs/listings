'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { Sparkles, Euro, Gauge, Zap, Car } from 'lucide-react';

export interface QuickFilter {
  id: string;
  label: string;
  icon: React.ReactNode;
  filter: Record<string, any>;
  count?: number;
}

interface CarQuickFiltersProps {
  selectedFilters: string[];
  onToggleFilter: (filterId: string, filter: Record<string, any>) => void;
  filterCounts?: Record<string, number>;
  className?: string;
}

export const CarQuickFilters: React.FC<CarQuickFiltersProps> = ({
  selectedFilters,
  onToggleFilter,
  filterCounts,
  className = '',
}) => {
  const t = useTranslations('cars');

  const quickFilters: QuickFilter[] = [
    {
      id: 'new',
      label: t('quickFilters.newCars'),
      icon: <Sparkles className="w-4 h-4" />,
      filter: { condition: 'new' },
    },
    {
      id: 'budget',
      label: t('quickFilters.under10k'),
      icon: <Euro className="w-4 h-4" />,
      filter: { priceMax: 10000 },
    },
    {
      id: 'lowMileage',
      label: t('quickFilters.lowMileage'),
      icon: <Gauge className="w-4 h-4" />,
      filter: { mileageMax: 50000 },
    },
    {
      id: 'electric',
      label: t('quickFilters.electric'),
      icon: <Zap className="w-4 h-4" />,
      filter: { fuelType: 'electric' },
    },
    {
      id: 'suv',
      label: t('quickFilters.suv'),
      icon: <Car className="w-4 h-4" />,
      filter: { bodyTypes: ['suv'] },
    },
  ];

  return (
    <div className={`flex flex-wrap gap-2 ${className}`}>
      {quickFilters.map((quickFilter) => {
        const isActive = selectedFilters.includes(quickFilter.id);
        const count = filterCounts?.[quickFilter.id];

        return (
          <button
            key={quickFilter.id}
            onClick={() => onToggleFilter(quickFilter.id, quickFilter.filter)}
            className={`btn btn-sm gap-2 ${
              isActive ? 'btn-primary' : 'btn-outline'
            }`}
          >
            {quickFilter.icon}
            <span>{quickFilter.label}</span>
            {count !== undefined && (
              <span className="badge badge-sm">{count}</span>
            )}
          </button>
        );
      })}
    </div>
  );
};

export default CarQuickFilters;
