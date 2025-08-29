'use client';

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Link from 'next/link';
// import Image from 'next/image';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/navigation';
import { PageTransition } from '@/components/ui/PageTransition';
import { useAuth } from '@/contexts/AuthContext';
import api from '@/services/api';
import CartIcon from '@/components/cart/CartIcon';
import { AuthButton } from '@/components/AuthButton';
// import { NestedCategorySelector } from '@/components/search/NestedCategorySelector';
import { useTranslations } from 'next-intl';
import configManager from '@/config';

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
          <p className="mt-2">Загрузка карты...</p>
        </div>
      </div>
    ),
  }
);

import {
  FiSearch,
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
import {
  BsHouseDoor,
  BsLaptop,
  BsBriefcase,
  BsPalette,
  BsTools,
  BsPhone,
  BsGem,
  BsHandbag,
} from 'react-icons/bs';
import { FaCar, FaTshirt } from 'react-icons/fa';
import { AiOutlineEye } from 'react-icons/ai';
import { HiOutlineSparkles } from 'react-icons/hi';
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
  const { user } = useAuth();
  const t = useTranslations('marketplace.home');
  const tCommon = useTranslations('common');
  const tFooter = useTranslations('common.footer');
  const [_mounted, setMounted] = useState(false);
  const [_selectedCategory] = useState<string | number>('all');
  const [currentBanner, setCurrentBanner] = useState(0);
  const [_showMobileMenu, _setShowMobileMenu] = useState(false);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [userLocation] = useState([44.7866, 20.4489]); // Координаты Белграда
  const [listings, setListings] = useState<any[]>([]);
  const [isLoadingListings, setIsLoadingListings] = useState(true);
  const [_categories, setCategories] = useState<any[]>([]);
  const [popularCategories, setPopularCategories] = useState<any[]>([]);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [officialStores, setOfficialStores] = useState<any[]>([]);
  const [_isLoadingStores, setIsLoadingStores] = useState(false);

  // Функция для получения URL объявления
  const getListingUrl = (deal: any) => {
    console.log('getListingUrl called with deal:', {
      id: deal.id,
      product_id: deal.product_id,
      listing_id: deal.listing_id,
      isStorefront: deal.isStorefront,
    });

    if (deal.isStorefront && deal.product_id) {
      // Для товаров витрин - используем product_id без префикса
      const url = `/${locale}/marketplace/${deal.product_id}`;
      console.log('Storefront URL:', url);
      return url;
    } else if (deal.listing_id) {
      // Для обычных объявлений - используем listing_id
      const url = `/${locale}/marketplace/${deal.listing_id}`;
      console.log('Listing URL:', url);
      return url;
    } else {
      // Fallback - извлекаем чистый ID из deal.id убрав префиксы
      const cleanId =
        typeof deal.id === 'string'
          ? deal.id.replace(/^(ml_|sp_)/, '')
          : deal.id;
      const url = `/${locale}/marketplace/${cleanId}`;
      console.log('Fallback URL:', url);
      return url;
    }
  };

  // Функция для открытия чата
  const handleStartChat = (deal: any) => {
    console.log('handleStartChat called with deal:', deal);

    if (!user) {
      // Если пользователь не авторизован, перенаправляем на страницу входа
      router.push('/login');
      return;
    }

    // Определяем URL для чата в зависимости от типа объявления
    if (deal.isStorefront && deal.storefront_id) {
      // B2C - чат с витриной, передаем storefront_product_id и seller_id (владелец витрины)
      console.log(
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
      console.log(
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

  // Устанавливаем mounted после гидрации для предотвращения hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

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

  // Загрузка категорий и популярных категорий
  useEffect(() => {
    const loadCategories = async () => {
      try {
        // Загружаем обычные категории для выпадающего списка
        const [categoriesResponse, popularResponse] = await Promise.all([
          api.get('/api/v1/marketplace/categories'),
          api.get(
            `/api/v1/marketplace/popular-categories?lang=${locale}&limit=8`
          ),
        ]);

        if (categoriesResponse.data.success) {
          setCategories(categoriesResponse.data.data);
        }

        if (popularResponse.data.success && popularResponse.data.data) {
          // Добавляем иконки для популярных категорий на основе их slug
          const iconMap: { [key: string]: any } = {
            'real-estate': BsHouseDoor,
            automotive: FaCar,
            electronics: BsLaptop,
            fashion: FaTshirt,
            jobs: BsBriefcase,
            services: BsTools,
            'hobbies-entertainment': BsPalette,
            'home-garden': BsHandbag,
            industrial: BsTools,
            'food-beverages': BsPhone,
            'books-stationery': BsGem,
            'antiques-art': BsPalette,
          };

          const colorMap: { [key: string]: string } = {
            'real-estate': 'text-blue-600',
            automotive: 'text-red-600',
            electronics: 'text-purple-600',
            fashion: 'text-pink-600',
            jobs: 'text-green-600',
            services: 'text-orange-600',
            'hobbies-entertainment': 'text-indigo-600',
            'home-garden': 'text-yellow-600',
            industrial: 'text-gray-600',
            'food-beverages': 'text-teal-600',
            'books-stationery': 'text-cyan-600',
            'antiques-art': 'text-rose-600',
          };

          const categoriesWithIcons = popularResponse.data.data.map(
            (cat: any) => ({
              ...cat,
              icon: iconMap[cat.slug] || BsHandbag,
              color: colorMap[cat.slug] || 'text-gray-600',
              count: cat.count ? `${cat.count}+` : '0',
            })
          );

          setPopularCategories(categoriesWithIcons);
          console.log('Popular categories loaded:', categoriesWithIcons);
        }
      } catch (error) {
        console.error('Failed to load categories:', error);
      } finally {
        setIsLoadingCategories(false);
      }
    };
    loadCategories();
  }, [locale]);

  // Загрузка витрин (официальных магазинов)
  useEffect(() => {
    const loadStorefronts = async () => {
      setIsLoadingStores(true);
      try {
        // Загружаем активные витрины
        const response = await api.get('/api/v1/storefronts', {
          params: {
            is_active: true,
            limit: 4,
            sort_by: 'products_count',
            sort_order: 'desc',
          },
        });

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
                category: store.category_name || 'Магазин',
                logo:
                  store.logo_url ||
                  `https://ui-avatars.com/api/?name=${initials}&background=${bgColor}&color=fff&size=128`,
                followers: store.followers_count
                  ? `${Math.floor(store.followers_count / 1000)}K`
                  : '0',
                products: store.products_count || 0,
                rating: store.rating || 0,
                verified: store.is_verified || false,
                discount: store.discount_text || '',
                bgImage: store.banner_url || bgImage,
                slug: store.slug,
                description: store.description,
              };
            }
          );

          setOfficialStores(formattedStores);
          console.log('Loaded storefronts:', formattedStores);
        } else {
          // Если нет реальных витрин, используем заглушки
          setOfficialStores([
            {
              id: 1,
              name: 'Агентство недвижимости',
              category: 'Недвижимость',
              logo: '/listings/storefronts/1/logo/10_2.jpeg',
              followers: '2K',
              products: 38,
              rating: 4.5,
              verified: true,
              discount: '',
              bgImage:
                'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400&h=200&fit=crop',
              slug: 'agenstvo',
              description:
                'Тут мы раскидаем по карте квартиры и будем их продавать',
            },
          ]);
        }
      } catch (error) {
        console.error('Failed to load storefronts:', error);
        // В случае ошибки тоже используем одну витрину из БД как заглушку
        setOfficialStores([
          {
            id: 1,
            name: 'Агентство недвижимости',
            category: 'Недвижимость',
            logo: '/listings/storefronts/1/logo/10_2.jpeg',
            followers: '2K',
            products: 38,
            rating: 4.5,
            verified: true,
            discount: '',
            bgImage:
              'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400&h=200&fit=crop',
            slug: 'agenstvo',
            description:
              'Тут мы раскидаем по карте квартиры и будем их продавать',
          },
        ]);
      } finally {
        setIsLoadingStores(false);
      }
    };

    loadStorefronts();
  }, []);

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
        searchParams.append('language', locale);
        searchParams.append('status', 'active');
        searchParams.append('product_types[]', 'marketplace');
        searchParams.append('product_types[]', 'storefront');

        console.log(
          'Request URL:',
          `/api/v1/search?${searchParams.toString()}`
        );
        const response = await api.get(
          `/api/v1/search?${searchParams.toString()}`
        );
        console.log('API Response:', response.data);

        if (
          response.data &&
          response.data.items &&
          response.data.items.length > 0
        ) {
          // Разделяем объявления на C2C и B2C для смешанного показа
          const allListings = response.data.items;
          console.log(
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

          // Создаем смешанную выборку: преимущественно C2C, но включаем B2C если есть
          let selectedListings = [];

          // Берем первые 6 C2C объявлений
          selectedListings.push(...c2cListings.slice(0, 6));

          // Добавляем 2 B2C объявления если есть
          if (b2cListings.length > 0) {
            selectedListings.push(...b2cListings.slice(0, 2));
          } else {
            // Если B2C нет, добавляем еще 2 C2C
            selectedListings.push(...c2cListings.slice(6, 8));
          }

          // Ограничиваем до 8 объявлений
          selectedListings = selectedListings.slice(0, 8);

          console.log(
            `Mixed selection: ${selectedListings.filter((l) => !l.storefrontId).length} C2C + ${selectedListings.filter((l) => l.storefrontId).length} B2C`
          );

          const apiListings = selectedListings.map((listing: any) => {
            // Вычисляем скидку если есть старая цена
            let discount = null;
            let oldPrice = null;

            if (
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

            const mappedListing = {
              id:
                listing.product_type === 'storefront'
                  ? `sp_${listing.product_id}` // Добавляем префикс для уникальности
                  : `ml_${listing.id}`,
              product_id:
                listing.product_type === 'storefront'
                  ? listing.product_id
                  : null,
              title: listing.name || listing.title,
              price: `${listing.price} ${listing.currency || 'РСД'}`,
              oldPrice,
              discount,
              location:
                listing.address_city ||
                listing.city ||
                listing.location?.city ||
                'Сербия',
              image:
                listing.images && listing.images.length > 0
                  ? configManager.buildImageUrl(
                      listing.images[0].url || listing.images[0].public_url
                    )
                  : 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=400&h=300&fit=crop', // fallback изображение
              rating: listing.rating || 4.0 + Math.random() * 1.0, // Используем настоящий рейтинг или генерируем
              reviews:
                listing.reviewCount || Math.floor(Math.random() * 500) + 10,
              isNew:
                new Date(listing.created_at || listing.createdAt) >
                new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // Новое если создано за последнюю неделю
              isPremium: listing.isPremium || false,
              isFavorite: false, // Это нужно будет получать из профиля пользователя
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
          });

          setListings(apiListings);
          console.log(
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
  }, [locale]);

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
            <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
              <HiOutlineSparkles className="w-6 h-6 text-warning" />
              {t('popularCategories')}
            </h2>
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
                  const Icon = cat.icon;
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
              >
                <FiGrid className="w-4 h-4" />
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`btn btn-sm ${viewMode === 'list' ? 'btn-primary' : 'btn-ghost'}`}
              >
                <FiList className="w-4 h-4" />
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
                      <img
                        src={deal.image}
                        alt={deal.title}
                        className="w-full h-full object-cover"
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
                          e.preventDefault();
                          e.stopPropagation();
                          console.log('Add to favorites:', deal.id);
                        }}
                      >
                        <FiHeart
                          className={`w-4 h-4 ${deal.isFavorite ? 'fill-error text-error' : ''}`}
                        />
                      </button>
                    </figure>
                    <div className="card-body p-4">
                      <h3 className="card-title text-base line-clamp-2">
                        {deal.title}
                      </h3>
                      <div className="flex items-center gap-2 text-sm">
                        <FiMapPin className="w-3 h-3" />
                        <span className="text-base-content/60">
                          {deal.location}
                        </span>
                      </div>
                      {deal.rating && (
                        <div className="flex items-center gap-1 text-sm">
                          <FiStar className="w-3 h-3 fill-warning text-warning" />
                          <span>{deal.rating}</span>
                          <span className="text-base-content/60">
                            ({deal.reviews})
                          </span>
                        </div>
                      )}
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
                        // B2C (витрина) - кнопка "В корзину" + "Написать в чат"
                        <div className="flex gap-2 mt-2">
                          <button
                            className="btn btn-primary btn-sm flex-1"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              console.log(
                                'Add to cart:',
                                deal.product_id || deal.id
                              );
                            }}
                          >
                            {t('addToCart')}
                          </button>
                          <button
                            className="btn btn-outline btn-sm"
                            onClick={(e) => {
                              e.preventDefault();
                              e.stopPropagation();
                              handleStartChat(deal);
                            }}
                          >
                            <FiMessageCircle className="w-4 h-4" />
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
                              e.preventDefault();
                              e.stopPropagation();
                              console.log('Add to favorites:', deal.id);
                            }}
                          >
                            <FiHeart className="w-4 h-4" />
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
                      title: item.title,
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
              <div
                key={store.id}
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
                        <div className="w-16 rounded-full ring ring-base-100 ring-offset-base-100 ring-offset-2">
                          <img src={store.logo} alt={store.name} />
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
                        <FiStar className="w-3 h-3 fill-warning text-warning" />
                        {store.rating}
                      </p>
                    </div>
                  </div>

                  <button className="btn btn-primary btn-sm mt-4 w-full">
                    {t('goToStore')}
                  </button>
                </div>
              </div>
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
                  <div className="card bg-base-100 w-64 hover:shadow-xl transition-all flex-shrink-0">
                    <figure className="h-40 overflow-hidden">
                      <img
                        src={deal.image}
                        alt={deal.title}
                        className="h-full w-full object-cover hover:scale-110 transition-transform duration-300"
                      />
                    </figure>
                    <div className="card-body p-4">
                      <h3 className="font-medium text-sm line-clamp-2">
                        {deal.title}
                      </h3>
                      <div className="flex items-center gap-2">
                        {deal.oldPrice && (
                          <span className="text-sm text-base-content/40 line-through">
                            {deal.oldPrice}
                          </span>
                        )}
                        <p className="text-lg font-bold text-primary">
                          {deal.price}
                        </p>
                      </div>
                      {deal.discount && (
                        <div className="badge badge-error badge-sm">
                          {deal.discount}
                        </div>
                      )}
                    </div>
                  </div>
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
              <p>{tFooter('copyright')}</p>
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

        {/* Мобильная навигация */}
        <div className="btm-nav lg:hidden">
          <button className="text-primary">
            <FiSearch className="w-5 h-5" />
            <span className="btm-nav-label">{t('search')}</span>
          </button>
          <button>
            <FiHeart className="w-5 h-5" />
            <span className="btm-nav-label">{t('favorites')}</span>
          </button>
          <div className="text-secondary">
            <CartIcon />
            <span className="btm-nav-label">{t('cart')}</span>
          </div>
          <div className="flex flex-col items-center justify-center">
            <AuthButton />
            <span className="btm-nav-label text-xs">{tCommon('profile')}</span>
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
