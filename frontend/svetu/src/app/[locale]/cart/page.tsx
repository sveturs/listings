'use client';

import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import SafeImage from '@/components/SafeImage';
import { motion, AnimatePresence } from 'framer-motion';
import {
  selectCartItems,
  selectCartTotal,
  selectCartItemsCount,
  removeItem as removeFromCart,
  updateQuantity,
  clearCart,
} from '@/store/slices/localCartSlice';
import type { AppDispatch } from '@/store';
import { PageTransition } from '@/components/ui/PageTransition';

export default function CartPage() {
  const t = useTranslations('cart');
  const locale = useLocale();
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();

  const items = useSelector(selectCartItems);
  const total = useSelector(selectCartTotal);
  const itemsCount = useSelector(selectCartItemsCount);
  const [isLoading, setIsLoading] = useState(false);

  // Группируем товары по витринам
  const itemsByStorefront = items.reduce(
    (acc, item) => {
      const storefrontId = item.storefrontId;
      if (!acc[storefrontId]) {
        acc[storefrontId] = {
          items: [],
          name: item.storefrontName || `Store ${storefrontId}`,
          slug: item.storefrontSlug || '',
          subtotal: 0,
        };
      }
      acc[storefrontId].items.push(item);
      acc[storefrontId].subtotal += item.price * item.quantity;
      return acc;
    },
    {} as Record<
      number,
      { items: typeof items; name: string; slug: string; subtotal: number }
    >
  );

  const handleRemove = (productId: number, variantId?: number) => {
    dispatch(removeFromCart({ productId, variantId }));
  };

  const handleQuantityChange = (
    productId: number,
    quantity: number,
    variantId?: number
  ) => {
    if (quantity > 0) {
      dispatch(updateQuantity({ productId, variantId, quantity }));
    }
  };

  const handleClearCart = () => {
    if (confirm(t('confirmClear'))) {
      dispatch(clearCart());
    }
  };

  const handleCheckout = () => {
    setIsLoading(true);
    router.push(`/${locale}/checkout`);
  };

  const calculateShipping = () => {
    // Бесплатная доставка от 5000 RSD
    return total >= 5000 ? 0 : 300;
  };

  const shipping = calculateShipping();
  const finalTotal = total + shipping;

  if (itemsCount === 0) {
    return (
      <PageTransition>
        <div className="min-h-screen bg-base-100 pt-24">
          <div className="container mx-auto px-4 py-16">
            <div className="max-w-md mx-auto text-center">
              <svg
                className="w-32 h-32 mx-auto mb-8 text-base-content/20"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1}
                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
              <h1 className="text-2xl font-bold mb-4">{t('emptyTitle')}</h1>
              <p className="text-base-content/60 mb-8">
                {t('emptyDescription')}
              </p>
              <Link href={`/${locale}`} className="btn btn-primary">
                {t('continueShopping')}
              </Link>
            </div>
          </div>
        </div>
      </PageTransition>
    );
  }

  return (
    <PageTransition>
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-8">
          {/* Header */}
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
            <h1 className="text-3xl font-bold">
              {t('title')} ({itemsCount})
            </h1>
            <button
              onClick={handleClearCart}
              className="btn btn-ghost btn-sm text-error"
            >
              <svg
                className="w-4 h-4 mr-2"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
              {t('clearCart')}
            </button>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Cart Items */}
            <div className="lg:col-span-2 space-y-6">
              {Object.entries(itemsByStorefront).map(
                ([storefrontId, group]) => (
                  <div key={storefrontId} className="card bg-base-200">
                    <div className="card-body">
                      {/* Store header */}
                      <div className="flex items-center justify-between mb-4">
                        <Link
                          href={`/${locale}/storefronts/${group.slug}`}
                          className="text-lg font-semibold hover:link"
                        >
                          {group.name}
                        </Link>
                        <span className="text-sm text-base-content/60">
                          {t('subtotal')}: {group.subtotal.toFixed(2)} RSD
                        </span>
                      </div>

                      {/* Items */}
                      <div className="divide-y divide-base-300">
                        <AnimatePresence>
                          {group.items.map((item) => (
                            <motion.div
                              key={`${item.productId}-${item.variantId || ''}`}
                              layout
                              exit={{ opacity: 0, height: 0 }}
                              className="py-4 first:pt-0 last:pb-0"
                            >
                              <div className="flex gap-4">
                                {/* Image */}
                                <Link
                                  href={`/${locale}/storefronts/${item.storefrontSlug}/products/${item.productId}`}
                                  className="relative w-24 h-24 rounded-lg overflow-hidden bg-base-100 flex-shrink-0"
                                >
                                  <SafeImage
                                    src={item.image}
                                    alt={item.name}
                                    fill
                                    className="object-cover"
                                    sizes="96px"
                                  />
                                </Link>

                                {/* Info */}
                                <div className="flex-1">
                                  <div className="flex justify-between">
                                    <div>
                                      <Link
                                        href={`/${locale}/storefronts/${item.storefrontSlug}/products/${item.productId}`}
                                        className="font-medium hover:link"
                                      >
                                        {item.name}
                                      </Link>
                                      {item.variantName && (
                                        <p className="text-sm text-base-content/60">
                                          {item.variantName}
                                        </p>
                                      )}
                                    </div>
                                    <button
                                      onClick={() =>
                                        handleRemove(
                                          item.productId,
                                          item.variantId
                                        )
                                      }
                                      className="btn btn-ghost btn-sm btn-circle"
                                    >
                                      <svg
                                        className="w-5 h-5"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                      >
                                        <path
                                          strokeLinecap="round"
                                          strokeLinejoin="round"
                                          strokeWidth={2}
                                          d="M6 18L18 6M6 6l12 12"
                                        />
                                      </svg>
                                    </button>
                                  </div>

                                  <div className="flex items-center justify-between mt-2">
                                    {/* Quantity selector */}
                                    <div className="join">
                                      <button
                                        onClick={() =>
                                          handleQuantityChange(
                                            item.productId,
                                            item.quantity - 1,
                                            item.variantId
                                          )
                                        }
                                        className="btn btn-sm join-item"
                                        disabled={item.quantity <= 1}
                                      >
                                        -
                                      </button>
                                      <input
                                        type="number"
                                        className="input input-sm join-item w-16 text-center"
                                        value={item.quantity}
                                        onChange={(e) =>
                                          handleQuantityChange(
                                            item.productId,
                                            parseInt(e.target.value) || 1,
                                            item.variantId
                                          )
                                        }
                                        min="1"
                                        max={item.stockQuantity || 99}
                                      />
                                      <button
                                        onClick={() =>
                                          handleQuantityChange(
                                            item.productId,
                                            item.quantity + 1,
                                            item.variantId
                                          )
                                        }
                                        className="btn btn-sm join-item"
                                        disabled={
                                          item.quantity >=
                                          (item.stockQuantity || 99)
                                        }
                                      >
                                        +
                                      </button>
                                    </div>

                                    {/* Price */}
                                    <div className="text-right">
                                      <div className="font-semibold">
                                        {(item.price * item.quantity).toFixed(
                                          2
                                        )}{' '}
                                        {item.currency}
                                      </div>
                                      {item.quantity > 1 && (
                                        <div className="text-xs text-base-content/60">
                                          {item.price.toFixed(2)} ×{' '}
                                          {item.quantity}
                                        </div>
                                      )}
                                    </div>
                                  </div>

                                  {/* Stock warning */}
                                  {item.stockQuantity &&
                                    item.stockQuantity < 5 && (
                                      <p className="text-xs text-warning mt-2">
                                        {t('lowStock', {
                                          count: item.stockQuantity,
                                        })}
                                      </p>
                                    )}
                                </div>
                              </div>
                            </motion.div>
                          ))}
                        </AnimatePresence>
                      </div>
                    </div>
                  </div>
                )
              )}

              {/* Promo code */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-semibold mb-3">{t('promoCode')}</h3>
                  <div className="flex gap-2">
                    <input
                      type="text"
                      placeholder={t('enterPromoCode')}
                      className="input input-bordered flex-1"
                    />
                    <button className="btn btn-primary">{t('apply')}</button>
                  </div>
                </div>
              </div>
            </div>

            {/* Order Summary */}
            <div className="lg:col-span-1">
              <div className="card bg-base-200 sticky top-24">
                <div className="card-body">
                  <h2 className="card-title mb-4">{t('orderSummary')}</h2>

                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>{t('subtotal')}</span>
                      <span>{total.toFixed(2)} RSD</span>
                    </div>
                    <div className="flex justify-between">
                      <span>{t('shipping')}</span>
                      <span>
                        {shipping === 0 ? (
                          <span className="text-success">
                            {t('freeShipping')}
                          </span>
                        ) : (
                          `${shipping} RSD`
                        )}
                      </span>
                    </div>
                    <div className="divider"></div>
                    <div className="flex justify-between text-lg font-bold">
                      <span>{t('total')}</span>
                      <span className="text-primary">
                        {finalTotal.toFixed(2)} RSD
                      </span>
                    </div>
                  </div>

                  {/* Free shipping progress */}
                  {shipping > 0 && (
                    <div className="mt-4">
                      <p className="text-sm text-base-content/60 mb-2">
                        {t('freeShippingAt', {
                          amount: (5000 - total).toFixed(2),
                        })}
                      </p>
                      <progress
                        className="progress progress-primary w-full"
                        value={total}
                        max={5000}
                      />
                    </div>
                  )}

                  <button
                    onClick={handleCheckout}
                    className="btn btn-primary btn-block mt-6"
                    disabled={isLoading}
                  >
                    {isLoading && <span className="loading loading-spinner" />}
                    {t('proceedToCheckout')}
                  </button>

                  <Link
                    href={`/${locale}`}
                    className="btn btn-ghost btn-block btn-sm"
                  >
                    {t('continueShopping')}
                  </Link>

                  {/* Security badges */}
                  <div className="mt-6 pt-6 border-t border-base-300">
                    <div className="flex items-center gap-2 text-sm text-base-content/60">
                      <svg
                        className="w-4 h-4 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span>{t('secureCheckout')}</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-base-content/60 mt-2">
                      <svg
                        className="w-4 h-4 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4z" />
                        <path
                          fillRule="evenodd"
                          d="M18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span>{t('securePayment')}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Recently viewed or recommendations */}
          <div className="mt-12">
            <h2 className="text-2xl font-bold mb-6">{t('recommendations')}</h2>
            {/* TODO: Add product recommendations */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="card bg-base-200 animate-pulse">
                  <div className="aspect-square bg-base-300"></div>
                  <div className="card-body p-4">
                    <div className="h-4 bg-base-300 rounded w-3/4 mb-2"></div>
                    <div className="h-4 bg-base-300 rounded w-1/2"></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
