'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';
import {
  Heart,
  Share2,
  MapPin,
  Clock,
  Eye,
  ChevronRight,
  Phone,
  MessageCircle,
  Scale,
} from 'lucide-react';
import CarBadges from './CarBadges';
import CarQuickInfo from './CarQuickInfo';

type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];

interface CarListingCardEnhancedProps {
  listing: MarketplaceListing;
  locale: string;
  onFavorite?: (listingId: number) => void;
  onShare?: (listingId: number) => void;
  onCompare?: (listingId: number) => void;
  onQuickView?: (listing: MarketplaceListing) => void;
  isGrid?: boolean;
  isFavorited?: boolean;
  isComparing?: boolean;
}

export const CarListingCardEnhanced: React.FC<CarListingCardEnhancedProps> = ({
  listing,
  locale,
  onFavorite,
  onShare,
  onCompare,
  onQuickView,
  isGrid = true,
  isFavorited = false,
  isComparing = false,
}) => {
  const t = useTranslations('cars');
  const [imageError, setImageError] = useState(false);
  const [isHovered, setIsHovered] = useState(false);

  const extractCarAttributes = () => {
    const attrs = listing.attributes || {};
    const result: Record<string, any> = {};

    Object.entries(attrs).forEach(([key, value]) => {
      if (key.includes('car_make') || key.includes('make')) {
        result.make = value;
      } else if (key.includes('car_model') || key.includes('model')) {
        result.model = value;
      } else if (key.includes('year')) {
        result.year = value;
      } else if (key.includes('mileage') || key.includes('kilometers')) {
        result.mileage = value;
      } else if (key.includes('fuel') || key.includes('fuel_type')) {
        result.fuelType = value;
      } else if (key.includes('transmission') || key.includes('gearbox')) {
        result.transmission = value;
      } else if (key.includes('engine') || key.includes('displacement')) {
        result.engineSize = value;
      } else if (key.includes('color') || key.includes('colour')) {
        result.color = value;
      } else if (key.includes('body') || key.includes('body_type')) {
        result.bodyType = value;
      } else if (key.includes('price_reduction')) {
        result.priceReduction = value;
      } else if (key.includes('is_featured')) {
        result.isFeatured = value;
      }
    });

    return result;
  };

  const carAttrs = extractCarAttributes();

  const formatTimeAgo = (date: string) => {
    const created = new Date(date);
    const now = new Date();
    const diffMs = now.getTime() - created.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffHours / 24);

    if (diffHours < 1) {
      return t('time.justNow');
    } else if (diffHours < 24) {
      return t('time.hoursAgo', { hours: diffHours });
    } else if (diffDays < 7) {
      return t('time.daysAgo', { days: diffDays });
    } else if (diffDays < 30) {
      return t('time.weeksAgo', { weeks: Math.floor(diffDays / 7) });
    } else {
      return t('time.monthsAgo', { months: Math.floor(diffDays / 30) });
    }
  };

  const mainImage = listing.images?.[0];
  const imageUrl = mainImage?.thumbnail_url;
  const hasMultipleImages = (listing.images?.length || 0) > 1;

  const carTitle =
    carAttrs.make && carAttrs.model
      ? `${carAttrs.make} ${carAttrs.model}${carAttrs.year ? ' ' + carAttrs.year : ''}`
      : listing.title;

  const handleActionClick = (
    e: React.MouseEvent,
    action: ((id: number) => void) | undefined
  ) => {
    e.preventDefault();
    e.stopPropagation();
    if (listing.id !== undefined && action) {
      action(listing.id);
    }
  };

  const handleQuickView = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onQuickView?.(listing);
  };

  if (!isGrid) {
    return (
      <div
        className="card card-side bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden"
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <figure className="relative w-80 h-56 flex-shrink-0">
          {imageUrl && !imageError ? (
            <>
              <Image
                src={imageUrl}
                alt={carTitle || ''}
                fill
                className="object-cover"
                sizes="320px"
                onError={() => setImageError(true)}
              />
              {hasMultipleImages && (
                <div className="absolute bottom-2 right-2 badge badge-neutral badge-sm">
                  +{(listing.images?.length || 1) - 1} {t('common.photos')}
                </div>
              )}
            </>
          ) : (
            <div className="w-full h-full bg-gradient-to-br from-base-200 to-base-300 flex items-center justify-center">
              <div className="text-center">
                <div className="text-6xl mb-2">🚗</div>
                <p className="text-xs text-base-content/50">
                  {t('common.noPhoto')}
                </p>
              </div>
            </div>
          )}
          <CarBadges
            isNew={carAttrs.isNew}
            priceReduction={carAttrs.priceReduction}
            isTop={carAttrs.isTop}
            isFeatured={carAttrs.isFeatured}
            createdAt={listing.created_at}
          />
        </figure>

        <div className="card-body p-4">
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <Link
                href={`/${locale}/listing/${listing.id || 0}`}
                className="hover:text-primary transition-colors"
              >
                <h3 className="card-title text-lg font-bold">{carTitle}</h3>
              </Link>
              {listing.description && (
                <p className="text-sm text-base-content/70 line-clamp-2 mt-1">
                  {listing.description}
                </p>
              )}
            </div>
            <div className="flex gap-1 ml-2">
              <button
                onClick={(e) => handleActionClick(e, onFavorite)}
                className={`btn btn-ghost btn-sm btn-circle ${isFavorited ? 'text-error' : ''}`}
                title={t('actions.favorite')}
              >
                <Heart
                  className={`w-5 h-5 ${isFavorited ? 'fill-current' : ''}`}
                />
              </button>
              <button
                onClick={(e) => handleActionClick(e, onShare)}
                className="btn btn-ghost btn-sm btn-circle"
                title={t('actions.share')}
              >
                <Share2 className="w-5 h-5" />
              </button>
              <button
                onClick={(e) => handleActionClick(e, onCompare)}
                className={`btn btn-ghost btn-sm btn-circle ${isComparing ? 'text-primary' : ''}`}
                title={t('actions.compare')}
              >
                <Scale
                  className={`w-5 h-5 ${isComparing ? 'fill-current' : ''}`}
                />
              </button>
            </div>
          </div>

          <div className="my-3">
            <CarQuickInfo
              year={carAttrs.year}
              mileage={carAttrs.mileage}
              fuelType={carAttrs.fuelType}
              transmission={carAttrs.transmission}
              engineSize={carAttrs.engineSize}
              compact={true}
            />
          </div>

          <div className="divider my-2"></div>

          <div className="flex justify-between items-end">
            <div>
              {listing.price && (
                <div className="flex items-baseline gap-2">
                  <span className="text-2xl font-bold text-primary">
                    €{listing.price.toLocaleString()}
                  </span>
                  {carAttrs.priceReduction && (
                    <span className="text-sm text-error line-through opacity-70">
                      €
                      {Math.round(
                        listing.price / (1 - carAttrs.priceReduction / 100)
                      ).toLocaleString()}
                    </span>
                  )}
                </div>
              )}
              <div className="flex items-center gap-3 mt-1">
                {(listing.city || listing.country) && (
                  <div className="flex items-center gap-1 text-sm text-base-content/60">
                    <MapPin className="w-3 h-3" />
                    <span>{listing.city || listing.country}</span>
                  </div>
                )}
                {listing.created_at && (
                  <div className="flex items-center gap-1 text-sm text-base-content/60">
                    <Clock className="w-3 h-3" />
                    <span>{formatTimeAgo(listing.created_at)}</span>
                  </div>
                )}
                <div className="flex items-center gap-1 text-sm text-base-content/60">
                  <Eye className="w-3 h-3" />
                  <span>{listing.views_count || 0}</span>
                </div>
              </div>
            </div>

            <div className="flex gap-2">
              {isHovered && (
                <>
                  <button
                    onClick={handleQuickView}
                    className="btn btn-sm btn-outline"
                  >
                    {t('actions.quickView')}
                  </button>
                  <Link
                    href={`/${locale}/listing/${listing.id || 0}`}
                    className="btn btn-sm btn-primary"
                  >
                    {t('actions.viewDetails')}
                    <ChevronRight className="w-4 h-4" />
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="card bg-base-100 shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden group"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <figure className="relative aspect-[4/3]">
        {imageUrl && !imageError ? (
          <>
            <Image
              src={imageUrl}
              alt={carTitle || ''}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-300"
              sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
              onError={() => setImageError(true)}
            />
            {hasMultipleImages && (
              <div className="absolute bottom-2 right-2 badge badge-neutral badge-sm">
                +{(listing.images?.length || 1) - 1}
              </div>
            )}
          </>
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-base-200 to-base-300 flex items-center justify-center">
            <div className="text-center">
              <div className="text-6xl mb-2">🚗</div>
              <p className="text-xs text-base-content/50">
                {t('common.noPhoto')}
              </p>
            </div>
          </div>
        )}

        <CarBadges
          isNew={carAttrs.isNew}
          priceReduction={carAttrs.priceReduction}
          isTop={carAttrs.isTop}
          isFeatured={carAttrs.isFeatured}
          createdAt={listing.created_at}
        />

        <div className="absolute top-2 right-2 flex flex-col gap-1">
          <button
            onClick={(e) => handleActionClick(e, onFavorite)}
            className={`btn btn-ghost btn-sm btn-circle bg-base-100/80 backdrop-blur ${
              isFavorited ? 'text-error' : ''
            }`}
            title={t('actions.favorite')}
          >
            <Heart className={`w-4 h-4 ${isFavorited ? 'fill-current' : ''}`} />
          </button>
          <button
            onClick={(e) => handleActionClick(e, onCompare)}
            className={`btn btn-ghost btn-sm btn-circle bg-base-100/80 backdrop-blur ${
              isComparing ? 'text-primary' : ''
            }`}
            title={t('actions.compare')}
          >
            <Scale className={`w-4 h-4 ${isComparing ? 'fill-current' : ''}`} />
          </button>
        </div>

        {isHovered && (
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center transition-opacity duration-300">
            <button
              onClick={handleQuickView}
              className="btn btn-primary btn-sm"
            >
              {t('actions.quickView')}
            </button>
          </div>
        )}
      </figure>

      <div className="card-body p-4">
        <Link
          href={`/${locale}/listing/${listing.id || 0}`}
          className="hover:text-primary transition-colors"
        >
          <h3 className="card-title text-base line-clamp-1">{carTitle}</h3>
        </Link>

        <div className="my-2">
          <CarQuickInfo
            year={carAttrs.year}
            mileage={carAttrs.mileage}
            fuelType={carAttrs.fuelType}
            transmission={carAttrs.transmission}
            compact={true}
          />
        </div>

        <div className="divider my-2"></div>

        <div className="space-y-2">
          {listing.price && (
            <div className="flex items-baseline justify-between">
              <span className="text-xl font-bold text-primary">
                €{listing.price.toLocaleString()}
              </span>
              {carAttrs.priceReduction && (
                <span className="text-xs text-error line-through opacity-70">
                  €
                  {Math.round(
                    listing.price / (1 - carAttrs.priceReduction / 100)
                  ).toLocaleString()}
                </span>
              )}
            </div>
          )}

          <div className="flex items-center justify-between text-xs text-base-content/60">
            <div className="flex items-center gap-1">
              <MapPin className="w-3 h-3" />
              <span>
                {listing.city || listing.country || t('common.noLocation')}
              </span>
            </div>
            <div className="flex items-center gap-1">
              <Clock className="w-3 h-3" />
              <span>
                {listing.created_at ? formatTimeAgo(listing.created_at) : ''}
              </span>
            </div>
          </div>

          <div className="flex justify-between items-center">
            <div className="flex items-center gap-1 text-xs text-base-content/60">
              <Eye className="w-3 h-3" />
              <span>
                {listing.views_count || 0} {t('common.views')}
              </span>
            </div>
            {isHovered && (
              <div className="flex gap-1">
                <button
                  className="btn btn-ghost btn-xs btn-circle"
                  title={t('actions.call')}
                >
                  <Phone className="w-3 h-3" />
                </button>
                <button
                  className="btn btn-ghost btn-xs btn-circle"
                  title={t('actions.message')}
                >
                  <MessageCircle className="w-3 h-3" />
                </button>
                <button
                  onClick={(e) => handleActionClick(e, onShare)}
                  className="btn btn-ghost btn-xs btn-circle"
                  title={t('actions.share')}
                >
                  <Share2 className="w-3 h-3" />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CarListingCardEnhanced;
