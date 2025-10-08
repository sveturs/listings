'use client';

import { useState, useEffect, useCallback } from 'react';
import Image from 'next/image';
import { useTranslations } from 'next-intl';
import configManager from '@/config';
import { useLocale } from 'next-intl';
import Link from 'next/link';

interface B2CProduct {
  id: number;
  storefront_id: number;
  name: string;
  description: string;
  price: number;
  currency: string;
  images?: Array<{
    id: number;
    image_url: string;
    is_default?: boolean;
  }>;
  storefront?: {
    slug: string;
    name: string;
  };
}

interface B2CProductQuickViewProps {
  productId: number;
  isOpen: boolean;
  onClose: () => void;
}

export default function B2CProductQuickView({
  productId,
  isOpen,
  onClose,
}: B2CProductQuickViewProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const [product, setProduct] = useState<B2CProduct | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProduct = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/b2c/products/${productId}`);
      if (!response.ok) {
        throw new Error('Failed to fetch product');
      }
      const data = await response.json();

      // Fetch storefront info if we have storefront_id
      if (data.storefront_id) {
        try {
          const storefrontResp = await fetch(
            `/api/v1/b2c/${data.storefront_id}`
          );
          if (storefrontResp.ok) {
            const storefrontData = await storefrontResp.json();
            data.storefront = storefrontData;
          }
        } catch (err) {
          console.error('Error fetching storefront:', err);
        }
      }

      setProduct(data);
    } catch (err) {
      console.error('Error fetching product:', err);
      setError(t('errors.loadingProduct'));
    } finally {
      setLoading(false);
    }
  }, [productId, t]);

  useEffect(() => {
    if (isOpen && productId) {
      fetchProduct();
    }
  }, [isOpen, productId, fetchProduct]);

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <button
          className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          onClick={onClose}
        >
          ✕
        </button>

        {loading && (
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        )}

        {error && (
          <div className="alert alert-error">
            <span>{error}</span>
          </div>
        )}

        {product && !loading && (
          <div className="space-y-4">
            {/* Product Image */}
            {product.images && product.images.length > 0 && (
              <div className="relative h-64 w-full rounded-lg overflow-hidden bg-base-200">
                <Image
                  src={configManager.buildImageUrl(product.images[0].image_url)}
                  alt={product.name}
                  fill
                  className="object-contain"
                />
              </div>
            )}

            {/* Product Info */}
            <div>
              <h3 className="text-2xl font-bold">{product.name}</h3>
              <p className="text-xl font-semibold text-primary mt-2">
                {product.price} {product.currency}
              </p>
            </div>

            {/* Description */}
            {product.description && (
              <div className="prose max-w-none">
                <p className="text-base-content/80">{product.description}</p>
              </div>
            )}

            {/* Actions */}
            <div className="modal-action">
              {product.storefront?.slug && (
                <Link
                  href={`/${locale}/b2c/${product.storefront.slug}/products/${product.id}`}
                  className="btn btn-primary"
                >
                  {t('viewFullDetails')}
                </Link>
              )}
              <button className="btn" onClick={onClose}>
                {t('close')}
              </button>
            </div>
          </div>
        )}
      </div>
      <div className="modal-backdrop bg-black/50" onClick={onClose}></div>
    </div>
  );
}
