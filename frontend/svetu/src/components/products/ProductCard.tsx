'use client';

import { memo } from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import { FiEdit2, FiEye, FiPackage, FiAlertCircle } from 'react-icons/fi';
import type { B2CProduct } from '@/types/b2c';
import SafeImage from '@/components/SafeImage';

interface ProductCardProps {
  product: B2CProduct;
  storefrontSlug: string;
  isSelected?: boolean;
  isSelectMode?: boolean;
  onToggleSelect?: (id: number) => void;
  viewMode?: 'grid' | 'list' | 'table';
}

export const ProductCard = memo(function ProductCard({
  product,
  storefrontSlug,
  isSelected = false,
  isSelectMode = false,
  onToggleSelect,
  viewMode = 'grid',
}: ProductCardProps) {
  const t = useTranslations('storefronts');

  const mainImage =
    product.images?.find((img) => img.is_default) || product.images?.[0];
  const stockStatusColor = product.stock_status
    ? {
        in_stock: 'text-success',
        low_stock: 'text-warning',
        out_of_stock: 'text-error',
      }[product.stock_status] || 'text-base-content'
    : 'text-base-content';

  const handleCardClick = (e: React.MouseEvent) => {
    // Не вызываем выбор если кликнули на чекбокс - он обработает сам
    if (
      isSelectMode &&
      onToggleSelect &&
      product.id &&
      e.target !== e.currentTarget
    ) {
      const target = e.target as HTMLElement;
      // Проверяем что клик не на чекбоксе или label
      if (
        !target.closest('input[type="checkbox"]') &&
        !target.closest('products.label')
      ) {
        e.preventDefault();
        onToggleSelect(product.id);
      }
    }
  };

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.stopPropagation();
    if (onToggleSelect && product.id) {
      onToggleSelect(product.id);
    }
  };

  if (viewMode === 'table') {
    return (
      <tr
        className={`hover cursor-pointer transition-colors ${
          isSelected ? 'bg-primary/10' : ''
        }`}
        onClick={handleCardClick}
      >
        {isSelectMode && (
          <td className="w-12">
            <label className="cursor-pointer">
              <input
                type="checkbox"
                className="checkbox checkbox-primary"
                checked={isSelected}
                onChange={handleCheckboxChange}
              />
            </label>
          </td>
        )}
        <td>
          <div className="flex items-center gap-3">
            {mainImage ? (
              <div className="avatar">
                <div className="mask mask-squircle w-12 h-12">
                  <SafeImage
                    src={mainImage.image_url}
                    alt={product.name || 'Product image'}
                    width={48}
                    height={48}
                  />
                </div>
              </div>
            ) : (
              <div className="avatar placeholder">
                <div className="bg-neutral text-neutral-content rounded-btn w-12 h-12">
                  <FiPackage className="w-6 h-6" />
                </div>
              </div>
            )}
            <div>
              <div className="font-bold">{product.name}</div>
              {product.sku && (
                <div className="text-sm opacity-50">SKU: {product.sku}</div>
              )}
            </div>
          </div>
        </td>
        <td>
          {product.stock_status && (
            <span className={`badge badge-sm ${stockStatusColor}`}>
              {t(`stock.${product.stock_status}`)}
            </span>
          )}
          <div className="text-sm opacity-50">
            {product.stock_quantity} {t('products.units')}
          </div>
        </td>
        <td>
          <div className="font-semibold">
            {product.price} {product.currency}
          </div>
        </td>
        <td>
          <span
            className={`badge ${product.is_active ? 'badge-success' : 'badge-ghost'}`}
          >
            {product.is_active ? t('products.active') : t('products.inactive')}
          </span>
        </td>
        <td>
          <div className="flex gap-2">
            <Link
              href={`/marketplace/storefronts/${storefrontSlug}/products/${product.id}`}
              className="btn btn-ghost btn-xs"
              onClick={(e) => e.stopPropagation()}
            >
              <FiEye className="w-4 h-4" />
            </Link>
            <Link
              href={`/marketplace/storefronts/${storefrontSlug}/products/${product.id}/edit`}
              className="btn btn-ghost btn-xs"
              onClick={(e) => e.stopPropagation()}
            >
              <FiEdit2 className="w-4 h-4" />
            </Link>
          </div>
        </td>
      </tr>
    );
  }

  return (
    <div
      className={`card bg-base-100 shadow-xl hover:shadow-2xl transition-all cursor-pointer ${
        isSelected ? 'ring-2 ring-primary ring-offset-2' : ''
      } ${viewMode === 'list' ? 'card-side' : ''}`}
      onClick={handleCardClick}
    >
      {/* Чекбокс выбора */}
      {isSelectMode && (
        <div className="absolute top-2 left-2 z-10">
          <label className="cursor-pointer">
            <input
              type="checkbox"
              className="checkbox checkbox-primary checkbox-sm bg-base-100"
              checked={isSelected}
              onChange={handleCheckboxChange}
            />
          </label>
        </div>
      )}

      {/* Изображение */}
      <figure className={viewMode === 'list' ? 'w-48' : 'relative pb-[75%]'}>
        {mainImage ? (
          <SafeImage
            src={mainImage.image_url}
            alt={product.name || 'Product image'}
            fill
            className="object-cover"
          />
        ) : (
          <div className="absolute inset-0 bg-base-200 flex items-center justify-center">
            <FiPackage className="w-12 h-12 text-base-content/30" />
          </div>
        )}

        {/* Статус активности */}
        {!product.is_active && (
          <div className="absolute top-2 right-2">
            <span className="badge badge-ghost">{t('products.inactive')}</span>
          </div>
        )}
      </figure>

      <div className="card-body">
        <h3 className="card-title text-lg">
          {product.name}
          {product.stock_status === 'out_of_stock' && (
            <FiAlertCircle className="w-5 h-5 text-error" />
          )}
        </h3>

        {/* SKU и категория */}
        <div className="text-sm text-base-content/70">
          {product.sku && <span>SKU: {product.sku}</span>}
          {product.category && (
            <span className="ml-2 badge badge-sm">{product.category.name}</span>
          )}
        </div>

        {/* Цена */}
        <div className="flex items-baseline gap-2 mt-2">
          <span className="text-2xl font-bold text-primary">
            {product.price} {product.currency}
          </span>
        </div>

        {/* Склад */}
        <div className="flex items-center justify-between mt-4">
          <div className="flex items-center gap-2">
            {product.stock_status && (
              <span className={`text-sm font-medium ${stockStatusColor}`}>
                {t(`stock.${product.stock_status}`)}
              </span>
            )}
            <span className="text-sm text-base-content/70">
              ({product.stock_quantity} {t('products.units')})
            </span>
          </div>

          {/* Действия */}
          {!isSelectMode && (
            <div className="card-actions">
              <Link
                href={`/marketplace/storefronts/${storefrontSlug}/products/${product.id}/edit`}
                className="btn btn-primary btn-sm"
                onClick={(e) => e.stopPropagation()}
              >
                <FiEdit2 className="w-4 h-4" />
                {t('products.edit')}
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
});
