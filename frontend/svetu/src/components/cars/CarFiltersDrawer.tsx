'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { X, Filter, RotateCcw, Check } from 'lucide-react';
import { CarFilters } from '@/components/marketplace/CarFilters';

interface CarFiltersDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  onApply: (filters: any) => void;
  onReset: () => void;
  currentFilters: any;
  activeFilterCount: number;
}

export const CarFiltersDrawer: React.FC<CarFiltersDrawerProps> = ({
  isOpen,
  onClose,
  onApply,
  onReset,
  currentFilters,
  activeFilterCount,
}) => {
  const t = useTranslations('cars');
  const [tempFilters, setTempFilters] = useState(currentFilters);
  const [hasChanges, setHasChanges] = useState(false);

  useEffect(() => {
    setTempFilters(currentFilters);
    setHasChanges(false);
  }, [currentFilters]);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  const handleFilterChange = (newFilters: any) => {
    setTempFilters(newFilters);
    setHasChanges(
      JSON.stringify(newFilters) !== JSON.stringify(currentFilters)
    );
  };

  const handleApply = () => {
    onApply(tempFilters);
    onClose();
  };

  const handleReset = () => {
    setTempFilters({});
    setHasChanges(true);
    onReset();
  };

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <>
      <div
        className={`fixed inset-0 bg-black/50 z-40 transition-opacity duration-300 lg:hidden ${
          isOpen
            ? 'opacity-100 pointer-events-auto'
            : 'opacity-0 pointer-events-none'
        }`}
        onClick={handleBackdropClick}
      />

      <div
        className={`fixed bottom-0 left-0 right-0 bg-base-100 z-50 transform transition-transform duration-300 ease-out lg:hidden ${
          isOpen ? 'translate-y-0' : 'translate-y-full'
        }`}
        style={{
          maxHeight: '85vh',
          borderTopLeftRadius: '1.5rem',
          borderTopRightRadius: '1.5rem',
        }}
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center justify-center pt-3 pb-2">
            <div className="w-12 h-1 bg-base-300 rounded-full" />
          </div>

          <div className="flex items-center justify-between px-4 pb-3 border-b">
            <div className="flex items-center gap-2">
              <Filter className="w-5 h-5 text-primary" />
              <h2 className="text-lg font-semibold">{t('filters.title')}</h2>
              {activeFilterCount > 0 && (
                <span className="badge badge-primary badge-sm">
                  {activeFilterCount}
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

          <div className="flex-1 overflow-y-auto px-4 py-4">
            <CarFilters onFiltersChange={handleFilterChange} />
          </div>

          <div className="border-t px-4 py-3">
            <div className="flex gap-2">
              <button
                onClick={handleReset}
                className="btn btn-ghost flex-1"
                disabled={activeFilterCount === 0}
              >
                <RotateCcw className="w-4 h-4" />
                {t('filters.reset')}
              </button>
              <button
                onClick={handleApply}
                className="btn btn-primary flex-1"
                disabled={!hasChanges}
              >
                <Check className="w-4 h-4" />
                {t('filters.apply')}
                {hasChanges && (
                  <span className="badge badge-sm badge-secondary">
                    {t('common.new')}
                  </span>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Floating filter button for mobile */}
      {!isOpen && (
        <button
          onClick={() => onClose()}
          className="fixed bottom-4 left-1/2 transform -translate-x-1/2 btn btn-primary shadow-lg lg:hidden z-30"
        >
          <Filter className="w-5 h-5" />
          {t('filters.showFilters')}
          {activeFilterCount > 0 && (
            <span className="badge badge-secondary badge-sm">
              {activeFilterCount}
            </span>
          )}
        </button>
      )}
    </>
  );
};

export default CarFiltersDrawer;
