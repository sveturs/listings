import React, {
  useState,
  useCallback,
  useRef,
  useEffect,
  useMemo,
} from 'react';
import Map, { Marker, Source, Layer } from 'react-map-gl';
import type { MapRef, MarkerDragEvent } from 'react-map-gl';
import {
  MapViewState,
  MapMarkerData,
  MapPopupData,
  MapControlsConfig,
} from '../types/gis';
import { useGeoSearch } from '../hooks/useGeoSearch';
import { useGeolocation } from '../hooks/useGeolocation';
import MapMarker from './MapMarker';
import MapPopup from './MapPopup';
import MapControls from './MapControls';
import { MapCluster } from './MapCluster';
import 'mapbox-gl/dist/mapbox-gl.css';
import type { ClusterData } from '../types/gis';

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
  loadClusters?: (
    bounds: {
      north: number;
      south: number;
      east: number;
      west: number;
    },
    zoom: number
  ) => Promise<ClusterData[]>;
  // Новые пропсы для маркера покупателя
  showBuyerMarker?: boolean;
  buyerLocation?: {
    longitude: number;
    latitude: number;
  };
  searchRadius?: number; // в метрах
  onBuyerLocationChange?: (location: {
    longitude: number;
    latitude: number;
  }) => void;
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
  loadClusters,
  showBuyerMarker = false,
  buyerLocation,
  searchRadius = 10000, // 10км по умолчанию
  onBuyerLocationChange,
}) => {
  console.log('[InteractiveMap] Received markers:', markers);

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
  const [selectedMarker, setSelectedMarker] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [useOpenStreetMap, setUseOpenStreetMap] = useState(false);

  // Состояние для кластеров
  const [clusters, setClusters] = useState<ClusterData[]>([]);
  const [_isLoadingClusters, setIsLoadingClusters] = useState(false);

  // Состояние для маркера покупателя
  const [internalBuyerLocation, setInternalBuyerLocation] = useState({
    longitude: buyerLocation?.longitude || viewState.longitude,
    latitude: buyerLocation?.latitude || viewState.latitude,
  });

  // Порог зума для переключения между кластерами и маркерами
  const CLUSTER_ZOOM_THRESHOLD = 14;

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
      console.info('Using Mapbox GL with provided token');
    }
  }, [accessToken, openStreetMapStyle]);

  // Обновление внутреннего состояния при изменении внешнего buyerLocation
  useEffect(() => {
    if (buyerLocation) {
      setInternalBuyerLocation(buyerLocation);
    }
  }, [buyerLocation]);

  // Загрузка кластеров при изменении области просмотра
  const loadClustersData = useCallback(
    async (
      bounds: { north: number; south: number; east: number; west: number },
      zoom: number
    ) => {
      if (!loadClusters) return;

      setIsLoadingClusters(true);
      try {
        const clustersData = await loadClusters(bounds, zoom);
        setClusters(clustersData || []);
      } catch (error) {
        console.error('Error loading clusters:', error);
        setClusters([]);
      } finally {
        setIsLoadingClusters(false);
      }
    },
    [loadClusters]
  );

  // Получение границ карты
  const getMapBounds = useCallback(() => {
    if (!mapRef.current) return null;

    const bounds = mapRef.current.getBounds();
    if (!bounds) return null;

    return {
      north: bounds.getNorth(),
      south: bounds.getSouth(),
      east: bounds.getEast(),
      west: bounds.getWest(),
    };
  }, []);

  // Загрузка кластеров при изменении viewport
  useEffect(() => {
    if (!loadClusters || !mapRef.current) return;

    const bounds = getMapBounds();
    if (!bounds) return;

    // Debounce для оптимизации
    const timeoutId = setTimeout(() => {
      loadClustersData(bounds, viewState.zoom);
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [viewState, loadClusters, loadClustersData, getMapBounds]);

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
      setSelectedMarker(marker.id);
      if (onMarkerClick) {
        onMarkerClick(marker);
      }
    },
    [onMarkerClick]
  );

  const handleMapClick = useCallback(
    (event: any) => {
      setSelectedMarker(null);
      if (onMapClick) {
        onMapClick(event);
      }
    },
    [onMapClick]
  );

  const handleStyleChange = useCallback((newStyle: string) => {
    setMapStyle(newStyle);
  }, []);

  const handleSearch = useCallback(
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
  //       setSelectedMarker(markerId);
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

  const handleBuyerMarkerDragEnd = useCallback(
    (event: MarkerDragEvent) => {
      const newLocation = {
        longitude: event.lngLat.lng,
        latitude: event.lngLat.lat,
      };
      setInternalBuyerLocation(newLocation);
      if (onBuyerLocationChange) {
        onBuyerLocationChange(newLocation);
      }
    },
    [onBuyerLocationChange]
  );

  // GeoJSON для круга радиуса поиска
  const radiusCircleGeoJSON = useMemo(() => {
    if (!showBuyerMarker) return null;

    // Создаем круг вокруг позиции покупателя
    const center = [
      internalBuyerLocation.longitude,
      internalBuyerLocation.latitude,
    ];
    const radiusInKm = searchRadius / 1000;
    const options = { steps: 64, units: 'kilometers' as const };

    // Простая аппроксимация круга полигоном
    const points = [];
    const numPoints = options.steps;
    for (let i = 0; i < numPoints; i++) {
      const angle = (i / numPoints) * 2 * Math.PI;
      const dx = radiusInKm * Math.cos(angle);
      const dy = radiusInKm * Math.sin(angle);

      // Приблизительное преобразование км в градусы
      const lat = center[1] + dy / 111.32;
      const lng =
        center[0] + dx / (111.32 * Math.cos((center[1] * Math.PI) / 180));
      points.push([lng, lat]);
    }
    points.push(points[0]); // Замыкаем полигон

    return {
      type: 'Feature' as const,
      geometry: {
        type: 'Polygon' as const,
        coordinates: [points],
      },
      properties: {},
    };
  }, [showBuyerMarker, internalBuyerLocation, searchRadius]);

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
      'fill-color': '#3B82F6',
      'fill-opacity': 0.1,
    },
  };

  const radiusLineLayer = {
    id: 'radius-line',
    type: 'line' as const,
    paint: {
      'line-color': '#3B82F6',
      'line-width': 2,
      'line-opacity': 0.8,
      'line-dasharray': [2, 2],
    },
  };

  // Обработчик клика по кластеру
  const handleClusterClick = useCallback(
    (cluster: ClusterData) => {
      if (!mapRef.current) return;

      // Увеличиваем масштаб и центрируем на кластере
      mapRef.current.flyTo({
        center: [cluster.center.lng, cluster.center.lat],
        zoom: cluster.zoom_expand || viewState.zoom + 2, // Используем рекомендованный zoom или увеличиваем на 2
        duration: 1000,
      });
    },
    [viewState.zoom]
  );

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
      >
        {/* Кластеры или маркеры в зависимости от уровня зума */}
        {viewState.zoom < CLUSTER_ZOOM_THRESHOLD && clusters.length > 0
          ? // Показываем кластеры
            clusters.map((cluster, index) => (
              <Marker
                key={`cluster-${index}`}
                longitude={cluster.center.lng}
                latitude={cluster.center.lat}
                anchor="center"
              >
                <MapCluster
                  count={cluster.count}
                  onClick={() => handleClusterClick(cluster)}
                />
              </Marker>
            ))
          : // Показываем обычные маркеры
            markers.map((marker) => (
              <MapMarker
                key={marker.id}
                marker={marker}
                selected={selectedMarker === marker.id}
                onClick={handleMarkerClick}
              />
            ))}

        {/* Всплывающее окно */}
        {popup && (
          <MapPopup popup={popup} onClose={() => setSelectedMarker(null)} />
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
          onSearch={handleSearch}
          isMobile={isMobile}
          useOpenStreetMap={useOpenStreetMap}
        />
      </Map>

      {/* Индикатор загрузки */}
      {isLoading && (
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
