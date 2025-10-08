'use client';

import React from 'react';
import type { components } from '@/types/generated/api';
import CarListingCard from './CarListingCard';
import Link from 'next/link';
import { OptimizedCarImage } from '@/components/common/OptimizedCarImage';

type C2CListing = components['schemas']['models.C2CListing'];

interface VirtualizedCarListProps {
  listings: C2CListing[];
  locale: string;
  isGrid?: boolean;
  onFavorite?: (listingId: number) => void;
  onShare?: (listingId: number) => void;
  itemsPerRow?: number;
}

/**
 * Компонент виртуализации для длинных списков автомобилей
 * Упрощенная версия без react-window для продакшн сборки
 */
export const VirtualizedCarList: React.FC<VirtualizedCarListProps> = ({
  listings,
  locale,
  isGrid = true,
  onFavorite,
  onShare,
}) => {
  // Простое отображение без виртуализации для production сборки
  return (
    <div
      className={
        isGrid
          ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
          : 'space-y-4'
      }
    >
      {listings.map((listing) => (
        <CarListingCard
          key={listing.id}
          listing={listing}
          locale={locale}
          onFavorite={onFavorite}
          onShare={onShare}
          isGrid={isGrid}
        />
      ))}
    </div>
  );
};

/**
 * Оптимизированная карточка для грид-режима с ленивой загрузкой изображений
 */
const _VirtualizedGridCard: React.FC<{
  listing: C2CListing;
  locale: string;
  onFavorite?: (listingId: number) => void;
  onShare?: (listingId: number) => void;
}> = ({ listing, locale, onFavorite }) => {
  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (listing.id !== undefined) {
      onFavorite?.(listing.id);
    }
  };

  return (
    <Link
      href={`/${locale}/listing/${listing.id || 0}`}
      className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow h-full"
    >
      <figure className="aspect-[4/3] relative">
        <OptimizedCarImage
          src={listing.images?.[0]?.thumbnail_url}
          alt={listing.title || ''}
        />
        <div className="absolute top-2 right-2">
          <button
            onClick={handleFavoriteClick}
            className="btn btn-ghost btn-sm btn-circle bg-base-100/80 backdrop-blur"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              />
            </svg>
          </button>
        </div>
      </figure>
      <div className="card-body p-4">
        <h3 className="card-title text-base line-clamp-1">{listing.title}</h3>
        {listing.price && (
          <p className="text-lg font-bold text-primary">
            €{listing.price.toLocaleString()}
          </p>
        )}
        <p className="text-sm text-base-content/60">
          {listing.city || listing.country}
        </p>
      </div>
    </Link>
  );
};

/**
 * Hook для определения необходимости виртуализации
 */
export const useVirtualization = (itemCount: number, threshold = 50) => {
  return {
    shouldVirtualize: itemCount > threshold,
    itemCount,
    threshold,
  };
};

export default VirtualizedCarList;
