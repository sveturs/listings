'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';
import configManager from '@/config';

// Use generated types from API
type ProductVariant =
  components['schemas']['backend_internal_domain_models.StorefrontProductVariant'];
type VariantMatrixResponse =
  components['schemas']['backend_internal_proj_storefront_types.VariantMatrixResponse'];
type BulkUpdateStockRequest =
  components['schemas']['backend_internal_proj_storefront_types.BulkUpdateStockRequest'];
type VariantAnalyticsResponse =
  components['schemas']['backend_internal_proj_storefront_types.VariantAnalyticsResponse'];

interface VariantManagerProps {
  productId: number;
  storefrontId: number;
  onSave: () => void;
  onCancel: () => void;
}

interface AttributeValue {
  id: number;
  value: string;
  display_name: string;
  color_hex?: string;
  image_url?: string;
  is_popular: boolean;
}

interface ProductVariantAttribute {
  id: number;
  name: string;
  display_name: string;
  type: string;
  is_required: boolean;
  affects_stock: boolean;
  sort_order: number;
}

export default function VariantManager({
  productId,
  storefrontId: _storefrontId,
  onSave,
  onCancel,
}: VariantManagerProps) {
  const t = useTranslations('storefronts');

  // State management
  const [variants, setVariants] = useState<ProductVariant[]>([]);
  const [attributes, setAttributes] = useState<ProductVariantAttribute[]>([]);
  const [attributeValues, setAttributeValues] = useState<
    Record<number, AttributeValue[]>
  >({});
  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<string, string[]>
  >({});
  const [_variantMatrix, setVariantMatrix] =
    useState<VariantMatrixResponse | null>(null);
  const [analytics, setAnalytics] = useState<VariantAnalyticsResponse | null>(
    null
  );

  // UI state
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [activeTab, setActiveTab] = useState<
    'attributes' | 'variants' | 'stock' | 'analytics'
  >('attributes');
  const [selectedVariants, setSelectedVariants] = useState<Set<number>>(
    new Set()
  );

  // Load data functions
  const loadVariants = useCallback(async () => {
    if (productId <= 0) return; // Skip for new products
    try {
      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefronts/storefront/products/${productId}/variants`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setVariants(data);
      }
    } catch (error) {
      console.error('Failed to load variants:', error);
    }
  }, [productId]);

  const loadAttributes = useCallback(async () => {
    try {
      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(`${apiUrl}/api/v1/public/variants/attributes`, {
        headers: {
          Authorization: token ? `Bearer ${token}` : '',
        },
      });
      if (response.ok) {
        const data = await response.json();
        setAttributes(
          data.filter((attr: ProductVariantAttribute) => attr.affects_stock)
        );

        // Load values for each attribute
        for (const attr of data) {
          await loadAttributeValues(attr.id);
        }
      }
    } catch (error) {
      console.error('Failed to load attributes:', error);
    }
  }, []);

  const loadAttributeValues = async (attributeId: number) => {
    try {
      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/public/variants/attributes/${attributeId}/values`,
        {
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        }
      );
      if (response.ok) {
        const values = await response.json();
        setAttributeValues((prev) => ({
          ...prev,
          [attributeId]: values,
        }));
      }
    } catch (error) {
      console.error(
        `Failed to load values for attribute ${attributeId}:`,
        error
      );
    }
  };

  const loadVariantMatrix = useCallback(async () => {
    if (productId <= 0) return; // Skip for new products
    try {
      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefronts/storefront/products/${productId}/variant-matrix`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setVariantMatrix(data);
      }
    } catch (error) {
      console.error('Failed to load variant matrix:', error);
    }
  }, [productId]);

  const loadAnalytics = useCallback(async () => {
    if (productId <= 0) return; // Skip for new products
    try {
      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefronts/storefront/products/${productId}/variants/analytics`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setAnalytics(data);
      }
    } catch (error) {
      console.error('Failed to load analytics:', error);
    }
  }, [productId]);

  // Initialize data
  useEffect(() => {
    const initializeData = async () => {
      setLoading(true);
      // Only load variants-related data if we have a valid productId
      if (productId > 0) {
        await Promise.all([
          loadVariants(),
          loadAttributes(),
          loadVariantMatrix(),
          loadAnalytics(),
        ]);
      } else {
        // For new products, only load attributes
        await loadAttributes();
      }
      setLoading(false);
    };

    initializeData();
  }, [
    productId,
    loadVariants,
    loadAttributes,
    loadVariantMatrix,
    loadAnalytics,
  ]);

  // Generate variants
  const generateVariants = async () => {
    try {
      setSaving(true);

      const request = {
        product_id: productId,
        attribute_matrix: selectedAttributes,
        default_stock: 10,
        base_price: 0,
      };

      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefronts/storefront/variants/generate`,
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
        await loadVariants();
        await loadVariantMatrix();
        setActiveTab('variants');
      } else {
        const error = await response.json();
        alert(`Failed to generate variants: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to generate variants:', error);
      alert('variants.Failed to generate variants');
    } finally {
      setSaving(false);
    }
  };

  // Bulk update stock
  const bulkUpdateStock = async (
    updates: { variant_id: number; stock_quantity: number }[]
  ) => {
    try {
      const request: BulkUpdateStockRequest = {
        updates: updates.map((update) => ({
          variant_id: update.variant_id,
          stock_quantity: update.stock_quantity,
        })),
      };

      const token = localStorage.getItem('access_token');
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefronts/storefront/products/${productId}/variants/bulk-update-stock`,
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
        await loadVariants();
        await loadAnalytics();
      } else {
        const error = await response.json();
        alert(`Failed to update stock: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to update stock:', error);
      alert('variants.Failed to update stock');
    }
  };

  // Attribute selection handlers
  const toggleAttributeValue = (attributeName: string, value: string) => {
    setSelectedAttributes((prev) => {
      const current = prev[attributeName] || [];
      const updated = current.includes(value)
        ? current.filter((v) => v !== value)
        : [...current, value];

      return {
        ...prev,
        [attributeName]: updated,
      };
    });
  };

  const getSelectedCombinationsCount = () => {
    const values = Object.values(selectedAttributes);
    return values.reduce((total, vals) => total * Math.max(vals.length, 1), 1);
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
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">{t('variants.title')}</h2>
          <p className="text-sm text-gray-600">
            {t('variants.subtitle', { count: variants.length })}
          </p>
        </div>
        <div className="flex space-x-2">
          <button onClick={onCancel} className="btn btn-outline">
            {t('variants.cancel')}
          </button>
          <button onClick={onSave} className="btn btn-primary">
            {t('variants.save')}
          </button>
        </div>
      </div>

      {/* Tab Navigation */}
      <div className="tabs tabs-boxed">
        <button
          className={`tab ${activeTab === 'attributes' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('attributes')}
        >
          {t('tabs.attributes')}
        </button>
        <button
          className={`tab ${activeTab === 'variants' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('variants')}
        >
          {t('tabs.variants')} ({variants.length})
        </button>
        <button
          className={`tab ${activeTab === 'stock' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('stock')}
        >
          {t('tabs.stock')}
        </button>
        <button
          className={`tab ${activeTab === 'analytics' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('analytics')}
        >
          {t('tabs.analytics')}
        </button>
      </div>

      {/* Tab Content */}
      {activeTab === 'attributes' && (
        <AttributeSelectionTab
          attributes={attributes}
          attributeValues={attributeValues}
          selectedAttributes={selectedAttributes}
          onToggleValue={toggleAttributeValue}
          onGenerate={generateVariants}
          generating={saving}
          combinationsCount={getSelectedCombinationsCount()}
          t={t}
        />
      )}

      {activeTab === 'variants' && (
        <VariantListTab
          variants={variants}
          selectedVariants={selectedVariants}
          onSelectionChange={setSelectedVariants}
          onVariantUpdate={loadVariants}
          t={t}
        />
      )}

      {activeTab === 'stock' && (
        <StockManagementTab
          variants={variants}
          onBulkUpdate={bulkUpdateStock}
          analytics={analytics}
          t={t}
        />
      )}

      {activeTab === 'analytics' && (
        <AnalyticsTab analytics={analytics} variants={variants} t={t} />
      )}
    </div>
  );
}

// Attribute Selection Tab Component
interface AttributeSelectionTabProps {
  attributes: ProductVariantAttribute[];
  attributeValues: Record<number, AttributeValue[]>;
  selectedAttributes: Record<string, string[]>;
  onToggleValue: (attributeName: string, value: string) => void;
  onGenerate: () => void;
  generating: boolean;
  combinationsCount: number;
  t: (key: string, params?: any) => string;
}

function AttributeSelectionTab({
  attributes,
  attributeValues,
  selectedAttributes,
  onToggleValue,
  onGenerate,
  generating,
  combinationsCount,
  t,
}: AttributeSelectionTabProps) {
  return (
    <div className="space-y-6">
      <div className="alert alert-info">
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
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <span>{t('variants.selectAttributesHelp')}</span>
      </div>

      {attributes.map((attribute) => {
        const values = attributeValues[attribute.id] || [];
        const selected = selectedAttributes[attribute.name] || [];

        return (
          <div key={attribute.id} className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <div className="flex justify-between items-center">
                <h3 className="card-title">
                  {attribute.display_name}
                  {attribute.is_required && (
                    <span className="badge badge-error">Required</span>
                  )}
                  {attribute.affects_stock && (
                    <span className="badge badge-success">Affects Stock</span>
                  )}
                </h3>
                <span className="badge badge-outline">
                  {selected.length} / {values.length} selected
                </span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 mt-4">
                {values.map((value) => {
                  const isSelected = selected.includes(value.value);

                  return (
                    <label
                      key={value.id}
                      className={`
                        flex items-center space-x-2 p-3 border-2 rounded-lg cursor-pointer transition-colors
                        ${
                          isSelected
                            ? 'border-primary bg-primary/10'
                            : 'border-base-300 hover:border-base-400'
                        }
                      `}
                    >
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() =>
                          onToggleValue(attribute.name, value.value)
                        }
                        className="checkbox checkbox-primary checkbox-sm"
                      />

                      {attribute.type === 'color' && value.color_hex && (
                        <div
                          className="w-4 h-4 rounded border border-base-300"
                          style={{ backgroundColor: value.color_hex }}
                        />
                      )}

                      <span className="text-sm font-medium">
                        {value.display_name}
                      </span>

                      {value.is_popular && (
                        <span className="badge badge-xs badge-warning">
                          Popular
                        </span>
                      )}
                    </label>
                  );
                })}
              </div>
            </div>
          </div>
        );
      })}

      {/* Generation Summary */}
      {combinationsCount > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h3 className="card-title">{t('variants.generationSummary')}</h3>
            <div className="stats shadow">
              <div className="stat">
                <div className="stat-title">
                  {t('variants.totalCombinations')}
                </div>
                <div className="stat-value text-primary">
                  {combinationsCount}
                </div>
                <div className="stat-desc">
                  {t('variants.variantsWillBeCreated', {
                    count: combinationsCount,
                  })}
                </div>
              </div>
            </div>

            <div className="card-actions justify-end">
              <button
                onClick={onGenerate}
                disabled={combinationsCount === 0 || generating}
                className={`btn btn-primary ${generating ? 'loading' : ''}`}
              >
                {generating
                  ? t('variants.generating')
                  : t('variants.generateVariants')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Variant List Tab Component
interface VariantListTabProps {
  variants: ProductVariant[];
  selectedVariants: Set<number>;
  onSelectionChange: (selected: Set<number>) => void;
  onVariantUpdate: () => void;
  t: (key: string, params?: any) => string;
}

function VariantListTab({
  variants,
  selectedVariants,
  onSelectionChange,
  onVariantUpdate,
  t,
}: VariantListTabProps) {
  const toggleVariantSelection = (variantId: number) => {
    const newSelected = new Set(selectedVariants);
    if (newSelected.has(variantId)) {
      newSelected.delete(variantId);
    } else {
      newSelected.add(variantId);
    }
    onSelectionChange(newSelected);
  };

  const selectAll = () => {
    onSelectionChange(new Set(variants.map((v) => v.id!)));
  };

  const clearSelection = () => {
    onSelectionChange(new Set());
  };

  return (
    <div className="space-y-4">
      {/* Bulk Actions */}
      <div className="flex justify-between items-center">
        <div className="flex space-x-2">
          <button onClick={selectAll} className="btn btn-sm btn-outline">
            {t('variants.selectAll')}
          </button>
          <button onClick={clearSelection} className="btn btn-sm btn-outline">
            {t('variants.clearSelection')}
          </button>
        </div>
        <span className="text-sm text-gray-600">
          {selectedVariants.size} / {variants.length} {t('variants.selected')}
        </span>
      </div>

      {/* Variants Table */}
      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th>
                <input
                  type="checkbox"
                  className="checkbox"
                  checked={
                    selectedVariants.size === variants.length &&
                    variants.length > 0
                  }
                  onChange={() =>
                    selectedVariants.size === variants.length
                      ? clearSelection()
                      : selectAll()
                  }
                />
              </th>
              <th>{t('variants.variant')}</th>
              <th>{t('variants.sku')}</th>
              <th>{t('variants.price')}</th>
              <th>{t('variants.stock')}</th>
              <th>{t('variants.status')}</th>
              <th>{t('variants.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {variants.map((variant) => (
              <VariantRow
                key={variant.id}
                variant={variant}
                isSelected={selectedVariants.has(variant.id!)}
                onToggleSelection={() => toggleVariantSelection(variant.id!)}
                onUpdate={onVariantUpdate}
                t={t}
              />
            ))}
          </tbody>
        </table>
      </div>

      {variants.length === 0 && (
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
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10"
              />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">
              {t('variants.noVariants')}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {t('variants.noVariantsDescription')}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}

// Variant Row Component
interface VariantRowProps {
  variant: ProductVariant;
  isSelected: boolean;
  onToggleSelection: () => void;
  onUpdate: () => void;
  t: (key: string, params?: any) => string;
}

function VariantRow({
  variant,
  isSelected,
  onToggleSelection,
  onUpdate,
  t,
}: VariantRowProps) {
  const [editing, setEditing] = useState(false);
  const [editValues, setEditValues] = useState({
    price: variant.price || 0,
    stock_quantity: variant.stock_quantity || 0,
    sku: variant.sku || '',
  });

  const saveChanges = async () => {
    try {
      const apiUrl = configManager.get('api.url');
      const response = await fetch(
        `${apiUrl}/api/v1/storefront/variants/${variant.id}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(editValues),
        }
      );

      if (response.ok) {
        setEditing(false);
        onUpdate();
      } else {
        const error = await response.json();
        alert(`Failed to update variant: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to update variant:', error);
      alert('variants.Failed to update variant');
    }
  };

  const getAttributesDisplay = () => {
    if (!variant.attributes) return '-';

    return Object.entries(variant.attributes)
      .map(([key, value]) => `${key}: ${value}`)
      .join(', ');
  };

  const getStockStatusBadge = () => {
    const stock = variant.stock_quantity || 0;
    if (stock === 0) {
      return (
        <span className="badge badge-error">{t('variants.outOfStock')}</span>
      );
    } else if (stock < 5) {
      return (
        <span className="badge badge-warning">{t('variants.lowStock')}</span>
      );
    } else {
      return (
        <span className="badge badge-success">{t('variants.inStock')}</span>
      );
    }
  };

  return (
    <tr>
      <td>
        <input
          type="checkbox"
          className="checkbox"
          checked={isSelected}
          onChange={onToggleSelection}
        />
      </td>
      <td>
        <div className="font-medium">
          {variant.name || getAttributesDisplay()}
        </div>
      </td>
      <td>
        {editing ? (
          <input
            type="text"
            value={editValues.sku}
            onChange={(e) =>
              setEditValues({ ...editValues, sku: e.target.value })
            }
            className="input input-sm input-bordered w-full max-w-xs"
          />
        ) : (
          <code className="text-sm">{variant.sku || '-'}</code>
        )}
      </td>
      <td>
        {editing ? (
          <input
            type="number"
            step="0.01"
            value={editValues.price}
            onChange={(e) =>
              setEditValues({
                ...editValues,
                price: parseFloat(e.target.value) || 0,
              })
            }
            className="input input-sm input-bordered w-full max-w-xs"
          />
        ) : (
          <span>{variant.price ? `${variant.price} RSD` : '-'}</span>
        )}
      </td>
      <td>
        {editing ? (
          <input
            type="number"
            min="0"
            value={editValues.stock_quantity}
            onChange={(e) =>
              setEditValues({
                ...editValues,
                stock_quantity: parseInt(e.target.value) || 0,
              })
            }
            className="input input-sm input-bordered w-full max-w-xs"
          />
        ) : (
          <span>{variant.stock_quantity || 0}</span>
        )}
      </td>
      <td>{getStockStatusBadge()}</td>
      <td>
        <div className="flex space-x-1">
          {editing ? (
            <>
              <button onClick={saveChanges} className="btn btn-sm btn-success">
                {t('variants.save')}
              </button>
              <button
                onClick={() => setEditing(false)}
                className="btn btn-sm btn-outline"
              >
                {t('variants.cancel')}
              </button>
            </>
          ) : (
            <button
              onClick={() => setEditing(true)}
              className="btn btn-sm btn-outline"
            >
              {t('variants.edit')}
            </button>
          )}
        </div>
      </td>
    </tr>
  );
}

// Stock Management Tab Component
interface StockManagementTabProps {
  variants: ProductVariant[];
  onBulkUpdate: (
    updates: { variant_id: number; stock_quantity: number }[]
  ) => void;
  analytics: VariantAnalyticsResponse | null;
  t: (key: string, params?: any) => string;
}

function StockManagementTab({
  variants,
  onBulkUpdate,
  analytics,
  t,
}: StockManagementTabProps) {
  const [bulkOperation, setBulkOperation] = useState<
    'add' | 'set' | 'subtract'
  >('add');
  const [bulkValue, setBulkValue] = useState(0);
  const [selectedForBulk, setSelectedForBulk] = useState<Set<number>>(
    new Set()
  );

  const applyBulkUpdate = () => {
    if (selectedForBulk.size === 0 || bulkValue === 0) return;

    const updates = Array.from(selectedForBulk)
      .map((variantId) => {
        const variant = variants.find((v) => v.id === variantId);
        if (!variant) return null;

        let newStock = 0;
        switch (bulkOperation) {
          case 'add':
            newStock = (variant.stock_quantity || 0) + bulkValue;
            break;
          case 'subtract':
            newStock = Math.max(0, (variant.stock_quantity || 0) - bulkValue);
            break;
          case 'set':
            newStock = bulkValue;
            break;
        }

        return {
          variant_id: variantId,
          stock_quantity: newStock,
        };
      })
      .filter(Boolean) as { variant_id: number; stock_quantity: number }[];

    onBulkUpdate(updates);
    setSelectedForBulk(new Set());
    setBulkValue(0);
  };

  const lowStockVariants = variants.filter((v) => (v.stock_quantity || 0) < 5);
  const outOfStockVariants = variants.filter(
    (v) => (v.stock_quantity || 0) === 0
  );

  return (
    <div className="space-y-6">
      {/* Stock Overview */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-figure text-primary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">{t('variants.totalStock')}</div>
          <div className="stat-value text-primary">
            {analytics?.total_stock ||
              variants.reduce((sum, v) => sum + (v.stock_quantity || 0), 0)}
          </div>
          <div className="stat-desc">{t('variants.acrossAllVariants')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-warning">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 18.5c-.77.833.192 2.5 1.732 2.5z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">{t('variants.lowStock')}</div>
          <div className="stat-value text-warning">
            {lowStockVariants.length}
          </div>
          <div className="stat-desc">{t('variants.variantsWithLowStock')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-error">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </div>
          <div className="stat-title">{t('variants.outOfStock')}</div>
          <div className="stat-value text-error">
            {outOfStockVariants.length}
          </div>
          <div className="stat-desc">{t('variants.variantsOutOfStock')}</div>
        </div>
      </div>

      {/* Bulk Operations */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h3 className="card-title">{t('variants.bulkStockOperations')}</h3>

          <div className="flex flex-wrap gap-4 items-end">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('variants.operation')}</span>
              </label>
              <select
                className="select select-bordered"
                value={bulkOperation}
                onChange={(e) =>
                  setBulkOperation(e.target.value as 'add' | 'set' | 'subtract')
                }
              >
                <option value="add">{t('variants.add')}</option>
                <option value="subtract">{t('variants.subtract')}</option>
                <option value="set">{t('variants.setTo')}</option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('variants.amount')}</span>
              </label>
              <input
                type="number"
                min="0"
                value={bulkValue}
                onChange={(e) => setBulkValue(parseInt(e.target.value) || 0)}
                className="input input-bordered"
                placeholder="0"
              />
            </div>

            <button
              onClick={applyBulkUpdate}
              disabled={selectedForBulk.size === 0 || bulkValue === 0}
              className="btn btn-primary"
            >
              {t('variants.applyToSelected', { count: selectedForBulk.size })}
            </button>
          </div>
        </div>
      </div>

      {/* Variants Stock Table */}
      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th>
                <input
                  type="checkbox"
                  className="checkbox"
                  onChange={(e) => {
                    if (e.target.checked) {
                      setSelectedForBulk(new Set(variants.map((v) => v.id!)));
                    } else {
                      setSelectedForBulk(new Set());
                    }
                  }}
                />
              </th>
              <th>{t('variants.variant')}</th>
              <th>{t('variants.currentStock')}</th>
              <th>{t('variants.status')}</th>
              <th>{t('variants.lastUpdated')}</th>
            </tr>
          </thead>
          <tbody>
            {variants.map((variant) => (
              <tr key={variant.id}>
                <td>
                  <input
                    type="checkbox"
                    className="checkbox"
                    checked={selectedForBulk.has(variant.id!)}
                    onChange={(e) => {
                      const newSelected = new Set(selectedForBulk);
                      if (e.target.checked) {
                        newSelected.add(variant.id!);
                      } else {
                        newSelected.delete(variant.id!);
                      }
                      setSelectedForBulk(newSelected);
                    }}
                  />
                </td>
                <td>
                  <div className="font-medium">
                    {variant.name ||
                      Object.entries(variant.attributes || {})
                        .map(([k, v]) => `${k}: ${v}`)
                        .join(', ')}
                  </div>
                  {variant.sku && (
                    <div className="text-sm text-gray-500">
                      SKU: {variant.sku}
                    </div>
                  )}
                </td>
                <td>
                  <span className="font-mono text-lg">
                    {variant.stock_quantity || 0}
                  </span>
                </td>
                <td>
                  {(variant.stock_quantity || 0) === 0 ? (
                    <span className="badge badge-error">
                      {t('variants.outOfStock')}
                    </span>
                  ) : (variant.stock_quantity || 0) < 5 ? (
                    <span className="badge badge-warning">
                      {t('variants.lowStock')}
                    </span>
                  ) : (
                    <span className="badge badge-success">
                      {t('variants.inStock')}
                    </span>
                  )}
                </td>
                <td>
                  <span className="text-sm text-gray-500">
                    {variant.updated_at
                      ? new Date(variant.updated_at).toLocaleDateString()
                      : '-'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// Analytics Tab Component
interface AnalyticsTabProps {
  analytics: VariantAnalyticsResponse | null;
  variants: ProductVariant[];
  t: (key: string, params?: any) => string;
}

function AnalyticsTab({
  analytics,
  variants: _variants,
  t,
}: AnalyticsTabProps) {
  if (!analytics) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
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
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            {t('variants.noAnalytics')}
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            {t('variants.noAnalyticsDescription')}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Overview Stats */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-title">{t('variants.totalVariants')}</div>
          <div className="stat-value">{analytics.total_variants}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.totalStock')}</div>
          <div className="stat-value">{analytics.total_stock}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.totalSold')}</div>
          <div className="stat-value">{analytics.total_sold}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('variants.lowStockVariants')}</div>
          <div className="stat-value text-warning">
            {analytics.low_stock_variants?.length || 0}
          </div>
        </div>
      </div>

      {/* Best Seller */}
      {analytics.best_seller && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h3 className="card-title">{t('variants.bestSeller')}</h3>
            <div className="flex justify-between items-center">
              <div>
                <div className="font-medium">
                  {Object.entries(
                    analytics.best_seller.variant_attributes || {}
                  )
                    .map(([k, v]) => `${k}: ${v}`)
                    .join(', ') || 'Variant'}
                </div>
                {analytics.best_seller.sku && (
                  <div className="text-sm text-gray-500">
                    SKU: {analytics.best_seller.sku}
                  </div>
                )}
              </div>
              <div className="text-right">
                <div className="text-2xl font-bold text-primary">
                  {analytics.best_seller.price} RSD
                </div>
                <div className="text-sm text-gray-500">
                  {t('variants.stock')}: {analytics.best_seller.stock_quantity}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Stock by Attribute */}
      {analytics.stock_by_attribute &&
        Object.keys(analytics.stock_by_attribute).length > 0 && (
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title">{t('variants.stockByAttribute')}</h3>
              <div className="space-y-4">
                {Object.entries(analytics.stock_by_attribute).map(
                  ([attribute, values]) => (
                    <div key={attribute}>
                      <h4 className="font-medium capitalize">{attribute}</h4>
                      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 mt-2">
                        {Object.entries(values).map(([value, stock]) => (
                          <div
                            key={value}
                            className="flex justify-between items-center p-2 bg-base-200 rounded"
                          >
                            <span className="text-sm">{value}</span>
                            <span className="font-mono font-medium">
                              {stock}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )
                )}
              </div>
            </div>
          </div>
        )}

      {/* Sales by Attribute */}
      {analytics.sales_by_attribute &&
        Object.keys(analytics.sales_by_attribute).length > 0 && (
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title">{t('variants.salesByAttribute')}</h3>
              <div className="space-y-4">
                {Object.entries(analytics.sales_by_attribute).map(
                  ([attribute, values]) => (
                    <div key={attribute}>
                      <h4 className="font-medium capitalize">{attribute}</h4>
                      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 mt-2">
                        {Object.entries(values).map(([value, sales]) => (
                          <div
                            key={value}
                            className="flex justify-between items-center p-2 bg-base-200 rounded"
                          >
                            <span className="text-sm">{value}</span>
                            <span className="font-mono font-medium text-success">
                              {sales}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )
                )}
              </div>
            </div>
          </div>
        )}
    </div>
  );
}
