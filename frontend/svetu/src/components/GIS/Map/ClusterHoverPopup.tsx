'use client';

import React from 'react';
import { Popup } from 'react-map-gl';
import Image from 'next/image';

interface ClusterListing {
  id: string;
  title: string;
  price: number;
  imageUrl?: string;
  category?: string;
  address?: string;
}

interface ClusterHoverPopupProps {
  coordinates: [number, number];
  listings: ClusterListing[];
  totalCount: number;
  onClose: () => void;
  onListingClick?: (listingId: string) => void;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
}

const ClusterHoverPopup: React.FC<ClusterHoverPopupProps> = ({
  coordinates,
  listings,
  totalCount,
  onClose,
  onListingClick,
  onMouseEnter,
  onMouseLeave,
}) => {
  // Форматирование цены
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ru-RU').format(price);
  };

  // Получение иконки категории
  const getCategoryIcon = (category?: string) => {
    if (!category) return '📦';

    const cat = category.toLowerCase();

    // Автомобили
    if (
      cat.includes('автомобил') ||
      cat.includes('car') ||
      cat.includes('vozilo')
    )
      return '🚗';
    // Недвижимость
    if (
      cat.includes('квартир') ||
      cat.includes('apartment') ||
      cat.includes('stan')
    )
      return '🏠';
    if (cat.includes('дом') || cat.includes('house') || cat.includes('kuća'))
      return '🏘️';
    if (cat.includes('комнат') || cat.includes('room') || cat.includes('soba'))
      return '🛏️';
    // Электроника
    if (
      cat.includes('телефон') ||
      cat.includes('phone') ||
      cat.includes('telefon')
    )
      return '📱';
    if (
      cat.includes('компьютер') ||
      cat.includes('computer') ||
      cat.includes('računar')
    )
      return '💻';
    // Работа
    if (cat.includes('работ') || cat.includes('job') || cat.includes('posao'))
      return '💼';
    // Услуги
    if (
      cat.includes('услуг') ||
      cat.includes('service') ||
      cat.includes('usluga')
    )
      return '🔧';

    return '📦';
  };

  const displayListings = listings.slice(0, 4); // Показываем максимум 4 объявления
  const remainingCount = totalCount - displayListings.length;

  return (
    <Popup
      longitude={coordinates[0]}
      latitude={coordinates[1]}
      closeButton={false}
      closeOnClick={false}
      onClose={onClose}
      anchor="bottom"
      offset={35}
      className="cluster-hover-popup"
      maxWidth="380px"
    >
      <div
        className="p-4"
        onMouseEnter={onMouseEnter}
        onMouseLeave={onMouseLeave}
      >
        {/* Заголовок */}
        <div className="mb-3 pb-3 border-b border-gray-200">
          <h3 className="text-lg font-bold text-gray-900">
            {totalCount} объявлени
            {totalCount === 1 ? 'е' : totalCount < 5 ? 'я' : 'й'}
          </h3>
          <p className="text-sm text-gray-600 mt-1">
            Приблизьте карту для просмотра
          </p>
        </div>

        {/* Список объявлений */}
        <div className="space-y-3">
          {displayListings.map((listing) => (
            <div
              key={listing.id}
              className="flex gap-3 p-2 -m-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              onClick={() => onListingClick?.(listing.id)}
            >
              {/* Изображение */}
              {listing.imageUrl ? (
                <div className="w-16 h-16 rounded-lg overflow-hidden flex-shrink-0">
                  <Image
                    src={listing.imageUrl}
                    alt={listing.title}
                    fill
                    className="object-cover"
                  />
                </div>
              ) : (
                <div className="w-16 h-16 rounded-lg bg-gray-100 flex items-center justify-center flex-shrink-0">
                  <span className="text-2xl opacity-50">
                    {getCategoryIcon(listing.category)}
                  </span>
                </div>
              )}

              {/* Информация */}
              <div className="flex-1 min-w-0">
                <h4 className="font-medium text-sm text-gray-900 line-clamp-1">
                  {listing.title}
                </h4>
                <div className="flex items-center gap-2 mt-1">
                  <span className="text-sm font-bold text-primary">
                    {formatPrice(listing.price)} RSD
                  </span>
                  {listing.category && (
                    <>
                      <span className="text-gray-300">•</span>
                      <span className="text-xs text-gray-600">
                        {listing.category}
                      </span>
                    </>
                  )}
                </div>
                {listing.address && (
                  <p className="text-xs text-gray-500 mt-1 line-clamp-1">
                    📍 {listing.address}
                  </p>
                )}
              </div>

              {/* Стрелка */}
              <div className="flex items-center">
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
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </div>
            </div>
          ))}
        </div>

        {/* Если есть еще объявления */}
        {remainingCount > 0 && (
          <div className="mt-3 pt-3 border-t border-gray-200">
            <p className="text-sm text-gray-600 text-center">
              и еще {remainingCount} объявлени
              {remainingCount === 1 ? 'е' : remainingCount < 5 ? 'я' : 'й'}
            </p>
          </div>
        )}

        {/* CTA */}
        <div className="mt-4 pt-3 border-t border-gray-200">
          <div className="flex items-center justify-between">
            <p className="text-sm text-gray-600">Нажмите на кластер</p>
            <div className="flex items-center gap-1 text-primary">
              <span className="text-sm font-medium">Приблизить</span>
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
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </Popup>
  );
};

export default ClusterHoverPopup;
