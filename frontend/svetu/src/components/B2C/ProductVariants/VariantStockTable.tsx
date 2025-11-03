'use client';

import React, { useState, useEffect, useRef } from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';
import { apiClient } from '@/services/api-client';

// Use generated types from API
type ProductVariant = components['schemas']['models.StorefrontProductVariant'];
type BulkUpdateStockRequest =
  components['schemas']['types.BulkUpdateStockRequest'];

interface VariantStockTableProps {
  productId: number;
  storefrontId: number;
  editable?: boolean;
}

interface StockUpdate {
  variant_id: number;
  stock_quantity: number;
  old_quantity: number;
}

interface FilterState {
  search: string;
  status: 'all' | 'in_stock' | 'low_stock' | 'out_of_stock';
  attribute: string;
  attributeValue: string;
}

interface SortState {
  field: 'name' | 'sku' | 'price' | 'stock' | 'updated_at';
  direction: 'asc' | 'desc';
}

export default function VariantStockTable({
  productId,
  storefrontId: _storefrontId,
  editable = true,
}: VariantStockTableProps) {
  const t = useTranslations('storefronts');

  // State management
  const [variants, setVariants] = useState<ProductVariant[]>([]);
  const [filteredVariants, setFilteredVariants] = useState<ProductVariant[]>(
    []
  );
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState(false);
  const [selectedVariants, setSelectedVariants] = useState<Set<number>>(
    new Set()
  );
  const [pendingUpdates, setPendingUpdates] = useState<
    Map<number, StockUpdate>
  >(new Map());

  // Filter and sort state
  const [filters, setFilters] = useState<FilterState>({
    search: '',
    status: 'all',
    attribute: '',
    attributeValue: '',
  });
  const [sort, setSort] = useState<SortState>({
    field: 'name',
    direction: 'asc',
  });

  // CSV import/export state
  const [importing, setImporting] = useState(false);
  const [exporting, setExporting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Bulk operations state
  const [bulkOperation, setBulkOperation] = useState<
    'add' | 'subtract' | 'set'
  >('add');
  const [bulkValue, setBulkValue] = useState<number>(0);

  // Available attributes for filtering
  const [availableAttributes, setAvailableAttributes] = useState<
    Record<string, Set<string>>
  >({});

  // Load variants
  const loadVariants = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get<ProductVariant[]>(
        `/marketplace/storefronts/storefront/products/${productId}/variants`
      );

      if (response.data) {
        const data = response.data;
        setVariants(data);

        // Extract available attributes for filtering
        const attrs: Record<string, Set<string>> = {};
        data.forEach((variant) => {
          if (variant.variant_attributes) {
            Object.entries(variant.variant_attributes).forEach(
              ([key, value]) => {
                if (!attrs[key]) {
                  attrs[key] = new Set();
                }
                attrs[key].add(String(value));
              }
            );
          }
        });
        setAvailableAttributes(attrs);
      } else if (response.error) {
        console.error('Failed to load variants:', response.error.message);
      }
    } catch (error) {
      console.error('Failed to load variants:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadVariants();
  }, [productId]); // eslint-disable-line react-hooks/exhaustive-deps

  // Apply filters and sorting
  useEffect(() => {
    let filtered = [...variants];

    // Apply search filter
    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      filtered = filtered.filter(
        (variant) =>
          (variant?.sku && variant?.sku.toLowerCase().includes(searchLower)) ||
          (variant.sku && variant.sku.toLowerCase().includes(searchLower)) ||
          (variant.variant_attributes &&
            Object.values(variant.variant_attributes).some((val) =>
              String(val).toLowerCase().includes(searchLower)
            ))
      );
    }

    // Apply status filter
    if (filters.status !== 'all') {
      filtered = filtered.filter((variant) => {
        const stock = variant.stock_quantity || 0;
        switch (filters.status) {
          case 'out_of_stock':
            return stock === 0;
          case 'low_stock':
            return stock > 0 && stock < 5;
          case 'in_stock':
            return stock >= 5;
          default:
            return true;
        }
      });
    }

    // Apply attribute filter
    if (filters.attribute && filters.attributeValue) {
      filtered = filtered.filter(
        (variant) =>
          variant.variant_attributes &&
          variant.variant_attributes[filters.attribute] ===
            filters.attributeValue
      );
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let aVal: any, bVal: any;

      switch (sort.field) {
        case 'name':
          aVal = a.sku || getAttributesDisplay(a) || '';
          bVal = b.sku || getAttributesDisplay(b) || '';
          break;
        case 'sku':
          aVal = a.sku || '';
          bVal = b.sku || '';
          break;
        case 'price':
          aVal = a.price || 0;
          bVal = b.price || 0;
          break;
        case 'stock':
          aVal = a.stock_quantity || 0;
          bVal = b.stock_quantity || 0;
          break;
        case 'updated_at':
          aVal = new Date(a.updated_at || 0).getTime();
          bVal = new Date(b.updated_at || 0).getTime();
          break;
        default:
          aVal = 0;
          bVal = 0;
      }

      if (aVal < bVal) return sort.direction === 'asc' ? -1 : 1;
      if (aVal > bVal) return sort.direction === 'asc' ? 1 : -1;
      return 0;
    });

    setFilteredVariants(filtered);
  }, [variants, filters, sort]);

  // Utility functions
  const getAttributesDisplay = (variant: ProductVariant): string => {
    if (!variant.variant_attributes) return '';
    return Object.entries(variant.variant_attributes)
      .map(([key, value]) => `${key}: ${value}`)
      .join(', ');
  };

  const getStockStatus = (
    stock: number
  ): 'in_stock' | 'low_stock' | 'out_of_stock' => {
    if (stock === 0) return 'out_of_stock';
    if (stock < 5) return 'low_stock';
    return 'in_stock';
  };

  const getStockStatusBadge = (stock: number) => {
    const status = getStockStatus(stock);
    switch (status) {
      case 'out_of_stock':
        return (
          <span className="badge badge-error">{t('variants.outOfStock')}</span>
        );
      case 'low_stock':
        return (
          <span className="badge badge-warning">{t('variants.lowStock')}</span>
        );
      case 'in_stock':
        return (
          <span className="badge badge-success">{t('variants.inStock')}</span>
        );
    }
  };

  // Selection handlers
  const toggleVariantSelection = (variantId: number) => {
    const newSelected = new Set(selectedVariants);
    if (newSelected.has(variantId)) {
      newSelected.delete(variantId);
    } else {
      newSelected.add(variantId);
    }
    setSelectedVariants(newSelected);
  };

  const selectAll = () => {
    setSelectedVariants(new Set(filteredVariants.map((v) => v.id!)));
  };

  const clearSelection = () => {
    setSelectedVariants(new Set());
  };

  // Stock update handlers
  const updateStock = (variantId: number, newStock: number) => {
    const variant = variants.find((v) => v.id === variantId);
    if (!variant) return;

    const update: StockUpdate = {
      variant_id: variantId,
      stock_quantity: Math.max(0, newStock),
      old_quantity: variant.stock_quantity || 0,
    };

    setPendingUpdates((prev) => new Map(prev).set(variantId, update));
  };

  const applyPendingUpdates = async () => {
    if (pendingUpdates.size === 0) return;

    try {
      setUpdating(true);

      const updates = Array.from(pendingUpdates.values()).map((update) => ({
        variant_id: update.variant_id,
        stock_quantity: update.stock_quantity,
      }));

      const request: BulkUpdateStockRequest = { updates };

      const response = await apiClient.post(
        `/marketplace/storefronts/storefront/products/${productId}/variants/bulk-update-stock`,
        request
      );

      if (response.data) {
        setPendingUpdates(new Map());
        await loadVariants();
      } else if (response.error) {
        alert(`Failed to update stock: ${response.error.message}`);
      }
    } catch (error) {
      console.error('Failed to update stock:', error);
      alert('variants.Failed to update stock');
    } finally {
      setUpdating(false);
    }
  };

  const cancelPendingUpdates = () => {
    setPendingUpdates(new Map());
  };

  // Bulk operations
  const applyBulkOperation = () => {
    if (selectedVariants.size === 0 || bulkValue === 0) return;

    selectedVariants.forEach((variantId) => {
      const variant = variants.find((v) => v.id === variantId);
      if (!variant) return;

      let newStock = variant.stock_quantity || 0;
      switch (bulkOperation) {
        case 'add':
          newStock += bulkValue;
          break;
        case 'subtract':
          newStock = Math.max(0, newStock - bulkValue);
          break;
        case 'set':
          newStock = bulkValue;
          break;
      }

      updateStock(variantId, newStock);
    });

    // Clear selection and bulk value
    setSelectedVariants(new Set());
    setBulkValue(0);
  };

  // CSV Import/Export
  const handleCSVImport = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      setImporting(true);
      const formData = new FormData();
      formData.append('file', file);

      const response = await apiClient.upload(
        `/marketplace/storefronts/storefront/products/${productId}/variants/import`,
        formData
      );

      if (response.data) {
        await loadVariants();
        alert(t('variants.csvImportSuccess'));
      } else if (response.error) {
        alert(`Import failed: ${response.error.message}`);
      }
    } catch (error) {
      console.error('Failed to import CSV:', error);
      alert('variants.Failed to import CSV');
    } finally {
      setImporting(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleCSVExport = async () => {
    try {
      setExporting(true);

      // Используем BFF proxy для безопасного доступа (credentials автоматически)
      const response = await fetch(
        `/api/v2/b2c/storefront/products/${productId}/variants/export`,
        {
          credentials: 'include', // JWT cookies передаются автоматически
        }
      );

      if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = `product-${productId}-variants.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
      } else {
        const error = await response.json();
        alert(`Export failed: ${error.error || error.message}`);
      }
    } catch (error) {
      console.error('Failed to export CSV:', error);
      alert('variants.Failed to export CSV');
    } finally {
      setExporting(false);
    }
  };

  // Sort handler
  const handleSort = (field: SortState['field']) => {
    setSort((prev) => ({
      field,
      direction:
        prev.field === field && prev.direction === 'asc' ? 'desc' : 'asc',
    }));
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2">{t('variants.loading')}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with stats */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-title">{t('variants.totalVariants')}</div>
          <div className="stat-value">{variants.length}</div>
          <div className="stat-desc">
            {t('variants.filteredCount', { count: filteredVariants.length })}
          </div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.totalStock')}</div>
          <div className="stat-value">
            {variants.reduce((sum, v) => sum + (v.stock_quantity || 0), 0)}
          </div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.lowStockCount')}</div>
          <div className="stat-value text-warning">
            {
              variants.filter(
                (v) =>
                  (v.stock_quantity || 0) > 0 && (v.stock_quantity || 0) < 5
              ).length
            }
          </div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.outOfStockCount')}</div>
          <div className="stat-value text-error">
            {variants.filter((v) => (v.stock_quantity || 0) === 0).length}
          </div>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {/* Search */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('variants.search')}</span>
              </label>
              <input
                type="text"
                placeholder={t('variants.searchPlaceholder')}
                className="input input-bordered"
                value={filters.search}
                onChange={(e) =>
                  setFilters((prev) => ({ ...prev, search: e.target.value }))
                }
              />
            </div>

            {/* Status Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('variants.status')}</span>
              </label>
              <select
                className="select select-bordered"
                value={filters.status}
                onChange={(e) =>
                  setFilters((prev) => ({
                    ...prev,
                    status: e.target.value as FilterState['status'],
                  }))
                }
              >
                <option value="all">{t('variants.allStatuses')}</option>
                <option value="in_stock">{t('variants.inStock')}</option>
                <option value="low_stock">{t('variants.lowStock')}</option>
                <option value="out_of_stock">{t('variants.outOfStock')}</option>
              </select>
            </div>

            {/* Attribute Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('variants.attribute')}</span>
              </label>
              <select
                className="select select-bordered"
                value={filters.attribute}
                onChange={(e) =>
                  setFilters((prev) => ({
                    ...prev,
                    attribute: e.target.value,
                    attributeValue: '',
                  }))
                }
              >
                <option value="">{t('variants.allAttributes')}</option>
                {Object.keys(availableAttributes).map((attr) => (
                  <option key={attr} value={attr}>
                    {attr}
                  </option>
                ))}
              </select>
            </div>

            {/* Attribute Value Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('variants.attributeValue')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={filters.attributeValue}
                onChange={(e) =>
                  setFilters((prev) => ({
                    ...prev,
                    attributeValue: e.target.value,
                  }))
                }
                disabled={!filters.attribute}
              >
                <option value="">{t('variants.allValues')}</option>
                {filters.attribute &&
                  availableAttributes[filters.attribute] &&
                  Array.from(availableAttributes[filters.attribute]).map(
                    (value) => (
                      <option key={value} value={value}>
                        {value}
                      </option>
                    )
                  )}
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Bulk Operations */}
      {editable && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <div className="flex flex-wrap justify-between items-center gap-4">
              <div className="flex flex-wrap items-center gap-4">
                {/* Selection controls */}
                <div className="flex space-x-2">
                  <button
                    onClick={selectAll}
                    className="btn btn-sm btn-outline"
                  >
                    {t('variants.selectAll')}
                  </button>
                  <button
                    onClick={clearSelection}
                    className="btn btn-sm btn-outline"
                  >
                    {t('variants.clearSelection')}
                  </button>
                  <span className="text-sm text-gray-600">
                    {selectedVariants.size} {t('variants.selected')}
                  </span>
                </div>

                {/* Bulk operations */}
                {selectedVariants.size > 0 && (
                  <div className="flex items-center space-x-2">
                    <select
                      className="select select-sm select-bordered"
                      value={bulkOperation}
                      onChange={(e) =>
                        setBulkOperation(
                          e.target.value as 'add' | 'subtract' | 'set'
                        )
                      }
                    >
                      <option value="add">{t('variants.add')}</option>
                      <option value="subtract">{t('variants.subtract')}</option>
                      <option value="set">{t('variants.setTo')}</option>
                    </select>
                    <input
                      type="number"
                      min="0"
                      className="input input-sm input-bordered w-20"
                      value={bulkValue}
                      onChange={(e) =>
                        setBulkValue(parseInt(e.target.value) || 0)
                      }
                    />
                    <button
                      onClick={applyBulkOperation}
                      disabled={bulkValue === 0}
                      className="btn btn-sm btn-primary"
                    >
                      {t('variants.apply')}
                    </button>
                  </div>
                )}
              </div>

              {/* CSV Operations */}
              <div className="flex space-x-2">
                <input
                  ref={fileInputRef}
                  type="file"
                  accept=".csv"
                  onChange={handleCSVImport}
                  className="hidden"
                />
                <button
                  onClick={() => fileInputRef.current?.click()}
                  disabled={importing}
                  className={`btn btn-sm btn-outline ${importing ? 'loading' : ''}`}
                >
                  {importing
                    ? t('variants.importing')
                    : t('variants.importCSV')}
                </button>
                <button
                  onClick={handleCSVExport}
                  disabled={exporting}
                  className={`btn btn-sm btn-outline ${exporting ? 'loading' : ''}`}
                >
                  {exporting
                    ? t('variants.exporting')
                    : t('variants.exportCSV')}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Pending Updates */}
      {pendingUpdates.size > 0 && (
        <div className="alert alert-warning">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 18.5c-.77.833.192 2.5 1.732 2.5z"
            ></path>
          </svg>
          <span>
            {t('variants.pendingUpdates', { count: pendingUpdates.size })}
          </span>
          <div className="flex space-x-2">
            <button
              onClick={applyPendingUpdates}
              disabled={updating}
              className={`btn btn-sm btn-success ${updating ? 'loading' : ''}`}
            >
              {updating ? t('variants.saving') : t('variants.saveChanges')}
            </button>
            <button
              onClick={cancelPendingUpdates}
              className="btn btn-sm btn-outline"
            >
              {t('variants.cancel')}
            </button>
          </div>
        </div>
      )}

      {/* Variants Table */}
      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              {editable && (
                <th>
                  <input
                    type="checkbox"
                    className="checkbox"
                    checked={
                      selectedVariants.size === filteredVariants.length &&
                      filteredVariants.length > 0
                    }
                    onChange={() =>
                      selectedVariants.size === filteredVariants.length
                        ? clearSelection()
                        : selectAll()
                    }
                  />
                </th>
              )}
              <th>
                <button
                  onClick={() => handleSort('name')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('variants.variant')}
                  {sort.field === 'name' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
              <th>
                <button
                  onClick={() => handleSort('sku')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('variants.sku')}
                  {sort.field === 'sku' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
              <th>
                <button
                  onClick={() => handleSort('price')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('variants.price')}
                  {sort.field === 'price' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
              <th>
                <button
                  onClick={() => handleSort('stock')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('variants.stock')}
                  {sort.field === 'stock' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
              <th>{t('variants.status')}</th>
              <th>
                <button
                  onClick={() => handleSort('updated_at')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('variants.lastUpdated')}
                  {sort.field === 'updated_at' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {filteredVariants.map((variant) => {
              const pendingUpdate = pendingUpdates.get(variant.id!);
              const currentStock =
                pendingUpdate?.stock_quantity ?? variant.stock_quantity ?? 0;

              return (
                <tr
                  key={variant.id}
                  className={pendingUpdate ? 'bg-warning/10' : ''}
                >
                  {editable && (
                    <td>
                      <input
                        type="checkbox"
                        className="checkbox"
                        checked={selectedVariants.has(variant.id!)}
                        onChange={() => toggleVariantSelection(variant.id!)}
                      />
                    </td>
                  )}
                  <td>
                    <div className="font-medium">
                      {variant?.sku || getAttributesDisplay(variant)}
                    </div>
                  </td>
                  <td>
                    <code className="text-sm">{variant.sku || '-'}</code>
                  </td>
                  <td>{variant.price ? `${variant.price} RSD` : '-'}</td>
                  <td>
                    {editable ? (
                      <div className="flex items-center space-x-2">
                        <input
                          type="number"
                          min="0"
                          value={currentStock}
                          onChange={(e) =>
                            updateStock(
                              variant.id!,
                              parseInt(e.target.value) || 0
                            )
                          }
                          className={`input input-sm input-bordered w-20 ${
                            pendingUpdate ? 'input-warning' : ''
                          }`}
                        />
                        {pendingUpdate && (
                          <span className="text-xs text-warning">
                            ({pendingUpdate.old_quantity} →{' '}
                            {pendingUpdate.stock_quantity})
                          </span>
                        )}
                      </div>
                    ) : (
                      <span className="font-mono">{currentStock}</span>
                    )}
                  </td>
                  <td>{getStockStatusBadge(currentStock)}</td>
                  <td>
                    <span className="text-sm text-gray-500">
                      {variant.updated_at
                        ? new Date(variant.updated_at).toLocaleDateString()
                        : '-'}
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {filteredVariants.length === 0 && (
        <div className="text-center py-8">
          <div className="text-gray-500">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2M4 13h2m13-8a2 2 0 11-4 0 2 2 0 014 0zM9 13h1m-1 4h1m6-4h1m-1 4h1"
              />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">
              {variants.length === 0
                ? t('variants.noVariants')
                : t('variants.noVariantsMatch')}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {variants.length === 0
                ? t('variants.noVariantsDescription')
                : t('variants.tryDifferentFilters')}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
