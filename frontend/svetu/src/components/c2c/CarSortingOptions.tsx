'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { ArrowUpDown, Calendar, Gauge, DollarSign } from 'lucide-react';

export type CarSortOption =
  | 'created_at_desc'
  | 'created_at_asc'
  | 'price_asc'
  | 'price_desc'
  | 'year_desc'
  | 'year_asc'
  | 'mileage_asc'
  | 'mileage_desc'
  | 'price_year_ratio';

interface CarSortingOptionsProps {
  value: CarSortOption;
  onChange: (value: CarSortOption) => void;
  className?: string;
}

export const CarSortingOptions: React.FC<CarSortingOptionsProps> = ({
  value,
  onChange,
  className = '',
}) => {
  const t = useTranslations('cars');

  const sortOptions: Array<{
    value: CarSortOption;
    label: string;
    icon: React.ReactNode;
  }> = [
    {
      value: 'created_at_desc',
      label: t('sorting.newest'),
      icon: <Calendar className="w-4 h-4" />,
    },
    {
      value: 'price_asc',
      label: t('sorting.priceLowest'),
      icon: <DollarSign className="w-4 h-4" />,
    },
    {
      value: 'price_desc',
      label: t('sorting.priceHighest'),
      icon: <DollarSign className="w-4 h-4" />,
    },
    {
      value: 'year_desc',
      label: t('sorting.yearNewest'),
      icon: <Calendar className="w-4 h-4" />,
    },
    {
      value: 'year_asc',
      label: t('sorting.yearOldest'),
      icon: <Calendar className="w-4 h-4" />,
    },
    {
      value: 'mileage_asc',
      label: t('sorting.mileageLowest'),
      icon: <Gauge className="w-4 h-4" />,
    },
    {
      value: 'mileage_desc',
      label: t('sorting.mileageHighest'),
      icon: <Gauge className="w-4 h-4" />,
    },
    {
      value: 'price_year_ratio',
      label: t('sorting.bestValue'),
      icon: <ArrowUpDown className="w-4 h-4" />,
    },
  ];

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <label className="flex items-center gap-2 text-sm font-medium">
        <ArrowUpDown className="w-4 h-4" />
        {t('sorting.sortBy')}:
      </label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value as CarSortOption)}
        className="select select-sm select-bordered"
      >
        {sortOptions.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default CarSortingOptions;
