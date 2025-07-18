'use client';

import React, { useState, useCallback } from 'react';
import Map from 'react-map-gl';
import { MapboxClusterLayer } from '../Map';
import { MapMarkerData, MapViewState } from '../types/gis';
import 'mapbox-gl/dist/mapbox-gl.css';

interface InteractiveMapWithClusteringProps {
  /** Токен доступа к Mapbox */
  mapboxAccessToken?: string;
  /** Начальное состояние карты */
  initialViewState?: Partial<MapViewState>;
  /** Массив маркеров */
  markers: MapMarkerData[];
  /** Обработчик клика по маркеру */
  onMarkerClick?: (marker: MapMarkerData) => void;
  /** Обработчик клика по кластеру */
  onClusterClick?: (clusterId: number, coordinates: [number, number]) => void;
  /** Показывать ли цены на маркерах */
  showPrices?: boolean;
  /** Настройки кластеризации */
  clusterOptions?: {
    radius?: number;
    maxZoom?: number;
    minPoints?: number;
  };
  /** Стили кластеров */
  clusterStyles?: {
    small?: { color?: string; size?: number; textColor?: string };
    medium?: { color?: string; size?: number; textColor?: string };
    large?: { color?: string; size?: number; textColor?: string };
  };
}

/**
 * Пример интеграции InteractiveMap с нативной кластеризацией MapboxClusterLayer
 *
 * Этот компонент демонстрирует как можно использовать новый MapboxClusterLayer
 * вместо серверной кластеризации для улучшения производительности.
 */
const InteractiveMapWithClustering: React.FC<
  InteractiveMapWithClusteringProps
> = ({
  mapboxAccessToken,
  initialViewState = {
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 10,
  },
  markers,
  onMarkerClick,
  onClusterClick,
  showPrices = false,
  clusterOptions = {
    radius: 50,
    maxZoom: 14,
    minPoints: 2,
  },
  clusterStyles,
}) => {
  // Состояние карты
  const [viewState, setViewState] = useState<MapViewState>({
    longitude: initialViewState.longitude || 20.4649,
    latitude: initialViewState.latitude || 44.8176,
    zoom: initialViewState.zoom || 10,
    pitch: initialViewState.pitch || 0,
    bearing: initialViewState.bearing || 0,
  });

  // Состояние выбранного маркера
  const [selectedMarker, setSelectedMarker] = useState<MapMarkerData | null>(
    null
  );

  // Получение токена
  const accessToken =
    mapboxAccessToken || process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  // Обработчики событий
  const handleViewStateChange = useCallback((newViewState: MapViewState) => {
    setViewState(newViewState);
  }, []);

  const handleMarkerClick = useCallback(
    (marker: MapMarkerData) => {
      setSelectedMarker(marker);
      if (onMarkerClick) {
        onMarkerClick(marker);
      }
    },
    [onMarkerClick]
  );

  const handleClusterClick = useCallback(
    (clusterId: number, coordinates: [number, number]) => {
      if (onClusterClick) {
        onClusterClick(clusterId, coordinates);
      } else {
        // Дефолтное поведение - зум к кластеру
        setViewState((prev) => ({
          ...prev,
          longitude: coordinates[0],
          latitude: coordinates[1],
          zoom: Math.min(prev.zoom + 3, 20),
        }));
      }
    },
    [onClusterClick]
  );

  const handleMapClick = useCallback(() => {
    setSelectedMarker(null);
  }, []);

  // Проверка токена
  if (!accessToken) {
    return (
      <div className="w-full h-96 flex items-center justify-center bg-gray-100 rounded-lg">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Mapbox Access Token Required
          </h3>
          <p className="text-sm text-gray-600">
            Пожалуйста, предоставьте токен доступа к Mapbox
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full h-96 relative">
      {/* Информация о выбранном маркере */}
      {selectedMarker && (
        <div className="absolute top-4 right-4 z-10 bg-white rounded-lg shadow-lg p-4 max-w-64">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-medium text-gray-900">
              {selectedMarker.title}
            </h4>
            <button
              onClick={() => setSelectedMarker(null)}
              className="text-gray-400 hover:text-gray-600 text-lg leading-none"
            >
              ×
            </button>
          </div>
          <p className="text-xs text-gray-600 mb-2">
            {selectedMarker.description}
          </p>
          <div className="text-xs text-gray-500 mb-2">
            Тип: {selectedMarker.type}
          </div>
          {selectedMarker.data?.price && (
            <div className="text-sm font-medium text-primary">
              {selectedMarker.data.price}€
            </div>
          )}
        </div>
      )}

      {/* Статистика */}
      <div className="absolute top-4 left-4 z-10 bg-white rounded-lg shadow-lg p-3">
        <div className="text-xs text-gray-600 mb-1">
          Маркеров: {markers.length}
        </div>
        <div className="text-xs text-gray-600">
          Зум: {viewState.zoom.toFixed(1)}
        </div>
      </div>

      {/* Карта */}
      <Map
        {...viewState}
        onMove={(evt) => handleViewStateChange(evt.viewState)}
        onClick={handleMapClick}
        mapboxAccessToken={accessToken}
        mapStyle="mapbox://styles/mapbox/streets-v12"
        style={{ width: '100%', height: '100%' }}
        attributionControl={false}
        logoPosition="bottom-right"
      >
        <MapboxClusterLayer
          markers={markers}
          clusterRadius={clusterOptions.radius}
          clusterMaxZoom={clusterOptions.maxZoom}
          clusterMinPoints={clusterOptions.minPoints}
          onClusterClick={handleClusterClick}
          onMarkerClick={handleMarkerClick}
          clusterStyles={clusterStyles}
        />
      </Map>

      {/* Кнопки управления */}
      <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-10">
        <div className="bg-white rounded-lg shadow-lg p-2 flex space-x-2">
          <button
            onClick={() =>
              setViewState((prev) => ({
                ...prev,
                zoom: Math.min(prev.zoom + 1, 20),
              }))
            }
            className="px-3 py-1 text-sm bg-primary text-white rounded hover:bg-primary-dark transition-colors"
          >
            +
          </button>
          <button
            onClick={() =>
              setViewState((prev) => ({
                ...prev,
                zoom: Math.max(prev.zoom - 1, 0),
              }))
            }
            className="px-3 py-1 text-sm bg-primary text-white rounded hover:bg-primary-dark transition-colors"
          >
            -
          </button>
          <button
            onClick={() =>
              setViewState({
                longitude: 20.4649,
                latitude: 44.8176,
                zoom: 10,
                pitch: 0,
                bearing: 0,
              })
            }
            className="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
          >
            Сброс
          </button>
        </div>
      </div>
    </div>
  );
};

export default InteractiveMapWithClustering;
