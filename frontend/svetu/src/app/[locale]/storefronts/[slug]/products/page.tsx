'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useDispatch } from 'react-redux';
import Link from 'next/link';
import { FiPlus, FiUpload, FiArrowLeft } from 'react-icons/fi';
import { toast } from 'react-hot-toast';
import { apiClient } from '@/services/api-client';
import { ProductList } from '@/components/products/ProductList';
import {
  setProducts,
  appendProducts,
  setPagination,
} from '@/store/slices/productSlice';
import type { AppDispatch } from '@/store';

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
  }>;
}

export default function ProductsPage({ params }: PageProps) {
  const [slug, setSlug] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [totalCount, setTotalCount] = useState(0);

  const t = useTranslations('storefronts');
  const locale = useLocale();
  const dispatch = useDispatch<AppDispatch>();

  const LIMIT = 20;

  useEffect(() => {
    params.then((p) => setSlug(p.slug));
  }, [params]);

  const loadProducts = useCallback(
    async (pageNum: number, append = false) => {
      if (!slug) return;

      try {
        setLoading(true);

        const params = new URLSearchParams({
          limit: LIMIT.toString(),
          offset: ((pageNum - 1) * LIMIT).toString(),
        });

        const response = await apiClient.get(
          `/api/v1/storefronts/slug/${slug}/products?${params}`
        );

        if (response.data) {
          const newProducts = Array.isArray(response.data)
            ? response.data
            : response.data.products || response.data.data || [];

          if (append) {
            dispatch(appendProducts(newProducts));
          } else {
            dispatch(setProducts(newProducts));
          }

          // Обработка метаданных пагинации
          const meta = response.data.meta || response.data.pagination;
          if (meta) {
            setTotalCount(meta.total || 0);
            const totalPages = Math.ceil((meta.total || 0) / LIMIT);
            setHasMore(pageNum < totalPages);

            dispatch(
              setPagination({
                page: pageNum,
                total: meta.total || 0,
                hasMore: pageNum < totalPages,
              })
            );
          } else {
            // Предполагаем, что есть еще данные, если получили полную страницу
            setHasMore(newProducts.length === LIMIT);
          }
        }
      } catch (error) {
        console.error('Failed to load products:', error);
        toast.error(t('products.errorLoadingProducts'));
      } finally {
        setLoading(false);
      }
    },
    [slug, dispatch, t]
  );

  useEffect(() => {
    if (slug) {
      setPage(1);
      loadProducts(1, false);
    }
  }, [slug, loadProducts]);

  const handleLoadMore = () => {
    if (hasMore && !loading) {
      const nextPage = page + 1;
      setPage(nextPage);
      loadProducts(nextPage, true);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Заголовок */}
      <div className="mb-8">
        <Link
          href={`/${locale}/storefronts/${slug}/dashboard`}
          className="inline-flex items-center text-primary hover:underline mb-4"
        >
          <FiArrowLeft className="w-4 h-4 mr-2" />
          {t('backToDashboard')}
        </Link>

        <div className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-4">
          <div>
            <h1 className="text-4xl font-bold">{t('products.title')}</h1>
            <p className="text-base-content/60 mt-2">
              {t('products.addProductDescription')}
            </p>
          </div>

          <div className="flex flex-wrap gap-3">
            <Link
              href={`/${locale}/storefronts/${slug}/products/import`}
              className="btn btn-outline gap-2"
            >
              <FiUpload className="w-5 h-5" />
              {t('products.importProducts')}
            </Link>
            <Link
              href={`/${locale}/storefronts/${slug}/products/new`}
              className="btn btn-primary gap-2"
            >
              <FiPlus className="w-5 h-5" />
              {t('products.addProduct')}
            </Link>
          </div>
        </div>
      </div>

      {/* Список товаров с массовыми операциями */}
      <ProductList
        storefrontSlug={slug}
        loading={loading}
        hasMore={hasMore}
        onLoadMore={handleLoadMore}
        totalCount={totalCount}
      />
    </div>
  );
}
