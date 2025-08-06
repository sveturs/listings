'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import SafeImage from '@/components/SafeImage';
import { useLocale, useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import config from '@/config';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import { useAuth } from '@/contexts/AuthContext';
import type { MarketplaceItem } from '@/types/marketplace';
import {
  formatAddressWithPrivacy,
  getFullLocalizedAddress,
  type LocationPrivacyLevel,
} from '@/utils/addressUtils';

interface SimilarListingsProps {
  listingId: number;
}

export default function SimilarListings({ listingId }: SimilarListingsProps) {
  const locale = useLocale();
  const tCommon = useTranslations('common');
  const t = useTranslations('marketplace');
  const router = useRouter();
  const { isAuthenticated, user } = useAuth();
  const [allListings, setAllListings] = useState<MarketplaceItem[]>([]);
  const [displayedListings, setDisplayedListings] = useState<MarketplaceItem[]>(
    []
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [displayCount, setDisplayCount] = useState(20);
  const [viewMode, setViewMode] = useViewPreference('grid');

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

        // Отладка API ответа для похожих объявлений
        console.log('SimilarListings API Response:', {
          totalItems: items.length,
          firstItem: items[0],
          sampleImagesData: items.slice(0, 3).map((item: any) => ({
            id: item.id,
            title: item.title,
            images: item.images,
            hasImages: !!(item.images && item.images.length > 0),
          })),
        });

        if (items.length > 0) {
          // Фильтруем текущее объявление из списка похожих
          const filteredListings = items.filter(
            (item: MarketplaceItem) => item.id !== listingId
          );

          console.log('SimilarListings Filtered:', {
            totalFiltered: filteredListings.length,
            withImages: filteredListings.filter(
              (item) => item.images && item.images.length > 0
            ).length,
            firstWithImages: filteredListings.find(
              (item) => item.images && item.images.length > 0
            ),
          });

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

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore: displayCount < allListings.length,
    onLoadMore: loadMore,
  });

  const handleChatClick = (e: React.MouseEvent, listing: MarketplaceItem) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      return;
    }

    if (listing.user_id === user?.id) {
      return;
    }

    router.push(
      `/${locale}/chat?listing_id=${listing.id}&seller_id=${listing.user_id}`
    );
  };

  const handleBuyClick = (e: React.MouseEvent, listing: MarketplaceItem) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(window.location.pathname)}`
      );
      return;
    }

    if (listing.user_id === user?.id) {
      return;
    }

    router.push(`/${locale}/marketplace/${listing.id}/buy`);
  };

  const isOnlinePurchaseAvailable = (listing: MarketplaceItem) => {
    return listing.storefront_id != null && listing.storefront_id > 0;
  };

  if (loading) {
    return (
      <div className="mt-12">
        <h2 className="text-2xl font-bold mb-6">
          {t('similarListings.title')}
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
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">{t('similarListings.title')}</h2>
        <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
      </div>

      <div
        className={
          viewMode === 'grid'
            ? 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4'
            : 'grid grid-cols-1 md:grid-cols-2 gap-4'
        }
      >
        {displayedListings.map((listing) =>
          viewMode === 'grid' ? (
            <div
              key={listing.id}
              className="group hover:shadow-lg transition-all duration-200"
            >
              <Link
                href={`/${locale}/marketplace/${listing.id}`}
                className="block"
              >
                <div className="relative aspect-square rounded-lg overflow-hidden bg-base-200 mb-2">
                  <SafeImage
                    src={(() => {
                      const imageUrl =
                        listing.images && listing.images.length > 0
                          ? config.buildImageUrl(listing.images[0].public_url)
                          : null;

                      // Отладка изображений в похожих объявлениях
                      if ((window as any).SimilarListingsImageDebug) {
                        console.log('SimilarListings SafeImage src debug:', {
                          listingId: listing.id,
                          title: listing.title,
                          hasImages: !!(
                            listing.images && listing.images.length > 0
                          ),
                          rawImageUrl: listing.images?.[0]?.public_url,
                          builtImageUrl: imageUrl,
                          imagesArray: listing.images,
                        });
                      }

                      return imageUrl;
                    })()}
                    alt={listing.title}
                    fill
                    className="object-cover group-hover:scale-105 transition-transform duration-200"
                    sizes="(max-width: 640px) 50vw, (max-width: 768px) 33vw, 20vw"
                    fallback={
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
                    }
                  />

                  {/* Price badge */}
                  <div className="absolute bottom-2 left-2 bg-base-100/90 backdrop-blur-sm px-2 py-1 rounded text-sm font-semibold">
                    {listing.price} $
                  </div>
                </div>

                <h3 className="font-medium text-sm line-clamp-2 mb-1 group-hover:text-primary transition-colors">
                  {listing.title}
                </h3>

                <div className="flex items-center gap-2 text-xs text-base-content/70 mb-2">
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
                      <span className="truncate">
                        {formatAddressWithPrivacy(
                          getFullLocalizedAddress(listing, locale),
                          listing.location_privacy as LocationPrivacyLevel
                        )}
                      </span>
                    </>
                  )}
                </div>
              </Link>

              {/* Action buttons */}
              <div className="flex gap-1 mt-2">
                {isOnlinePurchaseAvailable(listing) && (
                  <button
                    onClick={(e) => handleBuyClick(e, listing)}
                    className="btn btn-primary btn-sm flex-1 text-xs"
                  >
                    💳 {t('similarListings.buy')}
                  </button>
                )}
                <button
                  onClick={(e) => handleChatClick(e, listing)}
                  className="btn btn-outline btn-sm flex-1 text-xs"
                >
                  💬 {t('similarListings.chat')}
                </button>
              </div>
            </div>
          ) : (
            // List view
            <div
              key={listing.id}
              className="card bg-base-100 shadow-md hover:shadow-lg transition-shadow"
            >
              <div className="card-body p-4">
                <div className="flex gap-4">
                  <Link
                    href={`/${locale}/marketplace/${listing.id}`}
                    className="relative w-24 h-24 flex-shrink-0 rounded-lg overflow-hidden bg-base-200"
                  >
                    <SafeImage
                      src={(() => {
                        const imageUrl =
                          listing.images && listing.images.length > 0
                            ? config.buildImageUrl(listing.images[0].public_url)
                            : null;

                        // Отладка изображений в похожих объявлениях (List view)
                        if ((window as any).SimilarListingsImageDebug) {
                          console.log(
                            'SimilarListings SafeImage src debug (List view):',
                            {
                              listingId: listing.id,
                              title: listing.title,
                              hasImages: !!(
                                listing.images && listing.images.length > 0
                              ),
                              rawImageUrl: listing.images?.[0]?.public_url,
                              builtImageUrl: imageUrl,
                              imagesArray: listing.images,
                            }
                          );
                        }

                        return imageUrl;
                      })()}
                      alt={listing.title}
                      fill
                      className="object-cover"
                      sizes="96px"
                      fallback={
                        <div className="w-full h-full flex items-center justify-center">
                          <svg
                            className="w-8 h-8 text-base-content/20"
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
                      }
                    />
                  </Link>
                  <div className="flex-grow min-w-0">
                    <Link href={`/${locale}/marketplace/${listing.id}`}>
                      <h3 className="font-semibold line-clamp-1 mb-1 hover:text-primary transition-colors">
                        {listing.title}
                      </h3>
                    </Link>
                    <p className="text-lg font-bold text-primary mb-1">
                      {listing.price} $
                    </p>
                    {listing.location && (
                      <p className="text-sm text-base-content/60 flex items-center gap-1 mb-2">
                        <svg
                          className="w-3 h-3 flex-shrink-0"
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
                        <span className="truncate">
                          {formatAddressWithPrivacy(
                            getFullLocalizedAddress(listing, locale),
                            listing.location_privacy as LocationPrivacyLevel
                          )}
                        </span>
                      </p>
                    )}
                    {/* Action buttons */}
                    <div className="flex gap-2 mt-2">
                      {isOnlinePurchaseAvailable(listing) && (
                        <button
                          onClick={(e) => handleBuyClick(e, listing)}
                          className="btn btn-primary btn-sm"
                        >
                          💳 {t('similarListings.buy')}
                        </button>
                      )}
                      <button
                        onClick={(e) => handleChatClick(e, listing)}
                        className="btn btn-outline btn-sm"
                      >
                        💬 {t('similarListings.chat')}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )
        )}
      </div>

      <InfiniteScrollTrigger
        ref={loadMoreRef}
        loading={loading}
        hasMore={displayCount < allListings.length}
        onLoadMore={loadMore}
        showButton={false}
        loadMoreText={tCommon('loadMore')}
      />
    </div>
  );
}
