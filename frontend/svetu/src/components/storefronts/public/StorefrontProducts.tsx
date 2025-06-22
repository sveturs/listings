'use client';

import { useState, useEffect, useCallback } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import SafeImage from '@/components/SafeImage';
import { marketplaceApi } from '@/services/marketplaceApi';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import type { ListingCardData } from '@/types/marketplace';

interface StorefrontProductsProps {
  storefrontId: number;
}

export default function StorefrontProducts({ storefrontId }: StorefrontProductsProps) {
  const t = useTranslations();
  const locale = useLocale();
  const [products, setProducts] = useState<ListingCardData[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [_totalPages, setTotalPages] = useState(1);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [hasMore, setHasMore] = useState(false);

  const fetchProducts = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await marketplaceApi.getListings({
        storefront_id: storefrontId,
        limit: 12,
        offset: (currentPage - 1) * 12
      });
      
      // Безопасное извлечение данных из ответа API
      if (response && typeof response === 'object') {
        // API возвращает { data: { data: [...], success: true, meta: {...} } }
        const responseData = (response as any).data;
        const actualData = Array.isArray(responseData?.data) ? responseData.data : 
                          Array.isArray(responseData) ? responseData : [];
        
        if (currentPage === 1) {
          setProducts(actualData);
        } else {
          setProducts(prev => [...prev, ...actualData]);
        }
        
        // Безопасная проверка пагинации
        const pagination = responseData?.meta || (response as any).pagination;
        if (pagination && typeof pagination.total === 'number' && typeof pagination.limit === 'number' && pagination.limit > 0) {
          const calculatedTotalPages = Math.ceil(pagination.total / pagination.limit);
          setTotalPages(calculatedTotalPages);
          setHasMore(currentPage < calculatedTotalPages);
        }
      } else {
        setProducts([]);
      }
    } catch (error) {
      console.error('Error fetching products:', error);
      setProducts([]); // Устанавливаем пустой массив при ошибке
    } finally {
      setIsLoading(false);
    }
  }, [storefrontId, currentPage]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const handleLoadMore = useCallback(() => {
    if (hasMore && !isLoading) {
      setCurrentPage(prev => prev + 1);
    }
  }, [hasMore, isLoading]);

  const loadMoreRef = useInfiniteScroll({
    loading: isLoading,
    hasMore,
    onLoadMore: handleLoadMore,
  });

  if (isLoading && currentPage === 1) {
    return (
      <div className="card bg-base-200 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('storefronts.products.title')}</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="card bg-base-300 animate-pulse">
                <div className="aspect-[4/3] bg-base-content/10"></div>
                <div className="card-body">
                  <div className="h-4 bg-base-content/10 rounded w-3/4"></div>
                  <div className="h-6 bg-base-content/10 rounded w-1/2 mt-2"></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="card bg-base-200 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('storefronts.products.title')}</h2>
          <div className="text-center py-12">
            <svg className="w-16 h-16 mx-auto mb-4 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            <p className="text-base-content/60">{t('storefronts.noProducts')}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-200 shadow-xl">
      <div className="card-body">
        <div className="flex justify-between items-center mb-4">
          <h2 className="card-title text-2xl">{t('storefronts.products.title')}</h2>
          <div className="flex items-center gap-2">
            <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
            <Link 
              href={`/${locale}/marketplace?storefront_id=${storefrontId}`}
              className="btn btn-sm btn-ghost"
            >
              {t('common.viewAll')}
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>
          </div>
        </div>

        <div className={
          viewMode === 'grid' 
            ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
            : "space-y-4"
        }>
          {Array.isArray(products) ? products.map((product) => (
            viewMode === 'grid' ? (
              <Link 
                key={product.id} 
                href={`/${locale}/marketplace/${product.id}`}
                className="card bg-base-100 hover:shadow-lg transition-shadow"
              >
                <figure className="aspect-[4/3] relative">
                  <SafeImage
                    src={product.main_image}
                    alt={product.title}
                    fill
                    className="object-cover"
                    fallback={
                      <div className="w-full h-full bg-base-200 flex items-center justify-center">
                        <svg className="w-16 h-16 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
                        </svg>
                      </div>
                    }
                  />
                  {product.has_discount && product.discount_percentage && (
                    <div className="absolute top-2 right-2 badge badge-error">
                      -{product.discount_percentage}%
                    </div>
                  )}
                </figure>
                <div className="card-body p-4">
                  <h3 className="font-semibold line-clamp-2">{product.title}</h3>
                  <div className="flex items-baseline gap-2">
                    <span className="text-xl font-bold text-primary">
                      {product.price} $
                    </span>
                    {product.old_price && product.old_price > product.price && (
                      <span className="text-sm line-through text-base-content/50">
                        {product.old_price} $
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-base-content/60">
                    {product.location}
                  </p>
                </div>
              </Link>
            ) : (
              <Link 
                key={product.id} 
                href={`/${locale}/marketplace/${product.id}`}
                className="card bg-base-100 hover:shadow-lg transition-shadow"
              >
                <div className="card-body p-4">
                  <div className="flex gap-4">
                    <figure className="relative w-24 h-24 flex-shrink-0 rounded-lg overflow-hidden">
                      <SafeImage
                        src={product.main_image}
                        alt={product.title}
                        fill
                        className="object-cover"
                        fallback={
                          <div className="w-full h-full bg-base-200 flex items-center justify-center">
                            <svg className="w-8 h-8 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
                            </svg>
                          </div>
                        }
                      />
                      {product.has_discount && product.discount_percentage && (
                        <div className="absolute top-1 right-1 badge badge-error badge-sm">
                          -{product.discount_percentage}%
                        </div>
                      )}
                    </figure>
                    <div className="flex-grow">
                      <h3 className="font-semibold line-clamp-1">{product.title}</h3>
                      <div className="flex items-baseline gap-2 mt-1">
                        <span className="text-lg font-bold text-primary">
                          {product.price} $
                        </span>
                        {product.old_price && product.old_price > product.price && (
                          <span className="text-sm line-through text-base-content/50">
                            {product.old_price} $
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-base-content/60 mt-1">
                        {product.location}
                      </p>
                    </div>
                  </div>
                </div>
              </Link>
            )
          )) : (
            <div className="col-span-full text-center text-base-content/60">
              Нет данных о продуктах
            </div>
          )}
        </div>

        <InfiniteScrollTrigger
          ref={loadMoreRef}
          loading={isLoading}
          hasMore={hasMore}
          onLoadMore={handleLoadMore}
          showButton={true}
        />
      </div>
    </div>
  );
}