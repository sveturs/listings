'use client';

import React, { useEffect, useRef, useState, useMemo } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import SafeImage from '@/components/SafeImage';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuth } from '@/contexts/AuthContext';
import {
  selectCartItems as selectLocalCartItems,
  selectCartTotal as selectLocalCartTotal,
  selectCartItemsCount as selectLocalCartItemsCount,
  removeItem as removeFromLocalCart,
  updateQuantity as updateLocalQuantity,
} from '@/store/slices/localCartSlice';
import {
  selectAllCarts,
  selectAllCartsItemsCount,
  removeFromCart,
  updateCartItem,
} from '@/store/slices/cartSlice';
import type { AppDispatch } from '@/store';

interface CartItemDisplay {
  productId: number;
  variantId?: number;
  quantity: number;
  name: string;
  variantName?: string;
  price: number;
  currency: string;
  image?: string;
  storefrontId?: number;
  storefrontName?: string;
  storefrontSlug?: string;
  stockQuantity?: number;
  cartItemId?: number;
}

interface MiniCartProps {
  isOpen: boolean;
  onClose: () => void;
  anchorRef?: React.RefObject<HTMLElement | null>;
}

export default function MiniCart({
  isOpen,
  onClose,
  anchorRef,
}: MiniCartProps) {
  const t = useTranslations('cart');
  const locale = useLocale();
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();
  const { isAuthenticated } = useAuth();
  const cartRef = useRef<HTMLDivElement>(null);
  const [position, setPosition] = useState({ top: 0, right: 0 });

  // Выбираем данные из правильного slice в зависимости от авторизации
  const localItems = useSelector(selectLocalCartItems);
  const localTotal = useSelector(selectLocalCartTotal);
  const localItemsCount = useSelector(selectLocalCartItemsCount);

  // Для авторизованных пользователей получаем все корзины
  const allCarts = useSelector(selectAllCarts);
  const apiItemsCount = useSelector(selectAllCartsItemsCount);

  // Преобразуем товары из всех backend корзин в нужный формат
  const backendItems = useMemo((): CartItemDisplay[] => {
    if (!isAuthenticated || !allCarts || allCarts.length === 0) return [];

    const items: CartItemDisplay[] = [];

    allCarts.forEach((cart) => {
      if (cart.items && cart.items.length > 0) {
        cart.items.forEach((item) => {
          items.push({
            productId: item.product_id || 0,
            variantId: item.variant_id,
            quantity: item.quantity || 0,
            name: item.product?.name || 'Product',
            variantName: item.variant?.sku,
            // Преобразуем price_per_unit в число (может быть строкой из-за Decimal на backend)
            price:
              typeof item.price_per_unit === 'string'
                ? parseFloat(item.price_per_unit)
                : item.price_per_unit || 0,
            currency: 'RSD',
            image: item.product?.images?.[0]?.image_url,
            storefrontId: cart.storefront_id,
            storefrontName:
              cart.storefront?.name || `Store ${cart.storefront_id}`,
            storefrontSlug: cart.storefront?.slug || String(cart.storefront_id),
            stockQuantity: item.product?.stock_quantity || 99,
            cartItemId: item.id,
          });
        });
      }
    });

    return items;
  }, [allCarts, isAuthenticated]);

  const backendTotal = useMemo(() => {
    return backendItems.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );
  }, [backendItems]);

  // Используем данные в зависимости от авторизации
  const items: CartItemDisplay[] = isAuthenticated ? backendItems : localItems;
  const total = isAuthenticated ? backendTotal : localTotal;
  const itemsCount = isAuthenticated ? apiItemsCount : localItemsCount;

  // Группируем товары по витринам
  const itemsByStorefront = items.reduce(
    (acc, item) => {
      const storefrontId = item.storefrontId || 0; // Default to 0 if undefined
      if (!acc[storefrontId]) {
        acc[storefrontId] = {
          items: [],
          name: item.storefrontName || `Store ${storefrontId}`,
        };
      }
      acc[storefrontId].items.push(item);
      return acc;
    },
    {} as Record<number, { items: typeof items; name: string }>
  );

  // Вычисляем позицию относительно кнопки корзины
  useEffect(() => {
    if (isOpen && anchorRef?.current) {
      const rect = anchorRef.current.getBoundingClientRect();
      const cartWidth = 384; // w-96 = 24rem = 384px
      const viewportWidth = window.innerWidth;

      // Проверяем, не выходит ли корзина за правый край экрана
      let rightPos = viewportWidth - rect.right;

      // Если корзина выходит за левый край экрана, корректируем позицию
      if (rect.right - cartWidth < 0) {
        rightPos = viewportWidth - cartWidth - 16; // 16px отступ от края
      }

      setPosition({
        top: rect.bottom + 8, // 8px отступ от кнопки
        right: Math.max(16, rightPos), // минимум 16px от края
      });
    }
  }, [isOpen, anchorRef]);

  // Закрытие при клике вне компонента
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        cartRef.current &&
        !cartRef.current.contains(event.target as Node) &&
        anchorRef?.current &&
        !anchorRef.current.contains(event.target as Node)
      ) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onClose, anchorRef]);

  const handleRemove = async (productId: number, variantId?: number) => {
    if (isAuthenticated) {
      // Для авторизованных пользователей используем backend API
      // Находим товар, чтобы получить storefrontId и cartItemId
      const item = items.find(
        (i) => i.productId === productId && i.variantId === variantId
      );
      if (item && item.storefrontId && item.cartItemId) {
        await dispatch(
          removeFromCart({
            storefrontId: item.storefrontId,
            itemId: item.cartItemId,
          })
        );
      }
    } else {
      // Для неавторизованных используем localStorage
      dispatch(removeFromLocalCart({ productId, variantId }));
    }
  };

  const handleQuantityChange = async (
    productId: number,
    quantity: number,
    variantId?: number
  ) => {
    if (quantity <= 0) return;

    if (isAuthenticated) {
      // Для авторизованных пользователей используем backend API
      const item = items.find(
        (i) => i.productId === productId && i.variantId === variantId
      );
      if (item && item.storefrontId && item.cartItemId) {
        await dispatch(
          updateCartItem({
            storefrontId: item.storefrontId,
            itemId: item.cartItemId,
            data: { quantity },
          })
        );
      }
    } else {
      // Для неавторизованных используем localStorage
      dispatch(updateLocalQuantity({ productId, variantId, quantity }));
    }
  };

  const handleGoToCart = () => {
    onClose();
    router.push(`/${locale}/cart`);
  };

  const handleCheckout = () => {
    onClose();
    router.push(`/${locale}/checkout`);
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          ref={cartRef}
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -10 }}
          transition={{ duration: 0.2 }}
          className="fixed w-96 max-w-[calc(100vw-2rem)] bg-base-100 rounded-lg shadow-xl border border-base-300 z-[150]"
          style={{
            top: `${position.top}px`,
            right: `${position.right}px`,
          }}
        >
          {/* Header */}
          <div className="p-4 border-b border-base-300">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold">
                {t('cart.miniCart.title')} ({itemsCount})
              </h3>
              <button
                onClick={onClose}
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
          </div>

          {/* Content */}
          <div className="max-h-96 overflow-y-auto">
            {itemsCount === 0 ? (
              <div className="p-8 text-center">
                <svg
                  className="w-16 h-16 mx-auto mb-4 text-base-content/20"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                  />
                </svg>
                <p className="text-base-content/60">
                  {t('cart.miniCart.empty')}
                </p>
              </div>
            ) : (
              <div className="p-4 space-y-4">
                {Object.entries(itemsByStorefront).map(
                  ([storefrontId, group]) => (
                    <div key={storefrontId}>
                      {Object.keys(itemsByStorefront).length > 1 && (
                        <div className="text-xs font-medium text-base-content/60 mb-2">
                          {group.name}
                        </div>
                      )}
                      <div className="space-y-3">
                        {group.items.slice(0, 3).map((item) => (
                          <motion.div
                            key={`${item.productId}-${item.variantId || ''}`}
                            layout
                            className="flex gap-3"
                          >
                            {/* Image */}
                            <Link
                              href={`/${locale}/storefronts/${item.storefrontSlug}/products/${item.productId}`}
                              className="relative w-16 h-16 rounded-lg overflow-hidden bg-base-200 flex-shrink-0"
                            >
                              <SafeImage
                                src={item.image}
                                alt={item.name}
                                fill
                                className="object-cover"
                                sizes="64px"
                              />
                            </Link>

                            {/* Info */}
                            <div className="flex-1 min-w-0">
                              <Link
                                href={`/${locale}/storefronts/${item.storefrontSlug}/products/${item.productId}`}
                                className="text-sm font-medium hover:link line-clamp-1"
                              >
                                {item.name}
                              </Link>
                              {item.variantName && (
                                <p className="text-xs text-base-content/60">
                                  {item.variantName}
                                </p>
                              )}
                              <div className="flex items-center justify-between mt-1">
                                <span className="text-sm font-semibold">
                                  {(item.price * item.quantity).toFixed(2)}{' '}
                                  {item.currency}
                                </span>
                                <div className="flex items-center gap-1">
                                  <button
                                    onClick={() =>
                                      handleQuantityChange(
                                        item.productId,
                                        item.quantity - 1,
                                        item.variantId
                                      )
                                    }
                                    className="btn btn-ghost btn-xs"
                                    disabled={item.quantity <= 1}
                                  >
                                    -
                                  </button>
                                  <span className="text-xs w-8 text-center">
                                    {item.quantity}
                                  </span>
                                  <button
                                    onClick={() =>
                                      handleQuantityChange(
                                        item.productId,
                                        item.quantity + 1,
                                        item.variantId
                                      )
                                    }
                                    className="btn btn-ghost btn-xs"
                                    disabled={
                                      item.quantity >=
                                      (item.stockQuantity || 99)
                                    }
                                  >
                                    +
                                  </button>
                                </div>
                              </div>
                            </div>

                            {/* Remove button */}
                            <button
                              onClick={() =>
                                handleRemove(item.productId, item.variantId)
                              }
                              className="btn btn-ghost btn-xs btn-circle"
                            >
                              <svg
                                className="w-4 h-4"
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
                            </button>
                          </motion.div>
                        ))}
                      </div>
                      {group.items.length > 3 && (
                        <p className="text-xs text-base-content/60 mt-2">
                          {t('cart.miniCart.moreItems', {
                            count: group.items.length - 3,
                          })}
                        </p>
                      )}
                    </div>
                  )
                )}
              </div>
            )}
          </div>

          {/* Footer */}
          {itemsCount > 0 && (
            <div className="p-4 border-t border-base-300">
              <div className="flex items-center justify-between mb-4">
                <span className="font-medium">{t('cart.miniCart.total')}</span>
                <span className="text-xl font-bold text-primary">
                  {total.toFixed(2)} RSD
                </span>
              </div>
              <div className="space-y-2">
                <button
                  onClick={handleGoToCart}
                  className="btn btn-outline btn-block"
                >
                  {t('cart.miniCart.viewCart')}
                </button>
                <button
                  onClick={handleCheckout}
                  className="btn btn-primary btn-block"
                >
                  {t('cart.miniCart.checkout')}
                </button>
              </div>
            </div>
          )}

          {/* Free shipping progress */}
          {itemsCount > 0 && total < 5000 && (
            <div className="px-4 pb-4">
              <div className="text-xs text-base-content/60 mb-1">
                {t('cart.miniCart.freeShippingProgress', {
                  amount: (5000 - total).toFixed(2),
                })}
              </div>
              <progress
                className="progress progress-primary w-full h-2"
                value={total}
                max={5000}
              />
            </div>
          )}
        </motion.div>
      )}
    </AnimatePresence>
  );
}
