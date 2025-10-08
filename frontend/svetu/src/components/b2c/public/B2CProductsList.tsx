'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { storefrontProductsService } from '@/services/b2cProducts';
import { storefrontApi } from '@/services/b2cStoreApi';
import { UnifiedProductCard } from '@/components/common/UnifiedProductCard';
import { adaptB2CProduct } from '@/utils/product-adapters';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import type { components } from '@/types/generated/api';
import type { UnifiedStorefrontInfo } from '@/types/unified-product';

type B2CProduct = components['schemas']['models.B2CProduct'];

interface B2CProductsListProps {
  storefrontId: number;
  showTitle?: boolean;
  maxItems?: number;
}

export default function B2CProductsList({ 
  storefrontId, 
  showTitle = true,
  maxItems
}: B2CProductsListProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  
  const [products, setProducts] = useState<B2CProduct[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [_totalPages, setTotalPages] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [storefrontInfo, setStorefrontInfo] = useState<UnifiedStorefrontInfo | null>(null);

  const ITEMS_PER_PAGE = maxItems || 12;

  const fetchStorefrontInfo = useCallback(async () => {
    try {
      const storefront = await storefrontApi.getStorefrontById(storefrontId);
      setStorefrontInfo({
        id: storefront.id!,
        name: storefront.name!,
        slug: storefront.slug!,
      });
    } catch (err) {
      console.error('Error fetching storefront info:', err);
      // Используем fallback значения если не получилось загрузить
      setStorefrontInfo({
        id: storefrontId,
        name: 'Store',
        slug: 'store',
      });
    }
  }, [storefrontId]);

  const fetchProducts = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await storefrontProductsService.getProducts(String(storefrontId), {
        limit: ITEMS_PER_PAGE,
        offset: (currentPage - 1) * ITEMS_PER_PAGE,
        in_stock_only: false, // Показываем все товары, включая out_of_stock
      });

      if (currentPage === 1) {
        setProducts(response.products);
      } else {
        setProducts(prev => [...prev, ...response.products]);
      }

      const calculatedTotalPages = Math.ceil(response.total / ITEMS_PER_PAGE);
      setTotalPages(calculatedTotalPages);
      setHasMore(currentPage < calculatedTotalPages);
    } catch (err) {
      console.error('Error fetching products:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch products');
      setProducts([]);
    } finally {
      setIsLoading(false);
    }
  }, [storefrontId, currentPage, ITEMS_PER_PAGE]);

  useEffect(() => {
    fetchStorefrontInfo();
  }, [fetchStorefrontInfo]);

  useEffect(() => {
    if (storefrontInfo) {
      fetchProducts();
    }
  }, [fetchProducts, storefrontInfo]);

  const handleLoadMore = useCallback(() => {
    if (hasMore && !isLoading && !maxItems) {
      setCurrentPage(prev => prev + 1);
    }
  }, [hasMore, isLoading, maxItems]);

  const loadMoreRef = useInfiniteScroll({
    loading: isLoading,
    hasMore: hasMore && !maxItems,
    onLoadMore: handleLoadMore,
  });

  if ((isLoading && currentPage === 1) || !storefrontInfo) {
    return (
      <div className="space-y-4">
        {showTitle && (
          <div className="flex justify-between items-center">
            <h2 className="text-2xl font-bold">{t('products.title')}</h2>
            <div className="skeleton h-8 w-24"></div>
          </div>
        )}
        
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="card bg-base-300 animate-pulse">
              <div className="aspect-[4/3] bg-base-content/10"></div>
              <div className="card-body">
                <div className="h-4 bg-base-content/10 rounded w-3/4"></div>
                <div className="h-6 bg-base-content/10 rounded w-1/2 mt-2"></div>
                <div className="h-8 bg-base-content/10 rounded w-full mt-4"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        {showTitle && <h2 className="text-2xl font-bold">{t('products.title')}</h2>}
        <div className="alert alert-error">
          <span>{error}</span>
          <button 
            onClick={() => {
              setCurrentPage(1);
              fetchProducts();
            }}
            className="btn btn-sm"
          >
            {tCommon('tryAgain')}
          </button>
        </div>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="space-y-4">
        {showTitle && <h2 className="text-2xl font-bold">{t('products.title')}</h2>}
        <div className="text-center py-12">
          <svg className="w-16 h-16 mx-auto mb-4 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <p className="text-base-content/60">{t('products.noProducts')}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      {showTitle && (
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold">{t('products.title')}</h2>
          <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
        </div>
      )}

      {/* Products Grid/List */}
      <div className={
        viewMode === 'grid' 
          ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
          : "space-y-4"
      }>
        {products.map((product) => {
          if (!storefrontInfo) return null;
          
          const unifiedProduct = adaptB2CProduct(product, storefrontInfo);
          return (
            <UnifiedProductCard
              key={product.id}
              product={unifiedProduct}
              locale="en" // TODO: получать локаль из контекста
              viewMode={viewMode}
            />
          );
        })}
      </div>

      {/* Load More */}
      {!maxItems && (
        <InfiniteScrollTrigger
          ref={loadMoreRef}
          loading={isLoading}
          hasMore={hasMore}
          onLoadMore={handleLoadMore}
          showButton={true}
          loadMoreText={tCommon('loadMore')}
        />
      )}
    </div>
  );
}