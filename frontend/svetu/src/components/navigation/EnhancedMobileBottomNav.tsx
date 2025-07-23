'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Search, PlusCircle, MessageCircle, User } from 'lucide-react';
import { useTranslations } from 'next-intl';
import { useAuthContext } from '@/contexts/AuthContext';

interface NavItem {
  icon: React.ElementType;
  label: string;
  href: string;
  badge?: number;
  authRequired?: boolean;
  color?: string;
}

export const EnhancedMobileBottomNav: React.FC = () => {
  const pathname = usePathname();
  const t = useTranslations('navigation');
  const { isAuthenticated } = useAuthContext();
  const [mounted, setMounted] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    setMounted(true);
  }, []);

  const navItems: NavItem[] = [
    {
      icon: Home,
      label: t('home'),
      href: '/',
    },
    {
      icon: Search,
      label: t('search'),
      href: '/search',
    },
    {
      icon: PlusCircle,
      label: t('create'),
      href: '/create-listing',
      authRequired: true,
      color: 'text-primary',
    },
    {
      icon: MessageCircle,
      label: t('chats'),
      href: '/chat',
      authRequired: true,
      badge: 3, // TODO: подключить реальные уведомления
    },
    {
      icon: User,
      label: t('profile'),
      href: mounted && isAuthenticated ? '/profile' : '/auth/login',
    },
  ];

  // Фильтруем элементы, требующие авторизации
  const visibleItems = navItems.filter(
    (item) =>
      !item.authRequired || (item.authRequired && mounted && isAuthenticated)
  );

  // Функция для определения активности ссылки
  const isActive = (href: string) => {
    const pathWithoutLocale = pathname.replace(/^\/(en|ru|sr)/, '') || '/';
    const hrefWithoutLocale = href === '/' ? '/' : href;

    if (hrefWithoutLocale === '/') {
      return pathWithoutLocale === '/';
    }

    return pathWithoutLocale.startsWith(hrefWithoutLocale);
  };

  // Обновляем активный индекс при изменении pathname
  useEffect(() => {
    const currentIndex = visibleItems.findIndex((item) => isActive(item.href));
    if (currentIndex !== -1) {
      setActiveIndex(currentIndex);
    }
  }, [pathname, visibleItems]);

  // Добавляем локаль к href
  const getLocalizedHref = (href: string) => {
    const locale = pathname.split('/')[1];
    if (['en', 'ru', 'sr'].includes(locale)) {
      return href === '/' ? `/${locale}` : `/${locale}${href}`;
    }
    return href;
  };

  // Вычисляем позицию индикатора
  const indicatorStyle = {
    left: `${(activeIndex / visibleItems.length) * 100}%`,
    width: `${100 / visibleItems.length}%`,
  };

  return (
    <nav className="btm-nav btm-nav-sm md:hidden bg-base-100 border-t border-base-200 relative">
      {/* Анимированный индикатор */}
      <div
        className="absolute top-0 h-0.5 bg-primary transition-all duration-300 ease-out"
        style={indicatorStyle}
      />

      {visibleItems.map((item, index) => {
        const Icon = item.icon;
        const active = isActive(item.href);
        const localizedHref = getLocalizedHref(item.href);

        return (
          <Link
            key={item.href}
            href={localizedHref}
            className={`${active ? 'active' : ''} relative group`}
            onClick={() => setActiveIndex(index)}
          >
            <div className="flex flex-col items-center justify-center h-full">
              {/* Фоновая подсветка при наведении */}
              <div className="absolute inset-0 bg-primary/5 scale-0 group-hover:scale-100 transition-transform duration-200 rounded-lg" />

              {/* Иконка с анимацией */}
              <div className="relative">
                <Icon
                  className={`
                    w-5 h-5 transition-all duration-200
                    ${
                      active
                        ? 'text-primary scale-110'
                        : item.color || 'text-base-content/60'
                    }
                    group-hover:scale-110
                  `}
                />

                {/* Анимированный бейдж */}
                {item.badge !== undefined && item.badge > 0 && (
                  <span className="badge badge-error badge-xs absolute -top-1 -right-2 animate-pulse">
                    {item.badge > 9 ? '9+' : item.badge}
                  </span>
                )}

                {/* Ripple эффект для центральной кнопки */}
                {item.href === '/create-listing' && (
                  <div className="absolute inset-0 -m-2">
                    <div className="absolute inset-0 rounded-full bg-primary/20 animate-ping" />
                  </div>
                )}
              </div>

              {/* Подпись с анимацией */}
              <span
                className={`
                  text-xs mt-1 transition-all duration-200
                  ${
                    active
                      ? 'text-primary font-medium translate-y-0'
                      : 'text-base-content/60 translate-y-0.5'
                  }
                `}
              >
                {item.label}
              </span>

              {/* Точка для активного элемента */}
              {active && (
                <div className="absolute bottom-1 w-1 h-1 bg-primary rounded-full animate-pulse" />
              )}
            </div>
          </Link>
        );
      })}
    </nav>
  );
};
