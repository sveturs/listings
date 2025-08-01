'use client';

import React, { useState, useRef } from 'react';
import { useSelector } from 'react-redux';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { motion, AnimatePresence } from 'framer-motion';
import { selectCartItemsCount } from '@/store/slices/localCartSlice';
import MiniCart from './MiniCart';

export default function CartIcon() {
  const locale = useLocale();
  const router = useRouter();
  const itemsCount = useSelector(selectCartItemsCount);
  const [showMiniCart, setShowMiniCart] = useState(false);
  const iconRef = useRef<HTMLElement>(null);

  const handleClick = () => {
    // На мобильных устройствах сразу переходим в корзину
    if (window.innerWidth < 768) {
      router.push(`/${locale}/cart`);
    } else {
      setShowMiniCart(!showMiniCart);
    }
  };

  return (
    <div className="relative">
      <button
        ref={iconRef as React.RefObject<HTMLButtonElement>}
        onClick={handleClick}
        className="btn btn-ghost btn-circle relative"
        aria-label="Shopping cart"
      >
        <svg
          className="w-6 h-6"
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

        {/* Badge с количеством товаров */}
        <AnimatePresence>
          {itemsCount > 0 && (
            <motion.span
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0 }}
              className="absolute -top-1 -right-1 bg-primary text-primary-content rounded-full w-5 h-5 text-xs flex items-center justify-center font-bold"
            >
              {itemsCount > 99 ? '99+' : itemsCount}
            </motion.span>
          )}
        </AnimatePresence>
      </button>

      {/* Mini Cart Dropdown */}
      <MiniCart
        isOpen={showMiniCart}
        onClose={() => setShowMiniCart(false)}
        anchorRef={iconRef}
      />
    </div>
  );
}
