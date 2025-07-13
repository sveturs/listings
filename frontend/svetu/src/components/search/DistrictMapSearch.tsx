'use client';

import React, { useState, useCallback, useRef } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import { DistrictMapSelector } from './DistrictMapSelector';
import type { components } from '@/types/generated/api';
import { ListingPopup } from '../GIS/Map/MapPopup';
import type {
  MapMarkerData,
  MapPopupData,
  MapViewState,
} from '@/components/GIS/types/gis';

// Динамическая загрузка карты
const InteractiveMap = dynamic(
  () => import('@/components/GIS/Map/InteractiveMap'),
  { ssr: false }
);

type SpatialSearchResult = components['schemas']['types.SpatialSearchResult'];

export default function DistrictMapSearch() {
  const t = useTranslations();
  const mapRef = useRef<any>(null);

  const [viewState, setViewState] = useState<MapViewState>({
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 11,
  });

  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [popup, setPopup] = useState<MapPopupData | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Обработка результатов поиска
  const handleSearchResults = useCallback((results: SpatialSearchResult[]) => {
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

    setMarkers(newMarkers);
  }, []);

  // Обработка изменения границ района
  const handleDistrictBoundsChange = useCallback(
    (bounds: [number, number, number, number] | null) => {
      if (!bounds || !mapRef.current) return;

      const [minLng, minLat, maxLng, maxLat] = bounds;

      // Рассчитываем центр и масштаб для отображения всего района
      const centerLng = (minLng + maxLng) / 2;
      const centerLat = (minLat + maxLat) / 2;

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

  // Обработка клика по маркеру
  const handleMarkerClick = useCallback((marker: MapMarkerData) => {
    setPopup({
      longitude: marker.longitude,
      latitude: marker.latitude,
      content: (
        <ListingPopup
          title={marker.title}
          description={marker.description}
          price={marker.data?.price}
          currency={marker.data?.currency}
          imageUrl={marker.data?.imageUrl}
          link={`/listings/${marker.id}`}
          categoryName={marker.data?.categoryName}
          address={marker.data?.address}
        />
      ),
    });
  }, []);

  return (
    <div className="relative h-screen w-full">
      {/* Карта на весь экран */}
      <InteractiveMap
        ref={mapRef}
        viewState={viewState}
        onViewStateChange={setViewState}
        markers={markers}
        onMarkerClick={handleMarkerClick}
        popup={popup}
        onPopupClose={() => setPopup(null)}
        style={{ width: '100%', height: '100%' }}
      />

      {/* Панель выбора района */}
      <div className="absolute top-4 left-4 z-10 w-80 max-h-[calc(100vh-2rem)] overflow-y-auto">
        <DistrictMapSelector
          onSearchResults={handleSearchResults}
          onDistrictBoundsChange={handleDistrictBoundsChange}
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
      {isLoading && (
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      )}
    </div>
  );
}
