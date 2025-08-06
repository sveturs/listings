'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import MarketplaceList from './MarketplaceList';
import {
  UnifiedSearchItem,
  UnifiedSearchService,
} from '@/services/unifiedSearch';
import { RadiusSearchResult } from '@/components/GIS/types/gis';
import AdvancedGeoFilters from '@/components/GIS/AdvancedGeoFilters';
import { useAdvancedGeoFilters } from '@/hooks/useAdvancedGeoFilters';
import {
  MapIcon,
  AdjustmentsHorizontalIcon,
} from '@heroicons/react/24/outline';

// Редирект на основную карту
import { Link } from '@/i18n/routing';

const MapRedirectComponent = () => (
  <div className="mb-8 p-4 bg-base-200 rounded-lg text-center">
    <p className="text-base-content/70 mb-2">Переходим на полную карту...</p>
    <Link href="/map" className="btn btn-primary">
      Открыть карту
    </Link>
  </div>
);

interface HomePageWithFiltersProps {
  initialData: {
    items: UnifiedSearchItem[];
    total: number;
    page: number;
    limit: number;
    has_more: boolean;
  } | null;
  locale: string;
  error?: Error | null;
  paymentsEnabled?: boolean;
}

export default function HomePageWithFilters({
  initialData,
  locale,
  error,
  paymentsEnabled = false,
}: HomePageWithFiltersProps) {
  const t = useTranslations('marketplace.home');
  const [showMap, setShowMap] = useState(false);
  const [showFilters, setShowFilters] = useState(false);
  const [currentLocation, setCurrentLocation] = useState<
    { lat: number; lng: number } | undefined
  >();
  const [filteredData, setFilteredData] = useState(initialData);
  const [isFiltering, setIsFiltering] = useState(false);

  const { filters, setFilters } = useAdvancedGeoFilters();

  // Получение текущего местоположения пользователя
  useEffect(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setCurrentLocation({
            lat: position.coords.latitude,
            lng: position.coords.longitude,
          });
        },
        (error) => {
          console.error('Error getting location:', error);
        }
      );
    }
  }, []);

  // Применение расширенных геофильтров
  const applyFilters = useCallback(async () => {
    if (!initialData || Object.keys(filters).length === 0) {
      setFilteredData(initialData);
      return;
    }

    setIsFiltering(true);

    try {
      // Преобразуем фильтры в нужный формат
      const transformedFilters: any = {};

      if (filters.travelTime) {
        transformedFilters.travel_time = {
          center_lat: filters.travelTime.centerLat,
          center_lng: filters.travelTime.centerLng,
          max_minutes: filters.travelTime.maxMinutes,
          transport_mode: filters.travelTime.transportMode,
        };
      }

      if (filters.poiFilter) {
        transformedFilters.poi_filter = {
          poi_type: filters.poiFilter.poiType,
          max_distance: filters.poiFilter.maxDistance,
          min_count: filters.poiFilter.minCount,
        };
      }

      if (filters.densityFilter) {
        transformedFilters.density_filter = {
          avoid_crowded: filters.densityFilter.avoidCrowded,
          max_density: filters.densityFilter.maxDensity,
          min_density: filters.densityFilter.minDensity,
        };
      }

      // Получаем результаты с применением расширенных фильтров
      const result = await UnifiedSearchService.search({
        query: '',
        sort_by: 'date',
        sort_order: 'desc',
        page: 1,
        limit: 20,
        advanced_geo_filters: transformedFilters,
      });

      setFilteredData(result);
    } catch (error) {
      console.error('Error applying filters:', error);
      // При ошибке показываем исходные данные
      setFilteredData(initialData);
    } finally {
      setIsFiltering(false);
    }
  }, [filters, initialData]);

  // Слушаем события изменения фильтров
  useEffect(() => {
    const handleFiltersChanged = (_event: CustomEvent) => {
      applyFilters();
    };

    window.addEventListener(
      'advancedGeoFiltersChanged',
      handleFiltersChanged as EventListener
    );

    return () => {
      window.removeEventListener(
        'advancedGeoFiltersChanged',
        handleFiltersChanged as EventListener
      );
    };
  }, [applyFilters]);

  const _handleListingSelect = (listing: RadiusSearchResult) => {
    console.log('Selected listing:', listing);
  };

  return (
    <div className="flex gap-6">
      {/* Боковая панель с фильтрами для десктопа */}
      <aside
        className={`hidden lg:block w-80 ${showFilters ? '' : 'lg:hidden'}`}
      >
        <div className="sticky top-20">
          <AdvancedGeoFilters
            onFiltersChange={setFilters}
            currentLocation={currentLocation}
            className="shadow-lg"
          />
        </div>
      </aside>

      {/* Основной контент */}
      <div className="flex-1">
        {paymentsEnabled && (
          <div className="alert alert-info mb-8 shadow-lg">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <span>🎉 {t('paymentsNowAvailable')}</span>
          </div>
        )}

        {error && (
          <div className="alert alert-error mb-8 shadow-lg">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="stroke-current shrink-0 h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span>{t('errorLoadingData')}</span>
          </div>
        )}

        {/* Заголовок и контролы */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-base-content">
            {t('latestListings')}
            {isFiltering && (
              <span className="loading loading-spinner loading-sm ml-2"></span>
            )}
          </h2>

          <div className="flex items-center gap-2">
            {/* Кнопка фильтров для мобильных */}
            <button
              onClick={() => setShowFilters(!showFilters)}
              className="btn btn-ghost btn-sm lg:hidden"
            >
              <AdjustmentsHorizontalIcon className="h-5 w-5" />
              {t('filters')}
            </button>

            {/* Переключатель карта/список */}
            <div className="flex items-center space-x-2">
              <span className="text-sm text-base-content/70">{t('view')}:</span>
              <button
                onClick={() => setShowMap(!showMap)}
                className={`btn btn-sm ${showMap ? 'btn-primary' : 'btn-ghost'}`}
              >
                <MapIcon className="h-4 w-4 mr-1" />
                {t('map')}
              </button>
            </div>
          </div>
        </div>

        {/* Мобильная панель фильтров */}
        {showFilters && (
          <div className="lg:hidden mb-6">
            <AdvancedGeoFilters
              onFiltersChange={setFilters}
              currentLocation={currentLocation}
              className="shadow-lg"
            />
          </div>
        )}

        {/* Индикатор активных фильтров */}
        {Object.keys(filters).length > 0 && (
          <div className="alert alert-info mb-4">
            <span className="text-sm">
              {t('activeFilters')}: {Object.keys(filters).length}
            </span>
            <button
              onClick={() => setFilters({})}
              className="btn btn-ghost btn-xs"
            >
              {t('clearAll')}
            </button>
          </div>
        )}

        {/* Контент - карта или список */}
        {showMap ? (
          <MapRedirectComponent />
        ) : (
          <MarketplaceList initialData={filteredData} locale={locale} />
        )}
      </div>
    </div>
  );
}
