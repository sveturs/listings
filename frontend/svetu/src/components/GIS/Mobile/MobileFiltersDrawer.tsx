import React, { useState, useEffect, useRef, useCallback } from 'react';
import { SearchBar } from '@/components/SearchBar';
import WalkingAccessibilityControl from '../Map/WalkingAccessibilityControl';
// import { DistrictMapSelector } from '@/components/search';
import type { Feature, Polygon } from 'geojson';
import type { MapBounds } from '@/components/GIS/types/gis';
import { SmartFilters } from '@/components/marketplace/SmartFilters';
import { QuickFilters } from '@/components/marketplace/QuickFilters';
import { CategoryTreeSelector } from '@/components/common/CategoryTreeSelector';

interface MapFilters {
  categories: number[];
  priceFrom: number;
  priceTo: number;
  radius: number;
  accessibilityMode?: 'radius' | 'walking';
  walkingTime?: number;
  attributes?: Record<string, any>;
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
  // Props для DistrictMapSelector
  enableDistrictSearch?: boolean;
  onDistrictSearchResults?: (results: any[]) => void;
  onDistrictBoundsChange?: (
    bounds: [number, number, number, number] | null
  ) => void;
  onDistrictBoundaryChange?: (boundary: Feature<Polygon> | null) => void;
  currentViewport?: {
    bounds: MapBounds;
    center: { lat: number; lng: number };
  } | null;
  searchType?: 'address' | 'district';
  onSearchTypeChange?: (type: 'address' | 'district') => void;
  translations: {
    title: string;
    search: {
      address: string;
      placeholder: string;
      byAddress?: string;
      byDistrict?: string;
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
  enableDistrictSearch: _enableDistrictSearch,
  onDistrictSearchResults: _onDistrictSearchResults,
  onDistrictBoundsChange: _onDistrictBoundsChange,
  onDistrictBoundaryChange: _onDistrictBoundaryChange,
  currentViewport: _currentViewport,
  searchType: _searchType = 'address',
  onSearchTypeChange: _onSearchTypeChange,
  translations: t,
}) => {
  const [localFilters, setLocalFilters] = useState<MapFilters>({
    ...filters,
    accessibilityMode: filters.accessibilityMode || 'radius',
    walkingTime: filters.walkingTime || 15,
  });
  const [localSearchQuery, setLocalSearchQuery] = useState(searchQuery);

  // Swipe gesture state
  const [startX, setStartX] = useState(0);
  const [currentX, setCurrentX] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const drawerRef = useRef<HTMLDivElement>(null);

  // Синхронизация с внешними фильтрами
  useEffect(() => {
    setLocalFilters({
      ...filters,
      accessibilityMode: filters.accessibilityMode || 'radius',
      walkingTime: filters.walkingTime || 15,
    });
  }, [filters]);

  useEffect(() => {
    setLocalSearchQuery(searchQuery);
  }, [searchQuery]);

  const handleLocalFiltersChange = (newFilters: Partial<MapFilters>) => {
    setLocalFilters((prev) => ({ ...prev, ...newFilters }));
  };

  const handleQuickFilterSelect = (quickFilters: Record<string, any>) => {
    setLocalFilters((prev) => ({
      ...prev,
      attributes: {
        ...prev.attributes,
        ...quickFilters,
      },
    }));
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
      categories: [],
      priceFrom: 0,
      priceTo: 0,
      radius: 10000,
      accessibilityMode: 'radius' as const,
      walkingTime: 15,
      attributes: {},
    };
    setLocalFilters(resetFilters);
    setLocalSearchQuery('');
    onFiltersChange(resetFilters);
    onSearchChange('');
  };

  const hasActiveFilters =
    localFilters.categories.length > 0 ||
    localFilters.priceFrom > 0 ||
    localFilters.priceTo > 0;

  // Swipe gesture handlers
  const handleTouchStart = useCallback(
    (e: TouchEvent) => {
      if (!isOpen) return;
      setStartX(e.touches[0].clientX);
      setCurrentX(e.touches[0].clientX);
      setIsDragging(true);
    },
    [isOpen]
  );

  const handleTouchMove = useCallback(
    (e: TouchEvent) => {
      if (!isDragging) return;
      setCurrentX(e.touches[0].clientX);
    },
    [isDragging]
  );

  const handleTouchEnd = useCallback(() => {
    if (!isDragging) return;
    setIsDragging(false);

    const deltaX = currentX - startX;
    const threshold = 100; // Minimum distance for closing

    // Swipe right to close
    if (deltaX > threshold) {
      onClose();
    }
  }, [isDragging, startX, currentX, onClose]);

  // Add touch event listeners
  useEffect(() => {
    const drawer = drawerRef.current;
    if (!drawer) return;

    drawer.addEventListener('touchstart', handleTouchStart, { passive: false });
    drawer.addEventListener('touchmove', handleTouchMove, { passive: false });
    drawer.addEventListener('touchend', handleTouchEnd);

    return () => {
      drawer.removeEventListener('touchstart', handleTouchStart);
      drawer.removeEventListener('touchmove', handleTouchMove);
      drawer.removeEventListener('touchend', handleTouchEnd);
    };
  }, [handleTouchStart, handleTouchMove, handleTouchEnd]);

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
        ref={drawerRef}
        className={`fixed inset-y-0 right-0 z-[9999] w-full max-w-sm bg-white shadow-xl transform transition-transform duration-300 md:hidden ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
        style={{
          transform: isDragging
            ? `translateX(${Math.max(0, currentX - startX)}px)`
            : isOpen
              ? 'translateX(0)'
              : 'translateX(100%)',
        }}
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
              {/* Поиск по адресу - всегда показываем, так как районы отключены */}
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
                <CategoryTreeSelector
                  value={localFilters.categories}
                  onChange={(value) => {
                    const categories = Array.isArray(value)
                      ? value
                      : value
                        ? [value]
                        : [];
                    handleLocalFiltersChange({ categories });
                  }}
                  multiple={true}
                  placeholder={t.filters.allCategories}
                  showPath={true}
                  className="w-full"
                />
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

              {/* Быстрые фильтры */}
              {localFilters.categories &&
                localFilters.categories.length > 0 && (
                  <div className="mb-4">
                    <QuickFilters
                      categoryId={localFilters.categories[0].toString()}
                      onSelectFilter={handleQuickFilterSelect}
                      className="mb-4"
                    />
                  </div>
                )}

              {/* Динамические фильтры по атрибутам категории */}
              {localFilters.categories &&
                localFilters.categories.length > 0 && (
                  <div>
                    <SmartFilters
                      categoryId={localFilters.categories[0]}
                      onChange={(attributeFilters) =>
                        handleLocalFiltersChange({
                          attributes: attributeFilters,
                        })
                      }
                      lang={
                        typeof window !== 'undefined'
                          ? window.location.pathname.split('/')[1] || 'sr'
                          : 'sr'
                      }
                      className="space-y-3"
                    />
                  </div>
                )}

              {/* Радиус поиска с WalkingAccessibilityControl */}
              <div>
                <label className="block text-sm font-medium text-base-content mb-3">
                  {t.filters.radius}
                </label>
                <WalkingAccessibilityControl
                  mode={localFilters.accessibilityMode || 'radius'}
                  onModeChange={(mode) =>
                    handleLocalFiltersChange({ accessibilityMode: mode })
                  }
                  walkingTime={localFilters.walkingTime || 15}
                  onWalkingTimeChange={(time) =>
                    handleLocalFiltersChange({ walkingTime: time })
                  }
                  searchRadius={localFilters.radius}
                  onRadiusChange={(radius) =>
                    handleLocalFiltersChange({ radius })
                  }
                />
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
