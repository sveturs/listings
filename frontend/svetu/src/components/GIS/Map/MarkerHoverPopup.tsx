'use client';

import React from 'react';
import { Popup } from 'react-map-gl';
import type { MapMarkerData } from '../types/gis';

interface MarkerHoverPopupProps {
  marker: MapMarkerData;
  onClose: () => void;
}

const MarkerHoverPopup: React.FC<MarkerHoverPopupProps> = ({
  marker,
  onClose,
}) => {
  // Форматирование цены
  const formatPrice = (price?: number) => {
    if (!price) return '';
    return new Intl.NumberFormat('ru-RU').format(price);
  };

  // Получение иконки категории
  const getCategoryIcon = (category?: string) => {
    const categoryMap: { [key: string]: string } = {
      Недвижимость: '🏠',
      Транспорт: '🚗',
      Электроника: '💻',
      Одежда: '👕',
      Услуги: '🛠️',
      Работа: '💼',
    };
    return categoryMap[category || ''] || '📦';
  };

  return (
    <Popup
      longitude={marker.longitude}
      latitude={marker.latitude}
      closeButton={false}
      closeOnClick={false}
      onClose={onClose}
      anchor="bottom"
      offset={25}
      className="marker-hover-popup"
      maxWidth="320px"
    >
      <div className="relative overflow-hidden">
        {/* Фоновое изображение с blur эффектом */}
        {marker.imageUrl && (
          <div className="absolute inset-0 -z-10">
            <img
              src={marker.imageUrl}
              alt=""
              className="w-full h-full object-cover scale-110 blur-xl opacity-30"
            />
            <div className="absolute inset-0 bg-gradient-to-b from-white/70 via-white/85 to-white"></div>
          </div>
        )}

        {/* Основной контент */}
        <div className="relative">
          {/* Изображение товара */}
          {marker.imageUrl && (
            <div className="relative h-40 -m-3 mb-3 overflow-hidden">
              <img
                src={marker.imageUrl}
                alt={marker.title}
                className="w-full h-full object-cover"
                loading="lazy"
              />

              {/* Оверлей с градиентом */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent"></div>

              {/* Категория в углу */}
              {marker.metadata?.category && (
                <div className="absolute top-3 left-3">
                  <div className="flex items-center gap-1.5 px-3 py-1.5 bg-white/95 backdrop-blur-md rounded-full shadow-lg">
                    <span className="text-base">
                      {getCategoryIcon(marker.metadata.category)}
                    </span>
                    <span className="text-xs font-medium text-gray-800">
                      {marker.metadata.category}
                    </span>
                  </div>
                </div>
              )}

              {/* Цена внизу изображения */}
              {marker.metadata?.price && (
                <div className="absolute bottom-3 left-3 right-3">
                  <div className="flex items-baseline gap-1.5">
                    <span className="text-2xl font-bold text-white drop-shadow-lg">
                      {formatPrice(marker.metadata.price)}
                    </span>
                    <span className="text-sm text-white/90 font-medium">
                      {marker.metadata.currency || 'RSD'}
                    </span>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Информация о товаре */}
          <div className="px-4 pb-3 space-y-3">
            {/* Название */}
            <div>
              <h3 className="font-bold text-base text-gray-900 line-clamp-2 leading-tight">
                {marker.title}
              </h3>
            </div>

            {/* Дополнительная информация */}
            <div className="space-y-2">
              {/* Адрес */}
              {marker.data?.address && (
                <div className="flex items-start gap-2">
                  <svg
                    className="w-4 h-4 text-gray-400 mt-0.5 flex-shrink-0"
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
                  <span className="text-sm text-gray-600 line-clamp-1">
                    {marker.data.address}
                  </span>
                </div>
              )}

              {/* Быстрые характеристики, если есть */}
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-1">
                  <svg
                    className="w-4 h-4 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                  <span className="text-xs text-gray-500">234 просмотра</span>
                </div>
                <div className="flex items-center gap-1">
                  <svg
                    className="w-4 h-4 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span className="text-xs text-gray-500">2 часа назад</span>
                </div>
              </div>
            </div>

            {/* Разделитель */}
            <div className="border-t border-gray-100 pt-3">
              {/* CTA с анимацией */}
              <div className="flex items-center justify-between">
                <p className="text-sm text-gray-600 font-medium">
                  Нажмите для подробностей
                </p>
                <div className="flex items-center gap-1 text-primary">
                  <span className="text-sm font-medium">Открыть</span>
                  <svg
                    className="w-4 h-4 animate-pulse"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Popup>
  );
};

export default MarkerHoverPopup;
