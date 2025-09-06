'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import {
  Sparkles,
  TrendingUp,
  MapPin,
  DollarSign,
  Star,
  Zap,
  Package,
  Clock,
} from 'lucide-react';

interface FilterPreset {
  id: string;
  icon: React.ComponentType<{ className?: string }>;
  labelKey: string;
  filters: Record<string, any>;
  categoryId?: number;
}

interface QuickFilterPresetsProps {
  onPresetSelect: (filters: Record<string, any>, categoryId?: number) => void;
  className?: string;
}

const presets: FilterPreset[] = [
  {
    id: 'new-items',
    icon: Sparkles,
    labelKey: 'newItems',
    filters: {
      sort_by: 'created_at',
      sort_order: 'desc',
      condition: 'new',
    },
  },
  {
    id: 'hot-deals',
    icon: TrendingUp,
    labelKey: 'hotDeals',
    filters: {
      has_discount: true,
      sort_by: 'discount_percentage',
      sort_order: 'desc',
    },
  },
  {
    id: 'nearby',
    icon: MapPin,
    labelKey: 'nearby',
    filters: {
      distance: 5,
      sort_by: 'distance',
      sort_order: 'asc',
    },
  },
  {
    id: 'budget-cars',
    icon: DollarSign,
    labelKey: 'budgetCars',
    categoryId: 1003,
    filters: {
      price_max: 10000,
      sort_by: 'price',
      sort_order: 'asc',
    },
  },
  {
    id: 'premium-electronics',
    icon: Star,
    labelKey: 'premiumElectronics',
    categoryId: 1002,
    filters: {
      brand: 'apple',
      condition: 'new',
      hasWarranty: true,
    },
  },
  {
    id: 'quick-sale',
    icon: Zap,
    labelKey: 'quickSale',
    filters: {
      urgent_sale: true,
      sort_by: 'price',
      sort_order: 'asc',
    },
  },
  {
    id: 'free-shipping',
    icon: Package,
    labelKey: 'freeShipping',
    filters: {
      free_shipping: true,
      product_types: ['storefront'],
    },
  },
  {
    id: 'today-added',
    icon: Clock,
    labelKey: 'todayAdded',
    filters: {
      created_today: true,
      sort_by: 'created_at',
      sort_order: 'desc',
    },
  },
];

export const QuickFilterPresets: React.FC<QuickFilterPresetsProps> = ({
  onPresetSelect,
  className = '',
}) => {
  const t = useTranslations('search');

  const handlePresetClick = (preset: FilterPreset) => {
    onPresetSelect(preset.filters, preset.categoryId);
  };

  return (
    <div className={`${className}`}>
      <div className="flex items-center gap-2 mb-3">
        <Sparkles className="w-4 h-4 text-primary" />
        <span className="text-sm font-semibold text-base-content/70">
          {t('quickFilters')}
        </span>
      </div>

      <div className="flex flex-wrap gap-2">
        {presets.map((preset) => {
          const Icon = preset.icon;
          return (
            <button
              key={preset.id}
              onClick={() => handlePresetClick(preset)}
              className="btn btn-sm btn-outline gap-1 hover:btn-primary transition-all duration-200 hover:scale-105"
            >
              <Icon className="w-3.5 h-3.5" />
              <span className="text-xs">{t(`presets.${preset.labelKey}`)}</span>
            </button>
          );
        })}
      </div>

      <div className="divider text-xs text-base-content/40 mt-4 mb-2">
        {t('orCustomizeFilters')}
      </div>
    </div>
  );
};
