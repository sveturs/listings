'use client';

import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';
import {
  Car,
  Calendar,
  Gauge,
  Fuel,
  MapPin,
  Heart,
  Share2,
  Eye,
} from 'lucide-react';

type MarketplaceListing = components['schemas']['models.MarketplaceListing'];

interface CarListingCardProps {
  listing: MarketplaceListing;
  locale: string;
  onFavorite?: (listingId: number) => void;
  onShare?: (listingId: number) => void;
  isGrid?: boolean;
}

export const CarListingCard: React.FC<CarListingCardProps> = ({
  listing,
  locale,
  onFavorite,
  onShare,
  isGrid = true,
}) => {
  const t = useTranslations('cars');

  // Извлекаем автомобильные атрибуты из listing.attributes
  const extractCarAttributes = () => {
    const attrs = listing.attributes || {};
    const result: Record<string, any> = {};

    // Ищем атрибуты по ключам
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
      }
    });

    return result;
  };

  const carAttrs = extractCarAttributes();

  // Форматирование пробега
  const formatMileage = (mileage: any) => {
    if (!mileage) return null;
    const num = parseInt(mileage);
    if (isNaN(num)) return mileage;
    return num.toLocaleString() + ' km';
  };

  // Форматирование топлива
  const formatFuel = (fuel: string) => {
    if (!fuel) return null;
    const fuelTypes: Record<string, string> = {
      petrol: t('filters.petrol'),
      diesel: t('filters.diesel'),
      electric: t('filters.electric'),
      hybrid: t('filters.hybrid'),
      lpg: t('filters.lpg'),
    };
    return fuelTypes[fuel.toLowerCase()] || fuel;
  };

  // Форматирование коробки передач
  const formatTransmission = (transmission: string) => {
    if (!transmission) return null;
    const types: Record<string, string> = {
      manual: t('filters.manual'),
      automatic: t('filters.automatic'),
      semiAutomatic: t('filters.semiAutomatic'),
    };
    return types[transmission.toLowerCase()] || transmission;
  };

  // Определяем изображение
  const mainImage = listing.images?.[0];
  const imageUrl = mainImage?.thumbnail_url;

  // Создаем заголовок с маркой и моделью
  const carTitle =
    carAttrs.make && carAttrs.model
      ? `${carAttrs.make} ${carAttrs.model}`
      : listing.title;

  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (listing.id !== undefined) {
      onFavorite?.(listing.id);
    }
  };

  const handleShareClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (listing.id !== undefined) {
      onShare?.(listing.id);
    }
  };

  if (!isGrid) {
    // Горизонтальная карточка для списка
    return (
      <Link
        href={`/${locale}/listing/${listing.id || 0}`}
        className="card card-side bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300"
      >
        <figure className="relative w-72 h-48 flex-shrink-0">
          {imageUrl ? (
            <Image
              src={imageUrl}
              alt={carTitle || ''}
              fill
              className="object-cover"
              sizes="288px"
            />
          ) : (
            <div className="w-full h-full bg-base-200 flex items-center justify-center">
              <Car className="w-12 h-12 text-base-content/30" />
            </div>
          )}
        </figure>

        <div className="card-body">
          <div className="flex justify-between items-start">
            <div>
              <h3 className="card-title text-lg">{carTitle}</h3>
              {carAttrs.year && (
                <p className="text-sm text-base-content/60">{listing.title}</p>
              )}
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleFavoriteClick}
                className="btn btn-ghost btn-sm btn-circle"
              >
                <Heart className="w-5 h-5" />
              </button>
              <button
                onClick={handleShareClick}
                className="btn btn-ghost btn-sm btn-circle"
              >
                <Share2 className="w-5 h-5" />
              </button>
            </div>
          </div>

          <div className="grid grid-cols-3 gap-2 my-2">
            {carAttrs.year && (
              <div className="flex items-center gap-1 text-sm">
                <Calendar className="w-4 h-4 text-primary" />
                <span>{carAttrs.year}</span>
              </div>
            )}
            {carAttrs.mileage && (
              <div className="flex items-center gap-1 text-sm">
                <Gauge className="w-4 h-4 text-primary" />
                <span>{formatMileage(carAttrs.mileage)}</span>
              </div>
            )}
            {carAttrs.fuelType && (
              <div className="flex items-center gap-1 text-sm">
                <Fuel className="w-4 h-4 text-primary" />
                <span>{formatFuel(carAttrs.fuelType)}</span>
              </div>
            )}
          </div>

          {listing.description && (
            <p className="text-sm text-base-content/70 line-clamp-2">
              {listing.description}
            </p>
          )}

          <div className="card-actions justify-between items-center mt-auto">
            <div>
              {listing.price && (
                <div className="text-2xl font-bold text-primary">
                  €{listing.price.toLocaleString()}
                </div>
              )}
              {(listing.city || listing.country) && (
                <div className="flex items-center gap-1 text-sm text-base-content/60">
                  <MapPin className="w-3 h-3" />
                  <span>{listing.city || listing.country}</span>
                </div>
              )}
            </div>
            <div className="flex items-center gap-1 text-sm text-base-content/60">
              <Eye className="w-4 h-4" />
              <span>{listing.views_count || 0}</span>
            </div>
          </div>
        </div>
      </Link>
    );
  }

  // Вертикальная карточка для сетки
  return (
    <Link
      href={`/${locale}/listing/${listing.id || 0}`}
      className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300"
    >
      <figure className="relative aspect-[4/3]">
        {imageUrl ? (
          <Image
            src={imageUrl}
            alt={carTitle || ''}
            fill
            className="object-cover"
            sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
          />
        ) : (
          <div className="w-full h-full bg-base-200 flex items-center justify-center">
            <Car className="w-12 h-12 text-base-content/30" />
          </div>
        )}
        <div className="absolute top-2 right-2 flex gap-1">
          <button
            onClick={handleFavoriteClick}
            className="btn btn-ghost btn-sm btn-circle bg-base-100/80 backdrop-blur"
          >
            <Heart className="w-4 h-4" />
          </button>
        </div>
      </figure>

      <div className="card-body p-4">
        <h3 className="card-title text-base line-clamp-1">{carTitle}</h3>
        {carAttrs.year && (
          <p className="text-sm text-base-content/60 -mt-1">{listing.title}</p>
        )}

        <div className="grid grid-cols-2 gap-2 my-2">
          {carAttrs.year && (
            <div className="flex items-center gap-1 text-xs">
              <Calendar className="w-3 h-3 text-primary" />
              <span>{carAttrs.year}</span>
            </div>
          )}
          {carAttrs.mileage && (
            <div className="flex items-center gap-1 text-xs">
              <Gauge className="w-3 h-3 text-primary" />
              <span>{formatMileage(carAttrs.mileage)}</span>
            </div>
          )}
          {carAttrs.fuelType && (
            <div className="flex items-center gap-1 text-xs">
              <Fuel className="w-3 h-3 text-primary" />
              <span>{formatFuel(carAttrs.fuelType)}</span>
            </div>
          )}
          {carAttrs.transmission && (
            <div className="flex items-center gap-1 text-xs">
              <Car className="w-3 h-3 text-primary" />
              <span>{formatTransmission(carAttrs.transmission)}</span>
            </div>
          )}
        </div>

        <div className="card-actions justify-between items-center mt-auto">
          <div>
            {listing.price && (
              <div className="text-lg font-bold text-primary">
                €{listing.price.toLocaleString()}
              </div>
            )}
            {(listing.city || listing.country) && (
              <div className="flex items-center gap-1 text-xs text-base-content/60">
                <MapPin className="w-3 h-3" />
                <span>{listing.city || listing.country}</span>
              </div>
            )}
          </div>
          <div className="flex items-center gap-1 text-xs text-base-content/60">
            <Eye className="w-3 h-3" />
            <span>{listing.views_count || 0}</span>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default CarListingCard;
