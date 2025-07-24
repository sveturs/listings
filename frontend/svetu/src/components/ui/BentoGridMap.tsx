'use client';

import React, { useMemo } from 'react';
import Map, { Marker } from 'react-map-gl';
import type { ViewState } from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

interface BentoGridMapProps {
  listings?: Array<{
    id: string;
    latitude: number;
    longitude: number;
    price: number;
  }>;
  userLocation?: {
    latitude: number;
    longitude: number;
  };
}

export const BentoGridMap: React.FC<BentoGridMapProps> = ({
  listings = [],
  userLocation,
}) => {
  // Состояние для отслеживания доступности геолокации
  const [isGeolocationAvailable, setIsGeolocationAvailable] =
    React.useState(false);

  React.useEffect(() => {
    if ('geolocation' in navigator) {
      setIsGeolocationAvailable(true);
    }
  }, []);
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

    // Если нет локации пользователя, центрируем на первом объявлении
    // или на дефолтных координатах Белграда
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

  return (
    <div className="w-full h-full relative overflow-hidden rounded-lg">
      <Map
        initialViewState={initialViewState}
        style={{ width: '100%', height: '100%' }}
        mapStyle="mapbox://styles/mapbox/light-v11"
        mapboxAccessToken={process.env.NEXT_PUBLIC_MAPBOX_TOKEN}
        interactive={false} // Отключаем интерактивность для BentoGrid
        attributionControl={false}
      >
        {/* Маркер пользователя */}
        {userLocation && (
          <Marker
            longitude={userLocation.longitude}
            latitude={userLocation.latitude}
            anchor="bottom"
          >
            <div className="relative">
              <div className="absolute inset-0 bg-primary/30 rounded-full animate-ping" />
              <div className="w-4 h-4 bg-primary rounded-full border-2 border-white shadow-lg" />
            </div>
          </Marker>
        )}

        {/* Маркеры объявлений */}
        {listings.map((listing) => (
          <Marker
            key={listing.id}
            longitude={listing.longitude}
            latitude={listing.latitude}
            anchor="bottom"
          >
            <div className="relative group cursor-pointer">
              <div className="absolute inset-0 bg-secondary/20 rounded-full scale-150 opacity-0 group-hover:opacity-100 transition-opacity" />
              <div className="bg-white rounded-full px-2 py-1 text-xs font-semibold text-secondary shadow-md border border-secondary/20 group-hover:scale-110 transition-transform">
                €{listing.price}
              </div>
            </div>
          </Marker>
        ))}
      </Map>

      {/* Градиентная маска для эстетики */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute inset-x-0 top-0 h-8 bg-gradient-to-b from-white/20 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-8 bg-gradient-to-t from-white/20 to-transparent" />
      </div>

      {/* Кнопка геолокации */}
      {isGeolocationAvailable && !userLocation && (
        <button
          onClick={() => {
            if (navigator.geolocation) {
              navigator.geolocation.getCurrentPosition(
                (position) => {
                  console.log('User location:', position.coords);
                  // В будущем: передать координаты в родительский компонент
                },
                (error) => {
                  console.error('Geolocation error:', error);
                }
              );
            }
          }}
          className="absolute bottom-4 right-4 btn btn-sm btn-circle btn-primary shadow-lg hover:scale-110 transition-transform"
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
    </div>
  );
};
