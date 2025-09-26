'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import configManager from '@/config';
import type { components } from '@/types/generated/api';
import {
  X,
  ChevronLeft,
  ChevronRight,
  MapPin,
  Phone,
  MessageCircle,
  Heart,
  Share2,
  Scale,
  ExternalLink,
  Eye,
  Clock,
  User,
} from 'lucide-react';
import CarBadges from './CarBadges';
import CarQuickInfo from './CarQuickInfo';

type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];

interface CarQuickViewModalProps {
  isOpen: boolean;
  onClose: () => void;
  listing: MarketplaceListing | null;
  onNext?: () => void;
  onPrevious?: () => void;
  hasNext?: boolean;
  hasPrevious?: boolean;
  locale: string;
  onFavorite?: (listingId: number) => void;
  onCompare?: (listingId: number) => void;
  isFavorited?: boolean;
  isComparing?: boolean;
}

export const CarQuickViewModal: React.FC<CarQuickViewModalProps> = ({
  isOpen,
  onClose,
  listing,
  onNext,
  onPrevious,
  hasNext = false,
  hasPrevious = false,
  locale,
  onFavorite,
  onCompare,
  isFavorited = false,
  isComparing = false,
}) => {
  const t = useTranslations('cars');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [imageError, setImageError] = useState<Record<number, boolean>>({});

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      setCurrentImageIndex(0);
      setImageError({});
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen) return;

      switch (e.key) {
        case 'Escape':
          onClose();
          break;
        case 'ArrowLeft':
          if (hasPrevious) onPrevious?.();
          break;
        case 'ArrowRight':
          if (hasNext) onNext?.();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose, onNext, onPrevious, hasNext, hasPrevious]);

  if (!listing) return null;

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
      } else if (key.includes('seats')) {
        result.seats = value;
      } else if (key.includes('doors')) {
        result.doors = value;
      } else if (key.includes('power') || key.includes('horsepower')) {
        result.power = value;
      } else if (key.includes('drive') || key.includes('drivetrain')) {
        result.driveType = value;
      }
    });

    return result;
  };

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
    } else {
      return t('time.monthsAgo', { months: Math.floor(diffDays / 30) });
    }
  };

  const carAttrs = extractCarAttributes();
  const images = listing.images || [];
  const currentImage = images[currentImageIndex];
  const imageUrl = currentImage?.public_url
    ? configManager.buildImageUrl(currentImage.public_url)
    : null;

  const carTitle =
    carAttrs.make && carAttrs.model
      ? `${carAttrs.make} ${carAttrs.model}${carAttrs.year ? ' ' + carAttrs.year : ''}`
      : listing.title;

  const handleImageNavigation = (direction: 'prev' | 'next') => {
    if (direction === 'prev') {
      setCurrentImageIndex((prev) => (prev > 0 ? prev - 1 : images.length - 1));
    } else {
      setCurrentImageIndex((prev) => (prev < images.length - 1 ? prev + 1 : 0));
    }
  };

  return (
    <div
      className={`fixed inset-0 z-[100] flex items-center justify-center transition-opacity duration-300 ${
        isOpen
          ? 'opacity-100 pointer-events-auto'
          : 'opacity-0 pointer-events-none'
      }`}
    >
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />

      <div className="relative bg-base-100 rounded-lg shadow-2xl max-w-6xl w-full mx-4 max-h-[90vh] overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b">
          <h2 className="text-xl font-bold">{carTitle}</h2>
          <div className="flex items-center gap-2">
            {hasPrevious && (
              <button
                onClick={onPrevious}
                className="btn btn-ghost btn-sm btn-circle"
                title={t('navigation.previous')}
              >
                <ChevronLeft className="w-5 h-5" />
              </button>
            )}
            {hasNext && (
              <button
                onClick={onNext}
                className="btn btn-ghost btn-sm btn-circle"
                title={t('navigation.next')}
              >
                <ChevronRight className="w-5 h-5" />
              </button>
            )}
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        <div className="flex flex-col lg:flex-row overflow-y-auto max-h-[calc(90vh-4rem)]">
          <div className="lg:w-3/5 p-4">
            <div className="relative aspect-[16/10] bg-base-200 rounded-lg overflow-hidden">
              {imageUrl && !imageError[currentImageIndex] ? (
                <>
                  <Image
                    src={imageUrl}
                    alt={`${carTitle} - ${currentImageIndex + 1}`}
                    fill
                    className="object-contain"
                    sizes="800px"
                    onError={() =>
                      setImageError((prev) => ({
                        ...prev,
                        [currentImageIndex]: true,
                      }))
                    }
                  />
                  {images.length > 1 && (
                    <>
                      <button
                        onClick={() => handleImageNavigation('prev')}
                        className="absolute left-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80"
                      >
                        <ChevronLeft className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleImageNavigation('next')}
                        className="absolute right-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80"
                      >
                        <ChevronRight className="w-4 h-4" />
                      </button>
                      <div className="absolute bottom-2 left-1/2 -translate-x-1/2 badge badge-neutral">
                        {currentImageIndex + 1} / {images.length}
                      </div>
                    </>
                  )}
                </>
              ) : (
                <div className="w-full h-full flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-7xl mb-3">🚗</div>
                    <p className="text-base-content/50">
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
            </div>

            {images.length > 1 && (
              <div className="flex gap-2 mt-3 overflow-x-auto">
                {images.map((img, index) => (
                  <button
                    key={index}
                    onClick={() => setCurrentImageIndex(index)}
                    className={`relative w-20 h-20 flex-shrink-0 rounded-lg overflow-hidden border-2 transition-all ${
                      index === currentImageIndex
                        ? 'border-primary'
                        : 'border-transparent opacity-70 hover:opacity-100'
                    }`}
                  >
                    {img.public_url && !imageError[index] ? (
                      <Image
                        src={configManager.buildImageUrl(img.public_url)}
                        alt={`Thumbnail ${index + 1}`}
                        fill
                        className="object-cover"
                        sizes="80px"
                        onError={() =>
                          setImageError((prev) => ({ ...prev, [index]: true }))
                        }
                      />
                    ) : (
                      <div className="w-full h-full bg-base-200 flex items-center justify-center">
                        <div className="text-2xl">🚗</div>
                      </div>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>

          <div className="lg:w-2/5 p-4 lg:border-l">
            <div className="space-y-4">
              <div>
                {listing.price && (
                  <div className="text-3xl font-bold text-primary">
                    €{listing.price.toLocaleString()}
                  </div>
                )}
                <div className="flex items-center gap-3 mt-2 text-sm text-base-content/60">
                  {(listing.city || listing.country) && (
                    <div className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      <span>{listing.city || listing.country}</span>
                    </div>
                  )}
                  {listing.created_at && (
                    <div className="flex items-center gap-1">
                      <Clock className="w-4 h-4" />
                      <span>{formatTimeAgo(listing.created_at)}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-1">
                    <Eye className="w-4 h-4" />
                    <span>{listing.views_count || 0}</span>
                  </div>
                </div>
              </div>

              <div className="divider"></div>

              <CarQuickInfo
                year={carAttrs.year}
                mileage={carAttrs.mileage}
                fuelType={carAttrs.fuelType}
                transmission={carAttrs.transmission}
                engineSize={carAttrs.engineSize}
                color={carAttrs.color}
                bodyType={carAttrs.bodyType}
                seats={carAttrs.seats}
              />

              {listing.description && (
                <>
                  <div className="divider"></div>
                  <div>
                    <h3 className="font-semibold mb-2">
                      {t('quickView.description')}
                    </h3>
                    <p className="text-sm text-base-content/80 line-clamp-5">
                      {listing.description}
                    </p>
                  </div>
                </>
              )}

              {listing.user && (
                <>
                  <div className="divider"></div>
                  <div className="flex items-center gap-3">
                    <div className="avatar">
                      <div className="w-12 rounded-full bg-base-200">
                        <User className="w-12 h-12 p-3" />
                      </div>
                    </div>
                    <div>
                      <p className="font-semibold">
                        {listing.user.name || listing.user.email}
                      </p>
                      <p className="text-xs text-base-content/60">
                        {t('quickView.seller')}
                      </p>
                    </div>
                  </div>
                </>
              )}

              <div className="divider"></div>

              <div className="flex gap-2">
                <button className="btn btn-primary flex-1">
                  <Phone className="w-4 h-4" />
                  {t('actions.call')}
                </button>
                <button className="btn btn-outline flex-1">
                  <MessageCircle className="w-4 h-4" />
                  {t('actions.message')}
                </button>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => onFavorite?.(listing.id || 0)}
                  className={`btn btn-outline flex-1 ${isFavorited ? 'btn-error' : ''}`}
                >
                  <Heart
                    className={`w-4 h-4 ${isFavorited ? 'fill-current' : ''}`}
                  />
                  {t('actions.favorite')}
                </button>
                <button
                  onClick={() => onCompare?.(listing.id || 0)}
                  className={`btn btn-outline flex-1 ${isComparing ? 'btn-primary' : ''}`}
                >
                  <Scale
                    className={`w-4 h-4 ${isComparing ? 'fill-current' : ''}`}
                  />
                  {t('actions.compare')}
                </button>
                <button className="btn btn-outline btn-square">
                  <Share2 className="w-4 h-4" />
                </button>
              </div>

              <Link
                href={`/${locale}/listing/${listing.id || 0}`}
                className="btn btn-ghost btn-block"
                onClick={onClose}
              >
                <ExternalLink className="w-4 h-4" />
                {t('actions.viewFullDetails')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CarQuickViewModal;
