'use client';

import { useState, useEffect, useCallback } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import Image from 'next/image';
import { marketplaceApi } from '@/services/marketplaceApi';
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
  const [totalPages, setTotalPages] = useState(1);

  const fetchProducts = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await marketplaceApi.getListings({
        storefront_id: storefrontId,
        limit: 12,
        offset: (currentPage - 1) * 12
      });
      
      // Исправляем структуру ответа: API возвращает { data: { data: [...], success: true, meta: {...} } }
      const actualData = (response as any).data?.data || response.data || [];
      setProducts(actualData);
      
      // Проверяем пагинацию в правильном месте
      const pagination = (response as any).data?.meta || (response as any).pagination;
      if (pagination && pagination.total && pagination.limit) {
        setTotalPages(Math.ceil(pagination.total / pagination.limit));
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

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.isArray(products) ? products.map((product) => (
            <Link 
              key={product.id} 
              href={`/${locale}/marketplace/${product.id}`}
              className="card bg-base-100 hover:shadow-lg transition-shadow"
            >
              <figure className="aspect-[4/3] relative">
                <Image
                  src={product.main_image || '/placeholder-listing.jpg'}
                  alt={product.title}
                  fill
                  className="object-cover"
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
          )) : (
            <div className="col-span-full text-center text-base-content/60">
              Нет данных о продуктах
            </div>
          )}
        </div>

        {totalPages > 1 && (
          <div className="flex justify-center mt-6">
            <div className="join">
              <button 
                className="join-item btn btn-sm"
                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                disabled={currentPage === 1}
              >
                «
              </button>
              {[...Array(Math.min(5, totalPages))].map((_, i) => {
                const page = i + 1;
                return (
                  <button
                    key={page}
                    className={`join-item btn btn-sm ${currentPage === page ? 'btn-active' : ''}`}
                    onClick={() => setCurrentPage(page)}
                  >
                    {page}
                  </button>
                );
              })}
              <button 
                className="join-item btn btn-sm"
                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
              >
                »
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}