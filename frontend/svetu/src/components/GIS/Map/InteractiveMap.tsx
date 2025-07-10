import React, { useState, useCallback, useRef, useEffect } from 'react';
import Map, { MapRef } from 'react-map-gl';
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
import 'mapbox-gl/dist/mapbox-gl.css';

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
}) => {
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

  // Получение токена Mapbox из переменных окружения
  const accessToken =
    mapboxAccessToken || process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  useEffect(() => {
    if (!accessToken) {
      console.error('Mapbox access token is not provided');
    }
  }, [accessToken]);

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

  if (!accessToken) {
    return (
      <div
        className={`flex items-center justify-center bg-gray-100 ${className}`}
        style={style}
      >
        <div className="text-center">
          <div className="text-gray-500 mb-2">⚠️</div>
          <p className="text-gray-600">Mapbox access token not configured</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative ${className}`} style={style}>
      <Map
        ref={mapRef}
        {...viewState}
        onMove={(evt) => handleViewStateChange(evt.viewState)}
        onClick={handleMapClick}
        mapStyle={mapStyle}
        mapboxAccessToken={accessToken}
        attributionControl={false}
        logoPosition="bottom-left"
        style={{ width: '100%', height: '100%' }}
      >
        {/* Маркеры */}
        {markers.map((marker) => (
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

        {/* Контролы */}
        <MapControls
          config={controlsConfig}
          onStyleChange={handleStyleChange}
          onSearch={handleSearch}
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
