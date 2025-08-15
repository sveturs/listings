'use client';

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Link from 'next/link';
import dynamic from 'next/dynamic';

// Динамический импорт карты для избежания SSR проблем
const MapSection = dynamic(() => import('./components/MapSection'), {
  ssr: false,
  loading: () => (
    <div className="h-full w-full flex items-center justify-center bg-base-200 rounded-lg">
      <div className="text-center">
        <div className="loading loading-spinner loading-lg text-primary"></div>
        <p className="mt-2">Загрузка карты...</p>
      </div>
    </div>
  ),
});
import {
  FiSearch,
  FiMapPin,
  FiUser,
  FiShoppingCart,
  FiMenu,
  FiX,
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
import { AiOutlineThunderbolt, AiOutlineEye } from 'react-icons/ai';
import { HiOutlineSparkles } from 'react-icons/hi';

export default function IdealMarketplacePage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [currentBanner, setCurrentBanner] = useState(0);
  const [showMobileMenu, setShowMobileMenu] = useState(false);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [_userLocation, _setUserLocation] = useState('Белград');
  const [_cartCount, _setCartCount] = useState(3);

  // Баннеры для hero секции
  const banners = [
    {
      id: 1,
      title: '✅ Проверенная Черная пятница',
      subtitle: 'Только реальные скидки от 25%! Проверено историей цен',
      bgColor: 'bg-gradient-to-r from-purple-600 to-pink-600',
      cta: 'Смотреть акции',
      image: '🛍️',
      badge: 'AI проверка',
      details: '> 5% товаров со скидкой 25%+',
    },
    {
      id: 2,
      title: '🚚 Бесплатная доставка',
      subtitle: 'При покупке от €50',
      bgColor: 'bg-gradient-to-r from-blue-600 to-cyan-600',
      cta: 'Узнать больше',
      image: '📦',
    },
    {
      id: 3,
      title: '🛡️ Защита покупателя',
      subtitle: 'Безопасные сделки с эскроу',
      bgColor: 'bg-gradient-to-r from-green-600 to-teal-600',
      cta: 'Как работает',
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

  // Категории с иконками и количеством
  const categories = [
    {
      id: 'realestate',
      name: 'Недвижимость',
      icon: BsHouseDoor,
      count: '45K+',
      color: 'text-blue-600',
    },
    {
      id: 'auto',
      name: 'Транспорт',
      icon: FaCar,
      count: '28K+',
      color: 'text-red-600',
    },
    {
      id: 'electronics',
      name: 'Электроника',
      icon: BsLaptop,
      count: '67K+',
      color: 'text-purple-600',
    },
    {
      id: 'fashion',
      name: 'Одежда',
      icon: FaTshirt,
      count: '89K+',
      color: 'text-pink-600',
    },
    {
      id: 'job',
      name: 'Работа',
      icon: BsBriefcase,
      count: '12K+',
      color: 'text-green-600',
    },
    {
      id: 'services',
      name: 'Услуги',
      icon: BsTools,
      count: '35K+',
      color: 'text-orange-600',
    },
    {
      id: 'hobby',
      name: 'Хобби',
      icon: BsPalette,
      count: '23K+',
      color: 'text-indigo-600',
    },
    {
      id: 'home',
      name: 'Для дома',
      icon: BsHandbag,
      count: '54K+',
      color: 'text-yellow-600',
    },
  ];

  // Горячие предложения с реальными изображениями
  const hotDeals = [
    {
      id: 1,
      title: 'iPhone 15 Pro Max 256GB',
      price: '€1099',
      oldPrice: '€1399',
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
      id: 2,
      title: '2-комн квартира, центр, 65м²',
      price: '€85000',
      location: 'Нови Сад',
      image:
        'https://images.unsplash.com/photo-1512917774080-9991f1c4c750?w=400&h=300&fit=crop',
      rating: 4.9,
      reviews: 12,
      isNew: false,
      isPremium: true,
      isFavorite: true,
    },
    {
      id: 3,
      title: 'MacBook Air M3 13" 512GB',
      price: '€1299',
      oldPrice: '€1599',
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
    {
      id: 4,
      title: 'BMW X5 2021 xDrive30d',
      price: '€52900',
      location: 'Белград',
      image:
        'https://images.unsplash.com/photo-1555215858-9db736e8a7b8?w=400&h=300&fit=crop',
      rating: 5.0,
      reviews: 8,
      isNew: false,
      isPremium: true,
      isFavorite: true,
    },
    {
      id: 5,
      title: 'PlayStation 5 с играми',
      price: '€549',
      oldPrice: '€699',
      discount: '-21%',
      location: 'Белград',
      image:
        'https://images.unsplash.com/photo-1606813907291-d86efa9b94db?w=400&h=300&fit=crop',
      rating: 4.9,
      reviews: 445,
      isNew: false,
      isPremium: false,
      isFavorite: false,
    },
    {
      id: 6,
      title: 'Диван угловой, кожа',
      price: '€899',
      location: 'Нови Сад',
      image:
        'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400&h=300&fit=crop',
      rating: 4.7,
      reviews: 89,
      isNew: true,
      isPremium: false,
      isFavorite: true,
    },
    {
      id: 7,
      title: 'Nike Air Max 2024',
      price: '€149',
      oldPrice: '€199',
      discount: '-25%',
      location: 'Белград',
      image:
        'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400&h=300&fit=crop',
      rating: 4.8,
      reviews: 1234,
      isNew: true,
      isPremium: false,
      isFavorite: false,
    },
    {
      id: 8,
      title: 'Электросамокат Xiaomi Pro 2',
      price: '€449',
      location: 'Белград',
      image:
        'https://images.unsplash.com/photo-1593941966874-e9ec34e67d0e?w=400&h=300&fit=crop',
      rating: 4.6,
      reviews: 567,
      isNew: false,
      isPremium: false,
      isFavorite: false,
    },
  ];

  // Официальные магазины с реальными логотипами
  const stores = [
    {
      id: 1,
      name: 'TechnoWorld',
      category: 'Электроника',
      logo: 'https://ui-avatars.com/api/?name=TW&background=6366f1&color=fff&size=128',
      followers: '125K',
      products: 892,
      rating: 4.9,
      verified: true,
      discount: 'до -70%',
      bgImage:
        'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=400&h=200&fit=crop',
      blackFriday: true,
      realDiscount: '31% реальная скидка',
    },
    {
      id: 2,
      name: 'FashionHub',
      category: 'Одежда и обувь',
      logo: 'https://ui-avatars.com/api/?name=FH&background=ec4899&color=fff&size=128',
      followers: '89K',
      products: 1234,
      rating: 4.8,
      verified: true,
      discount: 'до -50%',
      bgImage:
        'https://images.unsplash.com/photo-1490481651871-ab68de25d43d?w=400&h=200&fit=crop',
    },
    {
      id: 3,
      name: 'HomeDecor',
      category: 'Дом и сад',
      logo: 'https://ui-avatars.com/api/?name=HD&background=10b981&color=fff&size=128',
      followers: '67K',
      products: 456,
      rating: 4.7,
      verified: true,
      discount: 'до -40%',
      bgImage:
        'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=400&h=200&fit=crop',
    },
    {
      id: 4,
      name: 'AutoParts',
      category: 'Автозапчасти',
      logo: 'https://ui-avatars.com/api/?name=AP&background=ef4444&color=fff&size=128',
      followers: '45K',
      products: 789,
      rating: 4.8,
      verified: true,
      discount: 'до -30%',
      bgImage:
        'https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?w=400&h=200&fit=crop',
    },
  ];

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

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-100 to-base-200">
      {/* Шапка сайта */}
      <header className="sticky top-0 z-50 bg-base-100/95 backdrop-blur-md border-b border-base-300">
        {/* Верхняя панель */}
        <div className="bg-primary text-primary-content py-1 text-sm">
          <div className="container mx-auto px-4 flex justify-between items-center">
            <div className="flex items-center gap-4">
              <span className="flex items-center gap-1">
                <FiMapPin className="w-3 h-3" />
                {_userLocation}
              </span>
              <Link href="/map" className="hover:underline">
                Выбрать другой город
              </Link>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/business" className="hover:underline">
                Для бизнеса
              </Link>
              <Link href="/help" className="hover:underline">
                Помощь
              </Link>
              <Link href="/app" className="hover:underline">
                📱 Приложение
              </Link>
            </div>
          </div>
        </div>

        {/* Основная шапка */}
        <div className="container mx-auto px-4 py-3">
          <div className="flex items-center gap-4">
            {/* Логотип */}
            <Link href="/" className="flex items-center gap-2">
              <img
                src="/logos/svetu-gradient-48x48.png"
                alt="SveTu Logo"
                className="w-10 h-10"
              />
              <div className="text-3xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                SveTu
              </div>
            </Link>

            {/* Кнопка каталога */}
            <button className="btn btn-primary hidden lg:flex items-center gap-2">
              <FiMenu className="w-5 h-5" />
              Каталог
            </button>

            {/* Поисковая строка */}
            <div className="flex-1 max-w-3xl">
              <div className="flex">
                <select
                  className="select select-bordered rounded-r-none w-40 hidden md:block"
                  value={selectedCategory}
                  onChange={(e) => setSelectedCategory(e.target.value)}
                >
                  <option value="all">Все категории</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
                <div className="relative flex-1">
                  <input
                    type="text"
                    placeholder="Поиск среди 2 млн товаров..."
                    className="input input-bordered w-full rounded-none"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                  <button className="absolute right-2 top-1/2 -translate-y-1/2">
                    <FiSearch className="w-5 h-5 text-base-content/50" />
                  </button>
                </div>
                <button className="btn btn-primary rounded-l-none">
                  Найти
                </button>
              </div>
            </div>

            {/* Действия пользователя */}
            <div className="flex items-center gap-2">
              <button className="btn btn-ghost btn-circle relative">
                <FiHeart className="w-5 h-5" />
                <span className="badge badge-sm badge-error absolute -top-1 -right-1">
                  2
                </span>
              </button>
              <button className="btn btn-ghost btn-circle relative">
                <FiShoppingCart className="w-5 h-5" />
                {_cartCount > 0 && (
                  <span className="badge badge-sm badge-error absolute -top-1 -right-1">
                    {_cartCount}
                  </span>
                )}
              </button>
              <Link
                href="/create"
                className="btn btn-secondary hidden lg:inline-flex"
              >
                Подать объявление
              </Link>
              <button className="btn btn-ghost btn-circle lg:btn lg:btn-ghost lg:btn-wide">
                <FiUser className="w-5 h-5" />
                <span className="hidden lg:inline ml-2">Войти</span>
              </button>
            </div>

            {/* Мобильное меню */}
            <button
              className="btn btn-ghost btn-circle lg:hidden"
              onClick={() => setShowMobileMenu(!showMobileMenu)}
            >
              {showMobileMenu ? (
                <FiX className="w-5 h-5" />
              ) : (
                <FiMenu className="w-5 h-5" />
              )}
            </button>
          </div>
        </div>

        {/* Категории под поиском */}
        <div className="border-t border-base-300 py-2 hidden lg:block">
          <div className="container mx-auto px-4">
            <div className="flex items-center gap-6 text-sm">
              {categories.slice(0, 8).map((cat) => {
                const Icon = cat.icon;
                return (
                  <Link
                    key={cat.id}
                    href={`/category/${cat.id}`}
                    className="flex items-center gap-2 hover:text-primary transition-colors"
                  >
                    <Icon className={`w-4 h-4 ${cat.color}`} />
                    <span>{cat.name}</span>
                    <span className="text-base-content/50">({cat.count})</span>
                  </Link>
                );
              })}
              <Link href="/categories" className="text-primary font-medium">
                Все категории →
              </Link>
            </div>
          </div>
        </div>
      </header>

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
                  className={`relative rounded-2xl p-8 lg:p-12 text-white overflow-hidden ${banners[currentBanner].bgColor}`}
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
                          idx === currentBanner ? 'w-8 bg-white' : 'bg-white/50'
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
              <div className="card bg-gradient-to-br from-orange-500 to-red-500 text-white">
                <div className="card-body">
                  <h3 className="card-title text-white">⚡ Молния-скидки</h3>
                  <p>Успей купить со скидкой до 90%</p>
                  <div className="text-2xl font-bold">02:45:18</div>
                  <button className="btn btn-white btn-sm">Смотреть</button>
                </div>
              </div>
              <div className="card bg-gradient-to-br from-green-500 to-teal-500 text-white">
                <div className="card-body">
                  <h3 className="card-title text-white">🎁 Подарок новым</h3>
                  <p>Скидка €10 на первый заказ</p>
                  <button className="btn btn-white btn-sm">Получить</button>
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
            Популярные категории
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
            {categories.map((cat) => {
              const Icon = cat.icon;
              return (
                <Link
                  key={cat.id}
                  href={`/category/${cat.id}`}
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
        </div>
      </section>

      {/* Горячие предложения */}
      <section className="py-8 bg-base-200/50">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <AiOutlineThunderbolt className="w-6 h-6 text-error" />
              Горячие предложения
            </h2>
            <div className="flex items-center gap-2">
              <button
                className={`btn btn-sm ${viewMode === 'grid' ? 'btn-primary' : 'btn-ghost'}`}
                onClick={() => setViewMode('grid')}
              >
                <FiGrid className="w-4 h-4" />
              </button>
              <button
                className={`btn btn-sm ${viewMode === 'list' ? 'btn-primary' : 'btn-ghost'}`}
                onClick={() => setViewMode('list')}
              >
                <FiList className="w-4 h-4" />
              </button>
              <Link href="/hot" className="btn btn-sm btn-ghost">
                Все предложения →
              </Link>
            </div>
          </div>

          <div
            className={`grid ${viewMode === 'grid' ? 'grid-cols-2 lg:grid-cols-4' : 'grid-cols-1'} gap-4`}
          >
            {hotDeals.map((deal) => (
              <motion.div
                key={deal.id}
                whileHover={{ scale: 1.02 }}
                className="card bg-base-100 hover:shadow-xl transition-all"
              >
                <figure className="relative h-48 overflow-hidden">
                  <img
                    src={deal.image}
                    alt={deal.title}
                    className="w-full h-full object-cover hover:scale-110 transition-transform duration-300"
                  />
                  {deal.discount && (
                    <div className="badge badge-error absolute top-2 left-2">
                      {deal.discount}
                    </div>
                  )}
                  {deal.isNew && (
                    <div className="badge badge-success absolute top-2 left-2">
                      NEW
                    </div>
                  )}
                  {deal.isPremium && (
                    <div className="badge badge-warning absolute top-2 right-2">
                      PREMIUM
                    </div>
                  )}
                  <button className="btn btn-circle btn-sm absolute top-2 right-2 bg-base-100/80 hover:bg-base-100">
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
                    <span className="text-xl font-bold text-primary">
                      {deal.price}
                    </span>
                  </div>
                  <button className="btn btn-primary btn-sm mt-2">
                    В корзину
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Товары рядом с вами */}
      <section className="py-8">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <FiMapPin className="w-6 h-6 text-info" />
              Рядом с вами
            </h2>
            <Link href="/map" className="btn btn-sm btn-ghost">
              Открыть карту →
            </Link>
          </div>

          <div className="grid lg:grid-cols-3 gap-6">
            {/* Карта */}
            <div className="lg:col-span-2">
              <div className="card bg-base-100 overflow-hidden">
                <div className="card-body p-0">
                  <div className="h-96 relative">
                    <MapSection />
                    {/* Фильтры на карте */}
                    <div className="absolute top-4 left-4 right-4 flex gap-2 z-[1000]">
                      <button className="btn btn-sm bg-base-100 shadow-lg">
                        До €100
                      </button>
                      <button className="btn btn-sm bg-base-100 shadow-lg">
                        Сегодня
                      </button>
                      <button className="btn btn-sm bg-base-100 shadow-lg">
                        С фото
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Статистика */}
            <div className="space-y-4">
              <div className="stats stats-vertical shadow w-full">
                <div className="stat">
                  <div className="stat-title">В вашем районе</div>
                  <div className="stat-value text-primary">1,234</div>
                  <div className="stat-desc">объявлений</div>
                </div>
                <div className="stat">
                  <div className="stat-title">Новых сегодня</div>
                  <div className="stat-value text-success">+89</div>
                  <div className="stat-desc">↗︎ больше чем вчера (57)</div>
                </div>
                <div className="stat">
                  <div className="stat-title">Средняя цена</div>
                  <div className="stat-value text-info">€450</div>
                  <div className="stat-desc">в радиусе 5 км</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Блок про систему проверки Черной пятницы */}
      <section className="py-8 bg-warning/5">
        <div className="container mx-auto px-4">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <div className="flex items-center gap-4 mb-4">
                <div className="badge badge-warning badge-lg">AI ПРОВЕРКА</div>
                <h3 className="text-2xl font-bold">
                  Как работает проверка Черной пятницы
                </h3>
              </div>
              <div className="grid md:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-3xl mb-2">📊</div>
                  <h4 className="font-bold mb-1">История цен</h4>
                  <p className="text-sm text-base-content/60">
                    Отслеживаем цены 60 дней до акции
                  </p>
                </div>
                <div className="text-center">
                  <div className="text-3xl mb-2">🔍</div>
                  <h4 className="font-bold mb-1">AI анализ</h4>
                  <p className="text-sm text-base-content/60">
                    Проверяем реальность скидки алгоритмом
                  </p>
                </div>
                <div className="text-center">
                  <div className="text-3xl mb-2">✅</div>
                  <h4 className="font-bold mb-1">Минимум 25%</h4>
                  <p className="text-sm text-base-content/60">
                    Только скидки от 25% на более чем 5% товаров
                  </p>
                </div>
                <div className="text-center">
                  <div className="text-3xl mb-2">🏆</div>
                  <h4 className="font-bold mb-1">Значок качества</h4>
                  <p className="text-sm text-base-content/60">
                    Получают только честные продавцы
                  </p>
                </div>
              </div>
              <div className="alert alert-info mt-4">
                <FiShield className="w-5 h-5" />
                <span>
                  <strong>Защита покупателей:</strong> Магазины с поддельными
                  скидками автоматически исключаются из программы Черная пятница
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Официальные магазины */}
      <section className="py-8 bg-gradient-to-r from-primary/5 to-secondary/5">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold flex items-center gap-2">
              <BsGem className="w-6 h-6 text-secondary" />
              Официальные магазины
            </h2>
            <Link href="/stores" className="btn btn-sm btn-ghost">
              Все магазины →
            </Link>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
            {stores.map((store) => (
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
                      <p className="text-base-content/60">Подписчики</p>
                      <p className="font-bold">{store.followers}</p>
                    </div>
                    <div className="text-center">
                      <p className="text-base-content/60">Товаров</p>
                      <p className="font-bold">{store.products}</p>
                    </div>
                    <div className="text-center">
                      <p className="text-base-content/60">Рейтинг</p>
                      <p className="font-bold flex items-center gap-1">
                        <FiStar className="w-3 h-3 fill-warning text-warning" />
                        {store.rating}
                      </p>
                    </div>
                  </div>

                  <button className="btn btn-primary btn-sm mt-4 w-full">
                    Перейти в магазин
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Рекомендации на основе просмотров */}
      <section className="py-8">
        <div className="container mx-auto px-4">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <AiOutlineEye className="w-6 h-6 text-info" />
            Рекомендуем на основе ваших просмотров
          </h2>

          <div className="carousel carousel-center space-x-4 pb-4">
            {hotDeals.map((deal, idx) => (
              <div key={`rec-${idx}`} className="carousel-item">
                <div className="card bg-base-100 w-64 hover:shadow-xl transition-all">
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
            Что сейчас ищут
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
            Почему выбирают SveTu?
          </h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-primary/10 rounded-full flex items-center justify-center">
                <FiTruck className="w-8 h-8 text-primary" />
              </div>
              <h3 className="font-bold mb-2">Быстрая доставка</h3>
              <p className="text-sm text-base-content/60">
                Доставка по всей Сербии от 1 дня
              </p>
            </div>
            <div className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-success/10 rounded-full flex items-center justify-center">
                <FiShield className="w-8 h-8 text-success" />
              </div>
              <h3 className="font-bold mb-2">Защита сделок</h3>
              <p className="text-sm text-base-content/60">
                Безопасные платежи через эскроу
              </p>
            </div>
            <div className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-warning/10 rounded-full flex items-center justify-center">
                <FiCreditCard className="w-8 h-8 text-warning" />
              </div>
              <h3 className="font-bold mb-2">Удобная оплата</h3>
              <p className="text-sm text-base-content/60">
                Все способы оплаты включая рассрочку
              </p>
            </div>
            <div className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-info/10 rounded-full flex items-center justify-center">
                <FiMessageCircle className="w-8 h-8 text-info" />
              </div>
              <h3 className="font-bold mb-2">Поддержка 24/7</h3>
              <p className="text-sm text-base-content/60">
                Помощь на каждом этапе сделки
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA секция */}
      <section className="py-12 bg-gradient-to-r from-primary to-secondary">
        <div className="container mx-auto px-4 text-center text-white">
          <h2 className="text-3xl font-bold mb-4">
            Начните покупать и продавать прямо сейчас!
          </h2>
          <p className="text-xl mb-8 opacity-90">
            Присоединяйтесь к 2 миллионам пользователей
          </p>
          <div className="flex gap-4 justify-center">
            <button className="btn btn-white btn-lg">Создать аккаунт</button>
            <button className="btn btn-outline btn-white btn-lg">
              Подать объявление
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
              <h3 className="text-2xl font-bold mb-4">SveTu</h3>
              <p className="text-base-content/60 mb-4">
                Крупнейшая площадка для покупки и продажи товаров в Сербии.
                Безопасные сделки, быстрая доставка, лучшие цены.
              </p>
              <div className="flex gap-4">
                <button className="btn btn-primary">
                  <BsPhone className="w-4 h-4 mr-2" />
                  App Store
                </button>
                <button className="btn btn-primary">
                  <BsPhone className="w-4 h-4 mr-2" />
                  Google Play
                </button>
              </div>
            </div>

            {/* Покупателям */}
            <div>
              <h4 className="font-bold mb-4">Покупателям</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link href="/how-to-buy" className="hover:text-primary">
                    Как купить
                  </Link>
                </li>
                <li>
                  <Link href="/delivery" className="hover:text-primary">
                    Доставка
                  </Link>
                </li>
                <li>
                  <Link href="/payment" className="hover:text-primary">
                    Оплата
                  </Link>
                </li>
                <li>
                  <Link href="/return" className="hover:text-primary">
                    Возврат
                  </Link>
                </li>
                <li>
                  <Link href="/warranty" className="hover:text-primary">
                    Гарантия
                  </Link>
                </li>
              </ul>
            </div>

            {/* Продавцам */}
            <div>
              <h4 className="font-bold mb-4">Продавцам</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link href="/how-to-sell" className="hover:text-primary">
                    Как продать
                  </Link>
                </li>
                <li>
                  <Link href="/tariffs" className="hover:text-primary">
                    Тарифы
                  </Link>
                </li>
                <li>
                  <Link href="/promotion" className="hover:text-primary">
                    Продвижение
                  </Link>
                </li>
                <li>
                  <Link href="/stores" className="hover:text-primary">
                    Магазины
                  </Link>
                </li>
                <li>
                  <Link href="/api" className="hover:text-primary">
                    API
                  </Link>
                </li>
              </ul>
            </div>

            {/* Помощь */}
            <div>
              <h4 className="font-bold mb-4">Помощь</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link href="/faq" className="hover:text-primary">
                    Частые вопросы
                  </Link>
                </li>
                <li>
                  <Link href="/support" className="hover:text-primary">
                    Поддержка
                  </Link>
                </li>
                <li>
                  <Link href="/rules" className="hover:text-primary">
                    Правила
                  </Link>
                </li>
                <li>
                  <Link href="/blog" className="hover:text-primary">
                    Блог
                  </Link>
                </li>
                <li>
                  <Link href="/contacts" className="hover:text-primary">
                    Контакты
                  </Link>
                </li>
              </ul>
            </div>
          </div>

          <div className="divider my-8"></div>

          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-base-content/60">
            <p>© 2025 SveTu. Все права защищены.</p>
            <div className="flex gap-4">
              <Link href="/terms" className="hover:text-primary">
                Условия использования
              </Link>
              <Link href="/privacy" className="hover:text-primary">
                Конфиденциальность
              </Link>
              <Link href="/cookies" className="hover:text-primary">
                Cookie
              </Link>
            </div>
          </div>
        </div>
      </footer>

      {/* Мобильная навигация */}
      <div className="btm-nav lg:hidden">
        <button className="text-primary">
          <FiSearch className="w-5 h-5" />
          <span className="btm-nav-label">Поиск</span>
        </button>
        <button>
          <FiHeart className="w-5 h-5" />
          <span className="btm-nav-label">Избранное</span>
        </button>
        <button className="text-secondary">
          <div className="indicator">
            <FiShoppingCart className="w-5 h-5" />
            {_cartCount > 0 && (
              <span className="badge badge-xs badge-error indicator-item">
                {_cartCount}
              </span>
            )}
          </div>
          <span className="btm-nav-label">Корзина</span>
        </button>
        <button>
          <FiUser className="w-5 h-5" />
          <span className="btm-nav-label">Профиль</span>
        </button>
      </div>
    </div>
  );
}
