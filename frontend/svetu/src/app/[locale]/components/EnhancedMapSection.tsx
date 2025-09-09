'use client';

import React, { useMemo, useCallback, useEffect } from 'react';
import Map, {
  Marker,
  Source,
  Layer,
  NavigationControl,
  GeolocateControl,
} from 'react-map-gl';
import type {
  ViewState,
  LayerProps,
  MapRef,
  MarkerDragEvent,
} from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { useTranslations } from 'next-intl';
import {
  getMapboxIsochrone,
  isPointInIsochrone,
} from '@/components/GIS/utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';
import {
  // FiMapPin,
  // FiShoppingBag,
  FiFilter,
  FiX,
  FiMaximize2,
  FiMinimize2,
} from 'react-icons/fi';
import { motion, AnimatePresence } from 'framer-motion';

interface EnhancedMapSectionProps {
  listings?: Array<{
    id: string | number;
    latitude: number;
    longitude: number;
    price: number;
    title?: string;
    category?: string;
    imageUrl?: string;
    isStorefront?: boolean;
    storeName?: string;
  }>;
  userLocation?: {
    latitude: number;
    longitude: number;
  };
  searchRadius?: number;
  showRadius?: boolean;
  enableClustering?: boolean;
  className?: string;
}

export const EnhancedMapSection: React.FC<EnhancedMapSectionProps> = ({
  listings = [],
  userLocation,
  searchRadius: initialSearchRadius = 5000,
  showRadius: initialShowRadius = true,
  enableClustering = true,
  className = '',
}) => {
  const _t = useTranslations('map');
  const mapRef = React.useRef<MapRef>(null);
  const [theme, setTheme] = React.useState<'light' | 'dark'>('light');
  const [mounted, setMounted] = React.useState(false);
  const [showRadius, setShowRadius] = React.useState(initialShowRadius);
  const [searchRadius, setSearchRadius] = React.useState(initialSearchRadius);
  const [selectedListing, setSelectedListing] = React.useState<
    string | number | null
  >(null);
  const [walkingMode, setWalkingMode] = React.useState<'radius' | 'walking'>(
    'radius'
  );
  const [walkingTime, setWalkingTime] = React.useState(15);
  const [userMarkerLocation, setUserMarkerLocation] = React.useState(() => {
    if (
      userLocation &&
      typeof userLocation.longitude === 'number' &&
      typeof userLocation.latitude === 'number' &&
      !isNaN(userLocation.longitude) &&
      !isNaN(userLocation.latitude)
    ) {
      return userLocation;
    }
    return { latitude: 44.7866, longitude: 20.4489 }; // Белград по умолчанию
  });
  const [isCompactControlExpanded, setIsCompactControlExpanded] =
    React.useState(false);
  const [isochroneData, setIsochroneData] =
    React.useState<Feature<Polygon> | null>(null);
  const [_isLoadingIsochrone, setIsLoadingIsochrone] = React.useState(false);
  const [isFullscreen, setIsFullscreen] = React.useState(false);
  const [showFilters, setShowFilters] = React.useState(false);
  const [selectedCategories, setSelectedCategories] = React.useState<string[]>(
    []
  );

  // Категории для фильтров
  const categories = [
    'Электроника',
    'Мебель',
    'Авто',
    'Одежда',
    'Недвижимость',
    'Услуги',
  ];

  useEffect(() => {
    setMounted(true);
    const getTheme = () => {
      const htmlTheme = document.documentElement.getAttribute('data-theme');
      return htmlTheme === 'dark' ? 'dark' : 'light';
    };
    setTheme(getTheme());

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

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme'],
    });

    return () => observer.disconnect();
  }, []);

  // Загрузка изохрона
  useEffect(() => {
    if (!userMarkerLocation || !showRadius || walkingMode !== 'walking') {
      setIsochroneData(null);
      return;
    }

    // Валидация координат для изохрона
    if (
      typeof userMarkerLocation.longitude !== 'number' ||
      typeof userMarkerLocation.latitude !== 'number' ||
      isNaN(userMarkerLocation.longitude) ||
      isNaN(userMarkerLocation.latitude) ||
      userMarkerLocation.latitude < -90 ||
      userMarkerLocation.latitude > 90 ||
      userMarkerLocation.longitude < -180 ||
      userMarkerLocation.longitude > 180
    ) {
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
        setIsochroneData(null);
      } finally {
        setIsLoadingIsochrone(false);
      }
    };

    const timer = setTimeout(loadIsochrone, 500);
    return () => clearTimeout(timer);
  }, [userMarkerLocation, showRadius, walkingMode, walkingTime]);

  const calculateDistance = useCallback(
    (lat1: number, lon1: number, lat2: number, lon2: number) => {
      const R = 6371e3;
      const φ1 = (lat1 * Math.PI) / 180;
      const φ2 = (lat2 * Math.PI) / 180;
      const Δφ = ((lat2 - lat1) * Math.PI) / 180;
      const Δλ = ((lon2 - lon1) * Math.PI) / 180;

      const a =
        Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
        Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

      return R * c;
    },
    []
  );

  // Фильтрация по категориям и радиусу
  const filteredListings = useMemo(() => {
    // Сначала фильтруем только listings с валидными координатами
    let filtered = listings.filter(
      (listing) =>
        typeof listing.latitude === 'number' &&
        typeof listing.longitude === 'number' &&
        !isNaN(listing.latitude) &&
        !isNaN(listing.longitude) &&
        listing.latitude >= -90 &&
        listing.latitude <= 90 &&
        listing.longitude >= -180 &&
        listing.longitude <= 180
    );

    // Фильтр по категориям
    if (selectedCategories.length > 0) {
      filtered = filtered.filter(
        (listing) =>
          listing.category && selectedCategories.includes(listing.category)
      );
    }

    // Фильтр по радиусу/изохрону
    if (!showRadius || !userMarkerLocation) return filtered;

    if (walkingMode === 'walking' && isochroneData) {
      return filtered.filter((listing) => {
        try {
          return isPointInIsochrone(
            [listing.longitude, listing.latitude],
            isochroneData
          );
        } catch (error) {
          console.warn(
            'Error checking isochrone for listing:',
            listing.id,
            error
          );
          return false;
        }
      });
    }

    const effectiveRadius =
      walkingMode === 'walking' ? walkingTime * 80 : searchRadius;

    return filtered.filter((listing) => {
      try {
        const distance = calculateDistance(
          userMarkerLocation.latitude,
          userMarkerLocation.longitude,
          listing.latitude,
          listing.longitude
        );
        return !isNaN(distance) && distance <= effectiveRadius;
      } catch (error) {
        console.warn(
          'Error calculating distance for listing:',
          listing.id,
          error
        );
        return false;
      }
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
    selectedCategories,
  ]);

  // GeoJSON для кластеризации
  const geoJsonData = useMemo(() => {
    return {
      type: 'FeatureCollection' as const,
      features: filteredListings
        .filter(
          (listing) =>
            typeof listing.latitude === 'number' &&
            typeof listing.longitude === 'number' &&
            !isNaN(listing.latitude) &&
            !isNaN(listing.longitude) &&
            listing.latitude >= -90 &&
            listing.latitude <= 90 &&
            listing.longitude >= -180 &&
            listing.longitude <= 180
        )
        .map((listing) => ({
          type: 'Feature' as const,
          geometry: {
            type: 'Point' as const,
            coordinates: [listing.longitude, listing.latitude],
          },
          properties: {
            id: listing.id,
            price: typeof listing.price === 'string' ? parseFloat(listing.price) : (listing.price || 0),
            title: listing.title,
            category: listing.category,
            isStorefront: listing.isStorefront || false,
            storeName: listing.storeName,
            imageUrl: listing.imageUrl,
          },
        })),
    };
  }, [filteredListings]);

  // Центр и масштаб карты
  const { center, zoom } = useMemo(() => {
    if (
      userLocation &&
      typeof userLocation.longitude === 'number' &&
      typeof userLocation.latitude === 'number' &&
      !isNaN(userLocation.longitude) &&
      !isNaN(userLocation.latitude)
    ) {
      return {
        center: {
          longitude: userLocation.longitude,
          latitude: userLocation.latitude,
        },
        zoom: 13,
      };
    }

    if (listings.length > 0) {
      const validListings = listings.filter(
        (l) =>
          typeof l.latitude === 'number' &&
          typeof l.longitude === 'number' &&
          !isNaN(l.latitude) &&
          !isNaN(l.longitude)
      );

      if (validListings.length > 0) {
        const avgLat =
          validListings.reduce((sum, l) => sum + l.latitude, 0) /
          validListings.length;
        const avgLng =
          validListings.reduce((sum, l) => sum + l.longitude, 0) /
          validListings.length;
        return {
          center: { longitude: avgLng, latitude: avgLat },
          zoom: 12,
        };
      }
    }

    // Белград по умолчанию
    return {
      center: { longitude: 20.4489, latitude: 44.7866 },
      zoom: 11,
    };
  }, [listings, userLocation]);

  const initialViewState: Partial<ViewState> = {
    ...center,
    zoom,
    pitch: 0,
    bearing: 0,
  };

  // Слои карты
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
          '#60a5fa',
          10,
          '#3b82f6',
          30,
          '#2563eb',
        ],
        'circle-radius': ['step', ['get', 'point_count'], 15, 10, 20, 30, 25],
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.8,
      },
    }),
    []
  );

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
      paint: { 'text-color': '#ffffff' },
    }),
    []
  );

  const unclusteredPointLayer: LayerProps = useMemo(
    () => ({
      id: 'unclustered-point',
      type: 'circle',
      source: 'listings',
      filter: ['!', ['has', 'point_count']],
      paint: {
        'circle-color': ['case', ['get', 'isStorefront'], '#f59e0b', '#3b82f6'],
        'circle-radius': 8,
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.9,
      },
    }),
    []
  );

  // Радиус поиска
  const radiusGeoJson = useMemo(() => {
    if (!userMarkerLocation || !showRadius) return null;

    // Валидация координат пользователя
    if (
      !userMarkerLocation ||
      typeof userMarkerLocation.longitude !== 'number' ||
      typeof userMarkerLocation.latitude !== 'number' ||
      isNaN(userMarkerLocation.longitude) ||
      isNaN(userMarkerLocation.latitude) ||
      userMarkerLocation.latitude < -90 ||
      userMarkerLocation.latitude > 90 ||
      userMarkerLocation.longitude < -180 ||
      userMarkerLocation.longitude > 180
    ) {
      return null;
    }

    if (walkingMode === 'walking' && isochroneData) {
      return isochroneData;
    }

    const numSides = 64;
    const angleStep = (2 * Math.PI) / numSides;
    const coordinates: [number, number][] = [];
    const lat = userMarkerLocation.latitude;
    const latRad = (lat * Math.PI) / 180;
    const metersPerDegreeLat = 111320;
    const metersPerDegreeLng = 111320 * Math.cos(latRad);

    // Проверяем что metersPerDegreeLng не NaN или Infinity
    if (isNaN(metersPerDegreeLng) || !isFinite(metersPerDegreeLng)) {
      return null;
    }

    for (let i = 0; i <= numSides; i++) {
      const angle = i * angleStep;
      const dx = (searchRadius * Math.cos(angle)) / metersPerDegreeLng;
      const dy = (searchRadius * Math.sin(angle)) / metersPerDegreeLat;

      // Проверяем валидность вычисленных координат
      const newLng = userMarkerLocation.longitude + dx;
      const newLat = userMarkerLocation.latitude + dy;

      if (
        isNaN(newLng) ||
        isNaN(newLat) ||
        !isFinite(newLng) ||
        !isFinite(newLat)
      ) {
        continue; // Пропускаем невалидные координаты
      }

      coordinates.push([newLng, newLat]);
    }

    // Проверяем что у нас есть минимум координат для полигона
    if (coordinates.length < 3) {
      return null;
    }

    return {
      type: 'Feature' as const,
      geometry: { type: 'Polygon' as const, coordinates: [coordinates] },
      properties: {},
    };
  }, [
    userMarkerLocation,
    searchRadius,
    showRadius,
    walkingMode,
    isochroneData,
  ]);

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

  const handleMapClick = useCallback((event: any) => {
    const features = event.features;
    if (!features || features.length === 0) return;

    const feature = features[0];

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
      setSelectedListing(feature.properties.id);
      // Можно добавить popup или модалку с деталями
    }
  }, []);

  const handleUserMarkerDrag = useCallback((event: MarkerDragEvent) => {
    const lng = event.lngLat.lng;
    const lat = event.lngLat.lat;

    // Проверяем валидность координат
    if (
      typeof lng === 'number' &&
      typeof lat === 'number' &&
      !isNaN(lng) &&
      !isNaN(lat) &&
      lat >= -90 &&
      lat <= 90 &&
      lng >= -180 &&
      lng <= 180
    ) {
      setUserMarkerLocation({
        longitude: lng,
        latitude: lat,
      });
    }
  }, []);

  if (!mounted) {
    return (
      <div
        className={`relative overflow-hidden rounded-lg bg-base-200 animate-pulse ${className}`}
      >
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`relative overflow-hidden rounded-lg ${isFullscreen ? 'fixed inset-0 z-50' : className}`}
    >
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
        {/* Контролы навигации */}
        <NavigationControl position="top-left" />
        <GeolocateControl position="top-left" />

        {/* Радиус поиска */}
        {radiusGeoJson && (
          <Source id="radius" type="geojson" data={radiusGeoJson}>
            <Layer {...radiusLayer} />
            <Layer {...radiusBorderLayer} />
          </Source>
        )}

        {/* Маркер пользователя */}
        {userMarkerLocation &&
          typeof userMarkerLocation.longitude === 'number' &&
          typeof userMarkerLocation.latitude === 'number' &&
          !isNaN(userMarkerLocation.longitude) &&
          !isNaN(userMarkerLocation.latitude) && (
            <Marker
              longitude={userMarkerLocation.longitude}
              latitude={userMarkerLocation.latitude}
              draggable
              onDrag={handleUserMarkerDrag}
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
                  <ellipse
                    cx="20"
                    cy="38"
                    rx="8"
                    ry="2"
                    fill="black"
                    fillOpacity="0.2"
                  />
                  <path
                    d="M20 36C20 36 32 24 32 16C32 9.37258 26.6274 4 20 4C13.3726 4 8 9.37258 8 16C8 24 20 36 20 36Z"
                    fill="#EF4444"
                    stroke="white"
                    strokeWidth="2"
                  />
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
            <Layer {...unclusteredPointLayer} />
          </Source>
        ) : (
          filteredListings
            .filter(
              (listing) =>
                typeof listing.latitude === 'number' &&
                typeof listing.longitude === 'number' &&
                !isNaN(listing.latitude) &&
                !isNaN(listing.longitude) &&
                listing.latitude >= -90 &&
                listing.latitude <= 90 &&
                listing.longitude >= -180 &&
                listing.longitude <= 180
            )
            .map((listing) => (
              <Marker
                key={listing.id}
                longitude={listing.longitude}
                latitude={listing.latitude}
                anchor="bottom"
                onClick={(e) => {
                  e.originalEvent.stopPropagation();
                  setSelectedListing(listing.id);
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
                    €{listing.price.toLocaleString()}
                  </div>
                </div>
              </Marker>
            ))
        )}
      </Map>

      {/* Панель управления */}
      <div className="absolute top-4 right-4 z-10 space-y-2">
        {/* Кнопка полноэкранного режима */}
        <button
          onClick={() => setIsFullscreen(!isFullscreen)}
          className="btn btn-sm btn-circle bg-white shadow-lg hover:shadow-xl"
        >
          {isFullscreen ? <FiMinimize2 /> : <FiMaximize2 />}
        </button>

        {/* Кнопка фильтров */}
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="btn btn-sm btn-circle bg-white shadow-lg hover:shadow-xl relative"
        >
          <FiFilter />
          {selectedCategories.length > 0 && (
            <span className="absolute -top-1 -right-1 badge badge-primary badge-xs">
              {selectedCategories.length}
            </span>
          )}
        </button>

        {/* Контрол радиуса */}
        {!isCompactControlExpanded ? (
          <button
            className="btn btn-sm btn-circle bg-white shadow-lg hover:shadow-xl"
            onClick={() => setIsCompactControlExpanded(true)}
          >
            <span className="text-lg">
              {walkingMode === 'walking' ? '🚶' : '📏'}
            </span>
            <div className="absolute -bottom-1 -right-1 badge badge-primary badge-xs">
              {walkingMode === 'walking'
                ? `${walkingTime}'`
                : searchRadius >= 1000
                  ? `${(searchRadius / 1000).toFixed(0)}км`
                  : `${searchRadius}м`}
            </div>
          </button>
        ) : (
          <div className="bg-white rounded-lg shadow-lg p-3 w-64">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm font-medium">
                {walkingMode === 'walking'
                  ? 'Пешая доступность'
                  : 'Радиус поиска'}
              </span>
              <button
                onClick={() => setIsCompactControlExpanded(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <FiX size={16} />
              </button>
            </div>

            <div className="space-y-3">
              <div className="flex gap-2">
                <button
                  onClick={() => setWalkingMode('radius')}
                  className={`btn btn-xs flex-1 ${walkingMode === 'radius' ? 'btn-primary' : 'btn-outline'}`}
                >
                  Радиус
                </button>
                <button
                  onClick={() => setWalkingMode('walking')}
                  className={`btn btn-xs flex-1 ${walkingMode === 'walking' ? 'btn-primary' : 'btn-outline'}`}
                >
                  Пешком
                </button>
              </div>

              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm checkbox-primary"
                  checked={showRadius}
                  onChange={(e) => setShowRadius(e.target.checked)}
                />
                <span className="text-xs">Показать зону</span>
              </label>

              {showRadius && (
                <div>
                  {walkingMode === 'walking' ? (
                    <>
                      <input
                        type="range"
                        min="5"
                        max="30"
                        value={walkingTime}
                        onChange={(e) => setWalkingTime(Number(e.target.value))}
                        className="range range-primary range-xs"
                      />
                      <div className="text-xs text-right text-primary mt-1">
                        {walkingTime} минут пешком
                      </div>
                    </>
                  ) : (
                    <>
                      <input
                        type="range"
                        min="100"
                        max="10000"
                        step="100"
                        value={searchRadius}
                        onChange={(e) =>
                          setSearchRadius(Number(e.target.value))
                        }
                        className="range range-primary range-xs"
                      />
                      <div className="text-xs text-right text-primary mt-1">
                        {searchRadius >= 1000
                          ? `${(searchRadius / 1000).toFixed(1)} км`
                          : `${searchRadius} м`}
                      </div>
                    </>
                  )}
                </div>
              )}
            </div>
          </div>
        )}
      </div>

      {/* Панель фильтров */}
      <AnimatePresence>
        {showFilters && (
          <motion.div
            initial={{ x: 300, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: 300, opacity: 0 }}
            transition={{ type: 'spring', damping: 25 }}
            className="absolute top-16 right-4 z-10 bg-white rounded-lg shadow-lg p-4 w-64"
          >
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold">Фильтры</h3>
              <button
                onClick={() => setShowFilters(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <FiX size={16} />
              </button>
            </div>

            <div className="space-y-2">
              <p className="text-sm font-medium mb-2">Категории:</p>
              {categories.map((cat) => (
                <label
                  key={cat}
                  className="flex items-center gap-2 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm checkbox-primary"
                    checked={selectedCategories.includes(cat)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedCategories([...selectedCategories, cat]);
                      } else {
                        setSelectedCategories(
                          selectedCategories.filter((c) => c !== cat)
                        );
                      }
                    }}
                  />
                  <span className="text-sm">{cat}</span>
                </label>
              ))}
              {selectedCategories.length > 0 && (
                <button
                  onClick={() => setSelectedCategories([])}
                  className="btn btn-xs btn-ghost w-full"
                >
                  Сбросить фильтры
                </button>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Легенда и счетчик */}
      <div className="absolute bottom-4 left-4 bg-white/90 backdrop-blur-sm rounded-lg shadow-lg p-3">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
            <span className="text-xs">{_t('listings')}</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-amber-500 rounded-full"></div>
            <span className="text-xs">{_t('storefronts')}</span>
          </div>
          <div className="pt-1 mt-1 border-t">
            <p className="text-xs font-medium">
              {_t('showing', {
                shown: filteredListings.length,
                total: listings.length,
              })}
            </p>
          </div>
        </div>
      </div>

      {/* Информация о выбранном объявлении */}
      <AnimatePresence>
        {selectedListing && (
          <motion.div
            initial={{ y: 100, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: 100, opacity: 0 }}
            className="absolute bottom-20 left-4 right-4 bg-white rounded-lg shadow-lg p-4 max-w-sm mx-auto"
          >
            {(() => {
              const listing = listings.find((l) => l.id === selectedListing);
              if (!listing) return null;

              return (
                <div className="flex gap-3">
                  {listing.imageUrl && (
                    <img
                      src={listing.imageUrl}
                      alt={listing.title}
                      className="w-20 h-20 rounded-lg object-cover"
                    />
                  )}
                  <div className="flex-1">
                    <h3 className="font-semibold">
                      {listing.title || _t('listing')}
                    </h3>
                    {listing.category && (
                      <p className="text-xs text-gray-500">
                        {listing.category}
                      </p>
                    )}
                    <p className="text-lg font-bold text-primary mt-1">
                      €{listing.price.toLocaleString()}
                    </p>
                    <div className="flex gap-2 mt-2">
                      <button className="btn btn-xs btn-primary">
                        Подробнее
                      </button>
                      <button
                        onClick={() => setSelectedListing(null)}
                        className="btn btn-xs btn-ghost"
                      >
                        Закрыть
                      </button>
                    </div>
                  </div>
                </div>
              );
            })()}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
