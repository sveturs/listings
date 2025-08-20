'use client';

import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Search, PlusCircle, MessageCircle, User, X } from 'lucide-react';
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

interface EnhancedMobileBottomNavProps {
  onClose?: () => void;
}

export const EnhancedMobileBottomNav: React.FC<
  EnhancedMobileBottomNavProps
> = ({ onClose }) => {
  const pathname = usePathname();
  const t = useTranslations('common');
  const { isAuthenticated } = useAuthContext();
  const [mounted, setMounted] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    setMounted(true);
  }, []);

  const navItems: NavItem[] = [
    {
      icon: Home,
      label: t('navigation.home'),
      href: '/',
    },
    {
      icon: Search,
      label: t('navigation.search'),
      href: '/search',
    },
    {
      icon: PlusCircle,
      label: t('navigation.create'),
      href: '/create-listing-choice',
      color: 'text-primary',
    },
    {
      icon: MessageCircle,
      label: t('navigation.chats'),
      href: '/chat',
      authRequired: true,
      badge: 3, // TODO: подключить реальные уведомления
    },
    {
      icon: User,
      label: t('navigation.profile'),
      href: mounted && isAuthenticated ? '/profile' : '/auth/login',
    },
  ];

  // Фильтруем элементы, требующие авторизации
  const visibleItems = navItems.filter(
    (item) =>
      !item.authRequired || (item.authRequired && mounted && isAuthenticated)
  );

  // Функция для определения активности ссылки
  const isActive = useCallback(
    (href: string) => {
      const pathWithoutLocale = pathname.replace(/^\/(en|ru|sr)/, '') || '/';
      const hrefWithoutLocale = href === '/' ? '/' : href;

      if (hrefWithoutLocale === '/') {
        return pathWithoutLocale === '/';
      }

      return pathWithoutLocale.startsWith(hrefWithoutLocale);
    },
    [pathname]
  );

  // Обновляем активный индекс при изменении pathname
  useEffect(() => {
    const currentIndex = visibleItems.findIndex((item) => isActive(item.href));
    if (currentIndex !== -1) {
      setActiveIndex(currentIndex);
    }
  }, [pathname, visibleItems, isActive]);

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
    <nav className="md:hidden bg-base-100 border-t border-base-200 relative">
      <div className="flex items-center justify-between px-2 py-2">
        {/* Навигационные элементы */}
        <div className="flex items-center justify-around flex-1">
          {visibleItems.map((item, index) => {
            const Icon = item.icon;
            const active = isActive(item.href);
            const localizedHref = getLocalizedHref(item.href);

            return (
              <Link
                key={item.href}
                href={localizedHref}
                className={`relative group flex items-center gap-2 px-3 py-2 rounded-lg transition-colors tooltip tooltip-top ${
                  active ? 'bg-primary/10 text-primary' : 'hover:bg-base-200'
                }`}
                data-tip={item.label}
                onClick={() => setActiveIndex(index)}
              >
                {/* Иконка */}
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

                  {/* Бейдж */}
                  {item.badge !== undefined && item.badge > 0 && (
                    <span className="badge badge-error badge-xs absolute -top-1 -right-2 animate-pulse">
                      {item.badge > 9 ? '9+' : item.badge}
                    </span>
                  )}

                  {/* Ripple эффект для центральной кнопки */}
                  {item.href === '/create-listing-choice' && (
                    <div className="absolute inset-0 -m-2">
                      <div className="absolute inset-0 rounded-full bg-primary/20 animate-ping" />
                    </div>
                  )}
                </div>

                {/* Подпись - скрываем на маленьких экранах */}
                <span
                  className={`
                    hidden sm:block text-xs transition-all duration-200
                    ${
                      active
                        ? 'text-primary font-medium'
                        : 'text-base-content/60'
                    }
                  `}
                >
                  {item.label}
                </span>
              </Link>
            );
          })}
        </div>

        {/* Кнопка закрытия */}
        {onClose && (
          <button
            onClick={onClose}
            className="btn btn-ghost btn-sm btn-circle ml-2"
            aria-label={t('navigation.close')}
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>

      {/* Анимированный индикатор активного элемента */}
      <div
        className="absolute top-0 h-0.5 bg-primary transition-all duration-300 ease-out"
        style={indicatorStyle}
      />
    </nav>
  );
};
