'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

interface ProductVariant {
  sku: string;
  price: number;
  stock_quantity: number;
  variant_attributes: Record<string, string>;
  is_default: boolean;
}

interface VariantStockManagerProps {
  variants: ProductVariant[];
  onVariantsChange: (variants: ProductVariant[]) => void;
}

export default function VariantStockManager({
  variants,
  onVariantsChange,
}: VariantStockManagerProps) {
  const t = useTranslations('storefronts');
  const [localVariants, setLocalVariants] =
    useState<ProductVariant[]>(variants);

  useEffect(() => {
    setLocalVariants(variants);
  }, [variants]);

  const updateVariantStock = (index: number, stockQuantity: number) => {
    const updatedVariants = [...localVariants];
    updatedVariants[index] = {
      ...updatedVariants[index],
      stock_quantity: stockQuantity,
    };
    setLocalVariants(updatedVariants);
    onVariantsChange(updatedVariants);
  };

  const updateVariantPrice = (index: number, price: number) => {
    const updatedVariants = [...localVariants];
    updatedVariants[index] = {
      ...updatedVariants[index],
      price: price,
    };
    setLocalVariants(updatedVariants);
    onVariantsChange(updatedVariants);
  };

  const formatVariantName = (variant: ProductVariant) => {
    return Object.entries(variant.variant_attributes)
      .map(([key, value]) => `${key}: ${value}`)
      .join(', ');
  };

  const getTotalStock = () => {
    return localVariants.reduce(
      (total, variant) => total + variant.stock_quantity,
      0
    );
  };

  const getAveragePrice = () => {
    if (localVariants.length === 0) return 0;
    const totalPrice = localVariants.reduce(
      (total, variant) => total + variant.price,
      0
    );
    return totalPrice / localVariants.length;
  };

  if (localVariants.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-base-content/70">{t('products.noVariantsToManage')}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Сводка */}
      <div className="stats stats-horizontal shadow w-full">
        <div className="stat">
          <div className="stat-title">{t('products.totalVariants')}</div>
          <div className="stat-value text-primary">{localVariants.length}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('products.totalStock')}</div>
          <div className="stat-value text-secondary">{getTotalStock()}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('products.averagePrice')}</div>
          <div className="stat-value text-accent">
            {Math.round(getAveragePrice())} RSD
          </div>
        </div>
      </div>

      {/* Таблица вариантов */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h3 className="card-title text-lg mb-4">{t('products.manageVariants')}</h3>

          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>{t('products.variant')}</th>
                  <th>SKU</th>
                  <th>{t('products.price')} (RSD)</th>
                  <th>{t('products.stock')}</th>
                  <th>{t('products.status')}</th>
                </tr>
              </thead>
              <tbody>
                {localVariants.map((variant, index) => (
                  <tr
                    key={variant.sku}
                    className={variant.is_default ? 'bg-primary/10' : ''}
                  >
                    <td>
                      <div className="flex items-center gap-2">
                        {variant.is_default && (
                          <div className="badge badge-primary badge-sm">
                            {t('products.main')}
                          </div>
                        )}
                        <span className="font-medium">
                          {formatVariantName(variant)}
                        </span>
                      </div>
                    </td>
                    <td>
                      <code className="text-sm bg-base-200 px-2 py-1 rounded">
                        {variant.sku}
                      </code>
                    </td>
                    <td>
                      <input
                        type="number"
                        className="input input-bordered input-sm w-24"
                        value={variant.price}
                        onChange={(e) =>
                          updateVariantPrice(
                            index,
                            parseFloat(e.target.value) || 0
                          )
                        }
                        min="0"
                        step="1"
                      />
                    </td>
                    <td>
                      <input
                        type="number"
                        className="input input-bordered input-sm w-20"
                        value={variant.stock_quantity}
                        onChange={(e) =>
                          updateVariantStock(
                            index,
                            parseInt(e.target.value) || 0
                          )
                        }
                        min="0"
                      />
                    </td>
                    <td>
                      {variant.stock_quantity > 0 ? (
                        <div className="badge badge-success">
                          {t('products.inStock')}
                        </div>
                      ) : (
                        <div className="badge badge-error">
                          {t('products.outOfStock')}
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Быстрые действия */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h4 className="font-medium mb-4">{t('products.quickActions')}</h4>
          <div className="flex flex-wrap gap-2">
            <button
              className="btn btn-outline btn-sm"
              onClick={() => {
                const stock = prompt(t('products.setStockPrompt'));
                if (stock) {
                  const stockNumber = parseInt(stock);
                  const updatedVariants = localVariants.map((variant) => ({
                    ...variant,
                    stock_quantity: stockNumber,
                  }));
                  setLocalVariants(updatedVariants);
                  onVariantsChange(updatedVariants);
                }
              }}
            >
              {t('products.setStockForAll')}
            </button>
            <button
              className="btn btn-outline btn-sm"
              onClick={() => {
                const price = prompt(t('products.setPricePrompt'));
                if (price) {
                  const priceNumber = parseFloat(price);
                  const updatedVariants = localVariants.map((variant) => ({
                    ...variant,
                    price: priceNumber,
                  }));
                  setLocalVariants(updatedVariants);
                  onVariantsChange(updatedVariants);
                }
              }}
            >
              {t('products.setPriceForAll')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
