'use client';

import React from 'react';
import { Popup } from 'react-map-gl';
import type { MapMarkerData } from '../types/gis';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import Image from 'next/image';

interface MarkerClickPopupProps {
  marker: MapMarkerData;
  onClose: () => void;
}

const MarkerClickPopup: React.FC<MarkerClickPopupProps> = ({
  marker,
  onClose,
}) => {
  const router = useRouter();
  const t = useTranslations('marketplace');

  // Форматирование адреса с учетом приватности
  const formatAddressWithPrivacy = (
    address: string,
    privacyLevel?: string
  ): string => {
    if (!address) return '';

    if (privacyLevel === 'exact') {
      return address;
    }

    const parts = address.split(',').map((part) => part.trim());

    switch (privacyLevel) {
      case 'approximate':
      case 'street':
        // Убираем номер дома
        if (parts.length > 2) {
          const streetPart = parts[0]
            .replace(/\d+[а-яА-Яa-zA-Z]?(\s|$)/g, '')
            .trim();
          return streetPart
            ? [streetPart, ...parts.slice(1)].join(', ')
            : parts.slice(1).join(', ');
        }
        return parts.slice(1).join(', ');

      case 'district':
        // Оставляем только район и город
        if (parts.length > 2) {
          return parts.slice(-2).join(', ');
        }
        return address;

      case 'city_only':
      case 'city':
        // Оставляем только город
        if (parts.length > 1) {
          return parts[parts.length - 1];
        }
        return address;

      case 'hidden':
        // Скрываем адрес полностью
        return 'Адрес скрыт';

      default:
        return address;
    }
  };

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

  // Форматирование цены
  const formatPrice = (price?: number) => {
    if (!price) return '';
    return new Intl.NumberFormat('ru-RU').format(price);
  };

  const handleViewDetails = () => {
    if (parsedData?.id || marker.id) {
      router.push(`/marketplace/${parsedData?.id || marker.id}`);
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
            <Image
              src={marker.imageUrl}
              alt={marker.title}
              fill
              className="object-cover"
            />

            {/* Категория */}
            {(parsedMetadata?.category || parsedData?.category) && (
              <div className="absolute top-3 left-3">
                <span className="text-sm px-3 py-1.5 bg-white/90 text-gray-800 rounded-full backdrop-blur-sm font-medium">
                  {parsedMetadata?.category || parsedData?.category}
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
          <h3
            className="font-bold text-lg text-gray-900"
            style={{
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
            }}
          >
            {marker.title}
          </h3>

          {/* Цена */}
          {(parsedMetadata?.price || parsedData?.price) && (
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-bold text-primary">
                {formatPrice(parsedMetadata?.price || parsedData?.price)}
              </span>
              <span className="text-base text-gray-500">
                {parsedMetadata?.currency || parsedData?.currency || 'RSD'}
              </span>
            </div>
          )}

          {/* Описание */}
          {marker.description && (
            <p
              className="text-sm text-gray-600"
              style={{
                display: '-webkit-box',
                WebkitLineClamp: 3,
                WebkitBoxOrient: 'vertical',
                overflow: 'hidden',
              }}
            >
              {marker.description}
            </p>
          )}

          {/* Адрес */}
          {(parsedData?.address ||
            parsedData?.location ||
            parsedMetadata?.address) && (
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
              <span>
                {formatAddressWithPrivacy(
                  parsedData?.address ||
                    parsedData?.location ||
                    parsedMetadata?.address,
                  parsedData?.locationPrivacy
                )}
              </span>
            </div>
          )}

          {/* Атрибуты и характеристики */}
          {(parsedData?.attributes || parsedMetadata || parsedData) && (
            <div className="space-y-2">
              <h4 className="font-medium text-gray-900">
                {t('characteristics')}
              </h4>
              <div className="grid grid-cols-2 gap-2 text-sm">
                {/* Показываем атрибуты из data.attributes или основные поля */}
                {parsedData?.attributes
                  ? Object.entries(parsedData.attributes)
                      .slice(0, 4)
                      .map(([key, value]) => (
                        <div key={key} className="flex justify-between">
                          <span className="text-gray-600">{key}:</span>
                          <span className="font-medium">{String(value)}</span>
                        </div>
                      ))
                  : // Если нет attributes, показываем основные поля
                    [
                      (parsedMetadata?.category || parsedData?.category) && [
                        t('category'),
                        parsedMetadata?.category || parsedData?.category,
                      ],
                      parsedData?.area && [t('area'), `${parsedData.area} м²`],
                      parsedData?.rooms && [t('rooms'), parsedData.rooms],
                      parsedData?.condition && [
                        t('condition'),
                        parsedData.condition,
                      ],
                      parsedData?.type && [t('type'), parsedData.type],
                      parsedData?.year && [t('year'), parsedData.year],
                    ]
                      .filter(Boolean)
                      .slice(0, 4)
                      .map(([key, value], index) => (
                        <div key={index} className="flex justify-between">
                          <span className="text-gray-600">{key}:</span>
                          <span className="font-medium">{String(value)}</span>
                        </div>
                      ))}
                {/* Показываем "еще X" если есть больше атрибутов */}
                {parsedData?.attributes &&
                  Object.keys(parsedData.attributes).length > 4 && (
                    <div className="col-span-2 text-center text-xs text-gray-500">
                      {t('andMore', {
                        count: Object.keys(parsedData.attributes).length - 4,
                      })}
                    </div>
                  )}
              </div>
            </div>
          )}

          {/* Дополнительная информация */}
          <div className="grid grid-cols-3 gap-3 pt-2 border-t border-gray-100">
            <div className="text-center">
              <p className="text-xs text-gray-500">{t('views')}</p>
              <p className="font-semibold text-gray-900">
                {/* Ищем views_count в разных местах */}
                {parsedData?.views_count ||
                  parsedData?.views ||
                  parsedData?.view_count ||
                  parsedMetadata?.views_count ||
                  parsedMetadata?.views ||
                  '—'}
              </p>
            </div>
            <div className="text-center">
              <p className="text-xs text-gray-500">{t('rating')}</p>
              <p className="font-semibold text-gray-900">
                {/* Ищем rating в разных местах */}
                {parsedData?.rating || parsedMetadata?.rating
                  ? `⭐ ${Number(parsedData?.rating || parsedMetadata?.rating).toFixed(1)}`
                  : '—'}
              </p>
            </div>
            <div className="text-center">
              <p className="text-xs text-gray-500">{t('published')}</p>
              <p className="font-semibold text-gray-900">
                {/* Ищем created_at в разных местах */}
                {parsedData?.created_at ||
                parsedData?.createdAt ||
                parsedMetadata?.created_at
                  ? new Date(
                      parsedData?.created_at ||
                        parsedData?.createdAt ||
                        parsedMetadata?.created_at
                    ).toLocaleDateString('sr-RS')
                  : '—'}
              </p>
            </div>
          </div>

          {/* Информация о продавце */}
          {parsedData?.seller && (
            <div className="pt-2 border-t border-gray-100">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                  {parsedData.seller.avatar ? (
                    <Image
                      src={parsedData.seller.avatar}
                      alt={parsedData.seller.name}
                      fill
                      className="rounded-full object-cover"
                    />
                  ) : (
                    <span className="text-gray-600 text-lg">👤</span>
                  )}
                </div>
                <div className="flex-1">
                  <p className="font-medium text-gray-900">
                    {parsedData.seller.name}
                  </p>
                  <p className="text-xs text-gray-500">
                    {t('seller')} •{' '}
                    {parsedData.seller.rating
                      ? `⭐ ${parsedData.seller.rating}`
                      : t('newSeller')}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Действия */}
          <div className="space-y-2 pt-3">
            {/* Основные кнопки */}
            <div className="flex gap-2">
              <button
                onClick={handleViewDetails}
                className="flex-1 btn btn-primary btn-sm"
              >
                {t('viewDetails')}
              </button>
              <button
                onClick={() => {
                  // TODO: Открыть чат с продавцом
                  console.log('Open chat with seller');
                }}
                className="flex-1 btn btn-outline btn-sm"
              >
                💬 {t('contactSeller')}
              </button>
            </div>

            {/* Дополнительные действия */}
            <div className="flex gap-2">
              <button
                className="btn btn-outline btn-sm flex-1"
                onClick={() => {
                  // TODO: Добавить в избранное
                  console.log('Add to favorites');
                }}
              >
                <svg
                  className="w-4 h-4 mr-1"
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
                {t('addToFavorites')}
              </button>
              <button
                className="btn btn-outline btn-sm flex-1"
                onClick={() => {
                  // TODO: Поделиться
                  console.log('Share listing');
                }}
              >
                <svg
                  className="w-4 h-4 mr-1"
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
                {t('share')}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Popup>
  );
};

export default MarkerClickPopup;
