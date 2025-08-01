'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Link, useRouter } from '@/i18n/routing';
import { usePathname, useSearchParams } from 'next/navigation';
import LanguageSwitcher from './LanguageSwitcher';
import { AuthButton } from './AuthButton';
import LoginModal from './LoginModal';
import { SearchBar } from './SearchBar';
import { useAuthContext } from '@/contexts/AuthContext';
import CartIcon from './cart/CartIcon';
import ShoppingCartModal from './cart/ShoppingCartModal';
import { ThemeToggle } from './ThemeToggle';
import Image from 'next/image';

export default function Header() {
  const t = useTranslations('header');
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const { isAuthenticated } = useAuthContext();
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const [isCartModalOpen, setIsCartModalOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  // Функция для извлечения ID витрины из пути
  const extractStorefrontIdFromPath = (path: string): number | null => {
    // Для страниц витрин вида /storefronts/tech-store-dmitry или /storefronts/tech-store-dmitry/products/1
    // можем временно извлечь ID из slug, если известно соответствие
    if (path.includes('/storefronts/')) {
      // Временное решение для витрины tech-store-dmitry = ID 4
      if (path.includes('tech-store-dmitry')) {
        return 4;
      }
      // Можно добавить другие известные витрины
    }
    return null;
  };

  // Определяем активную витрину из URL
  const currentStorefrontId = searchParams?.get('storefront')
    ? Number(searchParams?.get('storefront'))
    : pathname?.includes('/storefronts/')
      ? extractStorefrontIdFromPath(pathname)
      : null;

  // Не показываем мобильный поиск на странице поиска и главной
  const isSearchPage = pathname?.includes('/search');
  const isHomePage =
    pathname === '/' ||
    pathname === '/en' ||
    pathname === '/ru' ||
    pathname === '/sr';

  // Проверяем, что компонент смонтирован на клиенте
  useEffect(() => {
    setMounted(true);
  }, []);

  // Закрываем модалку логина при успешной аутентификации
  useEffect(() => {
    if (isAuthenticated && isLoginModalOpen) {
      setIsLoginModalOpen(false);
    }
  }, [isAuthenticated, isLoginModalOpen]);

  const handleCheckout = () => {
    if (currentStorefrontId) {
      router.push(`/checkout?storefront=${currentStorefrontId}`);
    }
  };

  const navItems = [{ href: '/map', label: t('nav.map') }];

  return (
    <>
      <header className="navbar bg-base-100 shadow-lg fixed top-0 left-0 right-0 z-[100] h-16">
        <div className="container mx-auto flex items-center">
          {/* Лого и мобильное меню */}
          <div className="flex-none">
            <div className="dropdown lg:hidden">
              <div tabIndex={0} role="button" className="btn btn-ghost btn-sm">
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
                    strokeWidth="2"
                    d="M4 6h16M4 12h8m-8 6h16"
                  />
                </svg>
              </div>
              <ul
                tabIndex={0}
                className="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
              >
                {navItems.map((item) => (
                  <li key={item.href}>
                    <Link href={item.href} className="flex items-center gap-2">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                        />
                      </svg>
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            <Link
              href="/"
              className="btn btn-ghost text-xl px-2 flex items-center gap-2"
            >
              <div className="w-8 h-8 flex items-center justify-center">
                <Image
                  src="/logos/svetu-gradient-48x48.png"
                  alt="SveTu Logo"
                  width={32}
                  height={32}
                  className="object-contain"
                />
              </div>
              SveTu
            </Link>
          </div>

          {/* Навигация для десктопа */}
          <div className="hidden lg:flex items-center ml-4">
            <ul className="menu menu-horizontal px-1">
              {navItems.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className="text-sm tooltip tooltip-bottom"
                    data-tip={item.label}
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
                        d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                      />
                    </svg>
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Правая часть */}
          <div className="flex-none flex items-center gap-2">
            {/* Корзина - показываем всегда для B2C товаров */}
            {mounted && <CartIcon />}

            {mounted && (
              <Link
                href="/create-listing-choice"
                className="btn btn-primary btn-sm hidden md:flex tooltip tooltip-bottom"
                data-tip={t('nav.createListing')}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-4 w-4"
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
                <span className="hidden md:inline">
                  {t('nav.createListing')}
                </span>
              </Link>
            )}
            <ThemeToggle />
            <LanguageSwitcher />
            <AuthButton onLoginClick={() => setIsLoginModalOpen(true)} />
          </div>
        </div>
      </header>

      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => setIsLoginModalOpen(false)}
      />

      {/* Модальное окно корзины */}
      {currentStorefrontId && (
        <ShoppingCartModal
          storefrontId={currentStorefrontId}
          isOpen={isCartModalOpen}
          onClose={() => setIsCartModalOpen(false)}
          onCheckout={handleCheckout}
        />
      )}

      {/* Мобильная поисковая строка - скрываем на странице поиска и главной */}
      {!isSearchPage && !isHomePage && (
        <div className="lg:hidden bg-base-100 border-t border-base-300 px-4 py-2 fixed top-16 left-0 right-0 z-[99]">
          <SearchBar className="w-full" placeholder={t('search.placeholder')} />
        </div>
      )}
    </>
  );
}
