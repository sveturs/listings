import React, { useState, useEffect } from 'react';
import { SearchBar } from '@/components/SearchBar';

interface MapFilters {
  category: string;
  priceFrom: number;
  priceTo: number;
  radius: number;
}

interface MobileFiltersDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  filters: MapFilters;
  onFiltersChange: (filters: Partial<MapFilters>) => void;
  searchQuery: string;
  onSearchChange: (query: string) => void;
  onSearch: (query: string) => void;
  isSearching?: boolean;
  markersCount: number;
  translations: {
    title: string;
    search: {
      address: string;
      placeholder: string;
    };
    filters: {
      category: string;
      allCategories: string;
      priceFrom: string;
      priceTo: string;
      radius: string;
    };
    categories: {
      realEstate: string;
      vehicles: string;
      electronics: string;
      clothing: string;
      services: string;
      jobs: string;
    };
    results: {
      showing: string;
      listings: string;
    };
    actions: {
      apply: string;
      reset: string;
    };
  };
}

const MobileFiltersDrawer: React.FC<MobileFiltersDrawerProps> = ({
  isOpen,
  onClose,
  filters,
  onFiltersChange,
  searchQuery,
  onSearchChange,
  onSearch,
  isSearching: _isSearching,
  markersCount,
  translations: t,
}) => {
  const [localFilters, setLocalFilters] = useState<MapFilters>(filters);
  const [localSearchQuery, setLocalSearchQuery] = useState(searchQuery);

  // Синхронизация с внешними фильтрами
  useEffect(() => {
    setLocalFilters(filters);
  }, [filters]);

  useEffect(() => {
    setLocalSearchQuery(searchQuery);
  }, [searchQuery]);

  const handleLocalFiltersChange = (newFilters: Partial<MapFilters>) => {
    setLocalFilters((prev) => ({ ...prev, ...newFilters }));
  };

  const handleApplyFilters = () => {
    onFiltersChange(localFilters);
    onSearchChange(localSearchQuery);
    if (localSearchQuery.trim()) {
      onSearch(localSearchQuery.trim());
    }
    onClose();
  };

  const handleResetFilters = () => {
    const resetFilters = {
      category: '',
      priceFrom: 0,
      priceTo: 0,
      radius: 10000,
    };
    setLocalFilters(resetFilters);
    setLocalSearchQuery('');
    onFiltersChange(resetFilters);
    onSearchChange('');
  };

  const hasActiveFilters =
    localFilters.category ||
    localFilters.priceFrom > 0 ||
    localFilters.priceTo > 0;

  return (
    <>
      {/* Backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-[9998] md:hidden"
          onClick={onClose}
          aria-label="Закрыть фильтры"
        />
      )}

      {/* Drawer */}
      <div
        className={`fixed inset-y-0 left-0 z-[9999] w-full max-w-sm bg-white shadow-xl transform transition-transform duration-300 md:hidden ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="flex flex-col h-full">
          {/* Заголовок */}
          <div className="flex items-center justify-between p-4 border-b border-base-300">
            <h2 className="text-lg font-semibold text-base-content">
              {t.title}
            </h2>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-200 rounded-full transition-all duration-200 flex-shrink-0 -mr-1 active:scale-95"
              aria-label="Закрыть фильтры"
            >
              <svg
                className="w-6 h-6 text-gray-700"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2.5}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          {/* Содержимое */}
          <div className="flex-1 overflow-y-auto overscroll-contain">
            {/* Поиск по адресу */}
            <div className="p-4 border-b border-base-300">
              <label className="block text-sm font-medium text-base-content mb-2">
                {t.search.address}
              </label>
              <SearchBar
                initialQuery={localSearchQuery}
                onSearch={(query) => setLocalSearchQuery(query)}
                placeholder={t.search.placeholder}
                className="w-full"
              />
            </div>

            {/* Фильтры */}
            <div className="p-4 space-y-6">
              {/* Категория */}
              <div>
                <label className="block text-sm font-medium text-base-content mb-2">
                  {t.filters.category}
                </label>
                <select
                  className="select select-bordered w-full"
                  value={localFilters.category}
                  onChange={(e) =>
                    handleLocalFiltersChange({ category: e.target.value })
                  }
                >
                  <option value="">{t.filters.allCategories}</option>
                  <option value="real-estate">{t.categories.realEstate}</option>
                  <option value="vehicles">{t.categories.vehicles}</option>
                  <option value="electronics">
                    {t.categories.electronics}
                  </option>
                  <option value="clothing">{t.categories.clothing}</option>
                  <option value="services">{t.categories.services}</option>
                  <option value="jobs">{t.categories.jobs}</option>
                </select>
              </div>

              {/* Цена от */}
              <div>
                <label className="block text-sm font-medium text-base-content mb-2">
                  {t.filters.priceFrom}
                </label>
                <input
                  type="number"
                  className="input input-bordered w-full"
                  value={localFilters.priceFrom || ''}
                  onChange={(e) =>
                    handleLocalFiltersChange({
                      priceFrom: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                />
              </div>

              {/* Цена до */}
              <div>
                <label className="block text-sm font-medium text-base-content mb-2">
                  {t.filters.priceTo}
                </label>
                <input
                  type="number"
                  className="input input-bordered w-full"
                  value={localFilters.priceTo || ''}
                  onChange={(e) =>
                    handleLocalFiltersChange({
                      priceTo: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="∞"
                />
              </div>

              {/* Радиус поиска */}
              <div>
                <label className="block text-sm font-medium text-base-content mb-2">
                  {t.filters.radius}: {Math.round(localFilters.radius / 1000)}{' '}
                  км
                </label>
                <input
                  type="range"
                  className="range range-primary w-full"
                  min="1000"
                  max="50000"
                  step="1000"
                  value={localFilters.radius}
                  onChange={(e) =>
                    handleLocalFiltersChange({
                      radius: parseInt(e.target.value),
                    })
                  }
                />
                <div className="flex justify-between text-xs text-base-content-secondary mt-1">
                  <span>1 км</span>
                  <span>50 км</span>
                </div>
              </div>
            </div>
          </div>

          {/* Нижняя панель с кнопками */}
          <div className="p-4 border-t border-base-300 bg-white shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)] flex-shrink-0">
            {/* Статистика */}
            <div className="text-sm text-base-content-secondary mb-4 text-center font-medium">
              {t.results.showing}:{' '}
              <span className="text-base-content">{markersCount}</span>{' '}
              {t.results.listings}
            </div>

            {/* Кнопки */}
            <div className="flex gap-3">
              <button
                onClick={handleResetFilters}
                className={`flex-1 px-4 py-3 text-sm font-medium rounded-lg transition-all duration-200 active:scale-95 ${
                  hasActiveFilters
                    ? 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                }`}
                disabled={!hasActiveFilters}
              >
                {t.actions.reset}
              </button>
              <button
                onClick={handleApplyFilters}
                className="flex-1 px-4 py-3 text-sm font-medium bg-primary text-white rounded-lg hover:bg-primary/90 transition-all duration-200 active:scale-95 shadow-sm"
              >
                {t.actions.apply}
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default MobileFiltersDrawer;
