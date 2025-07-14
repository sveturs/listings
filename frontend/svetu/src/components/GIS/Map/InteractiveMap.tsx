import React, {
  useState,
  useCallback,
  useRef,
  useEffect,
  useMemo,
} from 'react';

// Хук для детекции fullscreen режима
const useFullscreen = () => {
  const [isFullscreen, setIsFullscreen] = useState(false);

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  return isFullscreen;
};
import Map, { Marker, Source, Layer } from 'react-map-gl';
import type { MapRef, MarkerDragEvent } from 'react-map-gl';
import circle from '@turf/circle';
import {
  MapViewState,
  MapMarkerData,
  MapPopupData,
  MapControlsConfig,
} from '../types/gis';
import { generateStylizedIsochrone } from '../utils/isochrone';
import { getMapboxIsochrone } from '../utils/mapboxIsochrone';
import type { Feature, Polygon } from 'geojson';
import { useGeoSearch } from '../hooks/useGeoSearch';
import { useGeolocation } from '../hooks/useGeolocation';
import MapPopup from './MapPopup';
import MapControls from './MapControls';
import MapboxClusterLayer from './MapboxClusterLayer';
import MarkerHoverPopup from './MarkerHoverPopup';
// import NativeSliderControl from './NativeSliderControl';
import CompactSliderControl from './CompactSliderControl';
import FloatingSliderControl from './FloatingSliderControl';
import 'mapbox-gl/dist/mapbox-gl.css';
import '@/styles/map-popup.css';

interface InteractiveMapProps {
  initialViewState?: Partial<MapViewState>;
  markers?: MapMarkerData[];
  popup?: MapPopupData | null;
  onMarkerClick?: (marker: MapMarkerData) => void;
  onMapClick?: (event: any) => void;
  onViewStateChange?: (viewState: MapViewState) => void;
  controlsConfig?: MapControlsConfig;
  className?: string;
  style?: React.CSSProperties;
  mapboxAccessToken?: string;
  isMobile?: boolean;
  // Новые пропсы для маркера покупателя
  showBuyerMarker?: boolean;
  buyerLocation?: {
    longitude: number;
    latitude: number;
  };
  searchRadius?: number; // в метрах
  walkingMode?: 'radius' | 'walking';
  walkingTime?: number; // в минутах
  onBuyerLocationChange?: (location: {
    longitude: number;
    latitude: number;
  }) => void;
  onIsochroneChange?: (isochrone: Feature<Polygon> | null) => void;
  onWalkingModeChange?: (mode: 'radius' | 'walking') => void;
  onWalkingTimeChange?: (time: number) => void;
  onSearchRadiusChange?: (radius: number) => void;
  useNativeControl?: boolean; // Флаг для выбора типа контрола
  controlTranslations?: {
    walkingAccessibility: string;
    searchRadius: string;
    minutes: string;
    km: string;
    m: string;
    changeModeHint: string;
    holdForSettings: string;
    singleClickHint: string;
    mobileHint: string;
    desktopHint: string;
    updatingIsochrone: string;
  };
  // Визуализация границ районов
  districtBoundary?: Feature<Polygon> | null;
  onDistrictBoundaryChange?: (boundary: Feature<Polygon> | null) => void;
}

const InteractiveMap: React.FC<InteractiveMapProps> = ({
  initialViewState = {
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 12,
  },
  markers = [],
  popup = null,
  onMarkerClick,
  onMapClick,
  onViewStateChange,
  controlsConfig,
  className = '',
  style,
  mapboxAccessToken,
  isMobile = false,
  showBuyerMarker = false,
  buyerLocation,
  searchRadius = 10000, // 10км по умолчанию
  walkingMode = 'radius',
  walkingTime = 15,
  onBuyerLocationChange,
  onIsochroneChange,
  onWalkingModeChange,
  onWalkingTimeChange,
  onSearchRadiusChange,
  useNativeControl = false,
  controlTranslations,
  districtBoundary = null,
  onDistrictBoundaryChange: _onDistrictBoundaryChange,
}) => {
  // Детекция fullscreen режима
  const isFullscreen = useFullscreen();

  const mapRef = useRef<MapRef>(null);
  const { search } = useGeoSearch();
  const { getCurrentPosition } = useGeolocation();

  const [viewState, setViewState] = useState<MapViewState>({
    longitude: initialViewState.longitude || 20.4649,
    latitude: initialViewState.latitude || 44.8176,
    zoom: initialViewState.zoom || 12,
    pitch: initialViewState.pitch || 0,
    bearing: initialViewState.bearing || 0,
  });

  const [mapStyle, setMapStyle] = useState(
    'mapbox://styles/mapbox/streets-v12'
  );
  const [isLoading, setIsLoading] = useState(false);
  const [useOpenStreetMap, setUseOpenStreetMap] = useState(false);

  // Состояние для hover popup
  const [hoveredMarker, setHoveredMarker] = useState<MapMarkerData | null>(
    null
  );

  // Состояние для маркера покупателя
  const [internalBuyerLocation, setInternalBuyerLocation] = useState({
    longitude: buyerLocation?.longitude || viewState.longitude,
    latitude: buyerLocation?.latitude || viewState.latitude,
  });

  // Получение токена Mapbox из переменных окружения
  const accessToken =
    mapboxAccessToken || process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  // Создаем кастомный стиль для OpenStreetMap
  const openStreetMapStyle = useMemo(
    () => ({
      version: 8,
      sources: {
        'osm-tiles': {
          type: 'raster',
          tiles: [
            'https://a.tile.openstreetmap.org/{z}/{x}/{y}.png',
            'https://b.tile.openstreetmap.org/{z}/{x}/{y}.png',
            'https://c.tile.openstreetmap.org/{z}/{x}/{y}.png',
          ],
          tileSize: 256,
          attribution: '© OpenStreetMap contributors',
        },
      },
      layers: [
        {
          id: 'osm-tiles',
          type: 'raster',
          source: 'osm-tiles',
          minzoom: 0,
          maxzoom: 19,
        },
      ],
    }),
    []
  );

  // Автоматическое переключение на OpenStreetMap если токен не предоставлен
  useEffect(() => {
    if (!accessToken) {
      console.warn(
        'Mapbox access token is not provided, using OpenStreetMap as fallback'
      );
      setUseOpenStreetMap(true);
      setMapStyle(openStreetMapStyle as any);
    } else {
    }
  }, [accessToken, openStreetMapStyle]);

  // Обновление внутреннего состояния при изменении внешнего buyerLocation
  useEffect(() => {
    if (
      buyerLocation &&
      (buyerLocation.longitude !== internalBuyerLocation.longitude ||
        buyerLocation.latitude !== internalBuyerLocation.latitude)
    ) {
      setInternalBuyerLocation(buyerLocation);
    }
  }, [buyerLocation?.longitude, buyerLocation?.latitude]);

  // Логирование изменений границ района
  useEffect(() => {
    console.log('🗺️ District boundary in InteractiveMap:', districtBoundary);
    if (districtBoundary) {
      console.log('🗺️ District boundary type:', districtBoundary.type);
      console.log('🗺️ District boundary geometry:', districtBoundary.geometry);
      console.log(
        '🗺️ District boundary geometry type:',
        districtBoundary.geometry?.type
      );
      console.log(
        '🗺️ District boundary coordinates length:',
        districtBoundary.geometry?.coordinates?.length
      );
      if (districtBoundary.geometry?.coordinates?.[0]) {
        console.log(
          '🗺️ First coordinate ring length:',
          districtBoundary.geometry.coordinates[0].length
        );
        console.log(
          '🗺️ First few coordinates:',
          districtBoundary.geometry.coordinates[0].slice(0, 3)
        );
        console.log(
          '🗺️ Last few coordinates:',
          districtBoundary.geometry.coordinates[0].slice(-3)
        );

        // Проверяем, что полигон замкнут (первая и последняя точки должны совпадать)
        const coords = districtBoundary.geometry.coordinates[0];
        const first = coords[0];
        const last = coords[coords.length - 1];
        const isClosed = first[0] === last[0] && first[1] === last[1];
        console.log('🗺️ Polygon closed?', isClosed);

        // Проверяем валидность координат
        const hasValidCoords = coords.every(
          (coord) =>
            Array.isArray(coord) &&
            coord.length === 2 &&
            typeof coord[0] === 'number' &&
            typeof coord[1] === 'number' &&
            coord[0] >= -180 &&
            coord[0] <= 180 &&
            coord[1] >= -90 &&
            coord[1] <= 90
        );
        console.log('🗺️ Valid coordinates?', hasValidCoords);
      }
    }
  }, [districtBoundary]);

  const handleViewStateChange = useCallback(
    (newViewState: MapViewState) => {
      setViewState(newViewState);
      if (onViewStateChange) {
        onViewStateChange(newViewState);
      }
    },
    [onViewStateChange]
  );

  const handleMarkerClick = useCallback(
    (marker: MapMarkerData) => {
      if (onMarkerClick) {
        onMarkerClick(marker);
      }
    },
    [onMarkerClick]
  );

  const handleMapClick = useCallback(
    (event: any) => {
      if (onMapClick) {
        onMapClick(event);
      }
    },
    [onMapClick]
  );

  const handleStyleChange = useCallback((newStyle: string) => {
    setMapStyle(newStyle);
  }, []);

  const _handleSearch = useCallback(
    async (query: string) => {
      setIsLoading(true);
      try {
        const results = await search({
          query,
          limit: 1,
          language: 'ru',
        });

        if (results.length > 0) {
          const result = results[0];
          const newViewState = {
            longitude: parseFloat(result.lon),
            latitude: parseFloat(result.lat),
            zoom: 14,
            pitch: 0,
            bearing: 0,
          };

          setViewState(newViewState);

          // Анимация к результату поиска
          if (mapRef.current) {
            mapRef.current.flyTo({
              center: [newViewState.longitude, newViewState.latitude],
              zoom: newViewState.zoom,
              duration: 2000,
            });
          }
        }
      } catch (error) {
        console.error('Search error:', error);
      } finally {
        setIsLoading(false);
      }
    },
    [search]
  );

  const handleGeolocation = useCallback(async () => {
    try {
      const position = await getCurrentPosition();
      const newViewState = {
        longitude: position.longitude,
        latitude: position.latitude,
        zoom: 15,
        pitch: 0,
        bearing: 0,
      };

      setViewState(newViewState);

      if (mapRef.current) {
        mapRef.current.flyTo({
          center: [position.longitude, position.latitude],
          zoom: 15,
          duration: 2000,
        });
      }
    } catch (error) {
      console.error('Geolocation error:', error);
    }
  }, [getCurrentPosition]);

  // Функция для полета к маркеру (может использоваться в будущем)
  // const flyToMarker = useCallback(
  //   (markerId: string) => {
  //     const marker = markers.find((m) => m.id === markerId);
  //     if (marker && mapRef.current) {
  //       mapRef.current.flyTo({
  //         center: [marker.position[0], marker.position[1]],
  //         zoom: 16,
  //         duration: 1500,
  //       });
  //     }
  //   },
  //   [markers]
  // );

  // Обработчик перетаскивания маркера покупателя
  const handleBuyerMarkerDrag = useCallback((event: MarkerDragEvent) => {
    const newLocation = {
      longitude: event.lngLat.lng,
      latitude: event.lngLat.lat,
    };
    setInternalBuyerLocation(newLocation);
  }, []);

  const handleBuyerMarkerDragStart = useCallback(() => {
    setIsDragging(true);
  }, []);

  const handleBuyerMarkerDragEnd = useCallback(
    (event: MarkerDragEvent) => {
      const newLocation = {
        longitude: event.lngLat.lng,
        latitude: event.lngLat.lat,
      };
      setInternalBuyerLocation(newLocation);
      setIsDragging(false); // Теперь разрешаем пересчет изохрона
      if (onBuyerLocationChange) {
        onBuyerLocationChange(newLocation);
      }
    },
    [onBuyerLocationChange]
  );

  // Состояние для изохроны
  const [isochroneData, setIsochroneData] = useState<any>(null);
  const [isLoadingIsochrone, setIsLoadingIsochrone] = useState(false);
  const [_isMapLoaded, setIsMapLoaded] = useState(false);

  // Обработчики для hover
  const handleMarkerHover = useCallback((marker: MapMarkerData) => {
    setHoveredMarker(marker);
  }, []);

  const handleMarkerLeave = useCallback(() => {
    setHoveredMarker(null);
  }, []);

  // Состояние для отслеживания перетаскивания
  const [isDragging, setIsDragging] = useState(false);

  // Функция для загрузки изохроны
  const loadIsochrone = useCallback(async () => {
    if (!showBuyerMarker || walkingMode !== 'walking') {
      // Очищаем изохрон только при переключении режима
      setIsochroneData(null);
      if (onIsochroneChange) {
        onIsochroneChange(null);
      }
      return;
    }

    const center: [number, number] = [
      internalBuyerLocation.longitude,
      internalBuyerLocation.latitude,
    ];

    // НЕ очищаем старый изохрон до загрузки нового - это предотвращает мигание
    setIsLoadingIsochrone(true);

    // Используем переданное время ходьбы или 10 минут по умолчанию
    const timeInMinutes = walkingTime || 10;

    try {
      const isochrone = await getMapboxIsochrone({
        coordinates: center,
        minutes: timeInMinutes,
        profile: 'walking',
      });

      // Обновляем изохрон только после успешной загрузки
      setIsochroneData(isochrone);
      if (onIsochroneChange) {
        onIsochroneChange(isochrone);
      }
    } catch (error) {
      console.error('[InteractiveMap] Failed to fetch isochrone:', error);
      // При ошибке оставляем старый изохрон, не очищаем
    } finally {
      setIsLoadingIsochrone(false);
    }
  }, [
    showBuyerMarker,
    walkingMode,
    walkingTime,
    internalBuyerLocation.longitude,
    internalBuyerLocation.latitude,
    onIsochroneChange,
  ]);

  // Эффект для загрузки изохроны при изменении параметров (НЕ при перетаскивании)
  useEffect(() => {
    if (!isDragging) {
      loadIsochrone();
    }
  }, [loadIsochrone, isDragging]);

  // GeoJSON для радиуса поиска (круг или изохрона)
  const radiusCircleGeoJSON = useMemo(() => {
    if (!showBuyerMarker) return null;

    const center: [number, number] = [
      internalBuyerLocation.longitude,
      internalBuyerLocation.latitude,
    ];

    // Выбираем между радиусом и изохроной в зависимости от режима
    if (walkingMode === 'walking') {
      // Используем загруженную изохрону или fallback на локальную генерацию
      if (isochroneData) {
        return isochroneData;
      } else if (!isLoadingIsochrone) {
        // Если API недоступен и загрузка завершена, используем локальную генерацию
        const isochrone = generateStylizedIsochrone(center, 10); // 10 минут для пешехода
        return isochrone;
      }
      return null;
    } else {
      // Создаем обычный круг с помощью Turf.js
      const radiusInKm = searchRadius / 1000;
      const circleFeature = circle(center, radiusInKm, {
        steps: 64,
        units: 'kilometers',
      });
      return circleFeature;
    }
  }, [
    showBuyerMarker,
    internalBuyerLocation,
    searchRadius,
    walkingMode,
    isochroneData,
    isLoadingIsochrone,
  ]);

  // Стиль для слоя круга (закомментирован, не используется)
  // const radiusCircleLayer: CircleLayer = {
  //   id: 'radius-circle',
  //   type: 'circle',
  //   paint: {
  //     'circle-radius': 0,
  //     'circle-color': 'transparent',
  //   },
  // };

  const radiusFillLayer = {
    id: 'radius-fill',
    type: 'fill' as const,
    paint: {
      'fill-color': walkingMode === 'walking' ? '#10B981' : '#3B82F6', // Зеленый для пешеходного режима
      'fill-opacity': walkingMode === 'walking' ? 0.2 : 0.1, // Более заметная для пешеходного режима
      'fill-opacity-transition': {
        duration: 300,
        delay: 0,
      },
    },
  };

  const radiusLineLayer = {
    id: 'radius-line',
    type: 'line' as const,
    paint: {
      'line-color': walkingMode === 'walking' ? '#10B981' : '#3B82F6', // Зеленый для пешеходного режима
      'line-width': walkingMode === 'walking' ? 3 : 2, // Толще для пешеходного режима
      'line-opacity': 0.8,
      'line-dasharray': walkingMode === 'walking' ? [4, 4] : [2, 2], // Другой пунктир для пешеходного режима
      'line-opacity-transition': {
        duration: 300,
        delay: 0,
      },
    },
  };

  const fitBounds = useCallback(() => {
    if (markers.length > 0 && mapRef.current) {
      const bounds = markers.reduce(
        (acc, marker) => {
          return {
            minLng: Math.min(acc.minLng, marker.position[0]),
            maxLng: Math.max(acc.maxLng, marker.position[0]),
            minLat: Math.min(acc.minLat, marker.position[1]),
            maxLat: Math.max(acc.maxLat, marker.position[1]),
          };
        },
        {
          minLng: markers[0].position[0],
          maxLng: markers[0].position[0],
          minLat: markers[0].position[1],
          maxLat: markers[0].position[1],
        }
      );

      mapRef.current.fitBounds(
        [
          [bounds.minLng, bounds.minLat],
          [bounds.maxLng, bounds.maxLat],
        ],
        {
          padding: 50,
          duration: 1500,
        }
      );
    }
  }, [markers]);

  // Убираем проверку на отсутствие токена, так как теперь используем fallback

  return (
    <div className={`relative ${className}`} style={style}>
      <Map
        ref={mapRef}
        {...viewState}
        onMove={(evt) => handleViewStateChange(evt.viewState)}
        onClick={handleMapClick}
        mapStyle={mapStyle}
        mapboxAccessToken={useOpenStreetMap ? 'pk.dummy' : accessToken}
        attributionControl={false}
        logoPosition="bottom-left"
        style={{ width: '100%', height: '100%' }}
        onLoad={() => {
          setIsMapLoaded(true);
        }}
      >
        {/* Кластеризация маркеров с помощью MapboxClusterLayer */}
        {markers.length > 0 && (
          <MapboxClusterLayer
            markers={markers}
            onMarkerClick={handleMarkerClick}
            onMarkerHover={handleMarkerHover}
            onMarkerLeave={handleMarkerLeave}
            clusterRadius={50}
            clusterMaxZoom={14}
            clusterMinPoints={2}
            showPrices={false}
          />
        )}

        {/* Всплывающее окно */}
        {popup && <MapPopup popup={popup} onClose={() => {}} />}

        {/* Hover popup */}
        {hoveredMarker && (
          <MarkerHoverPopup
            marker={hoveredMarker}
            onClose={() => setHoveredMarker(null)}
          />
        )}

        {/* Границы выбранного района */}
        {districtBoundary && (
          <Source
            id="district-boundary-source"
            type="geojson"
            data={districtBoundary}
          >
            <Layer
              id="district-boundary-fill"
              type="fill"
              paint={{
                'fill-color': '#3b82f6',
                'fill-opacity': 0.3,
              }}
            />
            <Layer
              id="district-boundary-line"
              type="line"
              paint={{
                'line-color': '#3b82f6',
                'line-width': 4,
                'line-opacity': 1.0,
              }}
            />
          </Source>
        )}

        {/* Слой с радиусом поиска */}
        {showBuyerMarker && radiusCircleGeoJSON && (
          <Source type="geojson" data={radiusCircleGeoJSON}>
            <Layer {...radiusFillLayer} />
            <Layer {...radiusLineLayer} />
          </Source>
        )}

        {/* Маркер покупателя */}
        {showBuyerMarker && (
          <Marker
            longitude={internalBuyerLocation.longitude}
            latitude={internalBuyerLocation.latitude}
            draggable
            onDrag={handleBuyerMarkerDrag}
            onDragStart={handleBuyerMarkerDragStart}
            onDragEnd={handleBuyerMarkerDragEnd}
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

        {/* Контролы */}
        <MapControls
          config={controlsConfig}
          onStyleChange={handleStyleChange}
          isMobile={isMobile}
          useOpenStreetMap={useOpenStreetMap}
        />

        {/* Нативный контрол Mapbox */}
        {showBuyerMarker && useNativeControl && mapRef.current && (
          <CompactSliderControl
            map={mapRef.current.getMap()}
            mode={walkingMode}
            onModeChange={(mode) => {
              onWalkingModeChange?.(mode);
            }}
            walkingTime={walkingTime}
            onWalkingTimeChange={(time) => {
              onWalkingTimeChange?.(time);
            }}
            searchRadius={searchRadius}
            onRadiusChange={(radius) => {
              onSearchRadiusChange?.(radius);
            }}
            isFullscreen={isFullscreen}
            isMobile={isMobile}
            translations={controlTranslations}
          />
        )}
      </Map>

      {/* Плавающий контрол с выдвижным слайдером - вне MapBox контейнера */}
      {showBuyerMarker && !useNativeControl && (
        <FloatingSliderControl
          mode={walkingMode}
          isFullscreen={isFullscreen}
          isMobile={isMobile}
          onModeChange={(mode) => {
            onWalkingModeChange?.(mode);
          }}
          walkingTime={walkingTime}
          onWalkingTimeChange={(time) => {
            onWalkingTimeChange?.(time);
          }}
          searchRadius={searchRadius}
          onRadiusChange={(radius) => {
            onSearchRadiusChange?.(radius);
          }}
          translations={controlTranslations}
        />
      )}

      {/* Индикатор загрузки изохрона */}
      {isLoadingIsochrone && (
        <div className="absolute top-4 right-4 z-20">
          <div className="bg-white rounded-lg p-3 shadow-lg border border-gray-200">
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-green-500"></div>
              <span className="text-sm text-gray-600">
                {controlTranslations?.updatingIsochrone ||
                  'Updating accessibility zone...'}
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Индикатор общей загрузки */}
      {isLoading && !isLoadingIsochrone && (
        <div className="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center z-20">
          <div className="bg-white rounded-lg p-4 shadow-lg">
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
              <span className="text-sm text-gray-600">Поиск...</span>
            </div>
          </div>
        </div>
      )}

      {/* Панель быстрых действий */}
      <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-10">
        <div className="bg-white rounded-lg shadow-lg p-2 flex space-x-2">
          <button
            onClick={handleGeolocation}
            className="px-3 py-2 text-sm bg-primary text-white rounded-md hover:bg-primary-dark transition-colors"
            title="Мое местоположение"
          >
            📍 Где я?
          </button>

          {markers.length > 0 && (
            <button
              onClick={fitBounds}
              className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors"
              title="Показать все маркеры"
            >
              🔍 Показать все
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default InteractiveMap;
export { InteractiveMap };
