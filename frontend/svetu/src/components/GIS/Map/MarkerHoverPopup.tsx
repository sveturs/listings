'use client';

import React from 'react';
import { Popup } from 'react-map-gl';
import { useTranslations } from 'next-intl';
import type { MapMarkerData } from '../types/gis';
// import { getCategoryIcon } from '../../../utils/categoryIcons';
import SafeImage from '../../SafeImage';

interface MarkerHoverPopupProps {
  marker: MapMarkerData;
  onClose: () => void;
  onClick?: () => void;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
}

const MarkerHoverPopup: React.FC<MarkerHoverPopupProps> = ({
  marker,
  onClose,
  onClick,
  onMouseEnter,
  onMouseLeave,
}) => {
  const t = useTranslations('map');

  // Парсим JSON данные
  const parsedData = React.useMemo(() => {
    try {
      return typeof marker.data === 'string'
        ? JSON.parse(marker.data)
        : marker.data;
    } catch {
      return marker.data || {};
    }
  }, [marker.data]);

  const parsedMetadata = React.useMemo(() => {
    try {
      return typeof marker.metadata === 'string'
        ? JSON.parse(marker.metadata)
        : marker.metadata;
    } catch {
      return marker.metadata || {};
    }
  }, [marker.metadata]);

  // Получаем первое изображение
  const getFirstImage = () => {
    // Сначала проверяем marker.imageUrl
    if (marker.imageUrl) return marker.imageUrl;

    // Затем ищем в parsedData.image (может быть строкой или массивом)
    if (parsedData?.image) {
      if (Array.isArray(parsedData.image)) {
        return parsedData.image[0];
      }
      return parsedData.image;
    }

    // Ищем в parsedData.images (массив)
    if (
      parsedData?.images &&
      Array.isArray(parsedData.images) &&
      parsedData.images.length > 0
    ) {
      return parsedData.images[0];
    }

    return null;
  };

  const firstImage = getFirstImage();

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
      <div
        className="relative overflow-hidden cursor-pointer"
        onClick={onClick}
        onMouseEnter={onMouseEnter}
        onMouseLeave={onMouseLeave}
      >
        {/* Фоновое изображение с blur эффектом */}
        {firstImage && (
          <div className="absolute inset-0 -z-10">
            <SafeImage
              src={firstImage}
              alt=""
              fill
              className="object-cover scale-110 blur-xl opacity-30"
              sizes="400px"
            />
            <div className="absolute inset-0 bg-gradient-to-b from-white/70 via-white/85 to-white"></div>
          </div>
        )}

        {/* Основной контент */}
        <div className="relative">
          {/* Изображение товара */}
          {firstImage && (
            <div className="relative h-40 -m-3 mb-3 overflow-hidden">
              <SafeImage
                src={firstImage}
                alt={marker.title}
                fill
                className="object-cover"
                sizes="400px"
              />

              {/* Оверлей с градиентом */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent"></div>

              {/* Категория в углу */}
              {(parsedMetadata?.category || parsedData?.category) && (
                <div className="absolute top-3 left-3">
                  <div className="flex items-center gap-1.5 px-3 py-1.5 bg-white/95 backdrop-blur-md rounded-full shadow-lg">
                    <span className="text-base">
                      {parsedData?.icon ||
                        getCategoryIcon(
                          parsedMetadata?.category || parsedData?.category
                        )}
                    </span>
                    <span className="text-xs font-medium text-gray-800">
                      {parsedMetadata?.category || parsedData?.category}
                    </span>
                  </div>
                </div>
              )}

              {/* Цена внизу изображения */}
              {(parsedData?.price || parsedMetadata?.price) && (
                <div className="absolute bottom-3 left-3 right-3">
                  <div className="flex items-baseline gap-1.5">
                    <span className="text-2xl font-bold text-white drop-shadow-lg">
                      {formatPrice(parsedData?.price || parsedMetadata?.price)}
                    </span>
                    <span className="text-sm text-white/90 font-medium">
                      {parsedMetadata?.currency || 'RSD'}
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
              {(parsedData?.address || marker.data?.address) && (
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
                    {parsedData?.address || marker.data?.address}
                  </span>
                </div>
              )}

              {/* Быстрые характеристики, если есть */}
              <div className="flex items-center gap-3">
                {/* Просмотры */}
                {parsedData?.views_count !== undefined &&
                  parsedData?.views_count > 0 && (
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
                      <span className="text-xs text-gray-500">
                        {parsedData.views_count} {t('views')}
                      </span>
                    </div>
                  )}

                {/* Рейтинг */}
                {parsedData?.rating !== undefined && parsedData?.rating > 0 && (
                  <div className="flex items-center gap-1">
                    <svg
                      className="w-4 h-4 text-yellow-400"
                      fill="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                    </svg>
                    <span className="text-xs text-gray-500">
                      {Number(parsedData.rating).toFixed(1)}
                    </span>
                  </div>
                )}

                {/* Дата публикации */}
                {parsedData?.created_at && (
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
                    <span className="text-xs text-gray-500">
                      {new Date(parsedData.created_at).toLocaleDateString(
                        'sr-RS'
                      )}
                    </span>
                  </div>
                )}
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
