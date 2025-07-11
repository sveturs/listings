import React, { useState, useEffect } from 'react';
import { InteractiveMap } from '../index';
import { MapMarkerData, MapPopupData, MapViewState } from '../types/gis';
import { ListingPopup } from '../Map/MapPopup';

const MapDemo: React.FC = () => {
  const [viewState, setViewState] = useState<MapViewState>({
    longitude: 20.4649,
    latitude: 44.8176,
    zoom: 12,
  });

  const [markers, setMarkers] = useState<MapMarkerData[]>([]);
  const [popup, setPopup] = useState<MapPopupData | null>(null);

  // Демонстрационные данные маркеров
  useEffect(() => {
    const demoMarkers: MapMarkerData[] = [
      {
        id: '1',
        position: [20.4649, 44.8176],
        title: 'Хостел "Центр"',
        description: 'Уютный хостел в центре города',
        type: 'listing',
        data: {
          price: 25,
          currency: '€',
          rating: 4.5,
          reviews: 123,
          imageUrl: '/api/placeholder/200/150',
        },
      },
      {
        id: '2',
        position: [20.47, 44.82],
        title: 'Гостиница "Комфорт"',
        description: 'Современные номера с видом на реку',
        type: 'listing',
        data: {
          price: 45,
          currency: '€',
          rating: 4.2,
          reviews: 87,
          imageUrl: '/api/placeholder/200/150',
        },
      },
      {
        id: '3',
        position: [20.46, 44.815],
        title: 'Пользователь Марко',
        description: 'Активный путешественник',
        type: 'user',
        data: {
          avatar: '/api/placeholder/40/40',
          status: 'online',
        },
      },
      {
        id: '4',
        position: [20.475, 44.825],
        title: 'Крепость Калемегдан',
        description: 'Историческая достопримечательность',
        type: 'poi',
        data: {
          category: 'Достопримечательность',
          rating: 4.8,
        },
      },
    ];

    setMarkers(demoMarkers);
  }, []);

  const handleMarkerClick = (marker: MapMarkerData) => {
    if (marker.type === 'listing') {
      // Для объявлений показываем специальный popup
      setPopup({
        id: marker.id,
        position: marker.position,
        title: marker.title,
        description: marker.description,
        content: (
          <ListingPopup
            listing={{
              id: marker.id,
              title: marker.title,
              price: marker.data?.price || 0,
              currency: marker.data?.currency || '€',
              rating: marker.data?.rating,
              reviews: marker.data?.reviews,
              imageUrl: marker.data?.imageUrl,
              distance: 0.5,
              category: 'Размещение',
            }}
            position={marker.position}
            onClose={() => setPopup(null)}
            onViewDetails={(id) => {
              console.log('View details for:', id);
              setPopup(null);
            }}
          />
        ),
      });
    } else {
      // Для остальных типов показываем обычный popup
      setPopup({
        id: marker.id,
        position: marker.position,
        title: marker.title,
        description: marker.description,
        content: (
          <div className="space-y-3">
            {marker.type === 'user' && (
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                  👤
                </div>
                <div>
                  <p className="font-medium">{marker.title}</p>
                  <p className="text-sm text-green-600">В сети</p>
                </div>
              </div>
            )}

            {marker.type === 'poi' && (
              <div>
                <div className="flex items-center space-x-2 mb-2">
                  <span className="text-xl">🏛️</span>
                  <span className="text-sm bg-blue-100 text-blue-800 px-2 py-1 rounded">
                    {marker.data?.category}
                  </span>
                </div>
                <div className="flex items-center space-x-1">
                  <span className="text-yellow-400">★</span>
                  <span className="text-sm">{marker.data?.rating}</span>
                </div>
              </div>
            )}

            <button
              onClick={() => setPopup(null)}
              className="w-full bg-primary text-white py-2 px-4 rounded-md hover:bg-primary-dark transition-colors"
            >
              Подробнее
            </button>
          </div>
        ),
      });
    }
  };

  const handleMapClick = () => {
    setPopup(null);
  };

  return (
    <div className="w-full h-screen">
      <div className="absolute top-4 left-4 z-20 bg-white p-4 rounded-lg shadow-lg">
        <h2 className="text-lg font-semibold mb-2">Демонстрация карты</h2>
        <p className="text-sm text-gray-600 mb-3">
          Интерактивная карта с маркерами различных типов
        </p>
        <div className="space-y-2 text-sm">
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-blue-500 rounded-full"></div>
            <span>
              Размещение ({markers.filter((m) => m.type === 'listing').length})
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-orange-500 rounded-full"></div>
            <span>
              Пользователи ({markers.filter((m) => m.type === 'user').length})
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-4 h-4 bg-red-500 rounded-full"></div>
            <span>
              Достопримечательности (
              {markers.filter((m) => m.type === 'poi').length})
            </span>
          </div>
        </div>
      </div>

      <InteractiveMap
        initialViewState={viewState}
        markers={markers}
        popup={popup}
        onMarkerClick={handleMarkerClick}
        onMapClick={handleMapClick}
        onViewStateChange={setViewState}
        controlsConfig={{
          showZoom: true,
          showCompass: true,
          showFullscreen: true,
          showGeolocate: true,
          showNavigation: true,
          position: 'top-right',
        }}
        className="w-full h-full"
      />
    </div>
  );
};

export default MapDemo;
