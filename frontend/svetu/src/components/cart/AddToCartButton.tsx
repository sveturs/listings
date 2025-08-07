'use client';

import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useTranslations } from 'next-intl';
import { addToCart } from '@/store/slices/cartSlice';
import { addItem } from '@/store/slices/localCartSlice';
import { useAuth } from '@/contexts/AuthContext';
import type { AppDispatch } from '@/store';
import type { components } from '@/types/generated/api';

type StorefrontProductVariant =
  components['schemas']['backend_internal_domain_models.StorefrontProductVariant'];

interface ProductForCart {
  id: number;
  name: string;
  price: number;
  currency?: string;
  image?: string;
  storefrontId: number;
  storefrontName?: string;
  storefrontSlug?: string;
  stockQuantity?: number;
  stockStatus?: string;
  minOrderQuantity?: number;
  maxOrderQuantity?: number;
}

interface AddToCartButtonProps {
  product: ProductForCart;
  variant?: StorefrontProductVariant;
  quantity?: number;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
}

export default function AddToCartButton({
  product,
  variant,
  quantity = 1,
  className = '',
  size = 'md',
  disabled = false,
}: AddToCartButtonProps) {
  const t = useTranslations('cart');
  const dispatch = useDispatch<AppDispatch>();
  const { isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(false);
  const [justAdded, setJustAdded] = useState(false);

  const handleAddToCart = async () => {
    if (!product.id || disabled) return;

    try {
      setLoading(true);

      if (isAuthenticated && product.storefrontId) {
        // For authenticated users, save to backend
        await dispatch(
          addToCart({
            storefrontId: product.storefrontId,
            item: {
              product_id: product.id,
              variant_id: variant?.id,
              quantity,
            },
          })
        ).unwrap();
      } else {
        // For non-authenticated users, save to local storage
        dispatch(
          addItem({
            productId: product.id,
            variantId: variant?.id,
            name: product.name,
            variantName: variant?.name,
            price: variant?.price || product.price,
            currency: product.currency || 'RSD',
            quantity,
            stockQuantity: variant?.stock_quantity || product.stockQuantity,
            image: product.image,
            storefrontId: product.storefrontId,
            storefrontName: product.storefrontName || product.name,
            storefrontSlug:
              product.storefrontSlug || window.location.pathname.split('/')[3],
          })
        );
      }

      // Show success feedback
      setJustAdded(true);
      setTimeout(() => setJustAdded(false), 2000);
    } catch (error) {
      console.error('Failed to add to cart:', error);
    } finally {
      setLoading(false);
    }
  };

  const isOutOfStock = product.stockStatus === 'out_of_stock';
  const isDisabled = disabled || loading || isOutOfStock;

  const getButtonText = () => {
    if (justAdded) return t('added');
    if (isOutOfStock) return t('outOfStock');
    return t('addToCart');
  };

  const getButtonClass = () => {
    let baseClass = `btn ${className}`;

    if (justAdded) {
      baseClass += ' btn-success';
    } else if (isOutOfStock) {
      baseClass += ' btn-disabled';
    } else {
      baseClass += ' btn-primary';
    }

    switch (size) {
      case 'sm':
        baseClass += ' btn-sm';
        break;
      case 'lg':
        baseClass += ' btn-lg';
        break;
      default:
        break;
    }

    return baseClass;
  };

  return (
    <button
      onClick={handleAddToCart}
      disabled={isDisabled}
      className={getButtonClass()}
    >
      {loading && <span className="loading loading-spinner loading-sm mr-2" />}
      {justAdded && <span className="mr-2">✓</span>}
      {getButtonText()}
    </button>
  );
}
