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
}

export const MobileBottomNav: React.FC = () => {
  const pathname = usePathname();
  const t = useTranslations('navigation');
  const { isAuthenticated } = useAuthContext();
  const [mounted, setMounted] = useState(false);

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
    },
    {
      icon: MessageCircle,
      label: t('chats'),
      href: '/chat',
      authRequired: true,
      badge: 0, // TODO: добавить реальное количество непрочитанных
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
    // Извлекаем путь без локали
    const pathWithoutLocale = pathname.replace(/^\/(en|ru|sr)/, '') || '/';
    const hrefWithoutLocale = href === '/' ? '/' : href;

    if (hrefWithoutLocale === '/') {
      return pathWithoutLocale === '/';
    }

    return pathWithoutLocale.startsWith(hrefWithoutLocale);
  };

  // Добавляем локаль к href
  const getLocalizedHref = (href: string) => {
    const locale = pathname.split('/')[1];
    if (['en', 'ru', 'sr'].includes(locale)) {
      return href === '/' ? `/${locale}` : `/${locale}${href}`;
    }
    return href;
  };

  return (
    <nav className="btm-nav btm-nav-sm md:hidden bg-base-100 border-t border-base-200">
      {visibleItems.map((item) => {
        const Icon = item.icon;
        const active = isActive(item.href);
        const localizedHref = getLocalizedHref(item.href);

        return (
          <Link
            key={item.href}
            href={localizedHref}
            className={`${active ? 'active' : ''} relative`}
          >
            <div className="flex flex-col items-center justify-center h-full">
              {/* Активный индикатор */}
              {active && (
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-12 h-1 bg-primary rounded-b-full" />
              )}

              {/* Иконка с бейджем */}
              <div className="relative">
                <Icon
                  className={`w-5 h-5 ${active ? 'text-primary' : 'text-base-content/60'}`}
                />
                {item.badge !== undefined && item.badge > 0 && (
                  <span className="badge badge-error badge-xs absolute -top-1 -right-2">
                    {item.badge > 9 ? '9+' : item.badge}
                  </span>
                )}
              </div>

              {/* Подпись */}
              <span
                className={`text-xs mt-1 ${active ? 'text-primary font-medium' : 'text-base-content/60'}`}
              >
                {item.label}
              </span>
            </div>
          </Link>
        );
      })}
    </nav>
  );
};
