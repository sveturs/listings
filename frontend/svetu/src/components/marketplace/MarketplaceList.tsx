'use client';

import { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  UnifiedSearchService,
  UnifiedSearchItem,
} from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';
import MarketplaceCard from '@/components/MarketplaceCard';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';

// Адаптер для преобразования UnifiedSearchItem в MarketplaceItem
function convertToMarketplaceItem(
  unifiedItem: UnifiedSearchItem
): MarketplaceItem {
  return {
    id: unifiedItem.product_id,
    title: unifiedItem.name,
    description: unifiedItem.description,
    price: unifiedItem.price,
    images: unifiedItem.images.map((img) => ({
      id: 0,
      public_url: img.url,
      is_main: img.is_main,
    })),
    user: {
      id: 0,
      name: unifiedItem.storefront?.name || 'Пользователь',
      email: '',
      picture_url: '',
    },
    category: {
      id: unifiedItem.category.id,
      name: unifiedItem.category.name,
      slug: unifiedItem.category.slug || '',
    },
    city: unifiedItem.location?.city || '',
    country: unifiedItem.location?.country || '',
    location: unifiedItem.location?.city || '',
    created_at: unifiedItem.created_at || '',
    updated_at: '',
    user_id: 0,
    // Добавляем поля для определения типа товара и витрины
    product_type: unifiedItem.product_type,
    storefront_id: unifiedItem.storefront?.id,
  };
}

interface MarketplaceListProps {
  initialData: {
    items: UnifiedSearchItem[];
    total: number;
    page: number;
    limit: number;
    has_more: boolean;
  } | null;
  locale: string;
  productTypes?: ('marketplace' | 'storefront')[];
}

export default function MarketplaceList({
  initialData,
  locale,
  productTypes = ['marketplace', 'storefront'],
}: MarketplaceListProps) {
  // console.log('MarketplaceList render:', {
  //   initialData: !!initialData,
  //   locale,
  // });
  const t = useTranslations('home');
  const [items, setItems] = useState<UnifiedSearchItem[]>(
    initialData?.items || []
  );
  const [page, setPage] = useState(initialData?.page || 1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(initialData?.has_more ?? true);
  const [error, setError] = useState<string | null>(null);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [initialized, setInitialized] = useState(!!initialData);

  // console.log('MarketplaceList state:', {
  //   itemsCount: items.length,
  //   loading,
  //   hasMore,
  //   initialized,
  // });

  const manualLoad = useCallback(() => {
    // console.log('Manual load clicked!');
    setLoading(true);
    UnifiedSearchService.search({
      query: '',
      product_types: productTypes,
      sort_by: 'date',
      sort_order: 'desc',
      page: 1,
      limit: 20,
    })
      .then((response) => {
        // console.log('UnifiedSearchService.search response:', response);
        if (response && response.items) {
          setItems(response.items);
          setHasMore(response.has_more);
          setPage(response.page);
          setInitialized(true);
          // console.log('Loaded items:', response.items.length);
        }
      })
      .catch((err) => {
        console.error('Failed to load initial data:', err);
        setError(t('errorLoadingData'));
      })
      .finally(() => {
        setLoading(false);
      });
  }, [t, productTypes]);

  // Начальная загрузка данных, если они не были переданы через SSR
  useEffect(() => {
    if (!initialized && !loading) {
      manualLoad();
    }
  }, [initialized, loading, manualLoad]);

  // Перезагрузка данных при изменении типов товаров
  useEffect(() => {
    if (initialized) {
      setPage(1);
      setItems([]);
      setHasMore(true);
      manualLoad();
    }
  }, [productTypes, initialized, manualLoad]);

  const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    setError(null);

    try {
      const nextPage = page + 1;
      const response = await UnifiedSearchService.search({
        query: '',
        product_types: productTypes,
        sort_by: 'date',
        sort_order: 'desc',
        page: nextPage,
        limit: 20,
      });

      if (response && response.items) {
        setItems((prev) => {
          // Создаем Set с существующими ID для быстрой проверки
          const existingIds = new Set(prev.map((item) => item.id));
          // Фильтруем только новые объявления
          const newItems = response.items.filter(
            (item) => !existingIds.has(item.id)
          );
          return [...prev, ...newItems];
        });

        setPage(nextPage);
        setHasMore(response.has_more);
      }
    } catch (err) {
      console.error('Failed to load more items:', err);
      setError(t('errorLoadingData'));
    } finally {
      setLoading(false);
    }
  }, [loading, hasMore, page, t, productTypes]);

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore,
    onLoadMore: loadMore,
  });

  // Показываем состояние загрузки для начальной загрузки
  if (loading && !initialized && items.length === 0) {
    return (
      <div className="text-center py-16">
        <div className="loading loading-spinner loading-lg mb-4"></div>
        <p className="text-base-content/70">{t('loading')}</p>
      </div>
    );
  }

  // Показываем состояние "нет данных" только если мы уже загрузились и нет элементов
  if (initialized && items.length === 0 && !loading && !error) {
    return (
      <div className="text-center py-16">
        <div className="inline-flex items-center justify-center w-24 h-24 rounded-full bg-base-200 mb-6">
          <svg
            className="w-12 h-12 text-base-content/40"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
            />
          </svg>
        </div>
        <p className="text-xl text-base-content/70 mb-4">{t('noItems')}</p>
        <p className="text-base-content/50">
          Начните с создания первого объявления
        </p>
      </div>
    );
  }

  return (
    <>
      <div className="flex justify-between mb-4">
        <button onClick={manualLoad} className="btn btn-primary">
          Manual Load ({items.length} items)
        </button>
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
            item={convertToMarketplaceItem(item)}
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
        loadMoreText={t('loadMore')}
      />
    </>
  );
}
