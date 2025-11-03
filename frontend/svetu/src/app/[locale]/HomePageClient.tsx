'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Link from 'next/link';
// import Image from 'next/image';
import dynamic from 'next/dynamic';
import { useRouter, useSearchParams } from 'next/navigation';
import { logger } from '@/utils/logger';
import { PageTransition } from '@/components/ui/PageTransition';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient } from '@/services/api-client';
// import CartIcon from '@/components/cart/CartIcon';
// import { AuthButton } from '@/components/AuthButton';
import { NestedCategorySelector } from '@/components/search/NestedCategorySelector';
import { useTranslations } from 'next-intl';
import configManager, { buildImageUrl } from '@/config';
import { useDispatch, useSelector } from 'react-redux';
import { addToCart } from '@/store/slices/cartSlice';
import {
  fetchCategories,
  fetchPopularCategories,
} from '@/store/slices/categoriesSlice';
import { fetchProviders } from '@/store/slices/deliverySlice';
import type { AppDispatch, RootState } from '@/store';
import { toast } from 'react-hot-toast';
import { favoritesService } from '@/services/favorites';
import SafeImage from '@/components/SafeImage';

// Динамический импорт карты для избежания SSR проблем
const EnhancedMapSection = dynamic(
  () =>
    import('./components/EnhancedMapSection').then((mod) => ({
      default: mod.EnhancedMapSection,
    })),
  {
    ssr: false,
    loading: () => (
      <div className="h-full w-full flex items-center justify-center bg-base-200 rounded-lg">
        <div className="text-center">
          <div className="loading loading-spinner loading-lg text-primary"></div>
          <p className="mt-2 text-base-content">Загрузка карты...</p>
        </div>
      </div>
    ),
  }
);

import {
  // FiSearch,
  FiMapPin,
  // FiMenu,
  // FiX,
  FiChevronRight,
  FiTruck,
  FiShield,
  FiCreditCard,
  FiMessageCircle,
  FiStar,
  FiHeart,
  FiTrendingUp,
  FiGrid,
  FiList,
  FiShoppingBag,
} from 'react-icons/fi';
// Некоторые иконки используются на странице отдельно от категорий
import { BsGem, BsPhone } from 'react-icons/bs';
import { AiOutlineEye } from 'react-icons/ai';
import { HiOutlineSparkles } from 'react-icons/hi';
import { getCategoryIcon } from '@/utils/categoryIcons';
import NearbyStats from '@/components/home/NearbyStats';

interface HomePageClientProps {
  title: string;
  description: string;
  createListingText: string;
  homePageData: any;
  locale: string;
}

export default function HomePageClient({
  createListingText,
  locale,
}: HomePageClientProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user } = useAuth();
  const dispatch = useDispatch<AppDispatch>();
  const t = useTranslations('marketplace');
  const _tCommon = useTranslations('common');
  const tFooter = useTranslations('common');
  const [_mounted, setMounted] = useState(false);
  const [selectedCategory] = useState<string | number>('all');
  const [currentBanner, setCurrentBanner] = useState(0);
  const [_showMobileMenu, _setShowMobileMenu] = useState(false);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [userLocation] = useState([44.7866, 20.4489]); // Координаты Белграда
  const [listings, setListings] = useState<any[]>([]);
  const [isLoadingListings, setIsLoadingListings] = useState(true);
  const [categories, setCategories] = useState<any[]>([]);
  const [popularCategories, setPopularCategories] = useState<any[]>([]);
  const [favoriteIds, setFavoriteIds] = useState<Set<number>>(new Set());
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [officialStores, setOfficialStores] = useState<any[]>([]);
  const [_isLoadingStores, setIsLoadingStores] = useState(false);

  // Функция для получения URL объявления
  const getListingUrl = (deal: any) => {
    logger.debug('getListingUrl called with deal:', {
      id: deal.id,
      product_id: deal.product_id,
      listing_id: deal.listing_id,
      isStorefront: deal.isStorefront,
    });

    if (deal.isStorefront && deal.product_id) {
      // Для товаров витрин - используем product_id без префикса
      const url = `/${locale}/c2c/${deal.product_id}`;
      logger.debug('Storefront URL:', url);
      return url;
    } else if (deal.listing_id) {
      // Для обычных объявлений - используем listing_id
      const url = `/${locale}/c2c/${deal.listing_id}`;
      logger.debug('Listing URL:', url);
      return url;
    } else {
      // Fallback - извлекаем чистый ID из deal.id убрав префиксы
      const cleanId =
        typeof deal.id === 'string'
          ? deal.id.replace(/^(ml_|sp_)/, '')
          : deal.id;
      const url = `/${locale}/c2c/${cleanId}`;
      logger.debug('Fallback URL:', url);
      return url;
    }
  };

  // Функция для добавления в корзину
  const handleAddToCart = async (deal: any) => {
    logger.debug('handleAddToCart called with full deal object:', deal);

    try {
      // Если это не витрина, выходим
      if (!deal.isStorefront) {
        toast.error('Этот товар нельзя добавить в корзину');
        return;
      }

      // Пытаемся получить product_id
      const productId =
        deal.product_id ||
        (deal.id && typeof deal.id === 'string'
          ? parseInt(deal.id.replace('sp_', ''))
          : null);

      if (!productId) {
        toast.error('Не удалось определить товар');
        console.error('Cannot determine product_id from deal:', deal);
        return;
      }

      // Если нет storefront_id, пытаемся получить его через API
      let storefrontId = deal.storefront_id;

      if (!storefrontId) {
        logger.debug(
          'No storefront_id in deal, fetching product details from API...'
        );
        try {
          const response = await apiClient.get(`/marketplace/storefronts/products/${productId}`);
          if (response.data && response.data.storefront_id) {
            storefrontId = response.data.storefront_id;
            logger.debug('Got storefront_id from API:', storefrontId);
          }
        } catch (apiError) {
          console.error('Failed to fetch product details:', apiError);
        }
      }

      if (!storefrontId) {
        toast.error('Не удалось определить магазин');
        console.error('Cannot determine storefront_id for product:', productId);
        return;
      }

      logger.debug('Adding to cart:', { storefrontId, productId });

      // Добавляем товар в корзину
      await dispatch(
        addToCart({
          storefrontId: storefrontId,
          item: {
            product_id: productId,
            quantity: 1,
          },
        })
      ).unwrap();

      toast.success('Товар добавлен в корзину');
    } catch (error) {
      console.error('Error adding to cart:', error);
      toast.error('Ошибка при добавлении в корзину');
    }
  };

  // Функция для открытия чата
  const handleStartChat = (deal: any) => {
    logger.debug('handleStartChat called with deal:', deal);

    if (!user) {
      // Если пользователь не авторизован, перенаправляем на страницу входа
      router.push('/login');
      return;
    }

    // Определяем URL для чата в зависимости от типа объявления
    if (deal.isStorefront && deal.storefront_id) {
      // B2C - чат с витриной, передаем storefront_product_id и seller_id (владелец витрины)
      logger.debug(
        'Opening B2C chat with storefront_id:',
        deal.storefront_id,
        'product_id:',
        deal.product_id || deal.id,
        'seller_id:',
        deal.user_id
      );
      const productId = deal.product_id || deal.id;
      if (!deal.user_id) {
        console.error(
          'Missing seller_id for storefront product chat. Deal data:',
          deal
        );
        return;
      }
      router.push(
        `/${locale}/chat?storefront_product_id=${productId}&seller_id=${deal.user_id}`
      );
    } else if (deal.user_id) {
      // C2C - чат с продавцом обычного объявления
      const listingId = deal.listing_id || deal.id;
      logger.debug(
        'Opening C2C chat with user_id:',
        deal.user_id,
        'listing_id:',
        listingId
      );
      router.push(
        `/${locale}/chat?listing_id=${listingId}&seller_id=${deal.user_id}`
      );
    } else {
      console.error('Missing seller information for chat. Deal data:', deal);
    }
  };

  // Инициализация избранного
  const initializeFavorites = useCallback(async () => {
    if (user) {
      await favoritesService.initialize();
      const favorites = await favoritesService.getFavorites();
      setFavoriteIds(new Set(favorites.map((f) => Number(f.id))));
    }
  }, [user]);

  // Устанавливаем mounted после гидрации для предотвращения hydration mismatch
  useEffect(() => {
    setMounted(true);

    // OAuth авторизация теперь обрабатывается через BFF proxy с httpOnly cookies
    // Токен больше не передается через URL и не сохраняется в localStorage

    // Prefetch delivery providers для будущего использования в cart/checkout
    // Это не блокирует рендеринг страницы, но кэширует данные заранее
    dispatch(fetchProviders());
  }, [dispatch]);

  // Отдельный эффект для инициализации избранного при смене пользователя
  useEffect(() => {
    initializeFavorites();
  }, [initializeFavorites]);

  // Обработчик добавления/удаления из избранного
  const handleToggleFavorite = async (
    listingId: number | string,
    e: React.MouseEvent
  ) => {
    e.preventDefault();
    e.stopPropagation();

    if (!user) {
      toast.error('Войдите, чтобы добавить в избранное');
      router.push(`/${locale}/login`);
      return;
    }

    // Извлекаем числовой ID для проверки в favoriteIds
    const numericId =
      typeof listingId === 'string' && listingId.startsWith('sp_')
        ? parseInt(listingId.replace('sp_', ''))
        : Number(listingId);

    const isCurrentlyFavorite = favoriteIds.has(numericId);

    // Оптимистичное обновление UI
    if (isCurrentlyFavorite) {
      setFavoriteIds((prev) => {
        const newSet = new Set(prev);
        newSet.delete(numericId);
        return newSet;
      });
    } else {
      setFavoriteIds((prev) => new Set([...prev, numericId]));
    }

    // Выполняем запрос к API (передаем оригинальный listingId с префиксом)
    const success = await favoritesService.toggleFavorite(listingId);

    if (!success) {
      // Откатываем изменения если ошибка
      if (isCurrentlyFavorite) {
        setFavoriteIds((prev) => new Set([...prev, numericId]));
      } else {
        setFavoriteIds((prev) => {
          const newSet = new Set(prev);
          newSet.delete(numericId);
          return newSet;
        });
      }
    }
  };

  // Баннеры для hero секции
  const banners = [
    {
      id: 1,
      title: t('blackFridayTitle'),
      subtitle: t('blackFridaySubtitle'),
      bgColor: 'bg-gradient-to-r from-purple-600 to-pink-600',
      cta: t('blackFridayCta'),
      image: '🛍️',
      badge: t('blackFridayBadge'),
      details: t('blackFridayDetails'),
    },
    {
      id: 2,
      title: t('freeDeliveryTitle'),
      subtitle: t('freeDeliverySubtitle'),
      bgColor: 'bg-gradient-to-r from-blue-600 to-cyan-600',
      cta: t('freeDeliveryCta'),
      image: '📦',
    },
    {
      id: 3,
      title: t('buyerProtectionTitle'),
      subtitle: t('buyerProtectionSubtitle'),
      bgColor: 'bg-gradient-to-r from-green-600 to-teal-600',
      cta: t('buyerProtectionCta'),
      image: '🔒',
    },
  ];

  // OAuth токен обрабатывается в AuthContext.tsx
  // После успешной OAuth авторизации показываем уведомление
  useEffect(() => {
    // Проверяем, если пользователь только что вошел через OAuth
    // AuthContext обрабатывает токен из URL и устанавливает пользователя
    if (user && searchParams?.get('auth_token')) {
      // Токен уже обработан в AuthContext, просто показываем уведомление
      toast.success(t('loginSuccessful') || 'Successfully logged in!');
    }
  }, [user, searchParams, t]);

  // Загружаем избранное при изменении пользователя
  useEffect(() => {
    if (user) {
      initializeFavorites();
    } else {
      setFavoriteIds(new Set());
      favoritesService.clearCache();
    }
  }, [user, initializeFavorites]);

  // Автоматическая смена баннеров
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentBanner((prev) => (prev + 1) % banners.length);
    }, 5000);
    return () => clearInterval(interval);
  }, [banners.length]);

  // Популярные поисковые запросы
  const trendingSearches = [
    'iPhone 15',
    'PS5',
    'Квартира центр',
    'MacBook',
    'Электросамокат',
    'Диван',
    'AirPods',
    'Nike кроссовки',
    'Холодильник',
    'Велосипед',
  ];

  // Загружаем категории из Redux
  const {
    categories: reduxCategories,
    popularCategories: reduxPopularCategories,
    isLoadingCategories: reduxIsLoadingCategories,
    isLoadingPopular: reduxIsLoadingPopular,
  } = useSelector((state: RootState) => state.categories);

  // Обновляем локальное состояние из Redux
  useEffect(() => {
    if (reduxCategories.length > 0) {
      setCategories(reduxCategories);
    }
    if (reduxPopularCategories.length > 0) {
      setPopularCategories(reduxPopularCategories);
    }
    setIsLoadingCategories(reduxIsLoadingCategories || reduxIsLoadingPopular);
  }, [
    reduxCategories,
    reduxPopularCategories,
    reduxIsLoadingCategories,
    reduxIsLoadingPopular,
  ]);

  // Загрузка категорий через Redux
  useEffect(() => {
    // Загружаем обычные категории
    dispatch(fetchCategories());
    // Загружаем популярные категории с учетом локали
    dispatch(fetchPopularCategories({ locale: locale as string }));
  }, [dispatch, locale]);

  // Загрузка витрин (официальных магазинов)
  useEffect(() => {
    const loadStorefronts = async () => {
      setIsLoadingStores(true);
      try {
        // Загружаем активные витрины
        // Сначала загружаем больше витрин, чтобы выбрать те, у которых есть изображения
        const response = await apiClient.get(
          '/b2c?is_active=true&limit=10&sort_by=products_count&sort_order=desc'
        );

        if (response.data && response.data.storefronts) {
          // Форматируем данные витрин для отображения
          const formattedStores = response.data.storefronts.map(
            (store: any) => {
              // Генерируем цвет для аватара на основе имени
              const colors = [
                '6366f1',
                'ec4899',
                '10b981',
                'ef4444',
                'f59e0b',
                '8b5cf6',
              ];
              const colorIndex = store.id % colors.length;
              const bgColor = colors[colorIndex];

              // Берем первые 2 буквы названия для аватара
              const initials = store.name.substring(0, 2).toUpperCase();

              // Получаем случайное изображение для фона (можно заменить на реальные изображения категорий)
              const bgImages = [
                'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=400&h=200&fit=crop',
                'https://images.unsplash.com/photo-1490481651871-ab68de25d43d?w=400&h=200&fit=crop',
                'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400&h=200&fit=crop',
                'https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?w=400&h=200&fit=crop',
              ];
              const bgImage = bgImages[store.id % bgImages.length];

              return {
                id: store.id,
                name: store.name,
                category: store.category_name || t('defaultStoreCategory'),
                logo:
                  store.logo_url ||
                  `https://ui-avatars.com/api/?name=${initials}&background=${bgColor}&color=fff&size=128`,
                followers: store.followers_count
                  ? store.followers_count >= 1000
                    ? `${(store.followers_count / 1000).toFixed(1)}K`
                    : store.followers_count.toString()
                  : '-', // Показываем прочерк вместо 0
                products: store.products_count || '-', // Показываем прочерк если нет товаров
                rating: store.rating || null, // Используем реальный рейтинг из API
                viewsCount: store.views_count || 0, // Используем реальные просмотры из API
                verified: store.is_verified || false,
                discount: store.discount_text || '',
                bgImage: store.banner_url || bgImage,
                slug: store.slug,
                description: store.description,
              };
            }
          );

          // Приоритизируем витрины с изображениями
          const storesWithImages = formattedStores.filter(
            (store: any) => store.logo && store.logo.includes('svetu.rs')
          );
          const storesWithoutImages = formattedStores.filter(
            (store: any) => !store.logo || !store.logo.includes('svetu.rs')
          );

          // Комбинируем: сначала с изображениями, потом без
          const finalStores = [
            ...storesWithImages,
            ...storesWithoutImages,
          ].slice(0, 4);

          setOfficialStores(finalStores);
          logger.debug('Loaded storefronts:', finalStores);
        } else {
          // Если нет реальных витрин, показываем пустой массив
          setOfficialStores([]);
        }
      } catch (error) {
        console.error('Failed to load storefronts:', error);
        // В случае ошибки показываем пустой массив
        setOfficialStores([]);
      } finally {
        setIsLoadingStores(false);
      }
    };

    loadStorefronts();
  }, [t]);

  // Загрузка товаров через API поиска
  useEffect(() => {
    const loadListings = async () => {
      setIsLoadingListings(true);

      try {
        // Загружаем больше объявлений для смешанного показа C2C и B2C
        const searchParams = new URLSearchParams();
        searchParams.append('query', '');
        searchParams.append('size', '25');
        searchParams.append('page', '1');
        searchParams.append('sort', 'created_at');
        searchParams.append('sortDirection', 'desc');
        searchParams.append('lang', locale);
        searchParams.append('status', 'active');
        searchParams.append('product_types[]', 'marketplace');
        searchParams.append('product_types[]', 'storefront');

        logger.debug('Request URL:', `/search?${searchParams.toString()}`);
        const response = await apiClient.get(
          `/search?${searchParams.toString()}`
        );
        logger.debug('API Response:', response.data);

        if (
          response.data &&
          response.data.items &&
          response.data.items.length > 0
        ) {
          // Разделяем объявления на C2C и B2C для смешанного показа
          const allListings = response.data.items;
          logger.debug(
            'All listings product types:',
            JSON.stringify(
              allListings.map((l: any) => ({
                id: l.id,
                product_id: l.product_id,
                product_type: l.product_type,
                name: l.name || l.title,
              })),
              null,
              2
            )
          );
          const c2cListings = allListings.filter(
            (listing: any) => listing.product_type !== 'storefront'
          );
          const b2cListings = allListings.filter(
            (listing: any) => listing.product_type === 'storefront'
          );

          // Создаем смешанную выборку: смешиваем C2C и B2C объявления
          let selectedListings = [];

          // Смешиваем объявления для разнообразия
          // Берем поочередно из обоих списков
          const maxItems = 8;
          let c2cIndex = 0;
          let b2cIndex = 0;

          // Добавляем объявления поочередно: 2 C2C, 1 B2C, повторяем
          while (selectedListings.length < maxItems) {
            // Добавляем 2 C2C если есть
            for (
              let i = 0;
              i < 2 &&
              c2cIndex < c2cListings.length &&
              selectedListings.length < maxItems;
              i++
            ) {
              selectedListings.push(c2cListings[c2cIndex++]);
            }

            // Добавляем 1 B2C если есть
            if (
              b2cIndex < b2cListings.length &&
              selectedListings.length < maxItems
            ) {
              selectedListings.push(b2cListings[b2cIndex++]);
            }

            // Если закончились оба типа, выходим
            if (
              c2cIndex >= c2cListings.length &&
              b2cIndex >= b2cListings.length
            ) {
              break;
            }
          }

          // Если объявлений мало, берем все что есть
          if (selectedListings.length === 0) {
            selectedListings = allListings.slice(0, maxItems);
          }

          logger.debug(
            `Mixed selection: ${selectedListings.filter((l: any) => !l.storefrontId).length} C2C + ${selectedListings.filter((l: any) => l.storefrontId).length} B2C`
          );

          const apiListings = selectedListings.map(
            (listing: any, index: number) => {
              // Логируем структуру данных для отладки
              logger.debug('Processing listing:', {
                id: listing.id,
                name: listing.name,
                images: listing.images,
                hasImages: listing.images && listing.images.length > 0,
                firstImage: listing.images && listing.images[0],
              });

              // Вычисляем скидку если есть старая цена
              let discount = null;
              let oldPrice = null;

              // Проверяем наличие скидки из API или вычисляем из старой цены
              if (listing.has_discount && listing.old_price) {
                oldPrice = `${listing.old_price} ${listing.currency || 'РСД'}`;
                if (listing.discount_percentage) {
                  discount = `-${listing.discount_percentage}%`;
                } else if (listing.old_price > listing.price) {
                  const discountPercent = Math.round(
                    ((listing.old_price - listing.price) / listing.old_price) *
                      100
                  );
                  discount = `-${discountPercent}%`;
                }
              } else if (
                listing.originalPrice &&
                listing.price &&
                listing.originalPrice > listing.price
              ) {
                const discountPercent = Math.round(
                  ((listing.originalPrice - listing.price) /
                    listing.originalPrice) *
                    100
                );
                discount = `-${discountPercent}%`;
                oldPrice = `${listing.originalPrice} РСД`;
              }

              // Добавляем подробное логирование для отладки
              logger.debug('Mapping listing with storefront data:', {
                listing_id: listing.id,
                product_id: listing.product_id,
                product_type: listing.product_type,
                storefront_from_api: listing.storefront,
                storefront_id_direct: listing.storefront_id,
                user_from_api: listing.user,
                user_id_direct: listing.user_id,
              });

              // Формируем уникальный ключ для React
              // API уже возвращает listing.id с префиксом (sp_XXX или ml_XXX)
              // Добавляем index как fallback для гарантии уникальности
              const uniqueKey =
                listing.id || `${listing.product_type || 'item'}_${index}`;

              const mappedListing = {
                id: uniqueKey, // Используем уникальный ключ
                product_id:
                  listing.product_type === 'storefront'
                    ? listing.product_id
                    : null,
                title: listing.name || listing.title,
                price: `${listing.price} ${listing.currency || 'РСД'}`,
                oldPrice,
                discount,
                // Сохраняем все данные для локализации
                // Обрабатываем разные форматы location из search API и marketplace API
                location:
                  typeof listing.location === 'object' && listing.location
                    ? `${listing.location.city || ''}, ${listing.location.country || ''}`
                        .trim()
                        .replace(/^,\s*|\s*,$/, '')
                    : listing.location ||
                      listing.address_city ||
                      listing.city ||
                      'Сербия',
                city:
                  typeof listing.location === 'object' && listing.location
                    ? listing.location.city
                    : listing.city || listing.address_city,
                country:
                  typeof listing.location === 'object' && listing.location
                    ? listing.location.country
                    : listing.country || listing.address_country,
                address_multilingual:
                  listing.location?.address_multilingual ||
                  listing.address_multilingual,
                translations: listing.translations,
                image: (() => {
                  // Проверяем наличие изображений
                  if (listing.images && listing.images.length > 0) {
                    const firstImage = listing.images[0];

                    // Извлекаем URL из объекта изображения или используем как строку
                    let imageUrl: string;
                    if (typeof firstImage === 'object' && firstImage !== null) {
                      imageUrl = firstImage.url || firstImage.public_url || '';
                    } else if (typeof firstImage === 'string') {
                      imageUrl = firstImage;
                    } else {
                      imageUrl = '';
                    }

                    // Логируем для отладки
                    logger.debug(
                      'Building image URL for listing',
                      listing.id,
                      ':',
                      imageUrl,
                      'firstImage:',
                      firstImage
                    );

                    // Если у нас есть валидный URL
                    if (imageUrl) {
                      // Если URL уже полный (начинается с http), используем как есть
                      if (imageUrl.startsWith('http')) {
                        return imageUrl;
                      }

                      // Иначе строим URL через configManager
                      return configManager.buildImageUrl(imageUrl);
                    }
                  }

                  // Fallback изображение только если действительно нет изображений
                  logger.debug(
                    'No images for listing',
                    listing.id,
                    ', using fallback'
                  );
                  return 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=400&h=300&fit=crop';
                })(),
                rating: listing.rating || null, // Используем реальный рейтинг из API
                reviews: listing.reviewCount || listing.review_count || null, // Используем реальные отзывы из API
                viewsCount: listing.view_count ?? listing.views_count ?? null, // Используем реальные просмотры из API, включая 0
                isNew:
                  new Date(listing.created_at || listing.createdAt) >
                  new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // Новое если создано за последнюю неделю
                isPremium: listing.isPremium || false,
                isFavorite: favoriteIds.has(listing.id),
                category: listing.category?.name || listing.categoryName,
                isStorefront: listing.product_type === 'storefront',
                // Извлекаем user_id из объекта user (search API) или напрямую (marketplace API)
                user_id: listing.user?.id || listing.user_id,
                // Извлекаем storefront_id из объекта storefront (search API) или напрямую
                storefront_id: listing.storefront?.id || listing.storefront_id,
                // Сохраняем оригинальный listing_id для C2C товаров (удаляем префикс ml_ если есть)
                listing_id:
                  listing.product_type !== 'storefront'
                    ? typeof listing.id === 'string' &&
                      listing.id.startsWith('ml_')
                      ? listing.id.replace('ml_', '')
                      : listing.id
                    : null,
              };

              // Логирование для отладки
              if (!mappedListing.user_id && !mappedListing.storefront_id) {
                console.warn('Listing missing user_id and storefront_id:', {
                  original_listing: listing,
                  mapped_listing: mappedListing,
                });
              }

              return mappedListing;
            }
          );

          setListings(apiListings);
          logger.debug(
            'Loaded hot deals from API:',
            apiListings.length,
            'items'
          );
        } else {
          console.warn(
            'No listings data in API response, showing demo content for development'
          );
          // Fallback: показываем несколько демо объявлений когда API пуст
          setListings([
            {
              id: 'demo-1',
              title: 'iPhone 15 Pro Max 256GB',
              price: '130000 РСД',
              oldPrice: '167000 РСД',
              discount: '-21%',
              location: 'Белград',
              image:
                'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400&h=300&fit=crop',
              rating: 4.8,
              reviews: 234,
              viewsCount: 0,
              isNew: true,
              isPremium: false,
              isFavorite: false,
            },
            {
              id: 'demo-2',
              title: 'MacBook Air M3 13" 512GB',
              price: '155000 РСД',
              oldPrice: '190000 РСД',
              discount: '-19%',
              location: 'Белград',
              image:
                'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?w=400&h=300&fit=crop',
              rating: 4.9,
              reviews: 567,
              viewsCount: 0,
              isNew: true,
              isPremium: false,
              isFavorite: false,
            },
          ]);
        }
      } catch (error) {
        console.error('Failed to load hot deals from API:', error);

        // В случае ошибки показываем пустой массив вместо mock данных
        setListings([]);
      } finally {
        setIsLoadingListings(false);
      }
    };

    loadListings();
  }, [locale, favoriteIds]);

  return (
    <PageTransition mode="fade">
      <div className="min-h-screen bg-gradient-to-b from-base-100 to-base-200">
        {/* Hero секция с баннерами */}
        <section className="relative overflow-hidden">
          <div className="container mx-auto px-4 py-6">
            <div className="grid lg:grid-cols-3 gap-6">
              {/* Главный баннер */}
              <div className="lg:col-span-2">
                <AnimatePresence mode="wait">
                  <motion.div
                    key={currentBanner}
                    initial={{ opacity: 0, x: 100 }}
                    animate={{ opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: -100 }}
                    className={`relative rounded-2xl p-8 lg:p-12 text-white overflow-hidden h-[400px] ${banners[currentBanner].bgColor}`}
                    style={{
                      backgroundImage: `linear-gradient(rgba(0,0,0,0.3), rgba(0,0,0,0.3)), url('https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=1200&h=600&fit=crop')`,
                      backgroundSize: 'cover',
                      backgroundPosition: 'center',
                    }}
                  >
                    {banners[currentBanner].badge && (
                      <div className="absolute top-4 right-4 badge badge-warning badge-lg">
                        {banners[currentBanner].badge}
                      </div>
                    )}
                    <div className="relative z-10">
                      <h1 className="text-4xl lg:text-5xl font-bold mb-4 drop-shadow-lg">
                        {banners[currentBanner].title}
                      </h1>
                      <p className="text-xl mb-2 drop-shadow-lg">
                        {banners[currentBanner].subtitle}
                      </p>
                      {banners[currentBanner].details && (
                        <p className="text-sm mb-6 opacity-90 drop-shadow-lg">
                          {banners[currentBanner].details}
                        </p>
                      )}
                      <button className="btn btn-white btn-lg">
                        {banners[currentBanner].cta}
                        <FiChevronRight className="w-5 h-5 ml-2" />
                      </button>
                    </div>
                    <div className="absolute right-8 top-1/2 -translate-y-1/2 text-8xl opacity-20">
                      {banners[currentBanner].image}
                    </div>
                    {/* Индикаторы */}
                    <div className="absolute bottom-4 left-8 flex gap-2">
                      {banners.map((_, idx) => (
                        <button
                          key={idx}
                          className={`w-2 h-2 rounded-full transition-all ${
                            idx === currentBanner
                              ? 'w-8 bg-white'
                              : 'bg-white/50'
                          }`}
                          onClick={() => setCurrentBanner(idx)}
                          aria-label={_tCommon('carousel.goToSlide', {
                            index: (idx + 1).toString(),
                          })}
                          aria-current={idx === currentBanner}
                        />
                      ))}
                    </div>
                  </motion.div>
                </AnimatePresence>
              </div>

              {/* Боковые карточки */}
              <div className="space-y-4">
                <div className="card bg-gradient-to-br from-orange-500 to-red-500 text-white h-[190px]">
                  <div className="card-body">
                    <h3 className="card-title text-white">
                      {t('lightningDeals')}
                    </h3>
                    <p>{t('lightningDealsSubtitle')}</p>
                    <div className="text-2xl font-bold">02:45:18</div>
                    <button className="btn btn-white btn-sm">
                      {t('watch')}
                    </button>
                  </div>
                </div>
                <div className="card bg-gradient-to-br from-green-500 to-teal-500 text-white h-[190px]">
                  <div className="card-body">
                    <h3 className="card-title text-white">
                      {t('newUsersGift')}
                    </h3>
                    <p>{t('newUsersGiftSubtitle')}</p>
                    <button className="btn btn-white btn-sm">{t('get')}</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Популярные категории */}
        <section className="py-8">
          <div className="container mx-auto px-4">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold flex items-center gap-2">
                <HiOutlineSparkles className="w-6 h-6 text-warning" />
                {t('popularCategories')}
              </h2>
              <NestedCategorySelector
                categories={categories}
                selectedCategory={selectedCategory}
                onChange={(categoryId) => {
                  router.push(`/${locale}/search?category=${categoryId}`);
                }}
                placeholder={t('allCategories')}
                showCounts={true}
                className="btn btn-primary btn-sm gap-2"
              />
            </div>
            {isLoadingCategories ? (
              <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
                {[...Array(8)].map((_, i) => (
                  <div key={i} className="card bg-base-100">
                    <div className="card-body p-4">
                      <div className="skeleton h-14 w-14 rounded-full mx-auto mb-2"></div>
                      <div className="skeleton h-4 w-full"></div>
                      <div className="skeleton h-3 w-1/2 mx-auto"></div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
                {popularCategories.map((cat) => {
                  const Icon = getCategoryIcon(cat.iconName);
                  return (
                    <Link
                      key={cat.id}
                      href={`/search?category=${cat.id}`}
                      className="group"
                    >
                      <div className="card bg-base-100 hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
                        <div className="card-body p-4 text-center">
                          <div
                            className={`mx-auto mb-2 p-3 rounded-full bg-base-200 group-hover:bg-primary/10 transition-colors`}
                          >
                            <Icon className={`w-8 h-8 ${cat.color}`} />
                          </div>
                          <h3 className="font-medium text-sm">{cat.name}</h3>
                          <p className="text-xs text-base-content/60">
                            {cat.count}
                          </p>
                        </div>
                      </div>
                    </Link>
                  );
                })}
              </div>
            )}
          </div>
        </section>

        {/* Горячие предложения */}
        <section className="container mx-auto px-4 py-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <HiOutlineSparkles className="text-warning" />
              {t('hotDeals')}
            </h2>
            <div className="flex gap-2">
              <button
                onClick={() => setViewMode('grid')}
                className={`btn btn-sm ${viewMode === 'grid' ? 'btn-primary' : 'btn-ghost'}`}
                aria-label={_tCommon('view.gridView')}
                aria-pressed={viewMode === 'grid'}
              >
                <FiGrid className="w-4 h-4" aria-hidden="true" />
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`btn btn-sm ${viewMode === 'list' ? 'btn-primary' : 'btn-ghost'}`}
                aria-label={_tCommon('view.listView')}
                aria-pressed={viewMode === 'list'}
              >
                <FiList className="w-4 h-4" aria-hidden="true" />
              </button>
              <Link href="/hot" className="btn btn-sm btn-ghost">
                {t('allDeals')}
              </Link>
            </div>
          </div>

          {isLoadingListings ? (
            <div
              className={`grid ${viewMode === 'grid' ? 'grid-cols-2 lg:grid-cols-4' : 'grid-cols-1'} gap-4`}
            >
              {[...Array(8)].map((_, i) => (
                <div key={i} className="card bg-base-100">
                  <div className="skeleton h-48"></div>
                  <div className="card-body">
                    <div className="skeleton h-4 w-3/4"></div>
                    <div className="skeleton h-4 w-1/2"></div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div
              className={`grid ${viewMode === 'grid' ? 'grid-cols-2 lg:grid-cols-4' : 'grid-cols-1'} gap-4`}
            >
              {listings.map((deal) => (
                <Link
                  key={deal.id}
                  href={getListingUrl(deal)}
                  className="block"
                >
                  <motion.div
                    whileHover={{ scale: 1.02 }}
                    className="card bg-base-100 hover:shadow-xl transition-all"
                  >
                    <figure className="relative h-48 overflow-hidden">
                      <SafeImage
                        src={deal.image}
                        alt={deal.translations?.[locale]?.title || deal.title}
                        fill
                        className="object-cover"
                        sizes="(max-width: 640px) 50vw, (max-width: 1024px) 25vw, 25vw"
                      />

                      {/* Значок витрины для B2C объявлений */}
                      {deal.isStorefront && (
                        <div className="badge badge-info absolute top-2 left-2 flex items-center gap-1">
                          <FiShoppingBag className="w-3 h-3" />
                          {t('storefront')}
                        </div>
                      )}

                      {/* Остальные бейджи с учетом значка витрины */}
                      {deal.isNew && !deal.isStorefront && (
                        <div className="badge badge-secondary absolute top-2 left-2">
                          NEW
                        </div>
                      )}
                      {deal.isNew && deal.isStorefront && (
                        <div className="badge badge-secondary absolute top-12 left-2">
                          NEW
                        </div>
                      )}

                      {deal.discount && !deal.isStorefront && (
                        <div className="badge badge-error absolute top-2 left-2">
                          {deal.discount}
                        </div>
                      )}
                      {deal.discount && deal.isStorefront && (
                        <div className="badge badge-error absolute top-12 left-2">
                          {deal.discount}
                        </div>
                      )}

                      {deal.isPremium && (
                        <div className="badge badge-warning absolute top-2 right-2">
                          PREMIUM
                        </div>
                      )}

                      <button
                        className={`btn btn-circle btn-sm absolute ${deal.isPremium ? 'top-12 right-2' : 'top-2 right-2'} bg-base-100/80 hover:bg-base-100`}
                        onClick={(e) => {
                          // Используем listing_id если он есть
                          let listingId: number | string;
                          if (deal.listing_id) {
                            listingId =
                              typeof deal.listing_id === 'number'
                                ? deal.listing_id
                                : parseInt(String(deal.listing_id));
                          } else if (
                            typeof deal.id === 'string' &&
                            deal.id.startsWith('ml_')
                          ) {
                            listingId = parseInt(deal.id.replace('ml_', ''));
                          } else if (
                            typeof deal.id === 'string' &&
                            deal.id.startsWith('sp_')
                          ) {
                            // Сохраняем полный ID с префиксом для storefront продуктов
                            listingId = deal.id;
                          } else {
                            listingId = parseInt(String(deal.id));
                          }

                          // Проверяем валидность ID
                          const isValidId =
                            (typeof listingId === 'string' &&
                              listingId.startsWith('sp_')) ||
                            (typeof listingId === 'number' &&
                              !isNaN(listingId));

                          if (isValidId) {
                            handleToggleFavorite(listingId, e);
                          } else {
                            console.error(
                              'Invalid listing ID:',
                              deal.id,
                              deal.listing_id
                            );
                          }
                        }}
                        aria-label={
                          deal.isFavorite
                            ? _tCommon('product.removeFromFavorites')
                            : _tCommon('product.addToFavorites')
                        }
                      >
                        <FiHeart
                          className={`w-4 h-4 ${deal.isFavorite ? 'fill-error text-error' : ''}`}
                          aria-hidden="true"
                        />
                      </button>
                    </figure>
                    <div className="card-body p-4">
                      <h3 className="card-title text-base line-clamp-2">
                        {deal.translations?.[locale]?.title || deal.title}
                      </h3>
                      <div className="flex items-center gap-2 text-sm">
                        <FiMapPin className="w-3 h-3" />
                        <span className="text-base-content/60">
                          {(() => {
                            // Приоритет 1: Мультиязычный адрес из unified search
                            if (deal.address_multilingual?.[locale]) {
                              return deal.address_multilingual[locale];
                            }

                            // Приоритет 2: Перевод локации из translations
                            if (deal.translations?.[locale]?.location) {
                              return deal.translations[locale].location;
                            }

                            // Приоритет 3: Fallback на строковое значение location
                            if (typeof deal.location === 'string') {
                              return deal.location;
                            }

                            // Приоритет 4: Составляем из city и country
                            const city = deal.city || '';
                            const country = deal.country || '';
                            return (
                              `${city}${city && country ? ', ' : ''}${country}`.trim() ||
                              'Location not specified'
                            );
                          })()}
                        </span>
                      </div>
                      <div className="flex items-center justify-between mt-1">
                        <div className="flex items-center gap-3">
                          {deal.rating !== null && deal.rating > 0 && (
                            <div className="flex items-center gap-1 text-sm">
                              <FiStar className="w-3 h-3 fill-warning text-warning" />
                              <span className="font-medium">
                                {deal.rating.toFixed(1)}
                              </span>
                              {deal.reviews !== null && deal.reviews > 0 && (
                                <span className="text-base-content/60">
                                  ({deal.reviews})
                                </span>
                              )}
                            </div>
                          )}
                          {deal.viewsCount !== null &&
                            deal.viewsCount !== undefined && (
                              <div className="flex items-center gap-1 text-sm text-base-content/60">
                                <AiOutlineEye className="w-3.5 h-3.5" />
                                <span>
                                  {deal.viewsCount === 0
                                    ? '-'
                                    : deal.viewsCount.toLocaleString()}
                                </span>
                              </div>
                            )}
                        </div>
                      </div>
                      <div className="flex items-center gap-2 mt-2">
                        {deal.oldPrice && (
                          <span className="text-base-content/40 line-through text-sm">
                            {deal.oldPrice}
                          </span>
                        )}
                        <p className="text-xl font-bold text-primary">
                          {deal.price}
                        </p>
                      </div>

                      {/* Кнопки действий в зависимости от типа объявления */}
                      {deal.isStorefront ? (
                        // B2C (витрина) - кнопка "В корзину" + "В избранное" + "Написать в чат"
                        <div className="flex gap-2 mt-2">
                          <button
                            className="btn btn-primary btn-sm flex-1"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              handleAddToCart(deal);
                            }}
                          >
                            {t('addToCart')}
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              // Для товаров витрин используем ID с префиксом sp_
                              const productId =
                                deal.id &&
                                typeof deal.id === 'string' &&
                                deal.id.startsWith('sp_')
                                  ? deal.id
                                  : `sp_${deal.id || deal.product_id}`;
                              handleToggleFavorite(productId, e);
                            }}
                            aria-label={
                              favoriteIds.has(
                                parseInt(
                                  String(deal.id || deal.product_id).replace(
                                    'sp_',
                                    ''
                                  )
                                )
                              )
                                ? _tCommon('product.removeFromFavorites')
                                : _tCommon('product.addToFavorites')
                            }
                          >
                            <FiHeart
                              className={`w-4 h-4 ${
                                favoriteIds.has(
                                  parseInt(
                                    String(deal.id || deal.product_id).replace(
                                      'sp_',
                                      ''
                                    )
                                  )
                                )
                                  ? 'fill-current text-error'
                                  : ''
                              }`}
                              aria-hidden="true"
                            />
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              handleStartChat(deal);
                            }}
                            aria-label={_tCommon('chat.startChat')}
                          >
                            <FiMessageCircle
                              className="w-4 h-4"
                              aria-hidden="true"
                            />
                          </button>
                        </div>
                      ) : (
                        // C2C (обычное объявление) - "Написать в чат" + "В избранное"
                        <div className="flex gap-2 mt-2">
                          <button
                            className="btn btn-primary btn-sm flex-1"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              handleStartChat(deal);
                            }}
                          >
                            <FiMessageCircle className="w-4 h-4 mr-1" />
                            {t('writeToSeller')}
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={(e) => {
                              // Используем listing_id если он есть
                              let listingId: number | string;
                              if (deal.listing_id) {
                                listingId =
                                  typeof deal.listing_id === 'number'
                                    ? deal.listing_id
                                    : parseInt(String(deal.listing_id));
                              } else if (
                                typeof deal.id === 'string' &&
                                deal.id.startsWith('ml_')
                              ) {
                                listingId = parseInt(
                                  deal.id.replace('ml_', '')
                                );
                              } else if (
                                typeof deal.id === 'string' &&
                                deal.id.startsWith('sp_')
                              ) {
                                // Сохраняем полный ID с префиксом для storefront продуктов
                                listingId = deal.id;
                              } else {
                                listingId = parseInt(String(deal.id));
                              }

                              // Проверяем валидность ID
                              const isValidId =
                                (typeof listingId === 'string' &&
                                  listingId.startsWith('sp_')) ||
                                (typeof listingId === 'number' &&
                                  !isNaN(listingId));

                              if (isValidId) {
                                handleToggleFavorite(listingId, e);
                              } else {
                                console.error(
                                  'Invalid listing ID:',
                                  deal.id,
                                  deal.listing_id
                                );
                              }
                            }}
                            aria-label={
                              deal.isFavorite
                                ? _tCommon('product.removeFromFavorites')
                                : _tCommon('product.addToFavorites')
                            }
                          >
                            <FiHeart
                              className={`w-4 h-4 ${deal.isFavorite ? 'fill-error text-error' : ''}`}
                              aria-hidden="true"
                            />
                          </button>
                        </div>
                      )}
                    </div>
                  </motion.div>
                </Link>
              ))}
            </div>
          )}
        </section>

        {/* Товары рядом с вами */}
        <section className="container mx-auto px-4 py-8">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <FiMapPin className="text-primary" />
            {t('nearbyProducts')}
          </h2>
          <div className="grid lg:grid-cols-3 gap-6">
            {/* Карта */}
            <div className="lg:col-span-2">
              <div className="card bg-base-100 overflow-hidden">
                <div className="card-body p-0">
                  <EnhancedMapSection
                    className="h-96 w-full"
                    listings={listings.map((item) => ({
                      id: item.id,
                      latitude:
                        item.location?.lat ||
                        44.8125 + (Math.random() - 0.5) * 0.02,
                      longitude:
                        item.location?.lng ||
                        20.4612 + (Math.random() - 0.5) * 0.02,
                      price: item.price,
                      title: item.translations?.[locale]?.title || item.title,
                      category: item.category,
                      imageUrl: item.image,
                      isStorefront: item.isStorefront,
                    }))}
                    userLocation={
                      userLocation
                        ? {
                            latitude: userLocation[0],
                            longitude: userLocation[1],
                          }
                        : undefined
                    }
                    searchRadius={5000}
                    showRadius={true}
                    enableClustering={true}
                  />
                </div>
              </div>
            </div>

            {/* Статистика */}
            <div className="space-y-4">
              <NearbyStats />
            </div>
          </div>
        </section>

        {/* Блок про систему проверки Черной пятницы */}
        <section className="py-8 bg-warning/5">
          <div className="container mx-auto px-4">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="flex items-center gap-4 mb-4">
                  <div className="badge badge-warning badge-lg">
                    AI ПРОВЕРКА
                  </div>
                  <h3 className="text-2xl font-bold">
                    {t('howBlackFridayWorks')}
                  </h3>
                </div>
                <div className="grid md:grid-cols-4 gap-4">
                  <div className="text-center">
                    <div className="text-3xl mb-2">📊</div>
                    <h4 className="font-bold mb-1">{t('priceHistory')}</h4>
                    <p className="text-sm text-base-content/60">
                      {t('priceHistoryDesc')}
                    </p>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl mb-2">🤖</div>
                    <h4 className="font-bold mb-1">{t('aiAnalysis')}</h4>
                    <p className="text-sm text-base-content/60">
                      {t('aiAnalysisDesc')}
                    </p>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl mb-2">✅</div>
                    <h4 className="font-bold mb-1">{t('minimum25')}</h4>
                    <p className="text-sm text-base-content/60">
                      {t('minimum25Desc')}
                    </p>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl mb-2">🏆</div>
                    <h4 className="font-bold mb-1">{t('qualityBadge')}</h4>
                    <p className="text-sm text-base-content/60">
                      {t('qualityBadgeDesc')}
                    </p>
                  </div>
                </div>
                <div className="alert alert-info mt-4">
                  <FiShield className="w-5 h-5" />
                  <span>
                    <strong>{t('buyerProtectionNote')}</strong>{' '}
                    {t('buyerProtectionNoteDesc')}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Официальные магазины */}
        <section className="container mx-auto px-4 py-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <BsGem className="w-6 h-6 text-secondary" />
              {t('officialStores')}
            </h2>
            <Link href="/stores" className="btn btn-sm btn-ghost">
              {t('allStores')}
            </Link>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
            {officialStores.map((store) => (
              <Link
                key={store.id}
                href={`/${locale}/b2c/${store.slug || store.id}`}
                className="card bg-base-100 hover:shadow-xl transition-all overflow-hidden"
              >
                {/* Фоновое изображение магазина */}
                <div
                  className="h-24 relative"
                  style={{
                    backgroundImage: `linear-gradient(rgba(0,0,0,0.4), rgba(0,0,0,0.4)), url('${store.bgImage}')`,
                    backgroundSize: 'cover',
                    backgroundPosition: 'center',
                  }}
                >
                  {store.blackFriday && (
                    <div className="badge badge-warning absolute top-2 left-2">
                      ✅ Черная пятница
                    </div>
                  )}
                  {store.discount && (
                    <div className="badge badge-error absolute top-2 right-2">
                      {store.discount}
                    </div>
                  )}
                </div>

                <div className="card-body">
                  <div className="flex items-start justify-between -mt-8">
                    <div className="flex items-center gap-3">
                      <div className="avatar">
                        <div className="w-16 h-16 rounded-full ring ring-base-100 ring-offset-base-100 ring-offset-2 relative overflow-hidden">
                          <SafeImage
                            src={buildImageUrl(store.logo)}
                            alt={store.name}
                            fill
                            className="object-cover"
                            sizes="64px"
                          />
                        </div>
                      </div>
                      <div className="mt-8">
                        <h3 className="font-bold flex items-center gap-1">
                          {store.name}
                          {store.verified && (
                            <FiShield className="w-4 h-4 text-success" />
                          )}
                        </h3>
                        <p className="text-sm text-base-content/60">
                          {store.category}
                        </p>
                      </div>
                    </div>
                  </div>

                  {store.realDiscount && (
                    <div className="text-xs text-success font-medium mt-2">
                      {store.realDiscount}
                    </div>
                  )}

                  <div className="flex justify-between text-sm mt-4">
                    <div className="text-center">
                      <p className="text-base-content/60">{t('followers')}</p>
                      <p className="font-bold">{store.followers}</p>
                    </div>
                    <div className="text-center">
                      <p className="text-base-content/60">{t('products')}</p>
                      <p className="font-bold">{store.products}</p>
                    </div>
                    <div className="text-center">
                      <p className="text-base-content/60">{t('rating')}</p>
                      <p className="font-bold flex items-center gap-1">
                        {store.rating ? (
                          <>
                            <FiStar className="w-3 h-3 fill-warning text-warning" />
                            {store.rating.toFixed(1)}
                          </>
                        ) : (
                          '-'
                        )}
                      </p>
                    </div>
                  </div>

                  {store.viewsCount !== undefined && store.viewsCount > 0 && (
                    <div className="flex items-center justify-center gap-1 text-xs text-base-content/60 mt-3 pt-3 border-t border-base-200">
                      <AiOutlineEye className="w-3.5 h-3.5" />
                      <span>
                        {store.viewsCount.toLocaleString()} {t('views')}
                      </span>
                    </div>
                  )}

                  <div className="btn btn-primary btn-sm mt-4 w-full">
                    {t('goToStore')}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </section>

        {/* Рекомендации на основе просмотров */}
        <section className="py-8 overflow-hidden">
          <div className="container mx-auto px-4">
            <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
              <AiOutlineEye className="w-6 h-6 text-info" />
              {t('recommendedTitle')}
            </h2>

            <div className="carousel carousel-center w-full space-x-4 pb-4 overflow-x-auto">
              {listings.map((deal, idx) => (
                <div key={`rec-${idx}`} className="carousel-item">
                  <Link href={getListingUrl(deal)} className="block">
                    <div className="card bg-base-100 w-64 hover:shadow-xl transition-all flex-shrink-0">
                      <figure className="h-40 overflow-hidden relative">
                        <SafeImage
                          src={deal.image}
                          alt={deal.translations?.[locale]?.title || deal.title}
                          fill
                          className="object-cover hover:scale-110 transition-transform duration-300"
                          sizes="280px"
                        />
                        {/* Кнопка избранного */}
                        <button
                          className="btn btn-circle btn-sm absolute top-2 right-2 bg-base-100/80 hover:bg-base-100"
                          onClick={(e) => {
                            // Используем listing_id если он есть
                            let listingId: number | string;
                            if (deal.listing_id) {
                              listingId =
                                typeof deal.listing_id === 'number'
                                  ? deal.listing_id
                                  : parseInt(String(deal.listing_id));
                            } else if (
                              typeof deal.id === 'string' &&
                              deal.id.startsWith('ml_')
                            ) {
                              listingId = parseInt(deal.id.replace('ml_', ''));
                            } else if (
                              typeof deal.id === 'string' &&
                              deal.id.startsWith('sp_')
                            ) {
                              // Сохраняем полный ID с префиксом для storefront продуктов
                              listingId = deal.id;
                            } else {
                              listingId = parseInt(String(deal.id));
                            }

                            // Проверяем валидность ID
                            const isValidId =
                              (typeof listingId === 'string' &&
                                listingId.startsWith('sp_')) ||
                              (typeof listingId === 'number' &&
                                !isNaN(listingId));

                            if (isValidId) {
                              handleToggleFavorite(listingId, e);
                            } else {
                              console.error(
                                'Invalid listing ID:',
                                deal.id,
                                deal.listing_id
                              );
                            }
                          }}
                          aria-label={
                            deal.isFavorite
                              ? _tCommon('product.removeFromFavorites')
                              : _tCommon('product.addToFavorites')
                          }
                        >
                          <FiHeart
                            className={`w-4 h-4 ${deal.isFavorite ? 'fill-error text-error' : ''}`}
                            aria-hidden="true"
                          />
                        </button>
                      </figure>
                      <div className="card-body p-4">
                        <h3 className="font-medium text-sm line-clamp-2 mb-2">
                          {deal.translations?.[locale]?.title || deal.title}
                        </h3>

                        {/* Цена и скидка */}
                        <div className="flex items-center gap-2 mb-2">
                          {deal.oldPrice && (
                            <span className="text-sm text-base-content/40 line-through">
                              {deal.oldPrice}
                            </span>
                          )}
                          <p className="text-lg font-bold text-primary">
                            {deal.price}
                          </p>
                          {deal.discount && (
                            <div className="badge badge-error badge-sm">
                              {deal.discount}
                            </div>
                          )}
                        </div>

                        {/* Рейтинг и просмотры */}
                        <div className="flex items-center gap-3 text-xs">
                          {deal.rating !== null && deal.rating > 0 && (
                            <div className="flex items-center gap-1">
                              <FiStar className="w-3 h-3 fill-warning text-warning" />
                              <span className="font-medium">
                                {deal.rating.toFixed(1)}
                              </span>
                              {deal.reviews !== null && deal.reviews > 0 && (
                                <span className="text-base-content/60">
                                  ({deal.reviews})
                                </span>
                              )}
                            </div>
                          )}
                          {deal.viewsCount !== null &&
                            deal.viewsCount !== undefined && (
                              <div className="flex items-center gap-1 text-base-content/60">
                                <AiOutlineEye className="w-3 h-3" />
                                <span>
                                  {deal.viewsCount === 0
                                    ? '-'
                                    : deal.viewsCount.toLocaleString()}
                                </span>
                              </div>
                            )}
                        </div>
                      </div>
                    </div>
                  </Link>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Трендовые поиски */}
        <section className="py-8 bg-base-200/50">
          <div className="container mx-auto px-4">
            <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
              <FiTrendingUp className="w-6 h-6 text-success" />
              {t('trendingSearches')}
            </h2>
            <div className="flex flex-wrap gap-2">
              {trendingSearches.map((search) => (
                <button
                  key={search}
                  className="btn btn-sm btn-outline hover:btn-primary"
                >
                  {search}
                </button>
              ))}
            </div>
          </div>
        </section>

        {/* Преимущества */}
        <section className="py-12">
          <div className="container mx-auto px-4">
            <h2 className="text-2xl font-bold mb-8 text-center">
              {t('whyChooseUs')}
            </h2>
            <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 bg-primary/10 rounded-full flex items-center justify-center">
                  <FiTruck className="w-8 h-8 text-primary" />
                </div>
                <h3 className="font-bold mb-2">{t('fastDelivery')}</h3>
                <p className="text-sm text-base-content/60">
                  {t('fastDeliveryDesc')}
                </p>
              </div>
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 bg-success/10 rounded-full flex items-center justify-center">
                  <FiShield className="w-8 h-8 text-success" />
                </div>
                <h3 className="font-bold mb-2">{t('dealProtection')}</h3>
                <p className="text-sm text-base-content/60">
                  {t('dealProtectionDesc')}
                </p>
              </div>
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 bg-warning/10 rounded-full flex items-center justify-center">
                  <FiCreditCard className="w-8 h-8 text-warning" />
                </div>
                <h3 className="font-bold mb-2">{t('convenientPayment')}</h3>
                <p className="text-sm text-base-content/60">
                  {t('convenientPaymentDesc')}
                </p>
              </div>
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 bg-info/10 rounded-full flex items-center justify-center">
                  <FiMessageCircle className="w-8 h-8 text-info" />
                </div>
                <h3 className="font-bold mb-2">{t('support247')}</h3>
                <p className="text-sm text-base-content/60">
                  {t('support247Desc')}
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* CTA секция */}
        <section className="py-12 bg-gradient-to-r from-primary to-secondary">
          <div className="container mx-auto px-4 text-center text-white">
            <h2 className="text-3xl font-bold mb-4">{t('startNowTitle')}</h2>
            <p className="text-xl mb-8 opacity-90">{t('startNowSubtitle')}</p>
            <div className="flex gap-4 justify-center">
              <button className="btn btn-white btn-lg">
                {t('createAccount')}
              </button>
              <button className="btn btn-outline btn-white btn-lg">
                {t('postListing')}
              </button>
            </div>
          </div>
        </section>

        {/* Футер */}
        <footer className="bg-base-200">
          <div className="container mx-auto px-4 py-12">
            <div className="grid md:grid-cols-2 lg:grid-cols-5 gap-8">
              {/* О компании */}
              <div className="lg:col-span-2">
                <h3 className="text-2xl font-bold mb-4">
                  {tFooter('company')}
                </h3>
                <p className="text-base-content/60 mb-4">
                  {tFooter('companyDescription')}
                </p>
                <div className="flex gap-4">
                  <button className="btn btn-primary">
                    <BsPhone className="w-4 h-4 mr-2" />
                    {tFooter('appStore')}
                  </button>
                  <button className="btn btn-primary">
                    <BsPhone className="w-4 h-4 mr-2" />
                    {tFooter('googlePlay')}
                  </button>
                </div>
              </div>

              {/* Покупателям */}
              <div>
                <h4 className="font-bold mb-4">{tFooter('buyers')}</h4>
                <ul className="space-y-2 text-sm">
                  <li>
                    2.{' '}
                    <Link href="/how-to-buy" className="hover:text-primary">
                      {tFooter('howToBuy')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/delivery" className="hover:text-primary">
                      {tFooter('delivery')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/payment" className="hover:text-primary">
                      {tFooter('payment')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/return" className="hover:text-primary">
                      {tFooter('return')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/warranty" className="hover:text-primary">
                      {tFooter('warranty')}
                    </Link>
                  </li>
                </ul>
              </div>

              {/* Продавцам */}
              <div>
                <h4 className="font-bold mb-4">{tFooter('sellers')}</h4>
                <ul className="space-y-2 text-sm">
                  <li>
                    <Link href="/how-to-sell" className="hover:text-primary">
                      {tFooter('howToSell')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/tariffs" className="hover:text-primary">
                      {tFooter('tariffs')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/promotion" className="hover:text-primary">
                      {tFooter('promotion')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/stores" className="hover:text-primary">
                      {tFooter('stores')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/api" className="hover:text-primary">
                      {tFooter('api')}
                    </Link>
                  </li>
                </ul>
              </div>

              {/* Помощь */}
              <div>
                <h4 className="font-bold mb-4">{tFooter('help')}</h4>
                <ul className="space-y-2 text-sm">
                  <li>
                    <Link href="/faq" className="hover:text-primary">
                      {tFooter('frequentQuestions')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/support" className="hover:text-primary">
                      {tFooter('support')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/rules" className="hover:text-primary">
                      {tFooter('rules')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/blog" className="hover:text-primary">
                      {tFooter('blog')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/contacts" className="hover:text-primary">
                      {tFooter('contacts')}
                    </Link>
                  </li>
                </ul>
              </div>
            </div>

            <div className="divider my-8"></div>

            <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-base-content/60">
              <p>
                {tFooter('copyright')} • v{process.env.NEXT_PUBLIC_APP_VERSION}
              </p>
              <div className="flex gap-4">
                <Link href="/terms" className="hover:text-primary">
                  {tFooter('termsOfUse')}
                </Link>
                <Link href="/privacy" className="hover:text-primary">
                  {tFooter('confidentiality')}
                </Link>
                <Link href="/cookies" className="hover:text-primary">
                  {tFooter('cookie')}
                </Link>
              </div>
            </div>
          </div>
        </footer>

        {/* Плавающая кнопка создания объявления */}
        <Link
          href="/create-listing-choice"
          className="fixed bottom-6 right-6 btn btn-primary btn-circle btn-lg shadow-xl hover:shadow-2xl hover:scale-110 transition-all duration-200 z-50"
          title={createListingText}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
        </Link>
      </div>
    </PageTransition>
  );
}
