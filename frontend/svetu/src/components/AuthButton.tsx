'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import Image from 'next/image';
import { useChat } from '@/hooks/useChat';

interface AuthButtonProps {
  onLoginClick?: () => void;
}

export function AuthButton({ onLoginClick }: AuthButtonProps) {
  const {
    user,
    isAuthenticated,
    isLoading,
    isLoggingOut,
    error,
    login,
    logout,
    clearError,
  } = useAuth();
  const t = useTranslations('auth');
  const tCommon = useTranslations('common');
  const [imageError, setImageError] = useState(false);
  const [mounted, setMounted] = useState(false);

  // Получаем количество непрочитанных сообщений из useChat
  const { unreadCount } = useChat();

  // Временное логирование для отладки
  useEffect(() => {
    console.log('[AuthButton] Unread count:', unreadCount);
  }, [unreadCount]);

  // Ensure component is mounted before rendering dynamic content
  useEffect(() => {
    setMounted(true);
  }, []);

  // Reset image error when user changes
  useEffect(() => {
    setImageError(false);
  }, [user?.picture_url]);

  // Close error after 5 seconds
  useEffect(() => {
    if (error) {
      const timer = setTimeout(clearError, 5000);
      return () => clearTimeout(timer);
    }
  }, [error, clearError]);

  // Always show loading state during SSR and initial mount
  if (!mounted || isLoading) {
    return (
      <div
        className="skeleton h-10 w-24"
        aria-label="Loading authentication status"
      ></div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error py-2 px-4">
        <span className="text-sm">{error}</span>
        <button
          onClick={clearError}
          className="btn btn-ghost btn-xs"
          aria-label="Dismiss error"
        >
          ✕
        </button>
      </div>
    );
  }

  if (isAuthenticated && user) {
    return (
      <>
        {/* Иконка чата */}
        <Link href="/chat" className="btn btn-ghost btn-circle">
          <div className="indicator">
            <svg
              className="w-5 h-5"
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
            {/* Пульсирующий индикатор непрочитанных сообщений */}
            {unreadCount > 0 && (
              <span className="indicator-item indicator-top indicator-end">
                <span className="relative flex h-3 w-3">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                </span>
              </span>
            )}
          </div>
        </Link>

        {/* Меню пользователя */}
        <div className="dropdown dropdown-end">
          <div
            tabIndex={0}
            role="button"
            className="btn btn-ghost btn-circle avatar"
          >
            <div className="w-10 rounded-full">
              {user.picture_url && !imageError ? (
                <div className="relative w-10 h-10">
                  <Image
                    alt={`Profile picture of ${user.name}`}
                    src={user.picture_url}
                    fill
                    className="rounded-full object-cover"
                    onError={() => setImageError(true)}
                    sizes="40px"
                    priority={false}
                    placeholder="empty"
                  />
                </div>
              ) : (
                <div className="bg-neutral text-neutral-content w-full h-full flex items-center justify-center rounded-full">
                  <span className="text-lg font-semibold">
                    {user.name.charAt(0).toUpperCase()}
                  </span>
                </div>
              )}
            </div>
          </div>
          <ul
            tabIndex={0}
            className="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
          >
            <li className="menu-title">
              <span className="truncate">{user.name}</span>
              <span className="text-xs font-normal truncate">{user.email}</span>
            </li>
            <li>
              <Link href="/profile">{t('profile')}</Link>
            </li>
            <li>
              <Link href="/user-contacts">{t('myContacts')}</Link>
            </li>
            {user.is_admin && (
              <li>
                <Link href="/admin">{t('adminPanel')}</Link>
              </li>
            )}
            {user.is_admin && (
              <li>
                <Link href="/docs">
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
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                  {tCommon.has('documentation')
                    ? tCommon('documentation')
                    : 'Documentation'}
                </Link>
              </li>
            )}
            <li>
              <button
                onClick={logout}
                disabled={isLoggingOut}
                className={isLoggingOut ? 'loading' : ''}
              >
                {isLoggingOut ? t('loggingOut') : t('logout')}
              </button>
            </li>
          </ul>
        </div>
      </>
    );
  }

  return (
    <button
      onClick={() => (onLoginClick ? onLoginClick() : login())}
      className="btn btn-primary btn-sm"
      aria-label={t('login.enter')}
    >
      {t('login.enter')}
    </button>
  );
}
