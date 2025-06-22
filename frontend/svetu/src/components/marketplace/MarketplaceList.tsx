'use client';

import { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { MarketplaceService } from '@/services/marketplace';
import { MarketplaceItem } from '@/types/marketplace';
import MarketplaceCard from '@/components/MarketplaceCard';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';

interface MarketplaceListProps {
  initialData: {
    data: MarketplaceItem[];
    meta: {
      total: number;
      page: number;
      size: number;
      has_more: boolean;
    };
  };
  locale: string;
}

export default function MarketplaceList({
  initialData,
  locale,
}: MarketplaceListProps) {
  const t = useTranslations('home');
  const [items, setItems] = useState<MarketplaceItem[]>(initialData.data);
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(initialData.meta.has_more);
  const [error, setError] = useState<string | null>(null);
  const [viewMode, setViewMode] = useViewPreference('grid');

  const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    setError(null);

    try {
      const nextPage = page + 1;
      const response = await MarketplaceService.search({
        sort_by: 'date_desc',
        page: nextPage,
        size: 20,
      });

      if (response && response.data) {
        setItems((prev) => {
          // Создаем Set с существующими ID для быстрой проверки
          const existingIds = new Set(prev.map((item) => item.id));
          // Фильтруем только новые объявления
          const newItems = response.data.filter(
            (item) => !existingIds.has(item.id)
          );
          return [...prev, ...newItems];
        });

        setPage(nextPage);
        setHasMore(response.meta.has_more);
      }
    } catch (err) {
      console.error('Failed to load more items:', err);
      setError(t('errorLoadingData'));
    } finally {
      setLoading(false);
    }
  }, [loading, hasMore, page, t]);

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore,
    onLoadMore: loadMore,
  });

  return (
    <>
      <div className="flex justify-end mb-4">
        <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
      </div>

      <div
        className={
          viewMode === 'grid'
            ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'
            : 'space-y-4'
        }
      >
        {items.map((item) => (
          <MarketplaceCard
            key={item.id}
            item={item}
            locale={locale}
            viewMode={viewMode}
          />
        ))}
      </div>

      {error && (
        <div className="alert alert-error mt-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{error}</span>
        </div>
      )}

      <InfiniteScrollTrigger
        ref={loadMoreRef}
        loading={loading}
        hasMore={hasMore}
        onLoadMore={loadMore}
      />
    </>
  );
}
