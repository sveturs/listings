'use client';

import React, { useEffect, useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { DistrictSelector } from './DistrictSelector';
import {
  useVisibleCities,
  formatDistance,
} from '@/components/GIS/hooks/useVisibleCities';
import type { MapBounds } from '@/components/GIS/types/gis';
import type { components as _components } from '@/types/generated/api';
import type { Feature, Polygon } from 'geojson';

// Временные интерфейсы до исправления API типов
interface District {
  id: string;
  name: string;
  geometry?: any;
  boundary?: {
    coordinates: number[][][];
  };
  bounds?: [number, number, number, number];
  population?: number;
  area?: number;
  area_km2?: number;
}

interface SpatialSearchResult {
  id: string;
  title: string;
  description?: string;
  latitude: number;
  longitude: number;
  distance?: number;
  category?: string;
  price?: number;
  currency?: string;
  imageUrl?: string;
  first_image_url?: string;
  category_name?: string;
  address?: string;
  user_email?: string;
}

interface DistrictMapSelectorProps {
  onSearchResults?: (results: SpatialSearchResult[]) => void;
  onDistrictBoundsChange?: (
    bounds: [number, number, number, number] | null
  ) => void;
  onDistrictBoundaryChange?: (boundary: Feature<Polygon> | null) => void;
  onViewportChange?: (
    bounds: MapBounds,
    center: { lat: number; lng: number }
  ) => void;
  currentViewport?: {
    bounds: MapBounds;
    center: { lat: number; lng: number };
  } | null;
  className?: string;
}

export function DistrictMapSelector({
  onSearchResults,
  onDistrictBoundsChange,
  onDistrictBoundaryChange,
  onViewportChange,
  currentViewport,
  className = '',
}: DistrictMapSelectorProps) {
  const t = useTranslations('search');
  const [selectedDistrictId, setSelectedDistrictId] = useState<string | null>(
    null
  );
  const [selectedMunicipalityId, setSelectedMunicipalityId] = useState<
    string | null
  >(null);
  const [isSearching, setIsSearching] = useState(false);
  const [currentDistrict, setCurrentDistrict] = useState<District | null>(null);

  // Используем новый хук для отслеживания видимых городов
  const {
    visibleCities,
    closestCity,
    loading: citiesLoading,
    error: citiesError,
    updateViewport,
    refreshCities,
    shouldShowDistrictSearch,
    hasDistrictsInViewport,
  } = useVisibleCities();

  // Обработчик изменения viewport карты
  const _handleViewportChange = useCallback(
    (bounds: MapBounds, center: { lat: number; lng: number }) => {
      updateViewport(bounds, center);
      onViewportChange?.(bounds, center);
    },
    [updateViewport, onViewportChange]
  );

  // Отслеживание изменений viewport от MapPage
  useEffect(() => {
    if (currentViewport) {
      console.log('[DistrictMapSelector] Updating viewport from MapPage:', currentViewport);
      updateViewport(currentViewport.bounds, currentViewport.center);
    }
  }, [currentViewport, updateViewport]);

  // Загрузка информации о выбранном районе
  useEffect(() => {
    if (!selectedDistrictId) {
      setCurrentDistrict(null);
      onDistrictBoundsChange?.(null);
      return;
    }

    const fetchDistrictDetails = async () => {
      try {
        const response = await fetch(
          `/api/v1/gis/districts/${selectedDistrictId}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch district details');
        }
        const data = await response.json();
        const district = data.data as District;
        setCurrentDistrict(district);

        // Если у района есть границы, передаем их для отображения на карте
        if (district.boundary?.coordinates?.[0]) {
          const coords = district.boundary.coordinates[0];
          let minLng = 180,
            maxLng = -180,
            minLat = 90,
            maxLat = -90;

          coords.forEach((coord) => {
            minLng = Math.min(minLng, coord[0]);
            maxLng = Math.max(maxLng, coord[0]);
            minLat = Math.min(minLat, coord[1]);
            maxLat = Math.max(maxLat, coord[1]);
          });

          onDistrictBoundsChange?.([minLng, minLat, maxLng, maxLat]);
        }
      } catch (err) {
        console.error('Error fetching district details:', err);
      }
    };

    fetchDistrictDetails();
  }, [selectedDistrictId, onDistrictBoundsChange]);

  // Поиск объявлений при изменении выбора
  useEffect(() => {
    const searchListings = async () => {
      if (!selectedDistrictId && !selectedMunicipalityId) {
        onSearchResults?.([]);
        return;
      }

      setIsSearching(true);
      try {
        let url = '';
        if (selectedMunicipalityId) {
          url = `/api/v1/gis/search/by-municipality/${selectedMunicipalityId}`;
        } else if (selectedDistrictId) {
          url = `/api/v1/gis/search/by-district/${selectedDistrictId}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to search listings');
        }

        const data = await response.json();
        onSearchResults?.(data.data || []);
      } catch (err) {
        console.error('Error searching listings:', err);
        onSearchResults?.([]);
      } finally {
        setIsSearching(false);
      }
    };

    searchListings();
  }, [selectedDistrictId, selectedMunicipalityId, onSearchResults]);

  /**
   * Загружает границы района в формате GeoJSON
   */
  const loadDistrictBoundary = useCallback(
    async (districtId: string): Promise<Feature<Polygon> | null> => {
      try {
        const response = await fetch(
          `/api/v1/gis/districts/${encodeURIComponent(districtId)}/boundary`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (!data.success || !data.data) {
          throw new Error(data.error || 'Failed to fetch district boundary');
        }

        // data.data содержит объект с полем boundary
        const boundaryData = data.data.boundary;
        
        // Если boundary уже является объектом, используем его напрямую
        // Если это строка, парсим её
        const geoJson = typeof boundaryData === 'string' 
          ? JSON.parse(boundaryData) 
          : boundaryData;

        // Если это просто Polygon, оборачиваем в Feature
        if (geoJson.type === 'Polygon') {
          return {
            type: 'Feature',
            geometry: geoJson,
            properties: {
              id: data.data.id,
              name: data.data.name
            }
          } as Feature<Polygon>;
        }

        return geoJson as Feature<Polygon>;
      } catch (error) {
        console.error('Error loading district boundary:', error);
        return null;
      }
    },
    []
  );

  const handleDistrictChange = useCallback(
    async (districtId: string | null) => {
      setSelectedDistrictId(districtId);
      setSelectedMunicipalityId(null);

      // Загружаем границы выбранного района
      if (districtId && onDistrictBoundaryChange) {
        const boundary = await loadDistrictBoundary(districtId);
        onDistrictBoundaryChange(boundary);
      } else if (onDistrictBoundaryChange) {
        // Убираем границы если район не выбран
        onDistrictBoundaryChange(null);
      }
    },
    [loadDistrictBoundary, onDistrictBoundaryChange]
  );

  const handleMunicipalityChange = (municipalityId: string | null) => {
    setSelectedMunicipalityId(municipalityId);
  };

  // Не показываем компонент если нет городов с районами в viewport
  if (!shouldShowDistrictSearch && !citiesLoading) {
    return null;
  }

  return (
    <div className={`card bg-base-100 shadow-lg ${className}`}>
      <div className="card-body">
        {/* Заголовок с информацией о текущем городе */}
        <div className="flex items-center justify-between mb-4">
          <h3 className="card-title text-lg">
            {closestCity ? (
              <>
                {t('searchByDistrict')} — {closestCity.city.name}
                <div className="badge badge-outline badge-sm ml-2">
                  {formatDistance(closestCity.distance)}
                </div>
              </>
            ) : (
              t('searchByDistrict')
            )}
          </h3>

          {visibleCities.length > 1 && (
            <button
              className="btn btn-ghost btn-sm"
              onClick={refreshCities}
              disabled={citiesLoading}
            >
              {citiesLoading ? (
                <span className="loading loading-spinner loading-sm"></span>
              ) : (
                '🔄'
              )}
            </button>
          )}
        </div>

        {/* Информация о видимых городах */}
        {visibleCities.length > 0 && (
          <div className="alert alert-info mb-4">
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
            <span>
              {hasDistrictsInViewport
                ? t('citiesWithDistrictsFound', {
                    count: visibleCities.filter((c) => c.city.has_districts)
                      .length,
                  })
                : t('noCitiesWithDistricts')}
            </span>
          </div>
        )}

        {/* Селектор районов - показываем только если есть доступные районы */}
        {shouldShowDistrictSearch && (
          <DistrictSelector
            selectedDistrictId={selectedDistrictId || undefined}
            selectedMunicipalityId={selectedMunicipalityId || undefined}
            onDistrictChange={handleDistrictChange}
            onMunicipalityChange={handleMunicipalityChange}
          />
        )}

        {/* Информация о выбранном районе */}
        {currentDistrict && (
          <div className="stats stats-vertical lg:stats-horizontal shadow mt-4">
            {currentDistrict.population && (
              <div className="stat">
                <div className="stat-title">{t('population')}</div>
                <div className="stat-value text-base">
                  {currentDistrict.population.toLocaleString()}
                </div>
              </div>
            )}
            {currentDistrict.area_km2 && (
              <div className="stat">
                <div className="stat-title">{t('area')}</div>
                <div className="stat-value text-base">
                  {currentDistrict.area_km2} km²
                </div>
              </div>
            )}
          </div>
        )}

        {/* Индикаторы загрузки */}
        {(isSearching || citiesLoading) && (
          <div className="flex items-center justify-center mt-4">
            <span className="loading loading-spinner loading-md"></span>
            <span className="ml-2">
              {citiesLoading ? t('loadingCities') : t('searching')}
            </span>
          </div>
        )}

        {/* Ошибки */}
        {citiesError && (
          <div className="alert alert-error mt-4">
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
            <span>{citiesError}</span>
          </div>
        )}
      </div>
    </div>
  );
}
