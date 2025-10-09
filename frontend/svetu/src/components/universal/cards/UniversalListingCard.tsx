'use client';

import Link from 'next/link';
import { FC, useMemo } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import {
  FaHeart,
  FaRegHeart,
  FaExchangeAlt,
  FaMapMarkerAlt,
  FaClock,
} from 'react-icons/fa';
import { useAppSelector, useAppDispatch } from '@/store/hooks';
import {
  addToFavorites,
  removeFromFavorites,
} from '@/store/slices/favoritesSlice';
import { addItem, removeItem } from '@/store/slices/universalCompareSlice';
import {
  makeSelectCompareItemsByType,
  selectFavorites,
} from '@/store/selectors/compareSelectors';
import LazyImage from '@/components/common/LazyImage';

// Типы для универсальной карточки
export interface UniversalListingData {
  id: number;
  title: string;
  price: number;
  currency?: string;
  images?: string[];
  location?: {
    city?: string;
    district?: string;
    address?: string;
  };
  category?: string;
  categorySlug?: string;
  createdAt?: string;
  updatedAt?: string;
  // Универсальные атрибуты
  attributes?: Record<string, any>;
  // Специфичные поля для категорий
  customFields?: {
    label: string;
    value: string | number;
    icon?: React.ComponentType<any>;
  }[];
  // Бейджи
  badges?: {
    type: 'new' | 'discount' | 'top' | 'recommended' | 'verified' | 'urgent';
    label?: string;
    value?: string | number;
  }[];
  // Статистика
  stats?: {
    views?: number;
    favorites?: number;
    calls?: number;
  };
}

interface UniversalListingCardProps {
  data: UniversalListingData;
  type?: 'cars' | 'real_estate' | 'electronics' | 'marketplace' | 'storefront';
  layout?: 'grid' | 'list';
  showBadges?: boolean;
  showFavorite?: boolean;
  showCompare?: boolean;
  showStats?: boolean;
  onQuickView?: () => void;
  className?: string;
}

// Конфигурация для разных типов
const CARD_CONFIG = {
  cars: {
    linkPrefix: '/cars',
    compareEnabled: true,
    maxCompareItems: 3,
    defaultImage: '/images/car-placeholder.jpg',
  },
  real_estate: {
    linkPrefix: '/real-estate',
    compareEnabled: true,
    maxCompareItems: 4,
    defaultImage: '/images/property-placeholder.jpg',
  },
  electronics: {
    linkPrefix: '/electronics',
    compareEnabled: true,
    maxCompareItems: 5,
    defaultImage: '/images/product-placeholder.jpg',
  },
  marketplace: {
    linkPrefix: '/c2c',
    compareEnabled: false,
    maxCompareItems: 0,
    defaultImage: '/images/placeholder.jpg',
  },
  storefront: {
    linkPrefix: '/products',
    compareEnabled: false,
    maxCompareItems: 0,
    defaultImage: '/images/product-placeholder.jpg',
  },
};

const UniversalListingCard: FC<UniversalListingCardProps> = ({
  data,
  type = 'marketplace',
  layout = 'grid',
  showBadges = true,
  showFavorite = true,
  showCompare = true,
  showStats = false,
  onQuickView,
  className = '',
}) => {
  const t = useTranslations('marketplace');
  const locale = useLocale();
  const dispatch = useAppDispatch();

  const selectCompareItemsByType = useMemo(makeSelectCompareItemsByType, []);
  const favorites = useAppSelector(selectFavorites);
  const compareItems = useAppSelector((state) =>
    selectCompareItemsByType(state, type)
  );

  const config = CARD_CONFIG[type];
  const isFavorite = useMemo(
    () => favorites.some((fav) => fav.id === data.id),
    [favorites, data.id]
  );
  const isInCompare = useMemo(
    () => compareItems.some((item) => item.id === data.id),
    [compareItems, data.id]
  );
  const canCompare = config.compareEnabled && showCompare;

  // Форматирование цены
  const formattedPrice = useMemo(() => {
    if (!data.price) return t('priceOnRequest');
    const currency = data.currency || '€';
    return new Intl.NumberFormat(locale === 'ru' ? 'ru-RU' : 'en-US', {
      style: 'currency',
      currency: currency === '€' ? 'EUR' : 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(data.price);
  }, [data.price, data.currency, locale, t]);

  // Форматирование даты
  const formattedDate = useMemo(() => {
    if (!data.createdAt) return '';
    const date = new Date(data.createdAt);
    const now = new Date();
    const diffHours = Math.floor(
      (now.getTime() - date.getTime()) / (1000 * 60 * 60)
    );

    if (diffHours < 1) return t('justNow');
    if (diffHours < 24) return t('hoursAgo', { hours: diffHours });
    if (diffHours < 48) return t('yesterday');
    if (diffHours < 168)
      return t('daysAgo', { days: Math.floor(diffHours / 24) });

    return date.toLocaleDateString(locale);
  }, [data.createdAt, locale, t]);

  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isFavorite) {
      dispatch(
        removeFromFavorites({
          id: data.id,
          type: type === 'storefront' ? 'storefront' : 'marketplace',
        })
      );
    } else {
      dispatch(
        addToFavorites({
          id: data.id,
          type: type === 'storefront' ? 'storefront' : 'marketplace',
        })
      );
    }
  };

  const handleCompareClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (isInCompare) {
      dispatch(removeItem(data.id));
    } else {
      dispatch(
        addItem({
          id: data.id,
          category: type,
          title: data.title,
          price: data.price,
          image: data.images?.[0],
          attributes: data.attributes || {},
        })
      );
    }
  };

  const handleQuickView = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (onQuickView) onQuickView();
  };

  // Grid layout
  if (layout === 'grid') {
    return (
      <div
        className={`card card-compact bg-base-100 shadow-sm hover:shadow-lg transition-all duration-300 ${className}`}
      >
        <Link
          href={`/${locale}${config.linkPrefix}/${data.id}`}
          className="block"
        >
          <figure className="relative aspect-[4/3] overflow-hidden bg-base-200">
            {data.images && data.images[0] ? (
              <LazyImage
                src={data.images[0]}
                alt={data.title}
                fill
                className="object-cover hover:scale-105 transition-transform duration-300"
                sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                priority={false}
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-base-content/30">
                <span>{t('noImage')}</span>
              </div>
            )}

            {/* Бейджи */}
            {showBadges && data.badges && data.badges.length > 0 && (
              <div className="absolute top-2 left-2 flex flex-wrap gap-1">
                {data.badges.map((badge, index) => (
                  <span
                    key={index}
                    className={`badge ${
                      badge.type === 'new'
                        ? 'badge-primary'
                        : badge.type === 'discount'
                          ? 'badge-error'
                          : badge.type === 'top'
                            ? 'badge-warning'
                            : badge.type === 'recommended'
                              ? 'badge-success'
                              : badge.type === 'verified'
                                ? 'badge-info'
                                : 'badge-secondary'
                    } text-xs font-medium`}
                  >
                    {badge.label || badge.type}
                    {badge.value && ` ${badge.value}`}
                  </span>
                ))}
              </div>
            )}

            {/* Быстрые действия */}
            <div className="absolute top-2 right-2 flex flex-col gap-2">
              {showFavorite && (
                <button
                  onClick={handleFavoriteClick}
                  className="btn btn-circle btn-sm bg-base-100/80 backdrop-blur-sm hover:bg-base-100"
                  aria-label={
                    isFavorite ? t('removeFromFavorites') : t('addToFavorites')
                  }
                >
                  {isFavorite ? (
                    <FaHeart className="text-error" />
                  ) : (
                    <FaRegHeart />
                  )}
                </button>
              )}

              {canCompare && (
                <button
                  onClick={handleCompareClick}
                  className="btn btn-circle btn-sm bg-base-100/80 backdrop-blur-sm hover:bg-base-100"
                  aria-label={
                    isInCompare ? t('removeFromCompare') : t('addToCompare')
                  }
                >
                  <FaExchangeAlt
                    className={isInCompare ? 'text-primary' : ''}
                  />
                </button>
              )}
            </div>
          </figure>

          <div className="card-body p-4">
            {/* Заголовок и цена */}
            <h3 className="card-title text-base line-clamp-2 mb-2">
              {data.title}
            </h3>

            {/* Кастомные поля */}
            {data.customFields && data.customFields.length > 0 && (
              <div className="flex flex-wrap gap-2 text-sm text-base-content/70 mb-2">
                {data.customFields.slice(0, 3).map((field, index) => (
                  <span key={index} className="flex items-center gap-1">
                    {field.icon && <field.icon className="w-3 h-3" />}
                    <span>{field.value}</span>
                  </span>
                ))}
              </div>
            )}

            {/* Цена и локация */}
            <div className="flex justify-between items-end mt-auto">
              <div>
                <div className="text-lg font-bold text-primary">
                  {formattedPrice}
                </div>
                {data.location && (
                  <div className="flex items-center gap-1 text-xs text-base-content/60 mt-1">
                    <FaMapMarkerAlt className="w-3 h-3" />
                    <span>
                      {data.location.district && data.location.city
                        ? `${data.location.district}, ${data.location.city}`
                        : data.location.city || data.location.district}
                    </span>
                  </div>
                )}
              </div>

              {formattedDate && (
                <div className="flex items-center gap-1 text-xs text-base-content/60">
                  <FaClock className="w-3 h-3" />
                  <span>{formattedDate}</span>
                </div>
              )}
            </div>

            {/* Статистика */}
            {showStats && data.stats && (
              <div className="flex gap-3 text-xs text-base-content/60 mt-2 pt-2 border-t">
                {data.stats.views !== undefined && (
                  <span>
                    {t('views')}: {data.stats.views}
                  </span>
                )}
                {data.stats.favorites !== undefined &&
                  data.stats.favorites > 0 && (
                    <span>
                      {t('inFavorites')}: {data.stats.favorites}
                    </span>
                  )}
                {data.stats.calls !== undefined && data.stats.calls > 0 && (
                  <span>
                    {t('calls')}: {data.stats.calls}
                  </span>
                )}
              </div>
            )}

            {/* Кнопка быстрого просмотра */}
            {onQuickView && (
              <button
                onClick={handleQuickView}
                className="btn btn-sm btn-ghost w-full mt-2"
              >
                {t('quickView')}
              </button>
            )}
          </div>
        </Link>
      </div>
    );
  }

  // List layout
  return (
    <div
      className={`card bg-base-100 shadow-sm hover:shadow-lg transition-all duration-300 ${className}`}
    >
      <Link href={`/${locale}${config.linkPrefix}/${data.id}`} className="flex">
        <figure className="relative w-48 h-36 flex-shrink-0 overflow-hidden bg-base-200">
          {data.images && data.images[0] ? (
            <LazyImage
              src={data.images[0]}
              alt={data.title}
              fill
              className="object-cover hover:scale-105 transition-transform duration-300"
              sizes="200px"
              priority={false}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-base-content/30">
              <span>{t('noImage')}</span>
            </div>
          )}

          {/* Бейджи в list режиме */}
          {showBadges && data.badges && data.badges.length > 0 && (
            <div className="absolute top-1 left-1">
              <span
                className={`badge badge-sm ${
                  data.badges[0].type === 'new'
                    ? 'badge-primary'
                    : data.badges[0].type === 'discount'
                      ? 'badge-error'
                      : 'badge-secondary'
                }`}
              >
                {data.badges[0].label}
              </span>
            </div>
          )}
        </figure>

        <div className="card-body p-4 flex-1">
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <h3 className="font-semibold text-base line-clamp-1 mb-1">
                {data.title}
              </h3>

              {/* Кастомные поля в строку */}
              {data.customFields && data.customFields.length > 0 && (
                <div className="flex gap-3 text-sm text-base-content/70 mb-2">
                  {data.customFields.slice(0, 4).map((field, index) => (
                    <span key={index} className="flex items-center gap-1">
                      {field.icon && <field.icon className="w-3 h-3" />}
                      <span>{field.value}</span>
                    </span>
                  ))}
                </div>
              )}

              <div className="flex items-center gap-3 text-sm text-base-content/60">
                {data.location && (
                  <span className="flex items-center gap-1">
                    <FaMapMarkerAlt className="w-3 h-3" />
                    {data.location.city}
                  </span>
                )}
                {formattedDate && (
                  <span className="flex items-center gap-1">
                    <FaClock className="w-3 h-3" />
                    {formattedDate}
                  </span>
                )}
              </div>
            </div>

            <div className="flex flex-col items-end gap-2">
              <div className="text-lg font-bold text-primary">
                {formattedPrice}
              </div>

              <div className="flex gap-1">
                {showFavorite && (
                  <button
                    onClick={handleFavoriteClick}
                    className="btn btn-circle btn-sm"
                    aria-label={
                      isFavorite
                        ? t('removeFromFavorites')
                        : t('addToFavorites')
                    }
                  >
                    {isFavorite ? (
                      <FaHeart className="text-error" />
                    ) : (
                      <FaRegHeart />
                    )}
                  </button>
                )}

                {canCompare && (
                  <button
                    onClick={handleCompareClick}
                    className="btn btn-circle btn-sm"
                    aria-label={
                      isInCompare ? t('removeFromCompare') : t('addToCompare')
                    }
                  >
                    <FaExchangeAlt
                      className={isInCompare ? 'text-primary' : ''}
                    />
                  </button>
                )}

                {onQuickView && (
                  <button
                    onClick={handleQuickView}
                    className="btn btn-sm btn-ghost"
                  >
                    {t('quickView')}
                  </button>
                )}
              </div>
            </div>
          </div>

          {/* Статистика в list режиме */}
          {showStats && data.stats && (
            <div className="flex gap-3 text-xs text-base-content/60 mt-2 pt-2 border-t">
              {data.stats.views && (
                <span>
                  {t('views')}: {data.stats.views}
                </span>
              )}
              {data.stats.favorites && (
                <span>
                  {t('inFavorites')}: {data.stats.favorites}
                </span>
              )}
              {data.stats.calls && (
                <span>
                  {t('calls')}: {data.stats.calls}
                </span>
              )}
            </div>
          )}
        </div>
      </Link>
    </div>
  );
};

export default UniversalListingCard;
