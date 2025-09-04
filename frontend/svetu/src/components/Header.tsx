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
import { useAuthContext } from '@/contexts/AuthContext';
import CartIcon from './cart/CartIcon';
import ShoppingCartModal from './cart/ShoppingCartModal';
import { ThemeToggle } from './ThemeToggle';
import ChatIcon from './ChatIcon';

// Иконки
import { FiMapPin, FiHeart, FiMenu, FiX } from 'react-icons/fi';
import { BsHandbag } from 'react-icons/bs';

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
  const [popularCategories, setPopularCategories] = useState<any[]>([]);
  const [selectedCategory] = useState<string | number>('all');
  const [isVisible, setIsVisible] = useState(true);
  const [_lastScrollY, setLastScrollY] = useState(0);
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

    setLastScrollY((prevScrollY) => {
      const scrollDifference = Math.abs(currentScrollY - prevScrollY);

      // Увеличиваем threshold для предотвращения дрожания
      if (scrollDifference < 30) return prevScrollY;

      // Определяем направление скролла
      const isScrollingDown = currentScrollY > prevScrollY + 30;
      const isScrollingUp = currentScrollY < prevScrollY - 30;

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
      return currentScrollY;
    });
  }, []);

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

  // Загрузка популярных категорий для мобильного меню
  useEffect(() => {
    const loadCategories = async () => {
      try {
        const popularResponse = await api.get(
          `/api/v1/marketplace/popular-categories?lang=${locale}&limit=8`
        );

        if (popularResponse.data.success && popularResponse.data.data) {
          // Добавляем иконки для популярных категорий на основе их slug
          const iconMap: { [key: string]: any } = {
            'real-estate': BsHandbag,
            automotive: BsHandbag,
            electronics: BsHandbag,
            fashion: BsHandbag,
            jobs: BsHandbag,
            services: BsHandbag,
            'hobbies-entertainment': BsHandbag,
            'home-garden': BsHandbag,
            industrial: BsHandbag,
            'food-beverages': BsHandbag,
            'books-stationery': BsHandbag,
            'antiques-art': BsHandbag,
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
        console.error('Failed to load popular categories:', error);
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
    </>
  );
}
