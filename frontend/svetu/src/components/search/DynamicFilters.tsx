'use client';

import React, { useEffect, useState, useRef, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { CarFilters } from '@/components/marketplace/CarFilters';
import { RealEstateFilters } from './RealEstateFilters';
import { ElectronicsFilters } from './ElectronicsFilters';
import { GenericCategoryFilters } from './GenericCategoryFilters';
import { BaseFilters } from './BaseFilters';
import { ActiveFiltersChips } from './ActiveFiltersChips';
import { Filter } from 'lucide-react';

interface DynamicFiltersProps {
  categoryId?: number;
  onFiltersChange: (filters: Record<string, any>) => void;
  activeFilters: Record<string, any>;
  className?: string;
}

const CATEGORY_IDS = {
  AUTOMOTIVE: 1003,
  REAL_ESTATE: 1004,
  ELECTRONICS: 1002,
};

// Проверяем, является ли категория автомобильной (10100-10199)
const isAutomotiveCategory = (categoryId?: number) => {
  if (!categoryId) return false;
  return categoryId === CATEGORY_IDS.AUTOMOTIVE ||
         (categoryId >= 10100 && categoryId <= 10199);
};

export const DynamicFilters: React.FC<DynamicFiltersProps> = ({
  categoryId,
  onFiltersChange,
  activeFilters,
  className = '',
}) => {
  const t = useTranslations('search');
  const [categoryFilters, setCategoryFilters] = useState<Record<string, any>>(
    {}
  );
  const [baseFilters, setBaseFilters] = useState<Record<string, any>>({});
  const [isTransitioning, setIsTransitioning] = useState(false);
  const prevCategoryIdRef = useRef<number | undefined>(categoryId);

  // Use refs to store the latest values without causing re-renders
  const baseFiltersRef = useRef(baseFilters);
  const categoryFiltersRef = useRef(categoryFilters);

  // Update refs when state changes
  useEffect(() => {
    baseFiltersRef.current = baseFilters;
  }, [baseFilters]);

  useEffect(() => {
    categoryFiltersRef.current = categoryFilters;
  }, [categoryFilters]);

  useEffect(() => {
    if (prevCategoryIdRef.current !== categoryId && categoryId !== undefined) {
      setIsTransitioning(true);
      setCategoryFilters({});
      setTimeout(() => {
        setIsTransitioning(false);
      }, 300);
    }
    prevCategoryIdRef.current = categoryId;
  }, [categoryId]);

  const handleCategoryFiltersChange = useCallback(
    (filters: Record<string, any>) => {
      setCategoryFilters(filters);
      // Use ref to get latest baseFilters value
      const combinedFilters = { ...baseFiltersRef.current, ...filters };
      onFiltersChange(combinedFilters);
    },
    [onFiltersChange]
  );

  const handleBaseFiltersChange = useCallback(
    (filters: Record<string, any>) => {
      setBaseFilters(filters);
      // Use ref to get latest categoryFilters value
      const combinedFilters = { ...filters, ...categoryFiltersRef.current };
      onFiltersChange(combinedFilters);
    },
    [onFiltersChange]
  );

  const handleRemoveFilter = useCallback(
    (key: string) => {
      const newBaseFilters = { ...baseFiltersRef.current };
      const newCategoryFilters = { ...categoryFiltersRef.current };

      delete newBaseFilters[key];
      delete newCategoryFilters[key];

      setBaseFilters(newBaseFilters);
      setCategoryFilters(newCategoryFilters);

      // Notify parent of the change
      const combinedFilters = { ...newBaseFilters, ...newCategoryFilters };
      onFiltersChange(combinedFilters);
    },
    [onFiltersChange]
  );

  const renderCategoryFilters = () => {
    if (!categoryId) return null;

    // Проверяем автомобильные категории
    if (isAutomotiveCategory(categoryId)) {
      return (
        <CarFilters
          onFiltersChange={handleCategoryFiltersChange}
          className="space-y-4"
        />
      );
    }

    switch (categoryId) {
      case CATEGORY_IDS.REAL_ESTATE:
        return (
          <RealEstateFilters
            onFiltersChange={handleCategoryFiltersChange}
            className="space-y-4"
          />
        );
      case CATEGORY_IDS.ELECTRONICS:
        return (
          <ElectronicsFilters
            onFiltersChange={handleCategoryFiltersChange}
            className="space-y-4"
          />
        );
      default:
        return (
          <GenericCategoryFilters
            categoryId={categoryId}
            onFiltersChange={handleCategoryFiltersChange}
            className="space-y-4"
          />
        );
    }
  };

  const hasActiveFilters = Object.keys(activeFilters).length > 0;

  return (
    <div className={`space-y-4 ${className}`}>
      {hasActiveFilters && (
        <ActiveFiltersChips
          filters={activeFilters}
          onRemoveFilter={handleRemoveFilter}
        />
      )}

      <div className="border-b border-base-300 pb-4">
        <h3 className="text-sm font-semibold text-base-content/70 mb-3 flex items-center gap-2">
          <Filter className="w-4 h-4" />
          {t('baseFilters')}
        </h3>
        <BaseFilters onFiltersChange={handleBaseFiltersChange} />
      </div>

      {categoryId && (
        <div
          className={`border-b border-base-300 pb-4 transition-all duration-300 ease-in-out ${
            isTransitioning ? 'opacity-0 scale-95' : 'opacity-100 scale-100'
          }`}
        >
          <h3 className="text-sm font-semibold text-base-content/70 mb-3">
            {t('categorySpecificFilters')}
          </h3>
          <div className="animate-fadeIn">{renderCategoryFilters()}</div>
        </div>
      )}
    </div>
  );
};
