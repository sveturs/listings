'use client';

import React, { useState, useCallback, useMemo } from 'react';
import Map from 'react-map-gl';
import { MapboxClusterLayer } from '../Map';
import { MapMarkerData, MapViewState } from '../types/gis';
import 'mapbox-gl/dist/mapbox-gl.css';

interface MapboxClusterDemoProps {
  /** Токен доступа к Mapbox */
  mapboxAccessToken?: string;
  /** Начальное состояние карты */
  initialViewState?: Partial<MapViewState>;
  /** Показывать ли демонстрационные данные */
  showDemoData?: boolean;
}

const MapboxClusterDemo: React.FC<MapboxClusterDemoProps> = ({
  mapboxAccessToken,
  initialViewState = {
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 10,
  },
  showDemoData = true,
}) => {
  // Состояние карты
  const [viewState, setViewState] = useState<MapViewState>({
    longitude: initialViewState.longitude || 20.4649,
    latitude: initialViewState.latitude || 44.8176,
    zoom: initialViewState.zoom || 10,
    pitch: initialViewState.pitch || 0,
    bearing: initialViewState.bearing || 0,
  });

  // Состояние демонстрационных данных
  const [selectedMarker, setSelectedMarker] = useState<MapMarkerData | null>(
    null
  );
  const [showPrices, setShowPrices] = useState(false);

  // Получение токена
  const accessToken =
    mapboxAccessToken || process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  // Генерация демонстрационных данных
  const demoMarkers = useMemo(() => {
    if (!showDemoData) return [];

    const markers: MapMarkerData[] = [];
    const centerLat = 44.8176;
    const centerLng = 20.4649;
    const radius = 0.1; // Радиус для генерации точек

    // Генерируем случайные объявления
    for (let i = 0; i < 50; i++) {
      const angle = Math.random() * 2 * Math.PI;
      const distance = Math.random() * radius;
      const lat = centerLat + distance * Math.cos(angle);
      const lng = centerLng + distance * Math.sin(angle);

      markers.push({
        id: `listing-${i}`,
        position: [lng, lat],
        title: `Объявление ${i + 1}`,
        description: `Описание объявления ${i + 1}`,
        type: 'listing',
        data: {
          price: Math.floor(Math.random() * 500) + 50,
          currency: 'EUR',
          bedrooms: Math.floor(Math.random() * 5) + 1,
          bathrooms: Math.floor(Math.random() * 3) + 1,
        },
      });
    }

    // Добавляем несколько пользователей
    for (let i = 0; i < 5; i++) {
      const angle = Math.random() * 2 * Math.PI;
      const distance = Math.random() * radius;
      const lat = centerLat + distance * Math.cos(angle);
      const lng = centerLng + distance * Math.sin(angle);

      markers.push({
        id: `user-${i}`,
        position: [lng, lat],
        title: `Пользователь ${i + 1}`,
        description: `Пользователь ${i + 1}`,
        type: 'user',
        data: {
          username: `user${i + 1}`,
          avatar: `https://i.pravatar.cc/150?img=${i + 1}`,
        },
      });
    }

    // Добавляем несколько точек интереса
    for (let i = 0; i < 10; i++) {
      const angle = Math.random() * 2 * Math.PI;
      const distance = Math.random() * radius;
      const lat = centerLat + distance * Math.cos(angle);
      const lng = centerLng + distance * Math.sin(angle);

      markers.push({
        id: `poi-${i}`,
        position: [lng, lat],
        title: `Точка интереса ${i + 1}`,
        description: `Описание точки интереса ${i + 1}`,
        type: 'poi',
        data: {
          category: ['restaurant', 'shop', 'hospital', 'school', 'park'][
            Math.floor(Math.random() * 5)
          ],
          rating: Math.floor(Math.random() * 5) + 1,
        },
      });
    }

    return markers;
  }, [showDemoData]);

  // Обработчики событий
  const handleClusterClick = useCallback(
    (clusterId: number, coordinates: [number, number]) => {
      console.log('Cluster clicked:', clusterId, coordinates);
      // Увеличиваем зум к кластеру
      setViewState((prev) => ({
        ...prev,
        longitude: coordinates[0],
        latitude: coordinates[1],
        zoom: Math.min(prev.zoom + 2, 20),
      }));
    },
    []
  );

  const handleMarkerClick = useCallback((marker: MapMarkerData) => {
    console.log('Marker clicked:', marker);
    setSelectedMarker(marker);
  }, []);

  const handleViewStateChange = useCallback((newViewState: MapViewState) => {
    setViewState(newViewState);
  }, []);

  // Настройки стилей кластеров
  const clusterStyles = useMemo(
    () => ({
      small: {
        color: '#3b82f6',
        size: 40,
        textColor: '#ffffff',
      },
      medium: {
        color: '#059669',
        size: 55,
        textColor: '#ffffff',
      },
      large: {
        color: '#dc2626',
        size: 70,
        textColor: '#ffffff',
      },
    }),
    []
  );

  if (!accessToken) {
    return (
      <div className="w-full h-96 flex items-center justify-center bg-gray-100 rounded-lg">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Mapbox Access Token Required
          </h3>
          <p className="text-sm text-gray-600">
            Пожалуйста, предоставьте токен доступа к Mapbox для использования
            демонстрации
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full h-96 relative">
      {/* Панель управления */}
      <div className="absolute top-4 left-4 z-10 bg-white rounded-lg shadow-lg p-4 space-y-2">
        <h3 className="text-sm font-medium text-gray-900">
          Демонстрация кластеризации Mapbox
        </h3>
        <div className="flex items-center space-x-2">
          <input
            type="checkbox"
            id="show-prices"
            checked={showPrices}
            onChange={(e) => setShowPrices(e.target.checked)}
            className="rounded border-gray-300 text-primary focus:ring-primary"
          />
          <label htmlFor="show-prices" className="text-sm text-gray-700">
            Показывать цены
          </label>
        </div>
        <div className="text-xs text-gray-500">
          Маркеров: {demoMarkers.length}
        </div>
        <div className="text-xs text-gray-500">
          Зум: {viewState.zoom.toFixed(1)}
        </div>
      </div>

      {/* Информация о выбранном маркере */}
      {selectedMarker && (
        <div className="absolute top-4 right-4 z-10 bg-white rounded-lg shadow-lg p-4 max-w-64">
          <div className="flex items-center justify-between mb-2">
            <h4 className="text-sm font-medium text-gray-900">
              {selectedMarker.title}
            </h4>
            <button
              onClick={() => setSelectedMarker(null)}
              className="text-gray-400 hover:text-gray-600"
            >
              ×
            </button>
          </div>
          <p className="text-xs text-gray-600 mb-2">
            {selectedMarker.description}
          </p>
          <div className="text-xs text-gray-500">
            Тип: {selectedMarker.type}
          </div>
          {selectedMarker.data && (
            <div className="mt-2 text-xs text-gray-500">
              <pre className="whitespace-pre-wrap">
                {JSON.stringify(selectedMarker.data, null, 2)}
              </pre>
            </div>
          )}
        </div>
      )}

      {/* Карта */}
      <Map
        {...viewState}
        onMove={(evt) => handleViewStateChange(evt.viewState)}
        mapboxAccessToken={accessToken}
        mapStyle="mapbox://styles/mapbox/streets-v12"
        style={{ width: '100%', height: '100%' }}
        attributionControl={false}
        logoPosition="bottom-right"
      >
        <MapboxClusterLayer
          markers={demoMarkers}
          clusterRadius={50}
          clusterMaxZoom={14}
          clusterMinPoints={2}
          onClusterClick={handleClusterClick}
          onMarkerClick={handleMarkerClick}
          showPrices={showPrices}
          clusterStyles={clusterStyles}
        />
      </Map>

      {/* Легенда */}
      <div className="absolute bottom-4 left-4 z-10 bg-white rounded-lg shadow-lg p-3">
        <h4 className="text-xs font-medium text-gray-900 mb-2">Легенда</h4>
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-blue-500 rounded-full"></div>
            <span className="text-xs text-gray-600">Объявления</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-orange-500 rounded-full"></div>
            <span className="text-xs text-gray-600">Пользователи</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-red-500 rounded-full"></div>
            <span className="text-xs text-gray-600">Точки интереса</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-green-500 rounded-full"></div>
            <span className="text-xs text-gray-600">Кластеры</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MapboxClusterDemo;
