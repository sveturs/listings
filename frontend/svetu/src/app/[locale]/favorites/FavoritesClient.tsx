'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import { favoritesService, FavoriteItem } from '@/services/favorites';
import { buildImageUrl } from '@/config';
import { useAuth } from '@/contexts/AuthContext';
import {
  FiHeart,
  FiMapPin,
  FiMessageCircle,
  FiEye,
  FiTrash2,
} from 'react-icons/fi';
import { toast } from 'react-hot-toast';

interface FavoritesClientProps {
  locale: string;
  translations: {
    title: string;
    description: string;
    emptyTitle: string;
    emptyDescription: string;
    browseListings: string;
    removeFromFavorites: string;
    viewDetails: string;
    contactSeller: string;
    addToCart: string;
    loading: string;
  };
}

export default function FavoritesClient({
  locale,
  translations,
}: FavoritesClientProps) {
  const [favorites, setFavorites] = useState<FavoriteItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [removingIds, setRemovingIds] = useState<Set<number | string>>(
    new Set()
  );
  const { user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!user) {
      router.push(`/${locale}/login`);
      return;
    }

    loadFavorites();
  }, [user, locale, router]);

  const loadFavorites = async () => {
    setIsLoading(true);
    try {
      const items = await favoritesService.getFavorites();
      setFavorites(items);
    } catch (error) {
      console.error('Failed to load favorites:', error);
      toast.error('Ошибка при загрузке избранного');
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemoveFromFavorites = async (itemId: number | string) => {
    setRemovingIds((prev) => new Set([...prev, itemId]));

    const success = await favoritesService.removeFromFavorites(itemId);
    if (success) {
      setFavorites((prev) => prev.filter((item) => item.id !== itemId));
    }

    setRemovingIds((prev) => {
      const newSet = new Set(prev);
      newSet.delete(itemId);
      return newSet;
    });
  };

  const getListingUrl = (item: FavoriteItem) => {
    // Для товаров витрин используем другой URL
    if (typeof item.id === 'string' && item.id.startsWith('sp_')) {
      const productId = item.id.replace('sp_', '');
      return `/${locale}/b2c/product/${productId}`;
    }
    return `/${locale}/c2c/${item.id}`;
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 py-8">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-center h-96">
            <div className="text-center">
              <div className="loading loading-spinner loading-lg text-primary"></div>
              <p className="mt-4 text-lg">{translations.loading}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (favorites.length === 0) {
    return (
      <div className="min-h-screen bg-base-200 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-md mx-auto text-center py-16">
            <FiHeart className="w-24 h-24 mx-auto text-base-content/20 mb-4" />
            <h1 className="text-2xl font-bold mb-2">
              {translations.emptyTitle}
            </h1>
            <p className="text-base-content/60 mb-6">
              {translations.emptyDescription}
            </p>
            <Link href={`/${locale}/search`} className="btn btn-primary">
              {translations.browseListings}
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="container mx-auto px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2 flex items-center gap-2">
            <FiHeart className="text-error" />
            {translations.title}
          </h1>
          <p className="text-base-content/60">
            {translations.description} ({favorites.length})
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          <AnimatePresence>
            {favorites.map((item) => (
              <motion.div
                key={item.id}
                layout
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.9 }}
                transition={{ duration: 0.3 }}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow"
              >
                <figure className="relative h-48 overflow-hidden">
                  <Link href={getListingUrl(item)}>
                    <img
                      src={
                        item.image
                          ? buildImageUrl(item.image)
                          : '/images/placeholder.jpg'
                      }
                      alt={item.title}
                      className="w-full h-full object-cover hover:scale-110 transition-transform duration-300 cursor-pointer"
                      onError={(e) => {
                        e.currentTarget.src = '/images/placeholder.jpg';
                      }}
                    />
                  </Link>

                  {/* Кнопка удаления из избранного */}
                  <button
                    onClick={() => handleRemoveFromFavorites(item.id)}
                    disabled={removingIds.has(item.id)}
                    className="btn btn-circle btn-sm absolute top-2 right-2 bg-base-100/80 hover:bg-error hover:text-white"
                  >
                    {removingIds.has(item.id) ? (
                      <span className="loading loading-spinner loading-xs"></span>
                    ) : (
                      <FiTrash2 className="w-4 h-4" />
                    )}
                  </button>

                  {/* Новое объявление */}
                  {item.created_at &&
                    new Date(item.created_at) >
                      new Date(Date.now() - 7 * 24 * 60 * 60 * 1000) && (
                      <div className="badge badge-secondary absolute top-2 left-2">
                        NEW
                      </div>
                    )}
                </figure>

                <div className="card-body p-4">
                  <h3 className="card-title text-base line-clamp-2">
                    {item.title}
                  </h3>

                  {/* Локация */}
                  {item.location && (
                    <div className="flex items-center gap-1 text-sm text-base-content/60">
                      <FiMapPin className="w-3 h-3" />
                      <span>{item.location}</span>
                    </div>
                  )}

                  {/* Категория */}
                  {item.category && (
                    <div className="badge badge-ghost badge-sm">
                      {item.category.name}
                    </div>
                  )}

                  {/* Цена */}
                  <p className="text-xl font-bold text-primary mt-2">
                    {item.price} {item.currency || 'РСД'}
                  </p>

                  {/* Действия */}
                  <div className="card-actions justify-between mt-4">
                    <Link
                      href={getListingUrl(item)}
                      className="btn btn-sm btn-ghost"
                    >
                      <FiEye className="w-4 h-4" />
                      {translations.viewDetails}
                    </Link>

                    <button
                      onClick={() =>
                        router.push(
                          `/${locale}/chat?listing_id=${item.id}&seller_id=${item.user_id}`
                        )
                      }
                      className="btn btn-sm btn-primary"
                    >
                      <FiMessageCircle className="w-4 h-4" />
                      {translations.contactSeller}
                    </button>
                  </div>
                </div>
              </motion.div>
            ))}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}
