'use client';

import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useTranslations } from 'next-intl';
import {
  fetchCart,
  updateCartItem,
  removeFromCart,
  clearCart,
  selectCart,
  selectCartLoading,
  selectCartError,
  selectCartTotal,
} from '@/store/slices/cartSlice';
import SafeImage from '@/components/SafeImage';
import type { AppDispatch } from '@/store';
import type { components } from '@/types/generated/api';

type ShoppingCartItem =
  components['schemas']['backend_internal_domain_models.ShoppingCartItem'];

interface ShoppingCartModalProps {
  storefrontId: number;
  isOpen: boolean;
  onClose: () => void;
  onCheckout: () => void;
}

export default function ShoppingCartModal({
  storefrontId,
  isOpen,
  onClose,
  onCheckout,
}: ShoppingCartModalProps) {
  const t = useTranslations('cart');
  const dispatch = useDispatch<AppDispatch>();

  const cart = useSelector(selectCart);
  const loading = useSelector(selectCartLoading);
  const error = useSelector(selectCartError);
  const total = useSelector(selectCartTotal);

  useEffect(() => {
    if (isOpen && storefrontId) {
      dispatch(fetchCart(storefrontId));
    }
  }, [dispatch, isOpen, storefrontId]);

  const handleQuantityChange = async (itemId: number, quantity: number) => {
    if (quantity <= 0) {
      dispatch(removeFromCart({ storefrontId, itemId }));
    } else {
      dispatch(updateCartItem({ storefrontId, itemId, data: { quantity } }));
    }
  };

  const handleRemoveItem = (itemId: number) => {
    dispatch(removeFromCart({ storefrontId, itemId }));
  };

  const handleClearCart = () => {
    dispatch(clearCart(storefrontId));
  };

  const handleCheckout = () => {
    onCheckout();
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        <div
          className="fixed inset-0 bg-black bg-opacity-50"
          onClick={onClose}
        />

        <div className="relative w-full max-w-2xl rounded-lg bg-base-100 shadow-xl">
          {/* Header */}
          <div className="flex items-center justify-between border-b p-6">
            <h2 className="text-xl font-semibold">{t('title')}</h2>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              ✕
            </button>
          </div>

          {/* Content */}
          <div className="max-h-96 overflow-y-auto p-6">
            {loading && (
              <div className="flex justify-center py-8">
                <span className="loading loading-spinner loading-md"></span>
              </div>
            )}

            {error && (
              <div className="alert alert-error mb-4">
                <span>{error}</span>
              </div>
            )}

            {cart?.items && cart.items.length > 0 ? (
              <div className="space-y-4">
                {cart.items.map((item: ShoppingCartItem) => (
                  <div
                    key={item.id}
                    className="flex items-center gap-4 rounded-lg border p-4"
                  >
                    {/* Product Image */}
                    {item.product?.images?.[0] && (
                      <SafeImage
                        src={
                          item.product.images[0].thumbnail_url ||
                          item.product.images[0].image_url ||
                          ''
                        }
                        alt={item.product.name || 'Product'}
                        width={64}
                        height={64}
                        className="h-16 w-16 rounded-lg object-cover"
                      />
                    )}

                    {/* Product Info */}
                    <div className="flex-1">
                      <h4 className="font-medium">{item.product?.name}</h4>
                      {item.variant && (
                        <p className="text-sm text-base-content/70">
                          {Object.entries(item.variant.attributes || {})
                            .map(([key, value]) => `${key}: ${value}`)
                            .join(', ')}
                        </p>
                      )}
                      <p className="text-sm font-medium">
                        {Number(item.price_per_unit || 0).toFixed(2)}{' '}
                        {item.product?.currency}
                      </p>
                    </div>

                    {/* Quantity Controls */}
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() =>
                          handleQuantityChange(
                            item.id!,
                            (item.quantity || 1) - 1
                          )
                        }
                        className="btn btn-outline btn-sm btn-circle"
                        disabled={loading}
                      >
                        -
                      </button>
                      <span className="w-8 text-center">{item.quantity}</span>
                      <button
                        onClick={() =>
                          handleQuantityChange(
                            item.id!,
                            (item.quantity || 1) + 1
                          )
                        }
                        className="btn btn-outline btn-sm btn-circle"
                        disabled={loading}
                      >
                        +
                      </button>
                    </div>

                    {/* Item Total */}
                    <div className="text-right">
                      <p className="font-medium">
                        {Number(item.total_price || 0).toFixed(2)}{' '}
                        {item.product?.currency}
                      </p>
                      <button
                        onClick={() => handleRemoveItem(item.id!)}
                        className="btn btn-ghost btn-sm text-error"
                        disabled={loading}
                      >
                        {t('remove')}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              !loading && (
                <div className="py-8 text-center">
                  <p className="text-base-content/70">{t('empty')}</p>
                </div>
              )
            )}
          </div>

          {/* Footer */}
          {cart?.items && cart.items.length > 0 && (
            <div className="border-t p-6">
              <div className="mb-4 flex justify-between text-lg font-semibold">
                <span>{t('total')}</span>
                <span>
                  {total.toFixed(2)} {cart.items[0]?.product?.currency}
                </span>
              </div>

              <div className="flex gap-3">
                <button
                  onClick={handleClearCart}
                  className="btn btn-outline flex-1"
                  disabled={loading}
                >
                  {t('clear')}
                </button>
                <button
                  onClick={handleCheckout}
                  className="btn btn-primary flex-1"
                  disabled={loading || !cart.items.length}
                >
                  {t('checkout')}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
