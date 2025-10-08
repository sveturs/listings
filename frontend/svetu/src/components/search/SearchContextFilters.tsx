'use client';

import React, { Suspense } from 'react';
import { useTranslations } from 'next-intl';
import { SearchContextConfig } from '@/types/searchContext';
import { BaseFilters } from './BaseFilters';
import { CarFilters } from '@/components/c2c/CarFilters';
import { RealEstateFilters } from './RealEstateFilters';
import { ElectronicsFilters } from './ElectronicsFilters';

interface SearchContextFiltersProps {
  context: SearchContextConfig;
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

// Маппинг компонентов фильтров
const filterComponents: Record<string, React.ComponentType<any>> = {
  BaseFilters,
  CarFilters,
  RealEstateFilters,
  ElectronicsFilters,
  // Эти компоненты будут созданы позже
  ServiceFilters: BaseFilters, // Временно используем BaseFilters
  FashionFilters: BaseFilters,
  JobFilters: BaseFilters,
};

export const SearchContextFilters: React.FC<SearchContextFiltersProps> = ({
  context,
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations();

  // Получаем список компонентов для текущего контекста
  const componentsToRender = context.filterComponents || ['BaseFilters'];

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Быстрые фильтры контекста */}
      {context.quickFilters && context.quickFilters.length > 0 && (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title text-sm">Quick Filters</h3>
            <div className="flex flex-wrap gap-2">
              {context.quickFilters.map((filter) => (
                <button
                  key={filter.id}
                  onClick={() => onFiltersChange(filter.filters)}
                  className={`btn btn-sm btn-outline btn-${context.accentColor || 'primary'}`}
                >
                  {t(filter.label)}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Рендерим компоненты фильтров */}
      {componentsToRender.map((componentName) => {
        const Component = filterComponents[componentName];

        if (!Component) {
          console.warn(`Filter component ${componentName} not found`);
          return null;
        }

        // Специальная обработка для CarFilters в автомобильном контексте
        if (componentName === 'CarFilters' && context.id === 'automotive') {
          return (
            <Suspense
              key={componentName}
              fallback={
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="loading loading-spinner loading-md"></div>
                  </div>
                </div>
              }
            >
              <Component onFiltersChange={onFiltersChange} className="w-full" />
            </Suspense>
          );
        }

        // Для остальных компонентов
        return (
          <Suspense
            key={componentName}
            fallback={
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <div className="loading loading-spinner loading-md"></div>
                </div>
              </div>
            }
          >
            <Component onFiltersChange={onFiltersChange} className="w-full" />
          </Suspense>
        );
      })}
    </div>
  );
};
