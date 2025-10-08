'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import VariantSelector from '@/components/Storefront/ProductVariants/VariantSelector';
import type { components } from '@/types/generated/api';

type ProductVariant = components['schemas']['models.StorefrontProductVariant'];

interface VariantSelectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  productId: number;
  productName: string;
  productImage?: string;
  storefrontSlug: string;
  storefrontId?: number; // Добавляем storefrontId
  basePrice: number;
  baseCurrency: string;
  onAddToCart: (variant: ProductVariant | null, quantity: number) => void;
}

export default function VariantSelectionModal({
  isOpen,
  onClose,
  productId,
  productName,
  productImage,
  storefrontSlug,
  storefrontId,
  basePrice,
  baseCurrency,
  onAddToCart,
}: VariantSelectionModalProps) {
  const t = useTranslations('cart');
  const tCommon = useTranslations('common');
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  );
  const [quantity, setQuantity] = useState(1);
  const [hasVariants, setHasVariants] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [hasChecked, setHasChecked] = useState(false);

  const checkForVariants = async () => {
    if (hasChecked) {
      console.log('VariantSelectionModal: Already checked, skipping');
      return;
    }

    console.log(
      'VariantSelectionModal: checkForVariants started',
      storefrontSlug,
      storefrontId
    );
    setIsLoading(true);
    setHasChecked(true);
    try {
      let actualSlug = storefrontSlug;

      // Если slug пустой, загружаем реальный slug по ID витрины
      if ((!storefrontSlug || storefrontSlug.trim() === '') && storefrontId) {
        console.log(
          'VariantSelectionModal: Loading real slug by storefront ID',
          storefrontId
        );
        try {
          const storefrontResponse = await fetch(
            `/api/v1/storefronts/${storefrontId}`
          );
          if (storefrontResponse.ok) {
            const storefrontData = await storefrontResponse.json();
            actualSlug = storefrontData.slug || '';
            console.log('VariantSelectionModal: Found real slug:', actualSlug);
          }
        } catch (error) {
          console.error(
            'VariantSelectionModal: Failed to load storefront data:',
            error
          );
        }
      }

      // Если так и не получили slug, показываем модал без проверки вариантов
      if (!actualSlug || actualSlug.trim() === '') {
        console.log(
          'VariantSelectionModal: No valid slug found, showing modal without variant check'
        );
        setHasVariants(false);
        return;
      }

      // Используем публичный эндпоинт с реальным slug
      const url = `/api/v1/public/storefronts/${actualSlug}/products/${productId}`;

      const response = await fetch(url);
      console.log(
        'VariantSelectionModal: API response status:',
        response.status
      );
      if (response.ok) {
        const data = await response.json();
        console.log('VariantSelectionModal: API data:', data);
        console.log('VariantSelectionModal: data.variants:', data.variants);
        console.log(
          'VariantSelectionModal: variants length:',
          data.variants?.length
        );
        if (data.variants && data.variants.length > 0) {
          console.log('VariantSelectionModal: Setting hasVariants to true');
          setHasVariants(true);
        } else {
          console.log(
            'VariantSelectionModal: No variants found, showing quantity selector'
          );
          // No variants, but still show the modal for quantity selection
          setHasVariants(false);
        }
      }
    } catch (error) {
      console.error('VariantSelectionModal: Error checking variants:', error);
      // On error, show modal without variants
      setHasVariants(false);
    } finally {
      setIsLoading(false);
      console.log('VariantSelectionModal: checkForVariants finished');
    }
  };

  useEffect(() => {
    if (isOpen) {
      // Reset state when modal opens
      setSelectedVariant(null);
      setQuantity(1);
      setHasVariants(false);
      setHasChecked(false);
      checkForVariants();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, productId, storefrontSlug, storefrontId]);

  const handleAdd = () => {
    onAddToCart(selectedVariant, quantity);
    onClose();
  };

  if (!isOpen) return null;

  console.log('VariantSelectionModal: Rendering with state:', {
    isLoading,
    hasVariants,
    selectedVariant,
    shouldShowVariantSelector: !isLoading && hasVariants,
  });

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        {/* Header */}
        <div className="flex items-start gap-4 mb-6">
          {productImage && (
            <div className="w-20 h-20 rounded-lg overflow-hidden flex-shrink-0 relative">
              <Image
                src={productImage}
                alt={productName}
                fill
                sizes="80px"
                className="object-cover"
              />
            </div>
          )}
          <div className="flex-1">
            <h3 className="font-bold text-lg">{productName}</h3>
            <p className="text-2xl font-semibold text-primary mt-1">
              {selectedVariant?.price || basePrice} {baseCurrency}
            </p>
          </div>
          <button onClick={onClose} className="btn btn-sm btn-circle btn-ghost">
            ✕
          </button>
        </div>

        {/* Loading */}
        {isLoading && (
          <div className="flex justify-center py-8">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        )}

        {/* Content when not loading */}
        {!isLoading && (
          <div className="space-y-4">
            {/* Variant Selection - only if has variants */}
            {hasVariants && (
              <VariantSelector
                productId={productId}
                storefrontSlug={storefrontSlug}
                basePrice={basePrice}
                baseCurrency={baseCurrency}
                onVariantChange={setSelectedVariant}
              />
            )}

            {/* Quantity - always show */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('cart.quantity')}</span>
              </label>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="btn btn-sm btn-circle"
                >
                  -
                </button>
                <input
                  type="number"
                  value={quantity}
                  onChange={(e) =>
                    setQuantity(Math.max(1, parseInt(e.target.value) || 1))
                  }
                  className="input input-bordered w-20 text-center"
                  min="1"
                />
                <button
                  onClick={() => setQuantity(quantity + 1)}
                  className="btn btn-sm btn-circle"
                >
                  +
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Actions */}
        {!isLoading && (
          <div className="modal-action">
            <button onClick={onClose} className="btn btn-ghost">
              {tCommon('cancel')}
            </button>
            <button
              onClick={handleAdd}
              disabled={hasVariants && !selectedVariant}
              className="btn btn-primary"
            >
              {t('addToCart')}
            </button>
          </div>
        )}
      </div>
      <div className="modal-backdrop" onClick={onClose}></div>
    </div>
  );
}
