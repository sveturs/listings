'use client';

import React, { useEffect, useState, useCallback, useRef } from 'react';
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

  // Предыдущие значения для отслеживания изменений
  const prevDistrictIdRef = useRef<string | null>(null);
  const prevMunicipalityIdRef = useRef<string | null>(null);

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
      updateViewport(currentViewport.bounds, currentViewport.center);
    }
  }, [currentViewport, updateViewport]);

  // Кэш для районов, чтобы избежать повторных запросов
  const [districtCache, setDistrictCache] = useState<Map<string, District>>(
    new Map()
  );
  const [loadingDistricts, setLoadingDistricts] = useState<Set<string>>(
    new Set()
  );

  // Загрузка информации о выбранном районе (только для отображения статистики)
  useEffect(() => {
    if (!selectedDistrictId) {
      setCurrentDistrict(null);
      return;
    }

    // Проверяем кэш
    const cachedDistrict = districtCache.get(selectedDistrictId);
    if (cachedDistrict) {
      console.log(
        '📋 Using cached district data for stats:',
        selectedDistrictId
      );
      setCurrentDistrict(cachedDistrict);
      return;
    }

    // Проверяем, не загружается ли уже этот район
    if (loadingDistricts.has(selectedDistrictId)) {
      console.log('⏳ District already loading for stats:', selectedDistrictId);
      return;
    }

    const fetchDistrictDetails = async () => {
      console.log(
        '📡 Fetching district details for stats:',
        selectedDistrictId
      );

      // Добавляем в список загружаемых
      setLoadingDistricts((prev) => new Set(prev).add(selectedDistrictId));

      try {
        const response = await fetch(
          `/api/v1/gis/districts/${selectedDistrictId}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch district details');
        }
        const data = await response.json();
        const district = data.data as District;

        // Сохраняем в кэш
        setDistrictCache((prev) =>
          new Map(prev).set(selectedDistrictId, district)
        );
        setCurrentDistrict(district);
      } catch (err) {
        console.error('❌ Error fetching district details:', err);
      } finally {
        // Убираем из списка загружаемых
        setLoadingDistricts((prev) => {
          const newSet = new Set(prev);
          newSet.delete(selectedDistrictId);
          return newSet;
        });
      }
    };

    fetchDistrictDetails();
  }, [selectedDistrictId, districtCache, loadingDistricts]);

  // Debounce для поиска объявлений
  const searchTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Поиск объявлений при изменении выбора
  useEffect(() => {
    // Проверяем, действительно ли изменились значения
    const districtChanged = prevDistrictIdRef.current !== selectedDistrictId;
    const municipalityChanged =
      prevMunicipalityIdRef.current !== selectedMunicipalityId;

    if (!districtChanged && !municipalityChanged) {
      console.log('🔍 No actual changes in selection, skipping search');
      return;
    }

    // Обновляем предыдущие значения
    prevDistrictIdRef.current = selectedDistrictId;
    prevMunicipalityIdRef.current = selectedMunicipalityId;

    // Очищаем предыдущий таймер
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current);
    }

    const searchListings = async () => {
      if (!selectedDistrictId && !selectedMunicipalityId) {
        console.log('🔍 No district/municipality selected, clearing results');
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

        console.log('🔍 Searching listings:', url);
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to search listings');
        }

        const data = await response.json();
        console.log('📦 Search results:', data.data?.length || 0, 'items');
        onSearchResults?.(data.data || []);
      } catch (err) {
        console.error('Error searching listings:', err);
        onSearchResults?.([]);
      } finally {
        setIsSearching(false);
      }
    };

    // Debounce поиск на 500ms
    searchTimeoutRef.current = setTimeout(() => {
      searchListings();
    }, 500);

    // Cleanup функция
    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, [selectedDistrictId, selectedMunicipalityId, onSearchResults]);

  // Кэш для границ районов
  const [boundaryCache, setBoundaryCache] = useState<
    Map<string, Feature<Polygon> | null>
  >(new Map());
  const [loadingBoundaries, setLoadingBoundaries] = useState<Set<string>>(
    new Set()
  );

  /**
   * Загружает границы района в формате GeoJSON
   */
  const loadDistrictBoundary = useCallback(
    async (districtId: string): Promise<Feature<Polygon> | null> => {
      // Проверяем кэш
      const cachedBoundary = boundaryCache.get(districtId);
      if (cachedBoundary !== undefined) {
        console.log('📋 Using cached boundary data:', districtId);
        return cachedBoundary;
      }

      // Проверяем, не загружается ли уже эта граница
      if (loadingBoundaries.has(districtId)) {
        console.log('⏳ Boundary already loading:', districtId);
        return null;
      }

      console.log('📡 Fetching district boundary:', districtId);

      // Добавляем в список загружаемых
      setLoadingBoundaries((prev) => new Set(prev).add(districtId));

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
        const geoJson =
          typeof boundaryData === 'string'
            ? JSON.parse(boundaryData)
            : boundaryData;

        // Если это просто Polygon, оборачиваем в Feature
        let result: Feature<Polygon> | null = null;
        if (geoJson.type === 'Polygon') {
          result = {
            type: 'Feature',
            geometry: geoJson,
            properties: {
              id: data.data.id,
              name: data.data.name,
            },
          } as Feature<Polygon>;
        } else {
          result = geoJson as Feature<Polygon>;
        }

        // Сохраняем в кэш
        setBoundaryCache((prev) => new Map(prev).set(districtId, result));
        return result;
      } catch (error) {
        console.error('❌ Error fetching district boundary:', error);
        // Сохраняем null в кэш, чтобы не повторять неудачные запросы
        setBoundaryCache((prev) => new Map(prev).set(districtId, null));
        return null;
      } finally {
        // Убираем из списка загружаемых
        setLoadingBoundaries((prev) => {
          const newSet = new Set(prev);
          newSet.delete(districtId);
          return newSet;
        });
      }
    },
    [boundaryCache, loadingBoundaries]
  );

  const handleDistrictChange = useCallback(
    async (districtId: string | null) => {
      console.log('🌆 District changed:', districtId);
      setSelectedDistrictId(districtId);
      setSelectedMunicipalityId(null);

      if (!districtId) {
        // Очищаем границы если район не выбран
        console.log('🗺️ Clearing district boundary');
        onDistrictBoundaryChange?.(null);
        onDistrictBoundsChange?.(null);
        return;
      }

      // Загружаем границы выбранного района для отображения на карте
      if (onDistrictBoundaryChange) {
        console.log('📡 Loading district boundary for display...');
        const boundary = await loadDistrictBoundary(districtId);
        console.log('🗺️ District boundary loaded:', boundary);
        onDistrictBoundaryChange(boundary);

        // Если границы загружены, вычисляем bounding box для позиционирования карты
        if (boundary && boundary.geometry && onDistrictBoundsChange) {
          const coordinates = boundary.geometry.coordinates;
          if (coordinates && coordinates.length > 0) {
            // Для Polygon берем первый массив координат (внешний контур)
            const coords = coordinates[0];
            let minLng = 180,
              maxLng = -180,
              minLat = 90,
              maxLat = -90;

            coords.forEach((coord: number[]) => {
              minLng = Math.min(minLng, coord[0]);
              maxLng = Math.max(maxLng, coord[0]);
              minLat = Math.min(minLat, coord[1]);
              maxLat = Math.max(maxLat, coord[1]);
            });

            console.log('📍 District bounds calculated:', [
              minLng,
              minLat,
              maxLng,
              maxLat,
            ]);
            onDistrictBoundsChange([minLng, minLat, maxLng, maxLat]);
          }
        }
      }
    },
    [loadDistrictBoundary, onDistrictBoundaryChange, onDistrictBoundsChange]
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

          {(visibleCities?.length || 0) > 1 && (
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
        {(visibleCities?.length || 0) > 0 && (
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
                    count:
                      visibleCities?.filter((c) => c.city.has_districts)
                        ?.length || 0,
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
