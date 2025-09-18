'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { motion, AnimatePresence } from 'framer-motion';
import { FiHeart } from 'react-icons/fi';
import { useAuth } from '@/contexts/AuthContext';
import { favoritesService } from '@/services/favorites';

export default function FavoritesIcon() {
  const locale = useLocale();
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const [favoritesCount, setFavoritesCount] = useState(0);

  useEffect(() => {
    if (!isAuthenticated) {
      setFavoritesCount(0);
      return;
    }

    const loadCount = async () => {
      const count = await favoritesService.getFavoritesCount();
      setFavoritesCount(count);
    };

    loadCount();

    // Подписываемся на изменения избранного
    const handleFavoritesChange = () => {
      setFavoritesCount(favoritesService.getFavoritesIds().size);
    };

    // Добавляем слушатель изменений (через custom event)
    window.addEventListener('favoritesChanged', handleFavoritesChange);

    return () => {
      window.removeEventListener('favoritesChanged', handleFavoritesChange);
    };
  }, [isAuthenticated]);

  const handleClick = () => {
    router.push(`/${locale}/favorites`);
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <button
      onClick={handleClick}
      className="btn btn-ghost btn-circle relative hidden sm:inline-flex"
      aria-label="Favorites"
    >
      <FiHeart className="w-5 h-5" />

      {/* Badge с количеством избранного */}
      <AnimatePresence>
        {favoritesCount > 0 && (
          <motion.span
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            exit={{ scale: 0 }}
            className="absolute -top-1 -right-1 bg-error text-error-content rounded-full w-5 h-5 text-xs flex items-center justify-center font-bold"
          >
            {favoritesCount > 99 ? '99+' : favoritesCount}
          </motion.span>
        )}
      </AnimatePresence>
    </button>
  );
}
