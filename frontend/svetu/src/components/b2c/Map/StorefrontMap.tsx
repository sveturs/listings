'use client';

import { useEffect, useState, useCallback } from 'react';
import Image from 'next/image';
import dynamic from 'next/dynamic';
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  useMap,
  useMapEvents,
} from 'react-leaflet';
import L from 'leaflet';
import type { B2CStoreMapData } from '@/types/b2c';

// Исправляем иконки маркеров для Leaflet
import 'leaflet/dist/leaflet.css';

// Кастомные иконки
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: '/icons/marker-icon-2x.png',
  iconUrl: '/icons/marker-icon.png',
  shadowUrl: '/icons/marker-shadow.png',
});

// Создаем кастомную иконку для витрин
const storefrontIcon = new L.Icon({
  iconUrl: '/icons/storefront-marker.png',
  iconRetinaUrl: '/icons/storefront-marker-2x.png',
  shadowUrl: '/icons/marker-shadow.png',
  iconSize: [32, 32],
  iconAnchor: [16, 32],
  popupAnchor: [0, -32],
  shadowSize: [41, 41],
  shadowAnchor: [13, 41],
});

interface StorefrontMapProps {
  storefronts: B2CStoreMapData[];
  center?: { lat: number; lng: number };
  zoom?: number;
  height?: string;
  onStorefrontClick?: (storefront: B2CStoreMapData) => void;
  onBoundsChange?: (bounds: L.LatLngBounds) => void;
  className?: string;
  showSearch?: boolean;
  clustering?: boolean;
}

// Компонент для обработки событий карты
const MapEventHandler = ({
  onBoundsChange,
}: {
  onBoundsChange?: (bounds: L.LatLngBounds) => void;
}) => {
  const map = useMap();

  useMapEvents({
    moveend: () => {
      if (onBoundsChange) {
        onBoundsChange(map.getBounds());
      }
    },
    zoomend: () => {
      if (onBoundsChange) {
        onBoundsChange(map.getBounds());
      }
    },
  });

  return null;
};

// Компонент для автоматического позиционирования карты
const MapController = ({
  storefronts,
  center,
}: {
  storefronts: B2CStoreMapData[];
  center?: { lat: number; lng: number };
}) => {
  const map = useMap();

  useEffect(() => {
    if (storefronts.length > 0 && !center) {
      // Автоматически позиционируем карту по витринам
      const group = new L.FeatureGroup(
        storefronts
          .filter(
            (storefront) =>
              typeof storefront.latitude === 'number' &&
              typeof storefront.longitude === 'number'
          )
          .map((storefront) =>
            L.marker([
              storefront.latitude as number,
              storefront.longitude as number,
            ])
          )
      );

      if (group.getBounds().isValid()) {
        map.fitBounds(group.getBounds(), { padding: [20, 20] });
      }
    } else if (center) {
      map.setView([center.lat, center.lng], map.getZoom());
    }
  }, [map, storefronts, center]);

  return null;
};

const StorefrontMap: React.FC<StorefrontMapProps> = ({
  storefronts = [],
  center = { lat: 44.7866, lng: 20.4489 }, // Белград по умолчанию
  zoom = 12,
  height = '400px',
  onStorefrontClick,
  onBoundsChange,
  className = '',
  showSearch = false,
  // clustering = false,
}) => {
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  const handleMarkerClick = useCallback(
    (storefront: B2CStoreMapData) => {
      if (onStorefrontClick) {
        onStorefrontClick(storefront);
      }
    },
    [onStorefrontClick]
  );

  if (!isClient) {
    return (
      <div
        className={`bg-base-200 animate-pulse rounded-lg ${className}`}
        style={{ height }}
      >
        <div className="flex items-center justify-center h-full">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative ${className}`} style={{ height }}>
      <MapContainer
        center={[center.lat, center.lng]}
        zoom={zoom}
        style={{ height: '100%', width: '100%' }}
        className="rounded-lg z-0"
        scrollWheelZoom={true}
        zoomControl={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        {/* Обработчики событий карты */}
        <MapEventHandler onBoundsChange={onBoundsChange} />
        <MapController storefronts={storefronts} center={center} />

        {/* Маркеры витрин */}
        {storefronts
          .filter(
            (storefront) =>
              typeof storefront.latitude === 'number' &&
              typeof storefront.longitude === 'number'
          )
          .map((storefront) => (
            <Marker
              key={storefront.id}
              position={[
                storefront.latitude as number,
                storefront.longitude as number,
              ]}
              icon={storefrontIcon}
              eventHandlers={{
                click: () => handleMarkerClick(storefront),
              }}
            >
              <Popup>
                <div className="p-2 min-w-[200px]">
                  <h3 className="font-bold text-lg mb-2">{storefront.name}</h3>

                  {storefront.logo_url && (
                    <Image
                      src={storefront.logo_url || ''}
                      alt={storefront.name || 'Storefront logo'}
                      width={64}
                      height={64}
                      className="w-16 h-16 object-cover rounded mb-2"
                    />
                  )}

                  <p className="text-sm text-gray-600 mb-2">
                    {storefront.address}
                  </p>

                  {storefront.phone && (
                    <p className="text-sm mb-2">📞 {storefront.phone}</p>
                  )}

                  <div className="flex items-center gap-2 mb-2">
                    {storefront.rating && (
                      <div className="badge badge-primary">
                        ⭐ {storefront.rating.toFixed(1)}
                      </div>
                    )}

                    {storefront.accepts_cards && (
                      <div className="badge badge-success">💳 Карты</div>
                    )}

                    {storefront.working_now && (
                      <div className="badge badge-info">🕒 Открыто</div>
                    )}

                    {storefront.has_delivery && (
                      <div className="badge badge-accent">🚚 Доставка</div>
                    )}
                  </div>

                  <button
                    className="btn btn-primary btn-sm w-full"
                    onClick={() => handleMarkerClick(storefront)}
                  >
                    Открыть витрину
                  </button>
                </div>
              </Popup>
            </Marker>
          ))}
      </MapContainer>

      {/* Поиск на карте (если включен) */}
      {showSearch && (
        <div className="absolute top-4 left-4 z-[1000]">
          <div className="bg-white rounded-lg shadow-lg p-3">
            <input
              type="text"
              placeholder="Поиск витрин..."
              className="input input-bordered input-sm w-64"
            />
          </div>
        </div>
      )}

      {/* Счетчик витрин */}
      <div className="absolute bottom-4 right-4 z-[1000]">
        <div className="bg-white rounded-lg shadow-lg px-3 py-2">
          <span className="text-sm font-medium">
            Найдено: {storefronts.length} витрин
          </span>
        </div>
      </div>
    </div>
  );
};

// Экспортируем компонент с динамической загрузкой для SSR
export default dynamic(() => Promise.resolve(StorefrontMap), {
  ssr: false,
  loading: () => (
    <div
      className="bg-base-200 animate-pulse rounded-lg"
      style={{ height: '400px' }}
    >
      <div className="flex items-center justify-center h-full">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    </div>
  ),
});
