'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import {
  InteractiveMap,
  RadiusSearchControl,
  useRadiusSearch,
  createRadiusCircle,
  RadiusSearchResult,
} from '@/components/GIS';
import { MapViewState, MapMarkerData } from '@/components/GIS/types/gis';
import { Source, Layer } from 'react-map-gl';
import type { LayerProps } from 'react-map-gl';

// Динамический импорт для избежания SSR проблем
const DynamicMap = dynamic(() => Promise.resolve(InteractiveMap), {
  ssr: false,
  loading: () => (
    <div className="w-full h-96 bg-base-200 rounded-lg flex items-center justify-center">
      <div className="loading loading-spinner loading-lg"></div>
    </div>
  ),
});

interface MarketplaceMapWithRadiusSearchProps {
  className?: string;
  initialViewState?: Partial<MapViewState>;
  onListingSelect?: (listing: RadiusSearchResult) => void;
  showControls?: boolean;
}

const DEFAULT_VIEW_STATE: MapViewState = {
  longitude: 20.4649,
  latitude: 44.8176,
  zoom: 10,
  pitch: 0,
  bearing: 0,
};

export default function MarketplaceMapWithRadiusSearch({
  className = '',
  initialViewState = {},
  onListingSelect,
  showControls = true,
}: MarketplaceMapWithRadiusSearchProps) {
  const t = useTranslations('gis');

  // Состояние карты
  const [viewState, setViewState] = useState<MapViewState>({
    ...DEFAULT_VIEW_STATE,
    ...initialViewState,
  });

  // Состояние радиусного поиска
  const [searchCenter, setSearchCenter] = useState<{
    latitude: number;
    longitude: number;
  } | null>(null);
  const [searchRadius, setSearchRadius] = useState(5); // 5 км по умолчанию
  const [showSearchCircle, setShowSearchCircle] = useState(true);
  const [markers, setMarkers] = useState<MapMarkerData[]>([]);

  // Хук радиусного поиска
  const {
    results,
    loading: searchLoading,
    error: searchError,
    total,
    search,
    searchByAddress,
    searchByCurrentLocation,
    clearResults,
  } = useRadiusSearch();

  // Создание маркеров из результатов поиска
  const createMarkersFromResults = useCallback(
    (searchResults: RadiusSearchResult[]): MapMarkerData[] => {
      return searchResults.map((result) => ({
        id: result.id,
        position: [result.longitude, result.latitude] as [number, number],
        longitude: result.longitude,
        latitude: result.latitude,
        title: result.title,
        description: result.description,
        type: 'listing' as const,
        imageUrl: result.imageUrl,
        metadata: {
          price: result.price,
          currency: result.currency || 'RSD',
          category: result.category,
          distance: result.distance,
        },
        data: {
          title: result.title,
          price: result.price,
          category: result.category,
          image: result.imageUrl,
          id: result.id,
          distance: result.distance,
          metadata: result.metadata,
        },
      }));
    },
    []
  );

  // Обновление маркеров при изменении результатов поиска
  useEffect(() => {
    const newMarkers = createMarkersFromResults(results);
    setMarkers(newMarkers);
  }, [results, createMarkersFromResults]);

  // Обработка поиска
  const handleSearch = useCallback(
    async (
      searchResults: RadiusSearchResult[],
      center: { latitude: number; longitude: number },
      radius: number
    ) => {
      setSearchCenter(center);
      setSearchRadius(radius);

      // Обновляем view state чтобы показать область поиска
      setViewState((prev) => ({
        ...prev,
        latitude: center.latitude,
        longitude: center.longitude,
        zoom: Math.max(prev.zoom, 12),
      }));
    },
    []
  );

  // Обработка изменения радиуса
  const handleRadiusChange = useCallback((radius: number) => {
    setSearchRadius(radius);
  }, []);

  // Обработка изменения центра поиска
  const handleCenterChange = useCallback(
    (center: { latitude: number; longitude: number }) => {
      setSearchCenter(center);
    },
    []
  );

  // Обработка переключения показа круга
  const handleShowCircleToggle = useCallback((show: boolean) => {
    setShowSearchCircle(show);
  }, []);

  // Обработка клика по маркеру
  const handleMarkerClick = useCallback(
    (marker: MapMarkerData) => {
      if (onListingSelect && marker.data) {
        const listing: RadiusSearchResult = {
          id: marker.id,
          title: marker.title,
          description: marker.description,
          latitude: marker.latitude,
          longitude: marker.longitude,
          distance: marker.metadata?.distance || 0,
          category: marker.metadata?.category,
          price: marker.metadata?.price,
          currency: marker.metadata?.currency,
          imageUrl: marker.imageUrl,
          metadata: marker.data.metadata,
        };
        onListingSelect(listing);
      }
    },
    [onListingSelect]
  );

  // Обработка изменения view state
  const handleViewStateChange = useCallback((newViewState: MapViewState) => {
    setViewState(newViewState);
  }, []);

  // Создание данных для круга радиуса поиска
  const radiusCircleData = useMemo(() => {
    if (!searchCenter || !showSearchCircle) return null;

    return createRadiusCircle(searchCenter, searchRadius);
  }, [searchCenter, searchRadius, showSearchCircle]);

  // Стили для круга радиуса
  const circleLayerStyle: LayerProps = {
    id: 'radius-circle-fill',
    type: 'fill',
    paint: {
      'fill-color': '#3b82f6',
      'fill-opacity': 0.1,
    },
  };

  const circleOutlineLayerStyle: LayerProps = {
    id: 'radius-circle-outline',
    type: 'line',
    paint: {
      'line-color': '#3b82f6',
      'line-width': 2,
      'line-opacity': 0.6,
    },
  };

  return (
    <div className={`relative ${className}`}>
      {/* Контроль радиусного поиска */}
      {showControls && (
        <div className="absolute top-4 left-4 z-10 w-80 max-w-[calc(100vw-2rem)]">
          <RadiusSearchControl
            config={{
              minRadius: 0.1,
              maxRadius: 50,
              defaultRadius: 5,
              step: 0.1,
              showMyLocation: true,
              showAddressInput: true,
              showRadiusCircle: true,
              enableGeolocation: true,
            }}
            onSearch={handleSearch}
            onRadiusChange={handleRadiusChange}
            onCenterChange={handleCenterChange}
            onShowCircleToggle={handleShowCircleToggle}
            disabled={searchLoading}
          />
        </div>
      )}

      {/* Информация о результатах */}
      {total > 0 && (
        <div className="absolute top-4 right-4 z-10 bg-base-100 rounded-lg shadow-lg p-3">
          <div className="flex items-center space-x-2">
            <span className="w-2 h-2 bg-success rounded-full"></span>
            <span className="text-sm font-medium">
              {t('radius_search.results_found', {
                count: total,
                radius: `${searchRadius} км`,
              })}
            </span>
          </div>
        </div>
      )}

      {/* Ошибка поиска */}
      {searchError && (
        <div className="absolute top-4 right-4 z-10 alert alert-error shadow-lg max-w-sm">
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
          <span className="text-sm">{searchError}</span>
        </div>
      )}

      {/* Карта */}
      <div className="w-full h-96 relative">
        <DynamicMap
          initialViewState={viewState}
          markers={markers}
          onMarkerClick={handleMarkerClick}
          onViewStateChange={handleViewStateChange}
          className="w-full h-full"
          mapboxAccessToken={process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN}
          controlsConfig={{
            showNavigation: true,
            showFullscreen: true,
            showGeolocate: true,
            position: 'bottom-right',
          }}
        />

        {/* Пока что без круга - можно добавить позже через кастомизацию InteractiveMap */}
      </div>

      {/* Индикатор загрузки */}
      {searchLoading && (
        <div className="absolute inset-0 bg-base-100/50 flex items-center justify-center z-20">
          <div className="bg-base-100 rounded-lg shadow-lg p-4 flex items-center space-x-3">
            <div className="loading loading-spinner loading-md"></div>
            <span className="text-base-content">
              {t('radius_search.searching')}
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
