'use client';

import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import { DistrictMapSelector } from './DistrictMapSelector';
import type { components as _components } from '@/types/generated/api';
import { ListingPopup } from '../GIS/Map/MapPopup';
import type {
  MapMarkerData,
  MapPopupData,
  MapViewState,
} from '@/components/GIS/types/gis';
import type { Feature, Polygon } from 'geojson';
import {
  SearchModeProvider,
  useSearchMode,
} from '@/contexts/SearchModeContext';

// Динамическая загрузка карты
const InteractiveMap = dynamic(
  () => import('@/components/GIS/Map/InteractiveMap'),
  { ssr: false }
);

// Временный интерфейс до исправления API типов
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

function DistrictMapSearchInner() {
  const t = useTranslations();
  const mapRef = useRef<any>(null);
  const { setSearchMode } = useSearchMode();

  // Устанавливаем флаг блокировки радиусного поиска при монтировании компонента
  useEffect(() => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('blockRadiusSearch', 'true');
      // Устанавливаем глобальный флаг
      (window as any).__BLOCK_RADIUS_SEARCH__ = true;
      (window as any).__DISTRICT_PAGE_ACTIVE__ = true;
      console.log('🚫 Radius search blocked for district page (frontend)');

      // ВРЕМЕННО ОТКЛЮЧЕНО: Перехватчики fetch и XHR
      // Блокировка теперь работает на уровне backend
      console.log('ℹ️ Frontend interceptors disabled - using backend blocking');
    }

    return () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('blockRadiusSearch');
        delete (window as any).__BLOCK_RADIUS_SEARCH__;
        delete (window as any).__DISTRICT_PAGE_ACTIVE__;
        delete (window as any).__DISTRICT_MARKERS_SET__;
        console.log('✅ Radius search unblocked');
      }
    };
  }, []);

  const [viewState, setViewState] = useState<MapViewState>({
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 11,
  });

  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [popup, setPopup] = useState<MapPopupData | null>(null);
  const [_isLoading, _setIsLoading] = useState(false);
  const [districtBoundary, setDistrictBoundary] =
    useState<Feature<Polygon> | null>(null);
  const [_isDistrictSelected, setIsDistrictSelected] = useState(false);

  // Обработка результатов поиска
  const handleSearchResults = useCallback((results: SpatialSearchResult[]) => {
    console.log(
      '🔍 District search results received:',
      results.length,
      'items'
    );

    const newMarkers: MapMarkerData[] = results.map((result) => ({
      id: result.id,
      position: [result.longitude, result.latitude],
      longitude: result.longitude,
      latitude: result.latitude,
      title: result.title,
      description: result.description || '',
      type: 'listing' as const,
      data: {
        price: result.price,
        currency: result.currency,
        imageUrl: result.first_image_url || '/api/placeholder/200/150',
        categoryName: result.category_name,
        address: result.address,
        userEmail: result.user_email,
      },
    }));

    console.log('🗺️ Setting district markers:', newMarkers.length);
    setMarkers(newMarkers);
    setIsDistrictSelected(results.length > 0); // Устанавливаем флаг, что район выбран

    // Защита от очистки маркеров - устанавливаем флаг
    if (typeof window !== 'undefined') {
      (window as any).__DISTRICT_MARKERS_SET__ = true;
      setTimeout(() => {
        delete (window as any).__DISTRICT_MARKERS_SET__;
      }, 2000); // Защита на 2 секунды
    }
  }, []);

  // Обработка изменения границ района (для изменения viewport)
  const handleDistrictBoundsChange = useCallback(
    (bounds: [number, number, number, number] | null) => {
      if (!bounds || !mapRef.current) return;

      const [minLng, minLat, maxLng, maxLat] = bounds;

      // Рассчитываем центр и масштаб для отображения всего района
      const _centerLng = (minLng + maxLng) / 2;
      const _centerLat = (minLat + maxLat) / 2;

      // Добавляем небольшой отступ
      const padding = 0.01;
      const paddedBounds: [[number, number], [number, number]] = [
        [minLng - padding, minLat - padding],
        [maxLng + padding, maxLat + padding],
      ];

      // Используем fitBounds для плавного перехода к району
      mapRef.current.fitBounds(paddedBounds, {
        padding: 40,
        duration: 1000,
      });
    },
    []
  );

  // Обработка изменения границ района (для отображения на карте)
  const handleDistrictBoundaryChange = useCallback(
    (boundary: Feature<Polygon> | null) => {
      console.log('🗺️ District boundary changed:', boundary);
      setDistrictBoundary(boundary);
      setIsDistrictSelected(boundary !== null); // Устанавливаем флаг на основе наличия границ

      // Устанавливаем режим поиска
      setSearchMode(boundary !== null ? 'district' : 'none');
    },
    [setSearchMode]
  );

  // Обработка клика по маркеру
  const handleMarkerClick = useCallback((marker: MapMarkerData) => {
    setPopup({
      id: marker.id,
      position: [marker.longitude, marker.latitude],
      title: marker.title,
      description: marker.description,
      content: (
        <ListingPopup
          listing={{
            id: marker.id,
            title: marker.title,
            price: marker.data?.price || 0,
            currency: marker.data?.currency || 'RSD',
            imageUrl: marker.data?.imageUrl,
            category: marker.data?.categoryName,
          }}
          position={[marker.longitude, marker.latitude]}
          onClose={() => setPopup(null)}
        />
      ),
    });
  }, []);

  return (
    <div className="relative h-screen w-full">
      {/* Карта на весь экран */}
      <InteractiveMap
        initialViewState={viewState}
        onViewStateChange={setViewState}
        markers={markers}
        onMarkerClick={handleMarkerClick}
        popup={popup}
        style={{ width: '100%', height: '100%' }}
        districtBoundary={districtBoundary}
      />

      {/* Панель выбора района */}
      <div className="absolute top-4 left-4 z-10 w-80 max-h-[calc(100vh-2rem)] overflow-y-auto">
        <DistrictMapSelector
          onSearchResults={handleSearchResults}
          onDistrictBoundsChange={handleDistrictBoundsChange}
          onDistrictBoundaryChange={handleDistrictBoundaryChange}
          className="shadow-2xl"
        />
      </div>

      {/* Счетчик результатов */}
      {markers.length > 0 && (
        <div className="absolute bottom-4 left-4 z-10">
          <div className="badge badge-lg badge-primary">
            {t('search.found')}: {markers.length}
          </div>
        </div>
      )}

      {/* Индикатор загрузки */}
      {_isLoading && (
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      )}
    </div>
  );
}

export default function DistrictMapSearch() {
  // DISTRICT FUNCTIONALITY TEMPORARILY DISABLED
  return null;
  /*
  return (
    <SearchModeProvider>
      <DistrictMapSearchInner />
    </SearchModeProvider>
  );
  */
}
