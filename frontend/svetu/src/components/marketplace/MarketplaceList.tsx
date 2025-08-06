'use client';

import { useState, useCallback, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
  UnifiedSearchService,
  UnifiedSearchItem,
} from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';
import { UnifiedProductCard } from '@/components/common/UnifiedProductCard';
import { adaptMarketplaceItem } from '@/utils/product-adapters';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import GridColumnsToggle from '@/components/common/GridColumnsToggle';
import { useGridColumns } from '@/hooks/useGridColumns';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import { FilterValues } from './ListingFilters';
import { ListingGridSkeleton } from '@/components/ui/skeletons';

// Функция для конвертации фильтров в формат API
function convertFiltersToSearchParams(filters: FilterValues) {
  const searchParams: Record<string, any> = {};

  // Постоянные фильтры
  if (filters.priceMin) {
    searchParams.priceMin = filters.priceMin;
  }
  if (filters.priceMax) {
    searchParams.priceMax = filters.priceMax;
  }
  if (filters.condition) {
    searchParams.condition = filters.condition;
  }
  if (filters.sellerType) {
    // Для приватных продавцов не указываем storefrontID
    // Для компаний указываем, что нужны только товары из storefront
    searchParams.storefrontID =
      filters.sellerType === 'company' ? 'not_null' : null;
  }

  // Атрибуты категории
  if (
    filters.attributeFilters &&
    Object.keys(filters.attributeFilters).length > 0
  ) {
    searchParams.attributeFilters = filters.attributeFilters;
  }

  return searchParams;
}

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
    created_at: (unifiedItem as any).created_at || '',
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
  selectedCategoryId?: number | null;
  filters?: FilterValues;
}

export default function MarketplaceList({
  initialData,
  locale,
  productTypes = ['marketplace', 'storefront'],
  selectedCategoryId,
  filters = {},
}: MarketplaceListProps) {
  // console.log('MarketplaceList render:', {
  //   initialData: !!initialData,
  //   locale,
  // });
  const t = useTranslations('marketplace.home');
  const [items, setItems] = useState<UnifiedSearchItem[]>(
    initialData?.items || []
  );
  const [page, setPage] = useState(initialData?.page || 1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(initialData?.has_more ?? true);
  const [error, setError] = useState<string | null>(null);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [gridColumns, setGridColumns] = useGridColumns(1);
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
    setPage(1);
    setItems([]);
    setHasMore(true);

    const searchParams = convertFiltersToSearchParams(filters);

    UnifiedSearchService.search({
      query: '',
      product_types: productTypes,
      category_id: selectedCategoryId?.toString(),
      sort_by: 'date',
      sort_order: 'desc',
      page: 1,
      limit: 20,
      ...searchParams,
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
  }, [t, productTypes, selectedCategoryId, filters]);

  // Начальная загрузка данных, если они не были переданы через SSR
  useEffect(() => {
    if (!initialized && !loading) {
      manualLoad();
    }
  }, [initialized, loading, manualLoad]);

  // Создаем стабильный ключ для фильтров, чтобы избежать лишних рендеров
  const filtersKey = useMemo(() => {
    return JSON.stringify({
      priceMin: filters.priceMin,
      priceMax: filters.priceMax,
      condition: filters.condition,
      sellerType: filters.sellerType,
      hasDelivery: filters.hasDelivery,
      attributeFilters: filters.attributeFilters,
    });
  }, [filters]);

  // Перезагрузка данных при изменении типов товаров или категории
  useEffect(() => {
    if (initialized) {
      setPage(1);
      setItems([]);
      setHasMore(true);
      setLoading(true);

      const searchParams = convertFiltersToSearchParams(filters);

      UnifiedSearchService.search({
        query: '',
        product_types: productTypes,
        category_id: selectedCategoryId?.toString(),
        sort_by: 'date',
        sort_order: 'desc',
        page: 1,
        limit: 20,
        ...searchParams,
      })
        .then((response) => {
          if (response && response.items) {
            setItems(response.items);
            setHasMore(response.has_more);
            setPage(response.page);
          }
        })
        .catch((err) => {
          console.error('Failed to reload data:', err);
          setError(t('errorLoadingData'));
        })
        .finally(() => {
          setLoading(false);
        });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [productTypes, selectedCategoryId, filtersKey, initialized, t]);

  const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    setError(null);

    try {
      const nextPage = page + 1;
      const searchParams = convertFiltersToSearchParams(filters);

      const response = await UnifiedSearchService.search({
        query: '',
        product_types: productTypes,
        category_id: selectedCategoryId?.toString(),
        sort_by: 'date',
        sort_order: 'desc',
        page: nextPage,
        limit: 20,
        ...searchParams,
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
  }, [loading, hasMore, page, t, productTypes, selectedCategoryId, filters]);

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore,
    onLoadMore: loadMore,
  });

  // Показываем скелетоны для начальной загрузки
  if (loading && !initialized && items.length === 0) {
    return (
      <>
        <div className="flex justify-between mb-4">
          <div className="h-10 w-32 bg-base-300 rounded animate-pulse" />
          <div className="h-10 w-24 bg-base-300 rounded animate-pulse" />
        </div>
        <ListingGridSkeleton count={8} viewMode={viewMode} />
      </>
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
      <div className="flex justify-end mb-4">
        <div className="flex gap-2">
          {viewMode === 'grid' && (
            <div className="lg:hidden">
              <GridColumnsToggle
                currentColumns={gridColumns}
                onColumnsChange={setGridColumns}
              />
            </div>
          )}
          <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
        </div>
      </div>

      <div
        className={
          viewMode === 'grid'
            ? gridColumns === 1
              ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2 sm:gap-4 lg:gap-6'
              : gridColumns === 2
                ? 'grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-1.5 sm:gap-3'
                : 'grid grid-cols-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-1 sm:gap-2'
            : 'space-y-4'
        }
      >
        {items.map((item) => {
          const marketplaceItem = convertToMarketplaceItem(item);
          const unifiedProduct = adaptMarketplaceItem(marketplaceItem);
          return (
            <UnifiedProductCard
              key={item.id}
              product={unifiedProduct}
              locale={locale}
              viewMode={viewMode}
              gridColumns={viewMode === 'grid' ? gridColumns : undefined}
            />
          );
        })}
      </div>

      {/* Показываем скелетоны при загрузке следующей страницы */}
      {loading && initialized && items.length > 0 && (
        <div className="mt-6">
          <ListingGridSkeleton count={4} viewMode={viewMode} />
        </div>
      )}

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
