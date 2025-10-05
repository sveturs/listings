'use client';

import React from 'react';
import { useTranslations, useLocale } from 'next-intl';
import SafeImage from '@/components/SafeImage';
import AddToCartButton from '@/components/cart/AddToCartButton';
import type { components } from '@/types/generated/api';
import { MapPin } from 'lucide-react';
import {
  formatAddressWithPrivacy,
  getFullLocalizedAddress,
} from '@/utils/addressUtils';

type StorefrontProduct = components['schemas']['backend_internal_domain_models.StorefrontProduct'];

interface ProductCardProps {
  product: StorefrontProduct;
  storefrontId: number;
  viewMode?: 'grid' | 'list';
}

export default function ProductCard({
  product,
  storefrontId,
  viewMode = 'grid'
}: ProductCardProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();

  const mainImage = product.images?.find(img => img.is_default) || product.images?.[0];
  
  const stockStatusColor = {
    in_stock: 'text-success',
    low_stock: 'text-warning',
    out_of_stock: 'text-error',
  }[product.stock_status || 'in_stock'];

  if (viewMode === 'list') {
    return (
      <div className="card bg-base-100 shadow-xl hover:shadow-2xl border border-base-300 dark:border-base-600 transition-all card-side">
        {/* Image */}
        <figure className="w-48 h-32">
          {mainImage ? (
            <SafeImage
              src={mainImage.image_url || ''}
              alt={product.name || 'Product'}
              fill
              className="object-cover"
            />
          ) : (
            <div className="w-full h-full bg-base-200 flex items-center justify-center">
              <svg className="w-12 h-12 text-base-content/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          )}
        </figure>

        <div className="card-body flex-1">
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <h3 className="card-title text-lg">{product.name}</h3>
              
              {product.description && (
                <p className="text-sm text-base-content/70 line-clamp-2 mt-1">
                  {product.description}
                </p>
              )}

              {/* Location */}
              {product.has_individual_location && product.translations && (
                <div className="flex items-center gap-1 text-xs text-base-content/60 mt-1">
                  <MapPin className="w-3 h-3" />
                  <span>
                    {formatAddressWithPrivacy(
                      getFullLocalizedAddress(product as any, locale),
                      product.location_privacy as any
                    )}
                  </span>
                </div>
              )}

              <div className="flex items-center gap-4 mt-2">
                <span className="text-2xl font-bold text-primary">
                  {product.price} {product.currency}
                </span>
                
                <div className="flex items-center gap-2">
                  <span className={`text-sm font-medium ${stockStatusColor}`}>
                    {t(`stockStatus.${product.stock_status || 'in_stock'}`)}
                  </span>
                  <span className="text-sm text-base-content/70">
                    ({product.stock_quantity} {t('products.units')})
                  </span>
                </div>
              </div>
            </div>

            <div className="card-actions">
              <AddToCartButton
                product={{
                  id: product.id || 0,
                  name: product.name || '',
                  price: product.price || 0,
                  currency: product.currency,
                  image: product.images?.[0]?.image_url,
                  storefrontId: storefrontId,
                  stockQuantity: product.stock_quantity,
                  stockStatus: product.stock_status
                }}
                disabled={product.stock_status === 'out_of_stock' || !product.is_active}
                size="md"
              />
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-100 shadow-xl hover:shadow-2xl border border-base-300 dark:border-base-600 transition-all">
      {/* Image */}
      <figure className="relative pb-[75%]">
        {mainImage ? (
          <SafeImage
            src={mainImage.image_url || ''}
            alt={product.name || 'Product'}
            fill
            className="object-cover"
          />
        ) : (
          <div className="absolute inset-0 bg-base-200 flex items-center justify-center">
            <svg className="w-12 h-12 text-base-content/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
        )}

        {/* Out of stock overlay */}
        {product.stock_status === 'out_of_stock' && (
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
            <span className="badge badge-error badge-lg">{t('stockStatus.out_of_stock')}</span>
          </div>
        )}

        {/* Category badge */}
        {product.category && (
          <div className="absolute top-2 left-2">
            <span className="badge badge-primary badge-sm">{product.category.name}</span>
          </div>
        )}
      </figure>

      <div className="card-body">
        <h3 className="card-title text-lg line-clamp-2">{product.name}</h3>
        
        {product.description && (
          <p className="text-sm text-base-content/70 line-clamp-2">
            {product.description}
          </p>
        )}

        {/* Location */}
        {product.has_individual_location && product.translations && (
          <div className="flex items-center gap-1 text-xs text-base-content/60 mt-1">
            <MapPin className="w-3 h-3" />
            <span>
              {formatAddressWithPrivacy(
                getFullLocalizedAddress(product as any, locale),
                product.location_privacy as any
              )}
            </span>
          </div>
        )}

        {/* Price */}
        <div className="flex items-baseline gap-2 mt-2">
          <span className="text-2xl font-bold text-primary">
            {product.price} {product.currency}
          </span>
        </div>

        {/* Stock info */}
        <div className="flex items-center gap-2 mt-1">
          <span className={`text-sm font-medium ${stockStatusColor}`}>
            {t(`stockStatus.${product.stock_status || 'in_stock'}`)}
          </span>
          <span className="text-sm text-base-content/70">
            ({product.stock_quantity} {t('products.units')})
          </span>
        </div>

        {/* Add to cart button */}
        <div className="card-actions justify-end mt-4">
          <AddToCartButton
            product={{
              id: product.id || 0,
              name: product.name || '',
              price: product.price || 0,
              currency: product.currency,
              image: product.images?.[0]?.image_url,
              storefrontId: storefrontId,
              stockQuantity: product.stock_quantity,
              stockStatus: product.stock_status
            }}
            disabled={product.stock_status === 'out_of_stock' || !product.is_active}
            size="md"
            className="w-full"
          />
        </div>
      </div>
    </div>
  );
}