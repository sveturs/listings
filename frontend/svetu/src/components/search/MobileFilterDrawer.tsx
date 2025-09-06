'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { X, Filter } from 'lucide-react';
import { CategorySelector } from './CategorySelector';
import { DynamicFilters } from './DynamicFilters';
import { QuickFilterPresets } from './QuickFilterPresets';

interface MobileFilterDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  selectedCategoryId?: number;
  onCategorySelect: (categoryId: number | undefined) => void;
  filters: Record<string, any>;
  onFiltersChange: (filters: Record<string, any>) => void;
  activeFiltersCount: number;
}

export const MobileFilterDrawer: React.FC<MobileFilterDrawerProps> = ({
  isOpen,
  onClose,
  selectedCategoryId,
  onCategorySelect,
  filters,
  onFiltersChange,
  activeFiltersCount,
}) => {
  const t = useTranslations('search');

  if (!isOpen) return null;

  return (
    <>
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/50 z-40 lg:hidden animate-fadeIn"
        onClick={onClose}
      />

      {/* Drawer */}
      <div
        className={`fixed inset-y-0 right-0 w-full max-w-sm bg-base-100 z-50 transform transition-transform duration-300 ease-in-out lg:hidden ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-base-200">
          <div className="flex items-center gap-2">
            <Filter className="w-5 h-5 text-primary" />
            <h2 className="text-lg font-semibold">{t('filters')}</h2>
            {activeFiltersCount > 0 && (
              <span className="badge badge-primary badge-sm">
                {activeFiltersCount}
              </span>
            )}
          </div>
          <button
            onClick={onClose}
            className="btn btn-ghost btn-sm btn-circle"
            aria-label="Close filters"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="overflow-y-auto h-[calc(100vh-64px-56px)]">
          <div className="p-4 space-y-6">
            {/* Category Selector */}
            <div className="bg-base-200 rounded-lg p-4">
              <CategorySelector
                selectedCategoryId={selectedCategoryId}
                onCategorySelect={onCategorySelect}
              />
            </div>

            {/* Quick Presets */}
            <div className="bg-base-200 rounded-lg p-4">
              <QuickFilterPresets
                onPresetSelect={(presetFilters, categoryId) => {
                  if (categoryId) {
                    onCategorySelect(categoryId);
                  }
                  onFiltersChange({ ...filters, ...presetFilters });
                }}
              />
            </div>

            {/* Dynamic Filters */}
            <div className="bg-base-200 rounded-lg p-4">
              <DynamicFilters
                categoryId={selectedCategoryId}
                onFiltersChange={onFiltersChange}
                activeFilters={filters}
              />
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-base-200 bg-base-100">
          <div className="flex gap-2">
            <button
              onClick={() => {
                onFiltersChange({});
                onCategorySelect(undefined);
              }}
              className="btn btn-outline flex-1"
            >
              {t('clearAll')}
            </button>
            <button onClick={onClose} className="btn btn-primary flex-1">
              {t('apply')}
            </button>
          </div>
        </div>
      </div>
    </>
  );
};
