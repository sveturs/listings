'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useLocale } from 'next-intl';
import { apiClient } from '@/services/api-client';
import config from '@/config';
import type { MarketplaceItem } from '@/types/marketplace';

interface SimilarListingsProps {
  listingId: number;
}

export default function SimilarListings({ listingId }: SimilarListingsProps) {
  const locale = useLocale();
  const [allListings, setAllListings] = useState<MarketplaceItem[]>([]);
  const [displayedListings, setDisplayedListings] = useState<MarketplaceItem[]>(
    []
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [displayCount, setDisplayCount] = useState(20);
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadMoreRef = useRef<HTMLDivElement>(null);

  const ITEMS_PER_LOAD = 20;

  useEffect(() => {
    const fetchSimilarListings = async () => {
      try {
        setLoading(true);
        // Загружаем больше объявлений для постепенного показа
        const response = await apiClient.get<{
          data: MarketplaceItem[];
          success: boolean;
        }>(`/api/v1/marketplace/listings/${listingId}/similar?limit=100`);

        // API может вернуть либо массив напрямую, либо объект с полем data
        const items = Array.isArray(response.data)
          ? response.data
          : response.data?.data || [];

        if (items.length > 0) {
          // Фильтруем текущее объявление из списка похожих
          const filteredListings = items.filter(
            (item: MarketplaceItem) => item.id !== listingId
          );
          setAllListings(filteredListings);
          // Показываем первые 20 объявлений
          setDisplayedListings(filteredListings.slice(0, ITEMS_PER_LOAD));
        }
      } catch (err) {
        console.error('Error fetching similar listings:', err);
        setError('Failed to load similar listings');
      } finally {
        setLoading(false);
      }
    };

    fetchSimilarListings();
  }, [listingId]);

  // Обновляем отображаемые объявления при изменении displayCount
  useEffect(() => {
    setDisplayedListings(allListings.slice(0, displayCount));
  }, [displayCount, allListings]);

  const loadMore = useCallback(() => {
    if (displayCount < allListings.length) {
      setDisplayCount((prev) =>
        Math.min(prev + ITEMS_PER_LOAD, allListings.length)
      );
    }
  }, [displayCount, allListings.length]);

  // Настраиваем IntersectionObserver для автоматической подгрузки
  useEffect(() => {
    if (loading) return;

    const options = {
      root: null,
      rootMargin: '100px',
      threshold: 0.1,
    };

    observerRef.current = new IntersectionObserver((entries) => {
      const [entry] = entries;
      if (entry.isIntersecting && displayCount < allListings.length) {
        loadMore();
      }
    }, options);

    if (loadMoreRef.current) {
      observerRef.current.observe(loadMoreRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [loading, loadMore, displayCount, allListings.length]);

  if (loading) {
    return (
      <div className="mt-12">
        <h2 className="text-2xl font-bold mb-6">
          {locale === 'ru' ? 'Похожие объявления' : 'Similar listings'}
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="animate-pulse">
              <div className="bg-base-300 rounded-lg aspect-square mb-2"></div>
              <div className="h-4 bg-base-300 rounded w-3/4 mb-2"></div>
              <div className="h-4 bg-base-300 rounded w-1/2"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error || allListings.length === 0) {
    return null;
  }

  return (
    <div className="mt-12">
      <h2 className="text-2xl font-bold mb-6">
        {locale === 'ru' ? 'Похожие объявления' : 'Similar listings'}
      </h2>

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {displayedListings.map((listing) => (
          <Link
            key={listing.id}
            href={`/${locale}/marketplace/${listing.id}`}
            className="group hover:shadow-lg transition-all duration-200"
          >
            <div className="relative aspect-square rounded-lg overflow-hidden bg-base-200 mb-2">
              {listing.images && listing.images.length > 0 ? (
                <Image
                  src={config.buildImageUrl(listing.images[0].public_url)}
                  alt={listing.title}
                  fill
                  className="object-cover group-hover:scale-105 transition-transform duration-200"
                  sizes="(max-width: 640px) 50vw, (max-width: 768px) 33vw, 20vw"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center">
                  <svg
                    className="w-12 h-12 text-base-content/20"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                  </svg>
                </div>
              )}

              {/* Price badge */}
              <div className="absolute bottom-2 left-2 bg-base-100/90 backdrop-blur-sm px-2 py-1 rounded text-sm font-semibold">
                {listing.price} $
              </div>
            </div>

            <h3 className="font-medium text-sm line-clamp-2 mb-1 group-hover:text-primary transition-colors">
              {listing.title}
            </h3>

            <div className="flex items-center gap-2 text-xs text-base-content/70">
              {listing.location && (
                <>
                  <svg
                    className="w-3 h-3"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
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
                  <span className="truncate">{listing.location}</span>
                </>
              )}
            </div>
          </Link>
        ))}
      </div>

      {/* Элемент для отслеживания скролла */}
      {displayCount < allListings.length && (
        <div ref={loadMoreRef} className="flex justify-center mt-8">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      )}
    </div>
  );
}
