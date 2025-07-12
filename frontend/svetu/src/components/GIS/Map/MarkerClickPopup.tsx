'use client';

import React from 'react';
import { Popup } from 'react-map-gl';
import type { MapMarkerData } from '../types/gis';
import { useRouter } from 'next/navigation';

interface MarkerClickPopupProps {
  marker: MapMarkerData;
  onClose: () => void;
}

const MarkerClickPopup: React.FC<MarkerClickPopupProps> = ({
  marker,
  onClose,
}) => {
  const router = useRouter();

  // Форматирование цены
  const formatPrice = (price?: number) => {
    if (!price) return '';
    return new Intl.NumberFormat('ru-RU').format(price);
  };

  const handleViewDetails = () => {
    if (marker.data?.id) {
      router.push(`/marketplace/${marker.data.id}`);
    }
  };

  return (
    <Popup
      longitude={marker.longitude}
      latitude={marker.latitude}
      closeButton={true}
      closeOnClick={true}
      onClose={onClose}
      anchor="bottom"
      offset={35}
      className="marker-click-popup"
      maxWidth="400px"
    >
      <div className="p-0">
        {/* Изображение с галереей */}
        {marker.imageUrl && (
          <div className="relative h-48 -m-3 mb-3">
            <img
              src={marker.imageUrl}
              alt={marker.title}
              className="w-full h-full object-cover"
              loading="lazy"
            />

            {/* Категория */}
            {marker.metadata?.category && (
              <div className="absolute top-3 left-3">
                <span className="text-sm px-3 py-1.5 bg-white/90 text-gray-800 rounded-full backdrop-blur-sm font-medium">
                  {marker.metadata.category}
                </span>
              </div>
            )}

            {/* Индикатор галереи */}
            <div className="absolute bottom-3 right-3 bg-black/60 text-white px-2 py-1 rounded-full text-xs backdrop-blur-sm">
              📷 Больше фото
            </div>
          </div>
        )}

        {/* Контент */}
        <div className="space-y-3 p-4">
          {/* Название */}
          <h3 className="font-bold text-lg text-gray-900 line-clamp-2">
            {marker.title}
          </h3>

          {/* Цена */}
          {marker.metadata?.price && (
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-bold text-primary">
                {formatPrice(marker.metadata.price)}
              </span>
              <span className="text-base text-gray-500">
                {marker.metadata.currency || 'RSD'}
              </span>
            </div>
          )}

          {/* Описание */}
          {marker.description && (
            <p className="text-sm text-gray-600 line-clamp-3">
              {marker.description}
            </p>
          )}

          {/* Адрес */}
          {marker.data?.address && (
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
              <span>{marker.data.address}</span>
            </div>
          )}

          {/* Дополнительная информация */}
          <div className="grid grid-cols-2 gap-3 pt-2 border-t border-gray-100">
            <div className="text-center">
              <p className="text-xs text-gray-500">Просмотров</p>
              <p className="font-semibold text-gray-900">1,234</p>
            </div>
            <div className="text-center">
              <p className="text-xs text-gray-500">Добавлено</p>
              <p className="font-semibold text-gray-900">2 дня назад</p>
            </div>
          </div>

          {/* Действия */}
          <div className="flex gap-2 pt-3">
            <button
              onClick={handleViewDetails}
              className="flex-1 btn btn-primary btn-sm"
            >
              Подробнее
            </button>
            <button className="btn btn-outline btn-sm">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                />
              </svg>
            </button>
            <button className="btn btn-outline btn-sm">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m9.032 4.026a9.001 9.001 0 010-5.368m0 5.368a9.003 9.003 0 01-7.432 3.268 9.003 9.003 0 01-7.432-3.268m14.864 0a9.003 9.003 0 00-7.432-3.268 9.003 9.003 0 00-7.432 3.268"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </Popup>
  );
};

export default MarkerClickPopup;
