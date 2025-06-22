'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import Image from 'next/image';
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
import type { components } from '@/types/generated/api';

type StorefrontProduct = components['schemas']['models.StorefrontProduct'];

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
  const t = useTranslations();
  const locale = useLocale();
  const [products, setProducts] = useState<StorefrontProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [productToDelete, setProductToDelete] =
    useState<StorefrontProduct | null>(null);

  const loadProducts = async () => {
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

      const queryString = params.toString();
      const url = `/api/v1/storefronts/slug/${slug}/products${queryString ? `?${queryString}` : ''}`;
      const response = await apiClient.get(url);
      if (response.data) {
        setProducts(response.data);
      }
    } catch (error) {
      console.error('Failed to load products:', error);
      toast.error(t('storefronts.products.errorLoadingProducts'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [slug, search, filterStatus]);

  const handleDelete = async () => {
    if (!productToDelete || !slug) return;

    try {
      await apiClient.delete(
        `/api/v1/storefronts/slug/${slug}/products/${productToDelete.id}`
      );
      toast.success(t('storefronts.products.productDeleted'));
      setDeleteModalOpen(false);
      setProductToDelete(null);
      loadProducts();
    } catch (error) {
      console.error('Failed to delete product:', error);
      toast.error(t('storefronts.products.errorDeletingProduct'));
    }
  };

  const getStockStatusBadge = (product: StorefrontProduct) => {
    switch (product.stock_status) {
      case 'out_of_stock':
        return (
          <span className="badge badge-error">
            {t('storefronts.products.outOfStock')}
          </span>
        );
      case 'low_stock':
        return (
          <span className="badge badge-warning">
            {t('storefronts.products.lowStock')}
          </span>
        );
      case 'in_stock':
        return (
          <span className="badge badge-success">
            {t('storefronts.products.inStock')}
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
            href={`/${locale}/storefronts/${slug}/dashboard`}
            className="inline-flex items-center text-primary hover:underline mb-2"
          >
            <ArrowLeftIcon className="w-4 h-4 mr-2" />
            {t('storefronts.backToDashboard')}
          </Link>
          <h1 className="text-3xl font-bold">
            {t('storefronts.products.title')}
          </h1>
          <p className="text-base-content/60 mt-2">
            {t('storefronts.products.totalProducts', {
              count: products.length,
            })}
          </p>
        </div>
        <div className="flex gap-3">
          <Link
            href={`/${locale}/storefronts/${slug}/products/import`}
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
            {t('storefronts.products.importProducts')}
          </Link>
          <Link
            href={`/${locale}/storefronts/${slug}/products/new`}
            className="btn btn-primary"
          >
            <PlusIcon className="w-5 h-5" />
            {t('storefronts.products.addProduct')}
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
                  placeholder={t('storefronts.products.searchPlaceholder')}
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
                  {t('storefronts.products.allProducts')}
                </option>
                <option value="active">
                  {t('storefronts.products.activeOnly')}
                </option>
                <option value="inactive">
                  {t('storefronts.products.inactiveOnly')}
                </option>
                <option value="out_of_stock">
                  {t('storefronts.products.outOfStock')}
                </option>
                <option value="low_stock">
                  {t('storefronts.products.lowStock')}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Products Grid */}
      {loading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-12">
          <PhotoIcon className="w-16 h-16 mx-auto text-base-content/20 mb-4" />
          <p className="text-base-content/60 mb-4">
            {search || filterStatus !== 'all'
              ? t('storefronts.products.noProductsFound')
              : t('storefronts.products.noProducts')}
          </p>
          {!search && filterStatus === 'all' && (
            <Link
              href={`/${locale}/storefronts/${slug}/products/new`}
              className="btn btn-primary"
            >
              <PlusIcon className="w-5 h-5" />
              {t('storefronts.products.addFirstProduct')}
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {products.map((product) => (
            <div key={product.id} className="card bg-base-100 shadow-xl">
              {/* Product Image */}
              <figure className="px-4 pt-4">
                {product.images &&
                product.images.length > 0 &&
                product.images[0].thumbnail_url ? (
                  <Image
                    src={product.images[0].thumbnail_url}
                    alt={product.name || 'Product image'}
                    width={300}
                    height={200}
                    className="rounded-xl object-cover h-48 w-full"
                  />
                ) : (
                  <div className="bg-base-200 rounded-xl h-48 w-full flex items-center justify-center">
                    <PhotoIcon className="w-16 h-16 text-base-content/20" />
                  </div>
                )}
              </figure>

              <div className="card-body">
                {/* Product Name */}
                <h2 className="card-title text-lg">
                  {product.name}
                  {!product.is_active && (
                    <span className="badge badge-ghost">
                      {t('storefronts.products.inactive')}
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
                  {t('storefronts.products.stockCount', {
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
                    href={`/${locale}/storefronts/${slug}/products/${product.id}/edit`}
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
          ))}
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {deleteModalOpen && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {t('storefronts.products.deleteProductTitle')}
            </h3>
            <p className="py-4">
              {t('storefronts.products.deleteProductConfirm', {
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
                {t('common.cancel')}
              </button>
              <button className="btn btn-error" onClick={handleDelete}>
                {t('common.delete')}
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
