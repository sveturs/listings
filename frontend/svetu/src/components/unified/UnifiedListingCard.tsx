/**
 * Карточка unified listing (C2C или B2C)
 *
 * Универсальный компонент для отображения как C2C объявлений,
 * так и B2C товаров из витрин.
 */

'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import type { UnifiedListing } from '@/types/unified-listing';
import {
  getUnifiedListingDetailUrl,
  getMainImage,
} from '@/types/unified-listing';

export interface UnifiedListingCardProps {
  /** Unified listing (C2C или B2C) */
  listing: UnifiedListing;
  /** Текущая локаль */
  locale: string;
  /** Дополнительные CSS классы */
  className?: string;
  /** Приоритет загрузки изображения */
  priority?: boolean;
  /** Показывать бейдж типа (C2C/B2C) */
  showTypeBadge?: boolean;
  /** Показывать информацию о витрине (для B2C) */
  showStorefront?: boolean;
  /** Показывать количество просмотров */
  showViews?: boolean;
  /** Вариант отображения */
  variant?: 'grid' | 'list';
}

export function UnifiedListingCard({
  listing,
  locale,
  className = '',
  priority = false,
  showTypeBadge = true,
  showStorefront = true,
  showViews = false,
  variant = 'grid',
}: UnifiedListingCardProps) {
  const t = useTranslations('unified');

  // Получить главное изображение
  const mainImage = getMainImage(listing.images);

  // URL детальной страницы
  const detailUrl = getUnifiedListingDetailUrl(listing, locale);

  // Бейдж типа (C2C/B2C)
  const typeBadge = showTypeBadge && (
    <span
      className={`
        absolute top-2 right-2 z-10 px-2 py-1 rounded text-xs font-semibold
        shadow-md backdrop-blur-sm
        ${
          listing.source_type === 'c2c'
            ? 'bg-green-500/90 text-white'
            : 'bg-blue-500/90 text-white'
        }
      `}
      title={
        listing.source_type === 'c2c'
          ? t('badge.c2c_tooltip')
          : t('badge.b2c_tooltip')
      }
    >
      {listing.source_type === 'c2c' ? t('badge.c2c') : t('badge.b2c')}
    </span>
  );

  // Форматирование цены
  const formattedPrice = new Intl.NumberFormat(locale, {
    style: 'decimal',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(listing.price);

  if (variant === 'list') {
    return (
      <Link href={detailUrl} className={`block ${className}`}>
        <div className="bg-white rounded-lg shadow-sm hover:shadow-lg transition-shadow border border-gray-200 p-4">
          <div className="flex gap-4">
            {/* Изображение */}
            <div className="relative w-40 h-40 flex-shrink-0">
              {mainImage ? (
                <Image
                  src={mainImage.url}
                  alt={listing.title}
                  fill
                  className="object-cover rounded-lg"
                  priority={priority}
                />
              ) : (
                <div className="w-full h-full bg-gray-200 rounded-lg flex items-center justify-center">
                  <span className="text-gray-400 text-sm">{t('no_image')}</span>
                </div>
              )}
              {typeBadge}
            </div>

            {/* Информация */}
            <div className="flex-1 flex flex-col justify-between min-w-0">
              {/* Заголовок */}
              <div>
                <h3 className="font-semibold text-lg mb-2 line-clamp-2 text-gray-900">
                  {listing.title}
                </h3>

                <p className="text-gray-600 text-sm mb-3 line-clamp-2">
                  {listing.description}
                </p>

                {/* Местоположение */}
                {listing.city && (
                  <div className="flex items-center gap-1 text-sm text-gray-500 mb-2">
                    <span>📍</span>
                    <span>{listing.city}</span>
                  </div>
                )}
              </div>

              {/* Footer */}
              <div className="flex items-end justify-between">
                {/* Цена */}
                <div>
                  <span className="text-2xl font-bold text-blue-600">
                    {formattedPrice} RSD
                  </span>
                </div>

                {/* Метаинформация */}
                <div className="flex items-center gap-4 text-sm text-gray-500">
                  {showViews && (
                    <span className="flex items-center gap-1">
                      👁️ {listing.views_count}
                    </span>
                  )}
                  {listing.condition && (
                    <span>{t(`condition.${listing.condition}`)}</span>
                  )}
                </div>
              </div>

              {/* Витрина (только для B2C) */}
              {showStorefront &&
                listing.source_type === 'b2c' &&
                listing.storefront && (
                  <div className="mt-3 pt-3 border-t border-gray-200">
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-gray-500">
                        {t('from_storefront')}:
                      </span>
                      <span className="text-xs font-medium text-gray-700">
                        {listing.storefront.name}
                      </span>
                    </div>
                  </div>
                )}
            </div>
          </div>
        </div>
      </Link>
    );
  }

  // Grid variant (по умолчанию)
  return (
    <Link href={detailUrl} className={`block group ${className}`}>
      <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-all duration-300 group-hover:scale-105">
        {/* Изображение */}
        <div className="relative aspect-square">
          {mainImage ? (
            <Image
              src={mainImage.url}
              alt={listing.title}
              fill
              className="object-cover"
              priority={priority}
            />
          ) : (
            <div className="w-full h-full bg-gray-200 flex items-center justify-center">
              <span className="text-gray-400">{t('no_image')}</span>
            </div>
          )}
          {typeBadge}

          {/* Overlay gradient */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>

        {/* Информация */}
        <div className="p-4">
          {/* Заголовок */}
          <h3 className="font-semibold text-lg mb-2 line-clamp-2 text-gray-900 group-hover:text-blue-600 transition-colors">
            {listing.title}
          </h3>

          {/* Описание */}
          <p className="text-gray-600 text-sm mb-3 line-clamp-2">
            {listing.description}
          </p>

          {/* Местоположение */}
          {listing.city && (
            <div className="flex items-center gap-1 text-sm text-gray-500 mb-3">
              <span>📍</span>
              <span className="truncate">{listing.city}</span>
            </div>
          )}

          {/* Цена и состояние */}
          <div className="flex items-center justify-between mb-3">
            <span className="text-xl font-bold text-blue-600">
              {formattedPrice} RSD
            </span>

            {listing.condition && (
              <span className="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">
                {t(`condition.${listing.condition}`)}
              </span>
            )}
          </div>

          {/* Витрина (только для B2C) */}
          {showStorefront &&
            listing.source_type === 'b2c' &&
            listing.storefront && (
              <div className="pt-3 border-t border-gray-200">
                <div className="flex items-center gap-2">
                  <span className="text-xs text-gray-500">🏪</span>
                  <span className="text-xs font-medium text-gray-700 truncate">
                    {listing.storefront.name}
                  </span>
                </div>
              </div>
            )}

          {/* Просмотры */}
          {showViews && listing.views_count > 0 && (
            <div className="flex items-center gap-1 text-xs text-gray-500 mt-2">
              <span>👁️</span>
              <span>{listing.views_count}</span>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
}

/**
 * Компонент сетки unified listings
 */
export function UnifiedListingsGrid({
  listings,
  locale,
  className = '',
  emptyMessage,
}: {
  listings: UnifiedListing[];
  locale: string;
  className?: string;
  emptyMessage?: string;
}) {
  const t = useTranslations('unified');

  if (listings.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 text-lg">
          {emptyMessage || t('no_listings_found')}
        </p>
      </div>
    );
  }

  return (
    <div
      className={`grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 ${className}`}
    >
      {listings.map((listing) => (
        <UnifiedListingCard
          key={`${listing.source_type}-${listing.id}`}
          listing={listing}
          locale={locale}
        />
      ))}
    </div>
  );
}

/**
 * Компонент списка unified listings
 */
export function UnifiedListingsList({
  listings,
  locale,
  className = '',
  emptyMessage,
}: {
  listings: UnifiedListing[];
  locale: string;
  className?: string;
  emptyMessage?: string;
}) {
  const t = useTranslations('unified');

  if (listings.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 text-lg">
          {emptyMessage || t('no_listings_found')}
        </p>
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {listings.map((listing) => (
        <UnifiedListingCard
          key={`${listing.source_type}-${listing.id}`}
          listing={listing}
          locale={locale}
          variant="list"
        />
      ))}
    </div>
  );
}
