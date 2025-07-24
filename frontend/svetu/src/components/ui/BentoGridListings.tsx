'use client';

import React, { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { ShoppingBag, TrendingUp, ArrowRight, LayoutGrid, List } from 'lucide-react';
import { UnifiedSearchService, UnifiedSearchItem } from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';
import { EnhancedListingCard } from '@/components/marketplace/EnhancedListingCard';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import { ListingGridSkeleton } from '@/components/ui/skeletons';

// Адаптер для преобразования UnifiedSearchItem в MarketplaceItem
function convertToMarketplaceItem(unifiedItem: UnifiedSearchItem): MarketplaceItem {
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
    product_type: unifiedItem.product_type,
    storefront_id: unifiedItem.storefront?.id,
  };
}

interface BentoGridListingsProps {
  locale: string;
  productTypes?: ('marketplace' | 'storefront')[];
  selectedCategoryId?: number | null;
  filters?: Record<string, any>;
}

export const BentoGridListings: React.FC<BentoGridListingsProps> = ({
  locale,
  productTypes = ['marketplace', 'storefront'],
  selectedCategoryId,
  filters = {},
}) => {
  const t = useTranslations('home');
  const [items, setItems] = useState<UnifiedSearchItem[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(true);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [showAll, setShowAll] = useState(false);

  // Конвертация фильтров в формат API
  const convertFiltersToSearchParams = (filters: Record<string, any>) => {
    const searchParams: Record<string, any> = {};
    
    if (filters.priceMin) searchParams.priceMin = filters.priceMin;
    if (filters.priceMax) searchParams.priceMax = filters.priceMax;
    if (filters.condition) searchParams.condition = filters.condition;
    if (filters.sellerType) {
      searchParams.storefrontID = filters.sellerType === 'company' ? 'not_null' : null;
    }
    if (filters.attributeFilters && Object.keys(filters.attributeFilters).length > 0) {
      searchParams.attributeFilters = filters.attributeFilters;
    }
    
    return searchParams;
  };

  // Загрузка данных
  const loadData = useCallback(async (pageNum: number = 1) => {
    try {
      setLoading(true);
      const searchParams = convertFiltersToSearchParams(filters);
      
      const response = await UnifiedSearchService.search({
        query: '',
        product_types: productTypes,
        category_id: selectedCategoryId?.toString(),
        sort_by: 'date',
        sort_order: 'desc',
        page: pageNum,
        limit: 12,
        ...searchParams,
      });

      if (response && response.items) {
        if (pageNum === 1) {
          setItems(response.items);
        } else {
          setItems((prev) => {
            const existingIds = new Set(prev.map((item) => item.id));
            const newItems = response.items.filter(
              (item) => !existingIds.has(item.id)
            );
            return [...prev, ...newItems];
          });
        }
        setHasMore(response.has_more);
        setPage(pageNum);
      }
    } catch (err) {
      console.error('Failed to load listings:', err);
    } finally {
      setLoading(false);
    }
  }, [productTypes, selectedCategoryId, filters]);

  useEffect(() => {
    loadData(1);
  }, [productTypes, selectedCategoryId, filters]);

  const loadMore = useCallback(() => {
    if (!loading && hasMore) {
      loadData(page + 1);
    }
  }, [loading, hasMore, page, loadData]);

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore: hasMore && showAll,
    onLoadMore: loadMore,
  });

  // Показываем только первые 8 элементов в компактном режиме
  const displayItems = showAll ? items : items.slice(0, 8);

  return (
    <div className="col-span-3 row-span-3 bg-base-100 rounded-2xl shadow-xl p-6 overflow-hidden">
      {/* Заголовок */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-accent/10 rounded-xl">
            <ShoppingBag className="w-6 h-6 text-accent" />
          </div>
          <div>
            <h3 className="text-xl font-bold">{t('latestListings')}</h3>
            <p className="text-sm text-base-content/60">
              {selectedCategoryId ? 'В выбранной категории' : 'Все категории'}
            </p>
          </div>
        </div>

        {/* Переключатель вида */}
        <div className="flex items-center gap-2">
          <div className="join join-sm">
            <button
              className={`btn btn-sm join-item ${viewMode === 'grid' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('grid')}
            >
              <LayoutGrid className="w-4 h-4" />
            </button>
            <button
              className={`btn btn-sm join-item ${viewMode === 'list' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('list')}
            >
              <List className="w-4 h-4" />
            </button>
          </div>
          {!showAll && items.length > 8 && (
            <button
              onClick={() => setShowAll(true)}
              className="btn btn-primary btn-sm gap-1"
            >
              Все объявления
              <ArrowRight className="w-4 h-4" />
            </button>
          )}
        </div>
      </div>

      {/* Список объявлений */}
      <div className="overflow-y-auto max-h-[calc(100%-5rem)] pr-2 -mr-2">
        {loading && items.length === 0 ? (
          <ListingGridSkeleton count={6} viewMode={viewMode} />
        ) : items.length === 0 ? (
          <div className="text-center py-16">
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-base-200 mb-4">
              <ShoppingBag className="w-10 h-10 text-base-content/40" />
            </div>
            <p className="text-lg text-base-content/70">{t('noItems')}</p>
          </div>
        ) : (
          <>
            <div
              className={
                viewMode === 'grid'
                  ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'
                  : 'space-y-3'
              }
            >
              {displayItems.map((item) => (
                <EnhancedListingCard
                  key={item.id}
                  item={convertToMarketplaceItem(item)}
                  locale={locale}
                  viewMode={viewMode}
                />
              ))}
            </div>

            {/* Infinite scroll trigger */}
            {showAll && hasMore && (
              <div ref={loadMoreRef} className="h-20 flex items-center justify-center">
                {loading && (
                  <span className="loading loading-spinner loading-md"></span>
                )}
              </div>
            )}

            {/* Кнопка "Свернуть" */}
            {showAll && items.length > 8 && (
              <div className="mt-6 text-center">
                <button
                  onClick={() => setShowAll(false)}
                  className="btn btn-ghost btn-sm"
                >
                  Свернуть
                </button>
              </div>
            )}
          </>
        )}
      </div>

      {/* Статистика внизу */}
      {!loading && items.length > 0 && (
        <div className="mt-4 pt-4 border-t border-base-200 flex items-center justify-between text-sm text-base-content/60">
          <span>Показано {displayItems.length} из {items.length}</span>
          <div className="flex items-center gap-1">
            <TrendingUp className="w-4 h-4" />
            <span>Обновлено только что</span>
          </div>
        </div>
      )}
    </div>
  );
};