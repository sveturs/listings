'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import {
  PlusIcon,
  PencilIcon,
  TrashIcon,
  PhotoIcon,
  ArrowLeftIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';
import { toast } from '@/utils/toast';
import { apiClient } from '@/services/api-client';
import ViewToggle from '@/components/common/ViewToggle';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import SafeImage from '@/components/SafeImage';
import type { components } from '@/types/generated/api';

type B2CProduct = components['schemas']['models.StorefrontProduct'];

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
  }>;
}

export default function ProductsPage({ params }: PageProps) {
  const [slug, setSlug] = useState<string>('');

  useEffect(() => {
    params.then((p) => setSlug(p.slug));
  }, [params]);
  const t = useTranslations('storefronts');
  const tStorefronts = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const [products, setProducts] = useState<B2CProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState<B2CProduct | null>(
    null
  );
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [_totalProducts, setTotalProducts] = useState(0);
  const PRODUCTS_PER_PAGE = 20;

  const loadProducts = async (page = 1, append = false) => {
    if (!slug) return; // Wait for slug to be loaded

    try {
      setLoading(true);
      const params = new URLSearchParams();
      if (search) params.append('search', search);
      if (filterStatus !== 'all') {
        if (filterStatus === 'active') params.append('is_active', 'true');
        if (filterStatus === 'inactive') params.append('is_active', 'false');
        if (filterStatus === 'out_of_stock')
          params.append('stock_status', 'out_of_stock');
        if (filterStatus === 'low_stock')
          params.append('stock_status', 'low_stock');
      }

      // Add pagination params
      params.append('page', page.toString());
      params.append('limit', PRODUCTS_PER_PAGE.toString());

      const queryString = params.toString();
      const url = `/api/v1/b2c/slug/${slug}/products?${queryString}`;
      const response = await apiClient.get(url);
      if (response.data) {
        const responseData = response.data;
        const newProducts = Array.isArray(responseData)
          ? responseData
          : responseData.products || responseData.data || [];

        if (append) {
          setProducts((prev) => [...prev, ...newProducts]);
        } else {
          setProducts(newProducts);
        }

        // Handle pagination metadata
        const meta = responseData.meta || responseData.pagination;
        if (meta) {
          setTotalProducts(meta.total || 0);
          const totalPages = Math.ceil((meta.total || 0) / PRODUCTS_PER_PAGE);
          setHasMore(page < totalPages);
        } else {
          // Fallback: assume there are more if we got a full page
          setHasMore(newProducts.length === PRODUCTS_PER_PAGE);
        }
      }
    } catch (error) {
      console.error('Failed to load products:', error);
      toast.error(tStorefronts('products.errorLoadingProducts'));
    } finally {
      setLoading(false);
    }
  };

  const handleLoadMore = () => {
    if (hasMore && !loading) {
      const nextPage = currentPage + 1;
      setCurrentPage(nextPage);
      loadProducts(nextPage, true);
    }
  };

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore,
    onLoadMore: handleLoadMore,
  });

  useEffect(() => {
    setCurrentPage(1);
    loadProducts(1, false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [slug, search, filterStatus]);

  const handleDelete = async () => {
    if (!productToDelete || !slug) return;

    try {
      await apiClient.delete(
        `/api/v1/b2c/slug/${slug}/products/${productToDelete.id}`
      );
      toast.success(tStorefronts('products.productDeleted'));
      setDeleteModalOpen(false);
      setProductToDelete(null);
      setCurrentPage(1);
      loadProducts(1, false);
    } catch (error) {
      console.error('Failed to delete product:', error);
      toast.error(tStorefronts('products.errorDeletingProduct'));
    }
  };

  const getStockStatusBadge = (product: B2CProduct) => {
    switch (product.stock_status) {
      case 'out_of_stock':
        return (
          <span className="badge badge-error">
            {tStorefronts('products.outOfStock')}
          </span>
        );
      case 'low_stock':
        return (
          <span className="badge badge-warning">
            {tStorefronts('products.lowStock')}
          </span>
        );
      case 'in_stock':
        return (
          <span className="badge badge-success">
            {tStorefronts('products.inStock')}
          </span>
        );
      default:
        return null;
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <Link
            href={`/${locale}/b2c/${slug}/dashboard`}
            className="inline-flex items-center text-primary hover:underline mb-2"
          >
            <ArrowLeftIcon className="w-4 h-4 mr-2" />
            {tStorefronts('backToDashboard')}
          </Link>
          <h1 className="text-3xl font-bold">
            {tStorefronts('products.title')}
          </h1>
          <p className="text-base-content/60 mt-2">
            {t('products.totalProducts', {
              count: products.length,
            })}
          </p>
        </div>
        <div className="flex gap-3">
          <Link
            href={`/${locale}/b2c/${slug}/products/import`}
            className="btn btn-outline btn-secondary"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
            {tStorefronts('products.importProducts')}
          </Link>
          <Link
            href={`/${locale}/b2c/${slug}/products/new`}
            className="btn btn-primary"
          >
            <PlusIcon className="w-5 h-5" />
            {tStorefronts('products.addProduct')}
          </Link>
        </div>
      </div>

      {/* Filters */}
      <div className="card bg-base-100 shadow-xl mb-6">
        <div className="card-body">
          <div className="flex flex-col md:flex-row gap-4">
            {/* Search */}
            <div className="form-control flex-1">
              <div className="input-group">
                <input
                  type="text"
                  placeholder={tStorefronts('products.searchPlaceholder')}
                  className="input input-bordered w-full"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
                <button className="btn btn-square">
                  <MagnifyingGlassIcon className="w-5 h-5" />
                </button>
              </div>
            </div>

            {/* Status Filter */}
            <div className="form-control w-full md:w-auto">
              <select
                className="select select-bordered"
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
              >
                <option value="all">
                  {tStorefronts('products.allProducts')}
                </option>
                <option value="active">
                  {tStorefronts('products.activeOnly')}
                </option>
                <option value="inactive">
                  {tStorefronts('products.inactiveOnly')}
                </option>
                <option value="out_of_stock">
                  {tStorefronts('products.outOfStock')}
                </option>
                <option value="low_stock">
                  {tStorefronts('products.lowStock')}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* View Toggle */}
      {products.length > 0 && (
        <div className="flex justify-end mb-4">
          <ViewToggle currentView={viewMode} onViewChange={setViewMode} />
        </div>
      )}

      {/* Products Grid/List */}
      {loading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-12">
          <PhotoIcon className="w-16 h-16 mx-auto text-base-content/20 mb-4" />
          <p className="text-base-content/60 mb-4">
            {search || filterStatus !== 'all'
              ? tStorefronts('products.noProductsFound')
              : tStorefronts('products.noProducts')}
          </p>
          {!search && filterStatus === 'all' && (
            <Link
              href={`/${locale}/b2c/${slug}/products/new`}
              className="btn btn-primary"
            >
              <PlusIcon className="w-5 h-5" />
              {tStorefronts('products.addFirstProduct')}
            </Link>
          )}
        </div>
      ) : (
        <div
          className={
            viewMode === 'grid'
              ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'
              : 'space-y-4'
          }
        >
          {products.map((product) =>
            viewMode === 'grid' ? (
              <div key={product.id} className="card bg-base-100 shadow-xl">
                {/* Product Image */}
                <figure className="px-4 pt-4">
                  <SafeImage
                    src={product.images?.[0]?.thumbnail_url}
                    alt={product.name || 'Product image'}
                    width={300}
                    height={200}
                    className="rounded-xl object-cover h-48 w-full"
                    fallback={
                      <div className="bg-base-200 rounded-xl h-48 w-full flex items-center justify-center">
                        <PhotoIcon className="w-16 h-16 text-base-content/20" />
                      </div>
                    }
                  />
                </figure>

                <div className="card-body">
                  {/* Product Name */}
                  <h2 className="card-title text-lg">
                    {product.name}
                    {!product.is_active && (
                      <span className="badge badge-ghost">
                        {tStorefronts('products.inactive')}
                      </span>
                    )}
                  </h2>

                  {/* Category */}
                  {product.category && (
                    <p className="text-sm text-base-content/60">
                      {product.category.name}
                    </p>
                  )}

                  {/* Price and Stock */}
                  <div className="flex justify-between items-center mt-2">
                    <span className="text-xl font-bold">
                      {product.price} {product.currency}
                    </span>
                    {getStockStatusBadge(product)}
                  </div>

                  {/* Stock Quantity */}
                  <p className="text-sm text-base-content/60">
                    {t('products.stockCount', {
                      count: product.stock_quantity || 0,
                    })}
                  </p>

                  {/* SKU */}
                  {product.sku && (
                    <p className="text-xs text-base-content/40">
                      SKU: {product.sku}
                    </p>
                  )}

                  {/* Actions */}
                  <div className="card-actions justify-end mt-4">
                    <Link
                      href={`/${locale}/b2c/${slug}/products/${product.id}/edit`}
                      className="btn btn-sm btn-ghost"
                    >
                      <PencilIcon className="w-4 h-4" />
                    </Link>
                    <button
                      onClick={() => {
                        setProductToDelete(product);
                        setDeleteModalOpen(true);
                      }}
                      className="btn btn-sm btn-ghost text-error"
                    >
                      <TrashIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              // List view
              <div key={product.id} className="card bg-base-100 shadow-xl">
                <div className="card-body p-4">
                  <div className="flex gap-4">
                    {/* Product Image */}
                    <figure className="relative w-32 h-32 flex-shrink-0 rounded-lg overflow-hidden bg-base-200">
                      <SafeImage
                        src={product.images?.[0]?.thumbnail_url}
                        alt={product.name || 'Product image'}
                        fill
                        className="object-cover"
                        fallback={
                          <div className="w-full h-full flex items-center justify-center">
                            <PhotoIcon className="w-12 h-12 text-base-content/20" />
                          </div>
                        }
                      />
                    </figure>

                    {/* Product Info */}
                    <div className="flex-grow">
                      <div className="flex justify-between items-start gap-4">
                        <div className="flex-grow">
                          <h2 className="text-lg font-semibold line-clamp-1">
                            {product.name}
                            {!product.is_active && (
                              <span className="badge badge-ghost badge-sm ml-2">
                                {tStorefronts('products.inactive')}
                              </span>
                            )}
                          </h2>
                          {product.category && (
                            <p className="text-sm text-base-content/60">
                              {product.category.name}
                            </p>
                          )}
                          <div className="flex items-center gap-4 mt-2">
                            <span className="text-xl font-bold">
                              {product.price} {product.currency}
                            </span>
                            {getStockStatusBadge(product)}
                            <span className="text-sm text-base-content/60">
                              {t('products.stockCount', {
                                count: product.stock_quantity || 0,
                              })}
                            </span>
                            {product.sku && (
                              <span className="text-xs text-base-content/40">
                                SKU: {product.sku}
                              </span>
                            )}
                          </div>
                        </div>

                        {/* Actions */}
                        <div className="flex gap-2">
                          <Link
                            href={`/${locale}/b2c/${slug}/products/${product.id}/edit`}
                            className="btn btn-sm btn-ghost"
                          >
                            <PencilIcon className="w-4 h-4" />
                          </Link>
                          <button
                            onClick={() => {
                              setProductToDelete(product);
                              setDeleteModalOpen(true);
                            }}
                            className="btn btn-sm btn-ghost text-error"
                          >
                            <TrashIcon className="w-4 h-4" />
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )
          )}
        </div>
      )}

      {products.length > 0 && (
        <InfiniteScrollTrigger
          ref={loadMoreRef}
          loading={loading}
          hasMore={hasMore}
          onLoadMore={handleLoadMore}
          showButton={true}
          loadMoreText={tCommon('loadMore')}
        />
      )}

      {/* Delete Confirmation Modal */}
      {deleteModalOpen && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {tStorefronts('products.deleteProductTitle')}
            </h3>
            <p className="py-4">
              {t('products.deleteProductConfirm', {
                name: productToDelete?.name || '',
              })}
            </p>
            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => {
                  setDeleteModalOpen(false);
                  setProductToDelete(null);
                }}
              >
                {tCommon('cancel')}
              </button>
              <button className="btn btn-error" onClick={handleDelete}>
                {tCommon('delete')}
              </button>
            </div>
          </div>
          <div
            className="modal-backdrop"
            onClick={() => setDeleteModalOpen(false)}
          ></div>
        </div>
      )}
    </div>
  );
}
