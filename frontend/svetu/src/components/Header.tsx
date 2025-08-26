'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Link, useRouter } from '@/i18n/routing';
import { usePathname, useSearchParams } from 'next/navigation';
import { useParams } from 'next/navigation';
import Image from 'next/image';
import { motion, AnimatePresence } from 'framer-motion';
import api from '@/services/api';

// Компоненты
import LanguageSwitcher from './LanguageSwitcher';
import { AuthButton } from './AuthButton';
import LoginModal from './LoginModal';
import { SearchAutocomplete } from './search/SearchAutocomplete';
import { NestedCategorySelector } from './search/NestedCategorySelector';
import { useAuthContext } from '@/contexts/AuthContext';
import CartIcon from './cart/CartIcon';
import ShoppingCartModal from './cart/ShoppingCartModal';
import { ThemeToggle } from './ThemeToggle';
import ChatIcon from './ChatIcon';

// Иконки
import { FiMapPin, FiHeart, FiMenu, FiX } from 'react-icons/fi';
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

interface HeaderProps {
  locale?: string;
}

export default function Header({ locale: propsLocale }: HeaderProps = {}) {
  const t = useTranslations('common');
  const tMarketplace = useTranslations('marketplace.home');
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const params = useParams();
  const { isAuthenticated, user } = useAuthContext();

  // Получаем locale из params или props
  const locale = propsLocale || params?.locale || 'en';

  // Состояния
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const [isCartModalOpen, setIsCartModalOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [showMobileMenu, setShowMobileMenu] = useState(false);
  const [categories, setCategories] = useState<any[]>([]);
  const [popularCategories, setPopularCategories] = useState<any[]>([]);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [selectedCategory] = useState<string | number>('all');
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [isScrolled, setIsScrolled] = useState(false);

  // Функция для извлечения ID витрины из пути
  const extractStorefrontIdFromPath = (path: string): number | null => {
    if (path.includes('/storefronts/')) {
      if (path.includes('tech-store-dmitry')) {
        return 4;
      }
    }
    return null;
  };

  // Определяем активную витрину из URL
  const currentStorefrontId = searchParams?.get('storefront')
    ? Number(searchParams?.get('storefront'))
    : pathname?.includes('/storefronts/')
      ? extractStorefrontIdFromPath(pathname)
      : null;

  // Проверяем тип страницы
  const isSearchPage = pathname?.includes('/search');
  const isHomePage =
    pathname === '/' ||
    pathname === '/en' ||
    pathname === '/ru' ||
    pathname === '/sr' ||
    pathname === `/${locale}`;

  // Обработка скролла для анимации скрытия хедера
  const handleScroll = useCallback(() => {
    const currentScrollY = window.scrollY;
    const scrollDifference = Math.abs(currentScrollY - lastScrollY);

    // Увеличиваем threshold для предотвращения дрожания
    if (scrollDifference < 30) return;

    // Определяем направление скролла
    const isScrollingDown = currentScrollY > lastScrollY + 30;
    const isScrollingUp = currentScrollY < lastScrollY - 30;

    // Логика показа/скрытия хедера
    if (isScrollingDown && currentScrollY > 150) {
      // Скролл вниз и проскроллили больше 150px - скрываем хедер
      setIsVisible(false);
    } else if (isScrollingUp || currentScrollY <= 100) {
      // Скролл вверх или в начале страницы - показываем хедер
      setIsVisible(true);
    }

    // Определяем, проскроллена ли страница (для тени)
    setIsScrolled(currentScrollY > 50);
    setLastScrollY(currentScrollY);
  }, [lastScrollY]);

  useEffect(() => {
    let timeoutId: NodeJS.Timeout;

    const scrollHandler = () => {
      // Простой throttle без лишней сложности
      clearTimeout(timeoutId);
      timeoutId = setTimeout(handleScroll, 16); // ~60fps
    };

    window.addEventListener('scroll', scrollHandler, { passive: true });
    return () => {
      window.removeEventListener('scroll', scrollHandler);
      clearTimeout(timeoutId);
    };
  }, [handleScroll]);

  // Монтирование компонента
  useEffect(() => {
    setMounted(true);
  }, []);

  // Закрываем модалку логина при успешной аутентификации
  useEffect(() => {
    if (isAuthenticated && isLoginModalOpen) {
      setIsLoginModalOpen(false);
    }
  }, [isAuthenticated, isLoginModalOpen]);

  // Загрузка категорий
  useEffect(() => {
    const loadCategories = async () => {
      try {
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
        }
      } catch (error) {
        console.error('Failed to load categories:', error);
      } finally {
        setIsLoadingCategories(false);
      }
    };
    loadCategories();
  }, [locale]);

  const handleCheckout = () => {
    if (currentStorefrontId) {
      router.push(`/checkout?storefront=${currentStorefrontId}`);
    }
  };

  // Класс для хедера (убираем CSS transitions чтобы избежать конфликта с Framer Motion)
  const headerClass = `sticky top-0 z-50 ${isScrolled ? 'shadow-lg' : ''}`;

  return (
    <>
      <motion.header
        className={headerClass}
        initial={{ y: 0 }}
        animate={{ y: isVisible ? 0 : -100 }}
        transition={{
          duration: 0.2,
          ease: 'easeOut',
          type: 'tween',
        }}
        style={{
          willChange: 'transform',
        }}
      >
        <div className="bg-base-100/95 backdrop-blur-md border-b border-base-300">
          {/* Основная шапка */}
          <div className="container mx-auto px-4 py-3">
            <div className="flex items-center gap-4">
              {/* Логотип */}
              <Link href="/" className="flex items-center gap-2">
                <div className="text-2xl">
                  <Image
                    src="/logos/svetu-gradient-48x48.png"
                    alt="SveTu"
                    width={32}
                    height={32}
                  />
                </div>
                <span className="text-xl font-bold hidden md:inline">
                  SveTu
                </span>
              </Link>

              {/* Поисковая строка - скрываем на мобильных */}
              <div className="flex-1 max-w-3xl hidden lg:block">
                <SearchAutocomplete
                  placeholder={tMarketplace('searchPlaceholder')}
                  selectedCategory={selectedCategory}
                  locale={locale as string}
                  className="w-full"
                />
              </div>

              {/* Действия пользователя */}
              <div className="flex items-center gap-2 ml-auto">
                {/* Карта */}
                <Link
                  href="/map"
                  className="btn btn-ghost btn-circle tooltip tooltip-bottom hidden sm:inline-flex"
                  data-tip={t('header.nav.map')}
                >
                  <FiMapPin className="w-5 h-5" />
                </Link>

                {/* Избранное - показываем только для авторизованных */}
                {mounted && user && (
                  <Link
                    href="/favorites"
                    className="btn btn-ghost btn-circle relative hidden sm:inline-flex"
                  >
                    <FiHeart className="w-5 h-5" />
                    {/* TODO: Добавить счетчик избранного */}
                  </Link>
                )}

                {/* Чат - показываем только для авторизованных */}
                {mounted && user && <ChatIcon />}

                {/* Корзина */}
                {mounted && <CartIcon />}

                {/* Создать объявление */}
                <Link
                  href="/create-listing-choice"
                  className="btn btn-secondary hidden lg:inline-flex"
                >
                  {t('header.nav.createListing')}
                </Link>

                {/* Переключатель темы */}
                <ThemeToggle />

                {/* Переключатель языка */}
                <LanguageSwitcher />

                {/* Кнопка входа/профиль */}
                <AuthButton onLoginClick={() => setIsLoginModalOpen(true)} />

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

            {/* Мобильная поисковая строка - показываем только не на главной и поиске */}
            {!isHomePage && !isSearchPage && (
              <div className="mt-2 lg:hidden">
                <SearchAutocomplete
                  placeholder={tMarketplace('searchPlaceholder')}
                  selectedCategory={selectedCategory}
                  locale={locale as string}
                  className="w-full"
                />
              </div>
            )}
          </div>

          {/* Категории под поиском - только на десктопе */}
          <AnimatePresence>
            {isVisible && (
              <motion.div
                className="border-t border-base-300 py-2 hidden lg:block"
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: 'auto', opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2 }}
              >
                <div className="container mx-auto px-4">
                  <div className="flex items-center gap-4 text-sm">
                    {/* Селектор всех категорий - слева перед списком */}
                    <div className="flex-shrink-0">
                      <NestedCategorySelector
                        categories={categories}
                        selectedCategory={selectedCategory}
                        onChange={(categoryId) => {
                          router.push(
                            `/${locale}/search?category=${categoryId}`
                          );
                        }}
                        placeholder={tMarketplace('allCategories')}
                        showCounts={true}
                        className="btn btn-ghost btn-sm text-primary hover:bg-primary/10 font-medium gap-1 px-3"
                      />
                    </div>

                    {/* Разделитель после кнопки */}
                    <div className="h-4 w-px bg-base-300"></div>

                    {/* Список популярных категорий */}
                    {isLoadingCategories
                      ? [...Array(7)].map((_, i) => (
                          <div
                            key={i}
                            className="flex items-center gap-2 animate-pulse"
                          >
                            <div className="w-4 h-4 bg-base-300 rounded"></div>
                            <div className="w-20 h-4 bg-base-300 rounded"></div>
                            <div className="w-10 h-3 bg-base-300 rounded"></div>
                          </div>
                        ))
                      : popularCategories.slice(0, 7).map((cat) => {
                          const Icon = cat.icon || BsHandbag;
                          const count = cat.listing_count || cat.count || 0;
                          const formattedCount =
                            count > 1000
                              ? `${Math.floor(count / 1000)}K+`
                              : count;
                          return (
                            <Link
                              key={cat.id}
                              href={`/search?category=${cat.id}`}
                              className="flex items-center gap-2 hover:text-primary transition-colors"
                            >
                              <Icon className={`w-4 h-4 ${cat.color}`} />
                              <span>{cat.name}</span>
                              <span className="text-base-content/50">
                                ({formattedCount})
                              </span>
                            </Link>
                          );
                        })}
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </motion.header>

      {/* Мобильное боковое меню */}
      <AnimatePresence>
        {showMobileMenu && (
          <>
            <motion.div
              className="fixed inset-0 bg-black/50 z-[90] lg:hidden"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setShowMobileMenu(false)}
            />
            <motion.div
              className="fixed right-0 top-0 h-full w-64 bg-base-100 z-[91] lg:hidden shadow-xl"
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'tween' }}
            >
              <div className="p-4">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="text-xl font-bold">Меню</h2>
                  <button
                    className="btn btn-ghost btn-circle btn-sm"
                    onClick={() => setShowMobileMenu(false)}
                  >
                    <FiX className="w-5 h-5" />
                  </button>
                </div>

                <ul className="menu menu-lg">
                  <li>
                    <Link href="/map">
                      <FiMapPin className="w-5 h-5" />
                      {t('header.nav.map')}
                    </Link>
                  </li>
                  {user && (
                    <li>
                      <Link href="/favorites">
                        <FiHeart className="w-5 h-5" />
                        Избранное
                      </Link>
                    </li>
                  )}
                  <li>
                    <Link href="/create-listing-choice">
                      Создать объявление
                    </Link>
                  </li>
                </ul>

                {/* Категории в мобильном меню */}
                <div className="divider">Категории</div>
                <ul className="menu menu-sm max-h-96 overflow-y-auto">
                  {popularCategories.map((cat) => {
                    const Icon = cat.icon || BsHandbag;
                    return (
                      <li key={cat.id}>
                        <Link
                          href={`/search?category=${cat.id}`}
                          onClick={() => setShowMobileMenu(false)}
                        >
                          <Icon className={`w-4 h-4 ${cat.color}`} />
                          {cat.name}
                        </Link>
                      </li>
                    );
                  })}
                </ul>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Модалки */}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => setIsLoginModalOpen(false)}
      />

      {currentStorefrontId && (
        <ShoppingCartModal
          storefrontId={currentStorefrontId}
          isOpen={isCartModalOpen}
          onClose={() => setIsCartModalOpen(false)}
          onCheckout={handleCheckout}
        />
      )}

      {/* Мобильная нижняя навигация */}
      <div className="btm-nav lg:hidden z-50">
        <Link href="/" className={pathname === `/${locale}` ? 'active' : ''}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
            />
          </svg>
          <span className="btm-nav-label text-xs">Главная</span>
        </Link>
        <Link
          href="/search"
          className={pathname?.includes('/search') ? 'active' : ''}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <span className="btm-nav-label text-xs">Поиск</span>
        </Link>
        <Link href="/create-listing-choice">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
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
          <span className="btm-nav-label text-xs">Создать</span>
        </Link>
        {mounted && user && (
          <Link
            href="/chat"
            className={pathname?.includes('/chat') ? 'active' : ''}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
            </svg>
            <span className="btm-nav-label text-xs">Чаты</span>
          </Link>
        )}
        {mounted ? (
          <Link
            href={user ? '/profile' : '/login'}
            className={pathname?.includes('/profile') ? 'active' : ''}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              />
            </svg>
            <span className="btm-nav-label text-xs">Профиль</span>
          </Link>
        ) : (
          <Link href="/profile" className="">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              />
            </svg>
            <span className="btm-nav-label text-xs">Профиль</span>
          </Link>
        )}
      </div>
    </>
  );
}
