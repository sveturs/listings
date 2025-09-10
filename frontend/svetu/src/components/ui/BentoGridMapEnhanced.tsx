'use client';

import React, { useMemo, useState, useCallback, useEffect } from 'react';
import Map, { Marker, Source, Layer } from 'react-map-gl';
import type {
  ViewState,
  LayerProps,
  MapRef,
  MarkerDragEvent,
} from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { getMapboxIsochrone } from '@/components/GIS/utils/mapboxIsochrone';
import { isPointInIsochrone } from '@/components/GIS/utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';
import { useTranslations } from 'next-intl';

// Компонент перетаскиваемой иконки местоположения
const DraggableLocationIcon: React.FC<{
  mapRef: React.RefObject<MapRef>;
  onDropLocation: (lng: number, lat: number) => void;
}> = ({ mapRef, onDropLocation }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [dragPosition, setDragPosition] = useState({ x: 0, y: 0 });

  const handleDragStart = (e: React.DragEvent<HTMLDivElement>) => {
    e.dataTransfer.effectAllowed = 'move';
    setIsDragging(true);
    const dragImage = new Image();
    dragImage.src =
      'data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=';
    e.dataTransfer.setDragImage(dragImage, 0, 0);
  };

  const handleDragEnd = () => {
    setIsDragging(false);
  };

  useEffect(() => {
    if (!mapRef.current) return;

    const mapContainer = mapRef.current.getContainer();

    const handleDragOver = (e: DragEvent) => {
      if (isDragging) {
        e.preventDefault();
        e.dataTransfer!.dropEffect = 'move';
        setDragPosition({ x: e.clientX, y: e.clientY });
      }
    };

    const handleDrop = (e: DragEvent) => {
      if (isDragging) {
        e.preventDefault();
        const map = mapRef.current;
        if (!map) return;

        const rect = map.getContainer().getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        const lngLat = map.unproject([x, y]);
        onDropLocation(lngLat.lng, lngLat.lat);
      }
    };

    mapContainer.addEventListener('dragover', handleDragOver);
    mapContainer.addEventListener('drop', handleDrop);

    return () => {
      mapContainer.removeEventListener('dragover', handleDragOver);
      mapContainer.removeEventListener('drop', handleDrop);
    };
  }, [isDragging, mapRef, onDropLocation]);

  return (
    <>
      <div
        className={`absolute bottom-20 right-4 z-20 cursor-move transition-all duration-200 group ${
          isDragging ? 'opacity-50 scale-95' : 'opacity-100 hover:scale-105'
        }`}
        draggable
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className="bg-white rounded-lg shadow-lg p-3 hover:shadow-xl transition-shadow border-2 border-transparent hover:border-red-100">
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="text-red-500"
          >
            <circle cx="12" cy="12" r="3" fill="currentColor" />
            <path
              d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
              fill="currentColor"
            />
          </svg>
        </div>
        <div className="absolute top-full mt-2 right-0 bg-gray-800 text-white text-xs rounded-lg px-3 py-2 whitespace-nowrap pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity">
          Перетащите на карту
        </div>
      </div>
      {isDragging && (
        <div
          className="fixed pointer-events-none z-50"
          style={{ left: dragPosition.x - 12, top: dragPosition.y - 24 }}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="text-red-500"
          >
            <circle cx="12" cy="12" r="3" fill="currentColor" />
            <path
              d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
              fill="currentColor"
            />
          </svg>
        </div>
      )}
    </>
  );
};

interface BentoGridMapEnhancedProps {
  listings?: Array<{
    id: string;
    latitude: number;
    longitude: number;
    price: number;
    isStorefront?: boolean; // Для различения витрин и обычных объявлений
    storeName?: string;
    imageUrl?: string;
    category?: string;
  }>;
  userLocation?: {
    latitude: number;
    longitude: number;
  };
  searchRadius?: number; // В метрах
  showRadius?: boolean;
  enableClustering?: boolean;
}

export const BentoGridMapEnhanced: React.FC<BentoGridMapEnhancedProps> = ({
  listings = [],
  userLocation,
  searchRadius: initialSearchRadius = 5000, // 5km по умолчанию
  showRadius: initialShowRadius = true,
  enableClustering = true,
}) => {
  const t = useTranslations('map');
  const mapRef = React.useRef<MapRef>(null);
  const [theme, setTheme] = React.useState<'light' | 'dark'>('light');
  const [mounted, setMounted] = React.useState(false);
  const [isGeolocationAvailable, setIsGeolocationAvailable] =
    React.useState(false);
  const [showRadius, setShowRadius] = React.useState(initialShowRadius);
  const [searchRadius, setSearchRadius] = React.useState(initialSearchRadius);
  const [, setSelectedListing] = React.useState<string | null>(null);
  const [walkingMode, setWalkingMode] = React.useState<'radius' | 'walking'>(
    'radius'
  );
  const [walkingTime, setWalkingTime] = React.useState(15);
  const [userMarkerLocation, setUserMarkerLocation] = React.useState(
    userLocation || { latitude: 44.7866, longitude: 20.4489 }
  );
  const [, setIsDragging] = React.useState(false);
  const [isCompactControlExpanded, setIsCompactControlExpanded] =
    React.useState(false);
  const [showMobileHint, setShowMobileHint] = React.useState(false);
  const firstInteractionRef = React.useRef(true);
  const [isochroneData, setIsochroneData] =
    React.useState<Feature<Polygon> | null>(null);
  const [isLoadingIsochrone, setIsLoadingIsochrone] = React.useState(false);
  const [isLongPressing, setIsLongPressing] = React.useState(false);

  React.useEffect(() => {
    if ('geolocation' in navigator) {
      setIsGeolocationAvailable(true);
    }
  }, []);

  // Отслеживание темы
  React.useEffect(() => {
    setMounted(true);

    // Функция для получения текущей темы
    const getTheme = () => {
      const htmlTheme = document.documentElement.getAttribute('data-theme');
      return htmlTheme === 'dark' ? 'dark' : 'light';
    };

    // Устанавливаем начальную тему
    setTheme(getTheme());

    // Создаем observer для отслеживания изменений темы
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (
          mutation.type === 'attributes' &&
          mutation.attributeName === 'data-theme'
        ) {
          setTheme(getTheme());
        }
      });
    });

    // Наблюдаем за изменениями атрибута data-theme
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme'],
    });

    return () => observer.disconnect();
  }, []);

  // Загрузка изохрона при изменении параметров
  React.useEffect(() => {
    if (!userMarkerLocation || !showRadius || walkingMode !== 'walking') {
      setIsochroneData(null);
      return;
    }

    const loadIsochrone = async () => {
      setIsLoadingIsochrone(true);
      try {
        const isochrone = await getMapboxIsochrone({
          coordinates: [
            userMarkerLocation.longitude,
            userMarkerLocation.latitude,
          ],
          minutes: walkingTime,
          profile: 'walking',
        });
        setIsochroneData(isochrone);
      } catch (error) {
        console.error('Failed to load isochrone:', error);
        // В случае ошибки изохрон будет сгенерирован локально (fallback в getMapboxIsochrone)
      } finally {
        setIsLoadingIsochrone(false);
      }
    };

    // Debounce для избежания слишком частых запросов
    const timer = setTimeout(loadIsochrone, 500);
    return () => clearTimeout(timer);
  }, [userMarkerLocation, showRadius, walkingMode, walkingTime]);

  // Функция для расчета расстояния между двумя точками (формула Хаверсина)
  const calculateDistance = useCallback(
    (lat1: number, lon1: number, lat2: number, lon2: number) => {
      const R = 6371e3; // Радиус Земли в метрах
      const φ1 = (lat1 * Math.PI) / 180;
      const φ2 = (lat2 * Math.PI) / 180;
      const Δφ = ((lat2 - lat1) * Math.PI) / 180;
      const Δλ = ((lon2 - lon1) * Math.PI) / 180;

      const a =
        Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
        Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

      return R * c; // Расстояние в метрах
    },
    []
  );

  // Фильтруем объявления по радиусу или изохрону
  const filteredListings = useMemo(() => {
    if (!showRadius || !userMarkerLocation) return listings;

    // Если режим пешей доступности и есть изохрон
    if (walkingMode === 'walking' && isochroneData) {
      return listings.filter((listing) => {
        return isPointInIsochrone(
          [listing.longitude, listing.latitude],
          isochroneData
        );
      });
    }

    // Иначе используем радиус
    const effectiveRadius =
      walkingMode === 'walking' ? walkingTime * 80 : searchRadius;

    return listings.filter((listing) => {
      const distance = calculateDistance(
        userMarkerLocation.latitude,
        userMarkerLocation.longitude,
        listing.latitude,
        listing.longitude
      );
      return distance <= effectiveRadius;
    });
  }, [
    listings,
    userMarkerLocation,
    showRadius,
    searchRadius,
    walkingMode,
    walkingTime,
    calculateDistance,
    isochroneData,
  ]);

  // Преобразуем данные в GeoJSON для кластеризации
  const geoJsonData = useMemo(() => {
    return {
      type: 'FeatureCollection' as const,
      features: filteredListings.map((listing) => ({
        type: 'Feature' as const,
        geometry: {
          type: 'Point' as const,
          coordinates: [listing.longitude, listing.latitude],
        },
        properties: {
          id: listing.id,
          price:
            typeof listing.price === 'string'
              ? parseFloat(listing.price)
              : listing.price || 0,
          isStorefront: listing.isStorefront || false,
          storeName: listing.storeName,
          imageUrl: listing.imageUrl,
          category: listing.category,
        },
      })),
    };
  }, [filteredListings]);

  // Определяем центр карты и масштаб
  const { center, zoom } = useMemo(() => {
    if (userLocation) {
      return {
        center: {
          longitude: userLocation.longitude,
          latitude: userLocation.latitude,
        },
        zoom: 13,
      };
    }

    if (listings.length > 0) {
      const avgLat =
        listings.reduce((sum, l) => sum + l.latitude, 0) / listings.length;
      const avgLng =
        listings.reduce((sum, l) => sum + l.longitude, 0) / listings.length;
      return {
        center: { longitude: avgLng, latitude: avgLat },
        zoom: 12,
      };
    }

    return {
      center: { longitude: 20.4489, latitude: 44.7866 }, // Белград
      zoom: 11,
    };
  }, [listings, userLocation]);

  const initialViewState: Partial<ViewState> = {
    ...center,
    zoom,
    pitch: 0,
    bearing: 0,
  };

  // Слой кластеров
  const clusterLayer: LayerProps = useMemo(
    () => ({
      id: 'clusters',
      type: 'circle',
      source: 'listings',
      filter: ['has', 'point_count'],
      paint: {
        'circle-color': [
          'step',
          ['get', 'point_count'],
          '#60a5fa', // blue-400
          10,
          '#3b82f6', // blue-500
          30,
          '#2563eb', // blue-600
        ],
        'circle-radius': [
          'step',
          ['get', 'point_count'],
          15, // маленький кластер
          10,
          20, // средний кластер
          30,
          25, // большой кластер
        ],
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.8,
      },
    }),
    []
  );

  // Слой текста кластеров
  const clusterCountLayer: LayerProps = useMemo(
    () => ({
      id: 'cluster-count',
      type: 'symbol',
      source: 'listings',
      filter: ['has', 'point_count'],
      layout: {
        'text-field': '{point_count_abbreviated}',
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 12,
      },
      paint: {
        'text-color': '#ffffff',
      },
    }),
    []
  );

  // Слой диапазона цен под кластерами
  const clusterPriceLayer: LayerProps = useMemo(
    () => ({
      id: 'cluster-price',
      type: 'symbol',
      source: 'listings',
      filter: ['has', 'point_count'],
      layout: {
        'text-field': [
          'concat',
          [
            'number-format',
            ['get', 'minPrice'],
            { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
          ],
          '-',
          [
            'number-format',
            ['get', 'maxPrice'],
            { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
          ],
          ' RSD',
        ],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 10,
        'text-anchor': 'top',
        'text-offset': [0, 1.5],
      },
      paint: {
        'text-color': theme === 'dark' ? '#e5e7eb' : '#374151',
        'text-halo-color': theme === 'dark' ? '#1f2937' : '#ffffff',
        'text-halo-width': 2,
      },
    }),
    [theme]
  );

  // Слой индивидуальных маркеров
  const unclusteredPointLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-point',
      type: 'circle',
      source: 'listings',
      filter: ['!', ['has', 'point_count']],
      paint: {
        'circle-color': [
          'case',
          ['get', 'isStorefront'],
          '#f59e0b', // amber-500 для витрин
          '#3b82f6', // blue-500 для обычных
        ],
        'circle-radius': 8,
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.9,
      },
    }),
    []
  );

  // Слой цен под индивидуальными маркерами
  const unclusteredPriceLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-price',
      type: 'symbol',
      source: 'listings',
      filter: ['!', ['has', 'point_count']],
      layout: {
        'text-field': [
          'concat',
          [
            'number-format',
            ['get', 'price'],
            { 'min-fraction-digits': 0, 'max-fraction-digits': 0 },
          ],
          ' RSD',
        ],
        'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
        'text-size': 10,
        'text-anchor': 'top',
        'text-offset': [0, 1],
      },
      paint: {
        'text-color': theme === 'dark' ? '#e5e7eb' : '#374151',
        'text-halo-color': theme === 'dark' ? '#1f2937' : '#ffffff',
        'text-halo-width': 2,
      },
    }),
    [theme]
  );

  // Данные для радиуса поиска или изохрона
  const radiusGeoJson = useMemo(() => {
    if (!userMarkerLocation || !showRadius) return null;

    // Если есть изохрон для режима пешей доступности, используем его
    if (walkingMode === 'walking' && isochroneData) {
      return isochroneData;
    }

    // Иначе генерируем круг для радиуса
    const numSides = 64;
    const angleStep = (2 * Math.PI) / numSides;
    const coordinates: [number, number][] = [];

    const effectiveRadius = searchRadius;

    // Константы для конвертации
    const lat = userMarkerLocation.latitude;
    const latRad = (lat * Math.PI) / 180;

    // Метры в градус широты (примерно одинаково везде)
    const metersPerDegreeLat = 111320;

    // Метры в градус долготы (зависит от широты)
    const metersPerDegreeLng = 111320 * Math.cos(latRad);

    for (let i = 0; i <= numSides; i++) {
      const angle = i * angleStep;
      const dx = (effectiveRadius * Math.cos(angle)) / metersPerDegreeLng;
      const dy = (effectiveRadius * Math.sin(angle)) / metersPerDegreeLat;

      coordinates.push([
        userMarkerLocation.longitude + dx,
        userMarkerLocation.latitude + dy,
      ]);
    }

    return {
      type: 'Feature' as const,
      geometry: {
        type: 'Polygon' as const,
        coordinates: [coordinates],
      },
      properties: {},
    };
  }, [
    userMarkerLocation,
    searchRadius,
    showRadius,
    walkingMode,
    isochroneData,
  ]);

  // Слой радиуса
  const radiusLayer: LayerProps = {
    id: 'search-radius',
    type: 'fill',
    paint: {
      'fill-color': walkingMode === 'walking' ? '#10B981' : '#3b82f6',
      'fill-opacity': walkingMode === 'walking' ? 0.2 : 0.1,
    },
  };

  const radiusBorderLayer: LayerProps = {
    id: 'search-radius-border',
    type: 'line',
    paint: {
      'line-color': walkingMode === 'walking' ? '#10B981' : '#3b82f6',
      'line-width': walkingMode === 'walking' ? 3 : 2,
      'line-dasharray': walkingMode === 'walking' ? [4, 4] : [2, 2],
      'line-opacity': 0.8,
    },
  };

  const handleGeolocation = useCallback(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          setUserMarkerLocation({ latitude, longitude });
          mapRef.current?.flyTo({
            center: [longitude, latitude],
            zoom: 14,
            duration: 1000,
          });
        },
        (error) => {
          console.error('Geolocation error:', error);
        }
      );
    }
  }, []);

  // Обработчики для маркера пользователя
  const handleUserMarkerDrag = useCallback((event: MarkerDragEvent) => {
    setUserMarkerLocation({
      longitude: event.lngLat.lng,
      latitude: event.lngLat.lat,
    });
  }, []);

  const handleUserMarkerDragStart = useCallback(() => {
    setIsDragging(true);
  }, []);

  const handleUserMarkerDragEnd = useCallback((event: MarkerDragEvent) => {
    setUserMarkerLocation({
      longitude: event.lngLat.lng,
      latitude: event.lngLat.lat,
    });
    setIsDragging(false);
  }, []);

  // Обновляем расположение маркера при изменении userLocation
  useEffect(() => {
    if (userLocation && !userMarkerLocation) {
      setUserMarkerLocation(userLocation);
    }
  }, [userLocation, userMarkerLocation]);

  const handleMapClick = useCallback((event: any) => {
    const features = event.features;
    if (!features || features.length === 0) return;

    const feature = features[0];

    // Если это кластер, увеличиваем масштаб
    if (feature.properties.cluster) {
      const clusterId = feature.properties.cluster_id;
      const mapboxMap = mapRef.current?.getMap();
      const source = mapboxMap?.getSource('listings') as any;

      source?.getClusterExpansionZoom(clusterId, (err: any, zoom: number) => {
        if (err) return;

        mapRef.current?.easeTo({
          center: feature.geometry.coordinates,
          zoom,
          duration: 500,
        });
      });
    } else {
      // Если это маркер, показываем информацию о нём
      setSelectedListing(feature.properties.id);
      // Можно добавить popup или перенаправление
      window.open(`/marketplace/${feature.properties.id}`, '_blank');
    }
  }, []);

  // Предотвращаем гидратационные ошибки
  if (!mounted) {
    return (
      <div className="w-full h-full relative overflow-hidden rounded-lg bg-base-200 animate-pulse">
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full h-full relative overflow-hidden rounded-lg">
      <Map
        ref={mapRef}
        initialViewState={initialViewState}
        style={{ width: '100%', height: '100%' }}
        mapStyle={
          theme === 'dark'
            ? 'mapbox://styles/mapbox/dark-v11'
            : 'mapbox://styles/mapbox/light-v11'
        }
        mapboxAccessToken={
          process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN ||
          process.env.NEXT_PUBLIC_MAPBOX_TOKEN
        }
        interactive={true}
        attributionControl={false}
        onClick={handleMapClick}
        interactiveLayerIds={
          enableClustering ? ['clusters', 'unclustered-point'] : undefined
        }
      >
        {/* Радиус поиска */}
        {radiusGeoJson && (
          <Source id="radius" type="geojson" data={radiusGeoJson}>
            <Layer {...radiusLayer} />
            <Layer {...radiusBorderLayer} />
          </Source>
        )}

        {/* Маркер пользователя */}
        {userMarkerLocation && (
          <Marker
            longitude={userMarkerLocation.longitude}
            latitude={userMarkerLocation.latitude}
            draggable
            onDrag={handleUserMarkerDrag}
            onDragStart={handleUserMarkerDragStart}
            onDragEnd={handleUserMarkerDragEnd}
            anchor="bottom"
          >
            <div
              className="cursor-move hover:scale-110 transition-transform"
              style={{ width: 40, height: 40 }}
            >
              <svg
                width="40"
                height="40"
                viewBox="0 0 40 40"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                {/* Тень */}
                <ellipse
                  cx="20"
                  cy="38"
                  rx="8"
                  ry="2"
                  fill="black"
                  fillOpacity="0.2"
                />
                {/* Основной маркер */}
                <path
                  d="M20 36C20 36 32 24 32 16C32 9.37258 26.6274 4 20 4C13.3726 4 8 9.37258 8 16C8 24 20 36 20 36Z"
                  fill="#EF4444"
                  stroke="white"
                  strokeWidth="2"
                />
                {/* Иконка человека */}
                <circle cx="20" cy="13" r="3" fill="white" />
                <path
                  d="M15 20C15 18.3431 16.3431 17 18 17H22C23.6569 17 25 18.3431 25 20V24C25 24.5523 24.5523 25 24 25H16C15.4477 25 15 24.5523 15 24V20Z"
                  fill="white"
                />
              </svg>
            </div>
          </Marker>
        )}

        {/* Кластеризованные маркеры */}
        {enableClustering ? (
          <Source
            id="listings"
            type="geojson"
            data={geoJsonData}
            cluster={true}
            clusterMaxZoom={14}
            clusterRadius={50}
            clusterProperties={{
              minPrice: ['min', ['get', 'price']],
              maxPrice: ['max', ['get', 'price']],
            }}
          >
            <Layer {...clusterLayer} />
            <Layer {...clusterCountLayer} />
            <Layer {...clusterPriceLayer} />
            <Layer {...unclusteredPointLayer} />
            <Layer {...unclusteredPriceLayer} />
          </Source>
        ) : (
          // Некластеризованные маркеры
          filteredListings.map((listing) => (
            <Marker
              key={listing.id}
              longitude={listing.longitude}
              latitude={listing.latitude}
              anchor="bottom"
              onClick={(e) => {
                e.originalEvent.stopPropagation();
                window.open(`/marketplace/${listing.id}`, '_blank');
              }}
            >
              <div className="relative group cursor-pointer">
                <div className="absolute inset-0 bg-secondary/20 rounded-full scale-150 opacity-0 group-hover:opacity-100 transition-opacity" />
                <div
                  className={`rounded-full px-2 py-1 text-xs font-semibold shadow-md border-2 border-white group-hover:scale-110 transition-transform ${
                    listing.isStorefront
                      ? 'bg-amber-500 text-white'
                      : 'bg-white text-secondary'
                  }`}
                >
                  {listing.price.toLocaleString()} RSD
                </div>
              </div>
            </Marker>
          ))
        )}

        {/* Компактный контрол с линейкой - ВНУТРИ Map */}
        {userMarkerLocation && (
          <div
            className="absolute top-2 right-2 z-10"
            style={{
              width: isCompactControlExpanded ? '260px' : '32px',
              height: isCompactControlExpanded ? 'auto' : '32px',
            }}
          >
            {!isCompactControlExpanded ? (
              // Компактное состояние - иконка-кнопка
              <button
                className="w-full h-full flex items-center justify-center bg-white rounded-lg shadow-lg hover:bg-gray-50 transition-colors relative"
                onClick={(e) => {
                  e.preventDefault();
                  // Одиночный клик - переключение режима только если не было long press
                  if (!isLongPressing) {
                    setWalkingMode(
                      walkingMode === 'walking' ? 'radius' : 'walking'
                    );
                  }
                  setIsLongPressing(false);
                }}
                onMouseDown={(e) => {
                  // Для десктопа - начинаем отсчет long press
                  setIsLongPressing(false);
                  const timer = setTimeout(() => {
                    setIsLongPressing(true);
                    setIsCompactControlExpanded(true);
                    // Вибрация для обратной связи (если поддерживается)
                    if ('vibrate' in navigator) {
                      navigator.vibrate(50);
                    }
                  }, 500); // 500ms для long press
                  (e.currentTarget as any).longPressTimer = timer;
                }}
                onMouseUp={(e) => {
                  // Отменяем таймер, если отпустили раньше
                  const timer = (e.currentTarget as any).longPressTimer;
                  if (timer) {
                    clearTimeout(timer);
                    delete (e.currentTarget as any).longPressTimer;
                  }
                }}
                onMouseLeave={(e) => {
                  // Отменяем таймер, если мышь ушла с кнопки
                  const timer = (e.currentTarget as any).longPressTimer;
                  if (timer) {
                    clearTimeout(timer);
                    delete (e.currentTarget as any).longPressTimer;
                  }
                  setIsLongPressing(false);
                }}
                onTouchStart={(e) => {
                  e.preventDefault();

                  // Показываем подсказку при первом касании
                  if (firstInteractionRef.current && 'ontouchstart' in window) {
                    setShowMobileHint(true);
                    firstInteractionRef.current = false;
                    setTimeout(() => setShowMobileHint(false), 3000);
                  }

                  const timer = setTimeout(() => {
                    setIsCompactControlExpanded(true);
                    setShowMobileHint(false);
                    // Вибрация для обратной связи
                    if ('vibrate' in navigator) {
                      navigator.vibrate(50);
                    }
                  }, 500); // 500ms для long press
                  // Сохраняем таймер в data-атрибуте
                  (e.currentTarget as any).longPressTimer = timer;
                }}
                onTouchEnd={(e) => {
                  e.preventDefault();
                  const timer = (e.currentTarget as any).longPressTimer;
                  if (timer) {
                    clearTimeout(timer);
                    delete (e.currentTarget as any).longPressTimer;
                  }
                  // Если контрол не раскрылся, значит это был короткий тап - меняем режим
                  if (!isCompactControlExpanded) {
                    setWalkingMode(
                      walkingMode === 'walking' ? 'radius' : 'walking'
                    );
                  }
                }}
                onTouchMove={(e) => {
                  const timer = (e.currentTarget as any).longPressTimer;
                  if (timer) {
                    clearTimeout(timer);
                    delete (e.currentTarget as any).longPressTimer;
                  }
                }}
                title="Клик - сменить режим, долгое нажатие - настройки"
              >
                <span className="text-lg">
                  {walkingMode === 'walking' ? '🚶' : '📏'}
                </span>
                {/* Индикатор загрузки в компактном режиме */}
                {isLoadingIsochrone && walkingMode === 'walking' && (
                  <div className="absolute top-0 right-0 -mt-1 -mr-1">
                    <span className="loading loading-spinner loading-xs text-green-600"></span>
                  </div>
                )}
                {/* Индикатор значения */}
                <div
                  className="absolute text-white font-bold text-[9px] rounded px-1"
                  style={{
                    backgroundColor:
                      walkingMode === 'walking' ? '#10B981' : '#3b82f6',
                    bottom: '-2px',
                    right: '-2px',
                    lineHeight: '12px',
                  }}
                >
                  {walkingMode === 'walking'
                    ? `${walkingTime}'`
                    : searchRadius >= 1000
                      ? `${(searchRadius / 1000).toFixed(0)}км`
                      : `${searchRadius}м`}
                </div>

                {/* Мобильная подсказка */}
                {showMobileHint && (
                  <div className="absolute -bottom-10 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs px-3 py-2 rounded-lg shadow-lg whitespace-nowrap z-20">
                    {t('holdForSettings')}
                    <div className="absolute -top-1 left-1/2 transform -translate-x-1/2 w-2 h-2 bg-gray-800 transform rotate-45"></div>
                  </div>
                )}
              </button>
            ) : (
              // Развернутое состояние с настройками
              <div className="bg-base-100 dark:bg-base-200 rounded-lg shadow-lg p-3 space-y-3">
                {/* Заголовок с кнопкой закрытия */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <span className="text-lg">
                      {walkingMode === 'walking' ? '🚶' : '📏'}
                    </span>
                    <span className="text-sm font-medium text-base-content">
                      {walkingMode === 'walking'
                        ? t('walkingAccessibility')
                        : t('searchRadius')}
                    </span>
                  </div>
                  <button
                    onClick={() => setIsCompactControlExpanded(false)}
                    className="text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
                  >
                    ✕
                  </button>
                </div>

                {/* Переключатель видимости радиуса */}
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm checkbox-primary"
                    checked={showRadius}
                    onChange={(e) => setShowRadius(e.target.checked)}
                  />
                  <span className="text-xs text-base-content">
                    {t('showZone')}
                  </span>
                </label>

                {/* Слайдер для настройки */}
                {showRadius && (
                  <div className="space-y-2">
                    {walkingMode === 'walking' ? (
                      <>
                        <input
                          type="range"
                          min="5"
                          max="30"
                          value={walkingTime}
                          onChange={(e) =>
                            setWalkingTime(Number(e.target.value))
                          }
                          className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                          style={{
                            background: `linear-gradient(to right, #10B981 0%, #10B981 ${((walkingTime - 5) / 25) * 100}%, #e5e7eb ${((walkingTime - 5) / 25) * 100}%, #e5e7eb 100%)`,
                          }}
                        />
                        <div className="text-xs text-right text-green-600 font-medium">
                          {t('minutesWalking', { minutes: walkingTime })}
                        </div>
                      </>
                    ) : (
                      <>
                        <input
                          type="range"
                          min="100"
                          max="150000"
                          step="100"
                          value={searchRadius}
                          onChange={(e) =>
                            setSearchRadius(Number(e.target.value))
                          }
                          className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                          style={{
                            background: `linear-gradient(to right, #3b82f6 0%, #3b82f6 ${((searchRadius - 100) / (150000 - 100)) * 100}%, #e5e7eb ${((searchRadius - 100) / (150000 - 100)) * 100}%, #e5e7eb 100%)`,
                          }}
                        />
                        <div className="text-xs text-right text-blue-600 font-medium">
                          {searchRadius >= 1000
                            ? `${(searchRadius / 1000).toFixed(1)} км`
                            : `${searchRadius} м`}
                        </div>
                      </>
                    )}
                  </div>
                )}

                {/* Подсказка или индикатор загрузки */}
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {isLoadingIsochrone ? (
                    <span className="flex items-center gap-1">
                      <span className="loading loading-spinner loading-xs"></span>
                      {t('loadingIsochrone')}
                    </span>
                  ) : (
                    t('clickToChangeMode')
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </Map>

      {/* Легенда и счетчик - вне Map */}
      <div className="absolute top-2 left-2 bg-base-100/90 dark:bg-base-200/90 backdrop-blur-sm rounded-lg shadow-lg p-2 pointer-events-auto">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
            <span className="text-xs text-base-content">{t('listings')}</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-amber-500 rounded-full"></div>
            <span className="text-xs text-base-content">
              {t('storefronts')}
            </span>
          </div>
          {showRadius && (
            <div className="pt-1 mt-1 border-t border-base-300">
              <p className="text-xs font-medium text-base-content">
                {filteredListings.length === listings.length
                  ? t('total', { count: listings.length })
                  : t('showing', {
                      shown: filteredListings.length,
                      total: listings.length,
                    })}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Кнопка геолокации - вне Map */}
      {isGeolocationAvailable && (
        <button
          onClick={handleGeolocation}
          className="absolute bottom-4 right-4 btn btn-sm btn-circle btn-primary shadow-lg hover:scale-110 transition-transform pointer-events-auto"
          title="Определить мое местоположение"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2}
            stroke="currentColor"
            className="w-4 h-4"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
            />
          </svg>
        </button>
      )}

      {/* Градиентная маска для эстетики */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute inset-x-0 top-0 h-8 bg-gradient-to-b from-white/20 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-8 bg-gradient-to-t from-white/20 to-transparent" />
      </div>

      {/* Стили для ползунков */}
      <style jsx>{`
        input[type='range']::-webkit-slider-thumb {
          appearance: none;
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: ${walkingMode === 'walking' ? '#10B981' : '#3b82f6'};
          cursor: pointer;
          border: 2px solid white;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }

        input[type='range']::-moz-range-thumb {
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: ${walkingMode === 'walking' ? '#10B981' : '#3b82f6'};
          cursor: pointer;
          border: 2px solid white;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
      `}</style>

      {/* Перетаскиваемая иконка для установки местоположения */}
      <DraggableLocationIcon
        mapRef={mapRef as React.RefObject<any>}
        onDropLocation={(lng, lat) => {
          setUserMarkerLocation({ longitude: lng, latitude: lat });
        }}
      />
    </div>
  );
};
