'use client';

import React, { useState, useEffect, useRef } from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';

// Use generated types from API
type ProductVariant =
  components['schemas']['backend_internal_domain_models.StorefrontProductVariant'];
type BulkUpdateStockRequest =
  components['schemas']['backend_internal_proj_storefront_types.BulkUpdateStockRequest'];

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
  const t = useTranslations('storefronts.variants');

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
      const token = localStorage.getItem('access_token');
      const response = await fetch(
        `/api/v1/storefronts/storefront/products/${productId}/variants`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.ok) {
        const data: ProductVariant[] = await response.json();
        setVariants(data);

        // Extract available attributes for filtering
        const attrs: Record<string, Set<string>> = {};
        data.forEach((variant) => {
          if (variant.attributes) {
            Object.entries(variant.attributes).forEach(([key, value]) => {
              if (!attrs[key]) {
                attrs[key] = new Set();
              }
              attrs[key].add(String(value));
            });
          }
        });
        setAvailableAttributes(attrs);
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
          (variant.name && variant.name.toLowerCase().includes(searchLower)) ||
          (variant.sku && variant.sku.toLowerCase().includes(searchLower)) ||
          (variant.attributes &&
            Object.values(variant.attributes).some((val) =>
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
          variant.attributes &&
          variant.attributes[filters.attribute] === filters.attributeValue
      );
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let aVal: any, bVal: any;

      switch (sort.field) {
        case 'name':
          aVal = a.name || getAttributesDisplay(a) || '';
          bVal = b.name || getAttributesDisplay(b) || '';
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
    if (!variant.attributes) return '';
    return Object.entries(variant.attributes)
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
        return <span className="badge badge-error">{t('outOfStock')}</span>;
      case 'low_stock':
        return <span className="badge badge-warning">{t('lowStock')}</span>;
      case 'in_stock':
        return <span className="badge badge-success">{t('inStock')}</span>;
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

      const token = localStorage.getItem('access_token');
      const response = await fetch(
        `/api/v1/storefronts/storefront/products/${productId}/variants/bulk-update-stock`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.ok) {
        setPendingUpdates(new Map());
        await loadVariants();
      } else {
        const error = await response.json();
        alert(`Failed to update stock: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to update stock:', error);
      alert('Failed to update stock');
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

      const token = localStorage.getItem('access_token');
      const response = await fetch(
        `/api/v1/storefronts/storefront/products/${productId}/variants/import`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        }
      );

      if (response.ok) {
        await loadVariants();
        alert(t('csvImportSuccess'));
      } else {
        const error = await response.json();
        alert(`Import failed: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to import CSV:', error);
      alert('Failed to import CSV');
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
      const token = localStorage.getItem('access_token');
      const response = await fetch(
        `/api/v1/storefronts/storefront/products/${productId}/variants/export`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
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
        alert(`Export failed: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to export CSV:', error);
      alert('Failed to export CSV');
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
        <span className="ml-2">{t('loading')}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with stats */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-title">{t('totalVariants')}</div>
          <div className="stat-value">{variants.length}</div>
          <div className="stat-desc">
            {t('filteredCount', { count: filteredVariants.length })}
          </div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('totalStock')}</div>
          <div className="stat-value">
            {variants.reduce((sum, v) => sum + (v.stock_quantity || 0), 0)}
          </div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('lowStockCount')}</div>
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
          <div className="stat-title">{t('outOfStockCount')}</div>
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
                <span className="label-text">{t('search')}</span>
              </label>
              <input
                type="text"
                placeholder={t('searchPlaceholder')}
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
                <span className="label-text">{t('status')}</span>
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
                <option value="all">{t('allStatuses')}</option>
                <option value="in_stock">{t('inStock')}</option>
                <option value="low_stock">{t('lowStock')}</option>
                <option value="out_of_stock">{t('outOfStock')}</option>
              </select>
            </div>

            {/* Attribute Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('attribute')}</span>
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
                <option value="">{t('allAttributes')}</option>
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
                <span className="label-text">{t('attributeValue')}</span>
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
                <option value="">{t('allValues')}</option>
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
                    {t('selectAll')}
                  </button>
                  <button
                    onClick={clearSelection}
                    className="btn btn-sm btn-outline"
                  >
                    {t('clearSelection')}
                  </button>
                  <span className="text-sm text-gray-600">
                    {selectedVariants.size} {t('selected')}
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
                      <option value="add">{t('add')}</option>
                      <option value="subtract">{t('subtract')}</option>
                      <option value="set">{t('setTo')}</option>
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
                      {t('apply')}
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
                  {importing ? t('importing') : t('importCSV')}
                </button>
                <button
                  onClick={handleCSVExport}
                  disabled={exporting}
                  className={`btn btn-sm btn-outline ${exporting ? 'loading' : ''}`}
                >
                  {exporting ? t('exporting') : t('exportCSV')}
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
          <span>{t('pendingUpdates', { count: pendingUpdates.size })}</span>
          <div className="flex space-x-2">
            <button
              onClick={applyPendingUpdates}
              disabled={updating}
              className={`btn btn-sm btn-success ${updating ? 'loading' : ''}`}
            >
              {updating ? t('saving') : t('saveChanges')}
            </button>
            <button
              onClick={cancelPendingUpdates}
              className="btn btn-sm btn-outline"
            >
              {t('cancel')}
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
                  {t('variant')}
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
                  {t('sku')}
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
                  {t('price')}
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
                  {t('stock')}
                  {sort.field === 'stock' && (
                    <span className="ml-1">
                      {sort.direction === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </button>
              </th>
              <th>{t('status')}</th>
              <th>
                <button
                  onClick={() => handleSort('updated_at')}
                  className="btn btn-ghost btn-sm"
                >
                  {t('lastUpdated')}
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
                      {variant.name || getAttributesDisplay(variant)}
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
              {variants.length === 0 ? t('noVariants') : t('noVariantsMatch')}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {variants.length === 0
                ? t('noVariantsDescription')
                : t('tryDifferentFilters')}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
