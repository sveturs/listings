'use client';

import React, { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { format } from 'date-fns';
import { ru, enUS, sr } from 'date-fns/locale';
import { useLocale } from 'next-intl';
import { apiClient } from '@/services/api-client';

interface StorefrontProduct {
  id: number;
  name: string;
  description: string;
  price: number;
  currency: string;
  sku?: string;
  stock?: number;
  unlimited_stock?: boolean;
  status: string;
  is_active?: boolean;
  storefront_id: number;
  storefront_name?: string;
  storefront_slug?: string;
  storefront?: {
    id: number;
    name: string;
    slug: string;
  };
  category_id?: number;
  category_name?: string;
  category?: {
    id: number;
    name: string;
    slug: string;
  };
  images?: Array<
    | {
        id: number;
        product_id: number;
        file_path: string;
        file_name: string;
        public_url: string;
        is_main: boolean;
      }
    | string
  > | null;
  sales_count?: number;
  created_at: string;
  updated_at: string;
}

export default function AdminStorefrontProductsTable() {
  const t = useTranslations('admin');
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();

  const [products, setProducts] = useState<StorefrontProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [selectedProducts, setSelectedProducts] = useState<Set<number>>(
    new Set()
  );
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState<number | null>(null);

  // Фильтры и пагинация
  const [searchQuery, setSearchQuery] = useState(
    searchParams.get('query') || ''
  );
  const [storefrontFilter] = useState(searchParams.get('storefront') || '');
  const [_categoryFilter] = useState(searchParams.get('category') || '');
  const [stockFilter, setStockFilter] = useState(
    searchParams.get('stock') || ''
  );
  const [_statusFilter] = useState(searchParams.get('status') || '');
  const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'));
  const [limit] = useState(20);
  const [sortBy, setSortBy] = useState(searchParams.get('sort') || 'date_desc');

  // Загрузка данных
  const fetchProducts = useCallback(async () => {
    setLoading(true);
    try {
      // Для админ панели загружаем ВСЕ витрины (включая неактивные)
      const storefrontsResponse = await apiClient.get(
        `/storefronts?include_inactive=true`
      );

      if (!storefrontsResponse.data) {
        console.error('Failed to fetch storefronts');
        setProducts([]);
        setLoading(false);
        return;
      }

      const storefrontsData = storefrontsResponse.data;

      // API возвращает витрины в поле storefronts или data
      const storefronts =
        storefrontsData.storefronts || storefrontsData.data || [];

      if (!storefronts || storefronts.length === 0) {
        setProducts([]);
        setTotal(0);
        setLoading(false);
        return;
      }
      let allProducts: StorefrontProduct[] = [];

      // Получаем товары для каждой витрины
      console.log('Loading products for storefronts:', storefronts);
      for (const storefront of storefronts) {
        if (!storefront.slug) continue;
        console.log(
          'Fetching products for storefront:',
          storefront.slug,
          storefront
        );

        const params = new URLSearchParams();
        if (searchQuery) params.append('search', searchQuery);
        // if (categoryFilter) params.append('category_id', categoryFilter);
        // if (statusFilter) params.append('status', statusFilter);
        params.append('limit', '100'); // Получаем больше товаров за раз
        // Для админки явно запрашиваем ВСЕ товары (включая неактивные)
        // Передаем is_active=all чтобы получить и активные и неактивные товары
        params.append('is_active', 'all');

        // Используем правильный эндпоинт с slug в пути
        const response = await apiClient.get(
          `/storefronts/slug/${storefront.slug}/products?${params}`
        );

        if (response.data) {
          const data = response.data;
          // API возвращает массив товаров напрямую, а не в поле data
          // Обрабатываем случай когда API возвращает null для пустых результатов
          let products = [];

          if (data === null || data === undefined) {
            products = [];
          } else if (Array.isArray(data)) {
            products = data;
          } else if (data.data) {
            products = Array.isArray(data.data) ? data.data : [];
          }

          // Добавляем информацию о витрине к каждому товару
          const productsWithStorefront = products.map((p: any) => ({
            ...p,
            storefront_name: storefront.name,
            storefront_id: storefront.id,
            storefront_slug: storefront.slug,
            // Исправляем имена полей для совместимости
            stock: p.stock || p.stock_quantity,
            is_active: p.is_active !== false,
            sales_count: p.sales_count || p.sold_count || 0,
          }));
          console.log(
            `Found ${productsWithStorefront.length} products for storefront ${storefront.slug}`
          );
          allProducts = [...allProducts, ...productsWithStorefront];
        }
      }

      // Фильтруем по складу
      if (stockFilter) {
        if (stockFilter === 'out_of_stock') {
          allProducts = allProducts.filter(
            (p) => !p.unlimited_stock && (!p.stock || p.stock === 0)
          );
        } else if (stockFilter === 'low_stock') {
          allProducts = allProducts.filter(
            (p) => !p.unlimited_stock && p.stock && p.stock < 10 && p.stock > 0
          );
        } else if (stockFilter === 'in_stock') {
          allProducts = allProducts.filter(
            (p) => p.unlimited_stock || (p.stock && p.stock > 0)
          );
        }
      }

      // Фильтруем по конкретной витрине если выбрана
      if (storefrontFilter) {
        allProducts = allProducts.filter(
          (p) => p.storefront_id === parseInt(storefrontFilter)
        );
      }

      // Сортировка
      allProducts.sort((a, b) => {
        switch (sortBy) {
          case 'name_asc':
            return a.name.localeCompare(b.name);
          case 'name_desc':
            return b.name.localeCompare(a.name);
          case 'price_asc':
            return a.price - b.price;
          case 'price_desc':
            return b.price - a.price;
          case 'stock_asc':
            return (a.stock || 0) - (b.stock || 0);
          case 'stock_desc':
            return (b.stock || 0) - (a.stock || 0);
          case 'sales_desc':
            return (b.sales_count || 0) - (a.sales_count || 0);
          case 'date_asc':
            return (
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime()
            );
          case 'date_desc':
          default:
            return (
              new Date(b.created_at).getTime() -
              new Date(a.created_at).getTime()
            );
        }
      });

      // Применяем пагинацию
      console.log(`Total products after filtering: ${allProducts.length}`);
      console.log('Current filters:', {
        stockFilter,
        storefrontFilter,
        searchQuery,
      });
      const start = (page - 1) * limit;
      const paginatedProducts = allProducts.slice(start, start + limit);
      console.log(
        `Showing products ${start} to ${start + limit}:`,
        paginatedProducts
      );

      setProducts(paginatedProducts);
      setTotal(allProducts.length);
    } catch (error) {
      console.error('Error fetching products:', error);
      setProducts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [
    searchQuery,
    storefrontFilter,
    // categoryFilter,
    stockFilter,
    // statusFilter,
    page,
    limit,
    sortBy,
  ]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  // Обновление URL при изменении фильтров
  useEffect(() => {
    const params = new URLSearchParams();
    if (searchQuery) params.set('query', searchQuery);
    if (storefrontFilter) params.set('storefront', storefrontFilter);
    // if (categoryFilter) params.set('category', categoryFilter);
    if (stockFilter) params.set('stock', stockFilter);
    // if (statusFilter) params.set('status', statusFilter);
    if (page > 1) params.set('page', page.toString());
    if (sortBy !== 'date_desc') params.set('sort', sortBy);

    const newUrl = params.toString() ? `?${params}` : '';
    router.push(`/${locale}/admin/storefront-products${newUrl}`, {
      scroll: false,
    });
  }, [
    searchQuery,
    storefrontFilter,
    // categoryFilter,
    stockFilter,
    // statusFilter,
    page,
    sortBy,
    router,
    locale,
  ]);

  // Обработчики действий
  const handleDelete = async (id: number, storefrontSlug?: string) => {
    try {
      // Находим товар чтобы получить slug витрины
      const product = products.find((p) => p.id === id);
      const slug =
        storefrontSlug || product?.storefront_slug || product?.storefront?.slug;

      if (!slug) {
        console.error('Cannot delete product: storefront slug not found');
        return;
      }

      const response = await apiClient.delete(
        `/storefronts/slug/${slug}/products/${id}`
      );

      if (response.data) {
        await fetchProducts();
        router.refresh(); // Инвалидировать серверный кеш
        setDeleteModalOpen(false);
        setProductToDelete(null);
      }
    } catch (error) {
      console.error('Error deleting product:', error);
    }
  };

  const handleBulkDeactivate = async () => {
    if (selectedProducts.size === 0) return;

    const confirmed = window.confirm(
      `Вы уверены что хотите деактивировать ${selectedProducts.size} товаров?`
    );
    if (!confirmed) return;

    try {
      // Деактивируем все товары параллельно
      const deactivatePromises = Array.from(selectedProducts).map(
        async (id) => {
          const product = products.find((p) => p.id === id);
          const slug = product?.storefront_slug || product?.storefront?.slug;

          if (!slug) {
            console.error(
              `Cannot deactivate product ${id}: storefront slug not found`
            );
            return false;
          }

          try {
            const response = await apiClient.put(
              `/storefronts/slug/${slug}/products/${id}`,
              { is_active: false }
            );

            if (!response.data) {
              console.error(`Failed to deactivate product ${id}`);
              return false;
            }
            return true;
          } catch (error) {
            console.error(`Error deactivating product ${id}:`, error);
            return false;
          }
        }
      );

      const results = await Promise.all(deactivatePromises);
      const successCount = results.filter((r) => r === true).length;
      console.log(
        `Successfully deactivated ${successCount} out of ${selectedProducts.size} products`
      );

      // Очищаем выбранные товары и обновляем список
      setSelectedProducts(new Set());
      await fetchProducts();
      router.refresh(); // Инвалидировать серверный кеш
    } catch (error) {
      console.error('Error during bulk deactivate:', error);
    }
  };

  const handleBulkDelete = async () => {
    if (selectedProducts.size === 0) return;

    const confirmed = window.confirm(
      `Вы уверены что хотите УДАЛИТЬ НАВСЕГДА ${selectedProducts.size} товаров? Это действие необратимо!`
    );
    if (!confirmed) return;

    try {
      // Удаляем все товары параллельно с параметром ?hard=true
      const deletePromises = Array.from(selectedProducts).map(async (id) => {
        const product = products.find((p) => p.id === id);
        const slug = product?.storefront_slug || product?.storefront?.slug;

        if (!slug) {
          console.error(
            `Cannot delete product ${id}: storefront slug not found`
          );
          return;
        }

        try {
          const response = await apiClient.delete(
            `/storefronts/slug/${slug}/products/${id}?hard=true`
          );

          if (!response.data) {
            console.error(`Failed to delete product ${id}`);
          }
          return !!response.data;
        } catch (error) {
          console.error(`Error deleting product ${id}:`, error);
          return false;
        }
      });

      const results = await Promise.all(deletePromises);
      const successCount = results.filter((r) => r === true).length;
      console.log(
        `Successfully deleted ${successCount} out of ${selectedProducts.size} products`
      );

      // Очищаем выбранные товары и обновляем список
      setSelectedProducts(new Set());
      await fetchProducts();
      router.refresh(); // Инвалидировать серверный кеш
    } catch (error) {
      console.error('Error during bulk delete:', error);
    }
  };

  const handleToggleActive = async (id: number, isActive: boolean) => {
    try {
      // Находим товар чтобы получить slug витрины
      const product = products.find((p) => p.id === id);
      const slug = product?.storefront_slug || product?.storefront?.slug;

      if (!slug) {
        console.error('Cannot update product: storefront slug not found');
        return;
      }

      const response = await apiClient.put(
        `/storefronts/slug/${slug}/products/${id}`,
        { is_active: !isActive }
      );

      if (response.data) {
        await fetchProducts();
        router.refresh(); // Инвалидировать серверный кеш
      }
    } catch (error) {
      console.error('Error toggling product status:', error);
    }
  };

  const handleSelectAll = () => {
    if (selectedProducts.size === products.length) {
      setSelectedProducts(new Set());
    } else {
      setSelectedProducts(new Set(products.map((p) => p.id)));
    }
  };

  const handleSelectOne = (id: number) => {
    const newSelected = new Set(selectedProducts);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelectedProducts(newSelected);
  };

  const totalPages = Math.ceil(total / limit);

  const getStatusBadge = (status: string, isActive?: boolean) => {
    if (isActive === false)
      return (
        <span className="badge badge-ghost">
          {t('storefrontProducts.status.inactive')}
        </span>
      );

    switch (status) {
      case 'active':
        return (
          <span className="badge badge-success">
            {t('storefrontProducts.status.active')}
          </span>
        );
      case 'draft':
        return (
          <span className="badge badge-warning">
            {t('storefrontProducts.status.draft')}
          </span>
        );
      case 'archived':
        return (
          <span className="badge badge-ghost">
            {t('storefrontProducts.status.archived')}
          </span>
        );
      default:
        return <span className="badge">{status}</span>;
    }
  };

  const getStockBadge = (stock?: number, unlimited?: boolean) => {
    if (unlimited) {
      return (
        <span className="badge badge-info">
          {t('storefrontProducts.stockStatus.unlimited')}
        </span>
      );
    }

    if (!stock || stock === 0) {
      return (
        <span className="badge badge-error">
          {t('storefrontProducts.stockStatus.outOfStock')}
        </span>
      );
    }

    if (stock < 10) {
      return (
        <span className="badge badge-warning">
          {t('storefrontProducts.stockStatus.lowStock')} ({stock})
        </span>
      );
    }

    return (
      <span className="badge badge-success">
        {t('storefrontProducts.stockStatus.inStock')} ({stock})
      </span>
    );
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const dateLocale = locale === 'ru' ? ru : locale === 'sr' ? sr : enUS;
    return format(date, 'dd MMM yyyy', { locale: dateLocale });
  };

  if (loading && products.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Фильтры и поиск */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex flex-wrap gap-4">
            <div className="form-control flex-1 min-w-[200px]">
              <input
                type="text"
                placeholder={t('storefrontProducts.search.placeholder')}
                className="input input-bordered"
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setPage(1);
                }}
              />
            </div>

            <select
              className="select select-bordered"
              value={stockFilter}
              onChange={(e) => {
                setStockFilter(e.target.value);
                setPage(1);
              }}
            >
              <option value="">
                {t('storefrontProducts.filter.allStatuses')}
              </option>
              <option value="in_stock">
                {t('storefrontProducts.filter.inStock')}
              </option>
              <option value="out_of_stock">
                {t('storefrontProducts.filter.outOfStock')}
              </option>
              <option value="low_stock">
                {t('storefrontProducts.filter.lowStock')}
              </option>
            </select>

            <select
              className="select select-bordered"
              value={sortBy}
              onChange={(e) => {
                setSortBy(e.target.value);
                setPage(1);
              }}
            >
              <option value="date_desc">
                {t('storefrontProducts.sort.dateDesc')}
              </option>
              <option value="date_asc">
                {t('storefrontProducts.sort.dateAsc')}
              </option>
              <option value="name_asc">
                {t('storefrontProducts.sort.nameAsc')}
              </option>
              <option value="name_desc">
                {t('storefrontProducts.sort.nameDesc')}
              </option>
              <option value="price_desc">
                {t('storefrontProducts.sort.priceDesc')}
              </option>
              <option value="price_asc">
                {t('storefrontProducts.sort.priceAsc')}
              </option>
              <option value="stock_desc">
                {t('storefrontProducts.sort.stockDesc')}
              </option>
              <option value="stock_asc">
                {t('storefrontProducts.sort.stockAsc')}
              </option>
              <option value="sales_desc">
                {t('storefrontProducts.sort.salesDesc')}
              </option>
            </select>

            {selectedProducts.size > 0 && (
              <div className="dropdown dropdown-end">
                <label tabIndex={0} className="btn btn-outline">
                  {t('storefrontProducts.bulkActions.title')} (
                  {selectedProducts.size})
                </label>
                <ul
                  tabIndex={0}
                  className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
                >
                  <li>
                    <a
                      onClick={() => {
                        if (selectedProducts.size > 0) {
                          console.log('Activate selected');
                        }
                      }}
                    >
                      {t('storefrontProducts.bulkActions.activate')}
                    </a>
                  </li>
                  <li>
                    <a
                      onClick={() => {
                        if (selectedProducts.size > 0) {
                          handleBulkDeactivate();
                        }
                      }}
                    >
                      {t('storefrontProducts.bulkActions.deactivate')}
                    </a>
                  </li>
                  <li>
                    <a
                      className="text-error hover:bg-error hover:text-error-content"
                      onClick={() => {
                        if (selectedProducts.size > 0) {
                          handleBulkDelete();
                        }
                      }}
                    >
                      {t('storefrontProducts.bulkActions.delete')}
                    </a>
                  </li>
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Таблица товаров */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={
                        selectedProducts.size === products.length &&
                        products.length > 0
                      }
                      onChange={handleSelectAll}
                    />
                  </th>
                  <th>{t('storefrontProducts.table.id')}</th>
                  <th>{t('storefrontProducts.table.image')}</th>
                  <th>{t('storefrontProducts.table.name')}</th>
                  <th>{t('storefrontProducts.table.sku')}</th>
                  <th>{t('storefrontProducts.table.price')}</th>
                  <th>{t('storefrontProducts.table.stock')}</th>
                  <th>{t('storefrontProducts.table.storefront')}</th>
                  <th>{t('storefrontProducts.table.category')}</th>
                  <th>{t('storefrontProducts.table.status')}</th>
                  <th>{t('storefrontProducts.table.sales')}</th>
                  <th>{t('storefrontProducts.table.created')}</th>
                  <th>{t('storefrontProducts.table.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {products.length === 0 ? (
                  <tr>
                    <td
                      colSpan={13}
                      className="text-center py-8 text-base-content/60"
                    >
                      {t('storefrontProducts.noProducts')}
                    </td>
                  </tr>
                ) : (
                  products.map((product) => (
                    <tr key={product.id}>
                      <td>
                        <input
                          type="checkbox"
                          className="checkbox"
                          checked={selectedProducts.has(product.id)}
                          onChange={() => handleSelectOne(product.id)}
                        />
                      </td>
                      <td>{product.id}</td>
                      <td>
                        {product.images &&
                        product.images.length > 0 &&
                        product.images[0] ? (
                          <div className="avatar">
                            <div className="w-12 h-12 rounded bg-base-200">
                              <img
                                src={(() => {
                                  const firstImage = product.images![0];
                                  const imagePath =
                                    typeof firstImage === 'object' &&
                                    firstImage !== null
                                      ? (firstImage as any).public_url ||
                                        (firstImage as any).file_path
                                      : firstImage;

                                  if (!imagePath) return '/placeholder.jpg';
                                  if (typeof imagePath !== 'string')
                                    return '/placeholder.jpg';
                                  return imagePath.startsWith('http')
                                    ? imagePath
                                    : imagePath.startsWith('/')
                                      ? imagePath
                                      : `/listings/${imagePath}`;
                                })()}
                                alt={product.name}
                                className="w-full h-full object-cover rounded"
                                onError={(e) => {
                                  const target = e.target as HTMLImageElement;
                                  target.style.display = 'none';
                                  if (target.parentElement) {
                                    target.parentElement.innerHTML = `
                                      <svg class="w-6 h-6 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                      </svg>
                                    `;
                                  }
                                }}
                              />
                            </div>
                          </div>
                        ) : (
                          <div className="w-12 h-12 bg-base-200 rounded flex items-center justify-center">
                            <svg
                              className="w-6 h-6 text-base-content/20"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                              />
                            </svg>
                          </div>
                        )}
                      </td>
                      <td>
                        <div className="font-medium">{product.name}</div>
                        <div className="text-xs text-base-content/60 truncate max-w-xs">
                          {product.description}
                        </div>
                      </td>
                      <td>
                        {product.sku ? (
                          <code className="text-xs">{product.sku}</code>
                        ) : (
                          <span className="text-base-content/40">—</span>
                        )}
                      </td>
                      <td>
                        <div className="font-semibold">
                          {product.price} {product.currency || 'RSD'}
                        </div>
                      </td>
                      <td>
                        {getStockBadge(product.stock, product.unlimited_stock)}
                      </td>
                      <td>
                        <div className="text-sm">
                          {product.storefront?.name ||
                            product.storefront_name ||
                            `ID: ${product.storefront_id}`}
                        </div>
                      </td>
                      <td>
                        {product.category?.name ||
                          product.category_name ||
                          (product.category_id ? (
                            `ID: ${product.category_id}`
                          ) : (
                            <span className="text-base-content/40">—</span>
                          ))}
                      </td>
                      <td>
                        {getStatusBadge(product.status, product.is_active)}
                      </td>
                      <td>
                        <div className="flex items-center gap-1">
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                            />
                          </svg>
                          {product.sales_count || 0}
                        </div>
                      </td>
                      <td>
                        <div className="text-xs">
                          {formatDate(product.created_at)}
                        </div>
                      </td>
                      <td>
                        <div className="flex items-center gap-1">
                          <Link
                            href={`/${locale}/storefronts/${product.storefront_slug || product.storefront?.slug || product.storefront_id}/products/${product.id}`}
                            className="btn btn-ghost btn-xs"
                            target="_blank"
                          >
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                              />
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                              />
                            </svg>
                          </Link>

                          <Link
                            href={`/${locale}/admin/storefront-products/${product.id}/edit`}
                            className="btn btn-ghost btn-xs"
                          >
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                              />
                            </svg>
                          </Link>

                          <button
                            className="btn btn-ghost btn-xs"
                            onClick={() =>
                              handleToggleActive(
                                product.id,
                                product.is_active || false
                              )
                            }
                          >
                            {product.is_active ? (
                              <svg
                                className="w-4 h-4 text-warning"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                                />
                              </svg>
                            ) : (
                              <svg
                                className="w-4 h-4 text-success"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                />
                              </svg>
                            )}
                          </button>

                          <button
                            className="btn btn-ghost btn-xs text-error"
                            onClick={() => {
                              setProductToDelete(product.id);
                              setDeleteModalOpen(true);
                            }}
                          >
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                              />
                            </svg>
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {/* Пагинация */}
          {totalPages > 1 && (
            <div className="flex justify-center p-4">
              <div className="join">
                <button
                  className="join-item btn"
                  disabled={page === 1}
                  onClick={() => setPage(page - 1)}
                >
                  «
                </button>

                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  let pageNum;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (page <= 3) {
                    pageNum = i + 1;
                  } else if (page >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = page - 2 + i;
                  }

                  return (
                    <button
                      key={pageNum}
                      className={`join-item btn ${pageNum === page ? 'btn-active' : ''}`}
                      onClick={() => setPage(pageNum)}
                    >
                      {pageNum}
                    </button>
                  );
                })}

                <button
                  className="join-item btn"
                  disabled={page === totalPages}
                  onClick={() => setPage(page + 1)}
                >
                  »
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Модальное окно удаления */}
      {deleteModalOpen && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {t('storefrontProducts.deleteModal.title')}
            </h3>
            <p className="py-4">
              {t('storefrontProducts.deleteModal.message')}
            </p>
            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setDeleteModalOpen(false);
                  setProductToDelete(null);
                }}
              >
                {t('common.cancel')}
              </button>
              <button
                className="btn btn-error"
                onClick={() => {
                  if (productToDelete) {
                    handleDelete(productToDelete);
                  }
                }}
              >
                {t('common.delete')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button onClick={() => setDeleteModalOpen(false)}>close</button>
          </form>
        </dialog>
      )}
    </div>
  );
}
