'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useSelector, useDispatch } from 'react-redux';
import {
  FiGrid,
  FiList,
  FiMenu,
  FiFilter,
  FiSearch,
  FiX,
  FiCheck,
} from 'react-icons/fi';
import { ProductCard } from './ProductCard';
import { BulkActions } from './BulkActions';
import { InfiniteScrollTrigger } from '@/components/common/InfiniteScrollTrigger';
import type { RootState, AppDispatch } from '@/store';
import {
  toggleProductSelection,
  selectAll,
  clearSelection,
  toggleSelectMode,
  setViewMode,
  setFilters,
  bulkDeleteProducts,
  bulkUpdateStatus,
  exportProducts,
} from '@/store/slices/productSlice';

interface ProductListProps {
  storefrontSlug: string;
  loading?: boolean;
  hasMore?: boolean;
  onLoadMore?: () => void;
  totalCount?: number;
}

export function ProductList({
  storefrontSlug,
  loading = false,
  hasMore = false,
  onLoadMore,
  totalCount = 0,
}: ProductListProps) {
  const t = useTranslations('storefronts.products');
  const dispatch = useDispatch<AppDispatch>();

  const {
    products,
    selectedIds,
    ui: { isSelectMode, viewMode },
    bulkOperation: { isProcessing },
    filters,
  } = useSelector((state: RootState) => state.products);

  const [showFilters, setShowFilters] = useState(false);
  const [localSearch, setLocalSearch] = useState(filters.search);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      dispatch(setFilters({ search: localSearch }));
    }, 300);
    return () => clearTimeout(timer);
  }, [localSearch, dispatch]);

  const selectedCount = selectedIds.length;
  const allSelected =
    products.length > 0 &&
    products.every((p) => p.id && selectedIds.includes(p.id));

  const handleToggleSelectMode = () => {
    dispatch(toggleSelectMode());
  };

  const handleSelectAll = () => {
    if (allSelected) {
      dispatch(clearSelection());
    } else {
      dispatch(selectAll());
    }
  };

  const handleToggleSelect = useCallback(
    (id: number) => {
      dispatch(toggleProductSelection(id));
    },
    [dispatch]
  );

  const handleBulkDelete = async () => {
    await dispatch(
      bulkDeleteProducts({ storefrontSlug, productIds: selectedIds })
    );
    dispatch(clearSelection());
  };

  const handleBulkStatusChange = async (isActive: boolean) => {
    await dispatch(
      bulkUpdateStatus({ storefrontSlug, productIds: selectedIds, isActive })
    );
    dispatch(clearSelection());
  };

  const handleBulkExport = () => {
    const productIds = selectedCount > 0 ? selectedIds : undefined;
    dispatch(exportProducts({ storefrontSlug, productIds, format: 'csv' }));
  };

  return (
    <div className="space-y-4">
      {/* Панель управления */}
      <div className="bg-base-100 rounded-lg shadow-sm border border-base-300 p-4">
        <div className="flex items-center justify-between gap-4 flex-wrap">
          {/* Левая часть - поиск и фильтры */}
          <div className="flex items-center gap-2 flex-1 min-w-0">
            {/* Поиск */}
            <div className="relative flex-1 max-w-md">
              <input
                type="text"
                placeholder={t('searchPlaceholder')}
                className="input input-bordered w-full pl-10"
                value={localSearch}
                onChange={(e) => setLocalSearch(e.target.value)}
              />
              <FiSearch className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/50" />
              {localSearch && (
                <button
                  onClick={() => setLocalSearch('')}
                  className="absolute right-3 top-1/2 -translate-y-1/2"
                >
                  <FiX className="w-4 h-4" />
                </button>
              )}
            </div>

            {/* Фильтры */}
            <button
              onClick={() => setShowFilters(!showFilters)}
              className={`btn btn-ghost gap-2 ${showFilters ? 'btn-active' : ''}`}
            >
              <FiFilter className="w-4 h-4" />
              {t('filters')}
              {Object.values(filters).some(
                (v) => v !== null && v !== '' && v !== 'all'
              ) && <span className="badge badge-primary badge-sm">!</span>}
            </button>
          </div>

          {/* Правая часть - режимы отображения и массовый выбор */}
          <div className="flex items-center gap-2">
            {/* Счетчик */}
            {totalCount > 0 && (
              <span className="text-sm text-base-content/70">
                {t('totalProducts', { count: totalCount })}
              </span>
            )}

            {/* Режим выбора */}
            <button
              onClick={handleToggleSelectMode}
              className={`btn gap-2 ${isSelectMode ? 'btn-primary' : 'btn-ghost'}`}
            >
              {isSelectMode ? (
                <FiCheck className="w-4 h-4" />
              ) : (
                <FiMenu className="w-4 h-4" />
              )}
              {t('bulk.selectMode')}
            </button>

            {/* Выбрать все (в режиме выбора) */}
            {isSelectMode && products.length > 0 && (
              <label className="label cursor-pointer gap-2">
                <input
                  type="checkbox"
                  className="checkbox checkbox-primary"
                  checked={allSelected}
                  onChange={handleSelectAll}
                />
                <span className="label-text">{t('bulk.selectAll')}</span>
              </label>
            )}

            <div className="divider divider-horizontal mx-1"></div>

            {/* Режимы отображения */}
            <div className="join">
              <button
                onClick={() => dispatch(setViewMode('grid'))}
                className={`btn btn-sm join-item ${viewMode === 'grid' ? 'btn-active' : ''}`}
              >
                <FiGrid className="w-4 h-4" />
              </button>
              <button
                onClick={() => dispatch(setViewMode('list'))}
                className={`btn btn-sm join-item ${viewMode === 'list' ? 'btn-active' : ''}`}
              >
                <FiList className="w-4 h-4" />
              </button>
              <button
                onClick={() => dispatch(setViewMode('table'))}
                className={`btn btn-sm join-item ${viewMode === 'table' ? 'btn-active' : ''}`}
              >
                <FiMenu className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        {/* Расширенные фильтры */}
        {showFilters && (
          <div className="mt-4 pt-4 border-t border-base-300">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* Статус товара */}
              <div>
                <label className="label">
                  <span className="label-text">{t('bulk.status')}</span>
                </label>
                <select
                  className="select select-bordered w-full"
                  value={
                    filters.isActive === null
                      ? 'all'
                      : filters.isActive
                        ? 'active'
                        : 'inactive'
                  }
                  onChange={(e) => {
                    const value = e.target.value;
                    dispatch(
                      setFilters({
                        isActive: value === 'all' ? null : value === 'active',
                      })
                    );
                  }}
                >
                  <option value="all">{t('allProducts')}</option>
                  <option value="active">{t('activeOnly')}</option>
                  <option value="inactive">{t('inactiveOnly')}</option>
                </select>
              </div>

              {/* Статус склада */}
              <div>
                <label className="label">
                  <span className="label-text">{t('inventory')}</span>
                </label>
                <select
                  className="select select-bordered w-full"
                  value={filters.stockStatus}
                  onChange={(e) => {
                    dispatch(
                      setFilters({
                        stockStatus: e.target.value as any,
                      })
                    );
                  }}
                >
                  <option value="all">{t('allProducts')}</option>
                  <option value="in_stock">{t('stockStatus.in_stock')}</option>
                  <option value="low_stock">
                    {t('stockStatus.low_stock')}
                  </option>
                  <option value="out_of_stock">
                    {t('stockStatus.out_of_stock')}
                  </option>
                </select>
              </div>

              {/* Диапазон цен */}
              <div>
                <label className="label">
                  <span className="label-text">{t('priceRange')}</span>
                </label>
                <div className="flex gap-2">
                  <input
                    type="number"
                    placeholder={t('min')}
                    className="input input-bordered w-full"
                    value={filters.minPrice || ''}
                    onChange={(e) => {
                      dispatch(
                        setFilters({
                          minPrice: e.target.value
                            ? Number(e.target.value)
                            : null,
                        })
                      );
                    }}
                  />
                  <span className="self-center">-</span>
                  <input
                    type="number"
                    placeholder={t('max')}
                    className="input input-bordered w-full"
                    value={filters.maxPrice || ''}
                    onChange={(e) => {
                      dispatch(
                        setFilters({
                          maxPrice: e.target.value
                            ? Number(e.target.value)
                            : null,
                        })
                      );
                    }}
                  />
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Массовые действия */}
      {isSelectMode && selectedCount > 0 && (
        <BulkActions
          selectedCount={selectedCount}
          onBulkDelete={handleBulkDelete}
          onBulkStatusChange={handleBulkStatusChange}
          onBulkExport={handleBulkExport}
          onClearSelection={() => dispatch(clearSelection())}
          isProcessing={isProcessing}
        />
      )}

      {/* Список товаров */}
      {viewMode === 'table' ? (
        <div className="overflow-x-auto bg-base-100 rounded-lg shadow-sm border border-base-300">
          <table className="table">
            <thead>
              <tr>
                {isSelectMode && (
                  <th>
                    <label>
                      <input
                        type="checkbox"
                        className="checkbox"
                        checked={allSelected}
                        onChange={handleSelectAll}
                      />
                    </label>
                  </th>
                )}
                <th>{t('productName')}</th>
                <th>{t('inventory')}</th>
                <th>{t('price')}</th>
                <th>{t('bulk.status')}</th>
                <th>{t('actions')}</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <ProductCard
                  key={product.id}
                  product={product}
                  storefrontSlug={storefrontSlug}
                  isSelected={
                    product.id ? selectedIds.includes(product.id) : false
                  }
                  isSelectMode={isSelectMode}
                  onToggleSelect={handleToggleSelect}
                  viewMode="table"
                />
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div
          className={`
          grid gap-4
          ${viewMode === 'grid' ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4' : 'grid-cols-1'}
        `}
        >
          {products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              storefrontSlug={storefrontSlug}
              isSelected={product.id ? selectedIds.includes(product.id) : false}
              isSelectMode={isSelectMode}
              onToggleSelect={handleToggleSelect}
              viewMode={viewMode}
            />
          ))}
        </div>
      )}

      {/* Сообщение об отсутствии товаров */}
      {!loading && products.length === 0 && (
        <div className="text-center py-12">
          <p className="text-base-content/70">
            {filters.search ||
            Object.values(filters).some(
              (v) => v !== null && v !== '' && v !== 'all'
            )
              ? t('noProductsFound')
              : t('noProducts')}
          </p>
        </div>
      )}

      {/* Индикатор загрузки */}
      {loading && (
        <div className="text-center py-8">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      )}

      {/* Бесконечная прокрутка */}
      {hasMore && onLoadMore && (
        <InfiniteScrollTrigger
          loading={loading}
          hasMore={hasMore}
          onLoadMore={onLoadMore}
          loadMoreText={t('loadMore')}
        />
      )}
    </div>
  );
}
