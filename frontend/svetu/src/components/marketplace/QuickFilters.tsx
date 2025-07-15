'use client';

import React from 'react';
import { useTranslations } from 'next-intl';

interface QuickFilter {
  id: string;
  label: string;
  icon?: string;
  filters: Record<string, any>;
}

interface QuickFiltersProps {
  categoryId: string | null;
  onSelectFilter: (filters: Record<string, any>) => void;
  className?: string;
}

// Предопределенные быстрые фильтры для категорий
const quickFiltersByCategory: Record<string, QuickFilter[]> = {
  // Автомобили
  '2000': [
    {
      id: 'new-cars',
      label: 'filters.quickOptions.new',
      icon: '✨',
      filters: {
        condition: 'new',
      },
    },
    {
      id: 'low-mileage',
      label: 'filters.quickOptions.lowMileage',
      icon: '🚗',
      filters: {
        mileage: { max: 50000 },
      },
    },
    {
      id: 'no-damage',
      label: 'filters.quickOptions.noDamage',
      icon: '✅',
      filters: {
        damaged: false,
      },
    },
  ],
  // Недвижимость - квартиры
  '1100': [
    {
      id: 'with-photo',
      label: 'filters.quickOptions.withPhoto',
      icon: '📸',
      filters: {
        has_photos: true,
      },
    },
    {
      id: 'new-building',
      label: 'filters.quickOptions.newBuilding',
      icon: '🏗️',
      filters: {
        building_type: 'new',
      },
    },
    {
      id: 'with-balcony',
      label: 'filters.quickOptions.withBalcony',
      icon: '🌇',
      filters: {
        has_balcony: true,
      },
    },
  ],
  // Недвижимость - комнаты
  '1200': [
    {
      id: 'with-photo',
      label: 'filters.quickOptions.withPhoto',
      icon: '📸',
      filters: {
        has_photos: true,
      },
    },
    {
      id: 'furnished',
      label: 'filters.quickOptions.furnished',
      icon: '🛋️',
      filters: {
        furnished: true,
      },
    },
    {
      id: 'with-parking',
      label: 'filters.quickOptions.withParking',
      icon: '🚗',
      filters: {
        has_parking: true,
      },
    },
  ],
  // Недвижимость - дома
  '1300': [
    {
      id: 'with-photo',
      label: 'filters.quickOptions.withPhoto',
      icon: '📸',
      filters: {
        has_photos: true,
      },
    },
    {
      id: 'with-garden',
      label: 'filters.quickOptions.withGarden',
      icon: '🌳',
      filters: {
        has_garden: true,
      },
    },
    {
      id: 'with-garage',
      label: 'filters.quickOptions.withGarage',
      icon: '🚙',
      filters: {
        has_garage: true,
      },
    },
  ],
  // Электроника
  '3000': [
    {
      id: 'with-warranty',
      label: 'filters.quickOptions.withWarranty',
      icon: '🛡️',
      filters: {
        has_warranty: true,
      },
    },
    {
      id: 'like-new',
      label: 'filters.quickOptions.likeNew',
      icon: '⭐',
      filters: {
        condition: 'excellent',
      },
    },
    {
      id: 'in-box',
      label: 'filters.quickOptions.inBox',
      icon: '📦',
      filters: {
        original_packaging: true,
      },
    },
  ],
  // Работа
  '9000': [
    {
      id: 'full-time',
      label: 'filters.quickOptions.fullTime',
      icon: '💼',
      filters: {
        employment_type: 'full_time',
      },
    },
    {
      id: 'remote',
      label: 'filters.quickOptions.remote',
      icon: '🏠',
      filters: {
        remote_work: true,
      },
    },
    {
      id: 'with-experience',
      label: 'filters.quickOptions.noExperience',
      icon: '🎓',
      filters: {
        experience_required: 0,
      },
    },
  ],
};

export function QuickFilters({
  categoryId,
  onSelectFilter,
  className = '',
}: QuickFiltersProps) {
  const t = useTranslations('map');

  if (!categoryId) return null;

  const filters = quickFiltersByCategory[categoryId] || [];

  if (filters.length === 0) return null;

  return (
    <div className={`${className}`}>
      <h4 className="text-sm font-medium text-base-content mb-2">
        {t('filters.quickFilters')}
      </h4>
      <div className="flex flex-wrap gap-2">
        {filters.map((filter) => (
          <button
            key={filter.id}
            onClick={() => onSelectFilter(filter.filters)}
            className="btn btn-sm btn-outline hover:btn-primary group"
          >
            {filter.icon && <span className="text-base">{filter.icon}</span>}
            <span className="text-xs">{t(filter.label)}</span>
          </button>
        ))}
      </div>
    </div>
  );
}
