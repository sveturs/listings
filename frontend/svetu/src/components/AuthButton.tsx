'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { Link, useRouter } from '@/i18n/routing';
import Image from 'next/image';
import { useChat } from '@/hooks/useChat';
import { useAppDispatch } from '@/store/hooks';
import { clearCartOnLogout } from '@/store/slices/cartSlice';

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
    logout,
    clearError,
  } = useAuth();
  const router = useRouter();
  const t = useTranslations('auth');
  const tCommon = useTranslations('common');
  const [imageError, setImageError] = useState(false);
  const [mounted, setMounted] = useState(false);
  const dispatch = useAppDispatch();

  // Получаем количество непрочитанных сообщений из useChat
  const { unreadCount } = useChat();

  // Обработчик logout с очисткой корзины
  const handleLogout = async () => {
    // Сначала очищаем корзину в Redux
    dispatch(clearCartOnLogout());
    // Затем выполняем logout
    await logout();
  };

  // Временное логирование для отладки
  useEffect(() => {
    console.log('[AuthButton] Unread count:', unreadCount);
    console.log('[AuthButton] User data:', user);
    console.log('[AuthButton] Is admin?:', user?.is_admin);
  }, [unreadCount, user]);

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
                    {(user.name || user.email || 'U').charAt(0).toUpperCase()}
                  </span>
                </div>
              )}
            </div>
          </div>
          <ul
            tabIndex={0}
            className="menu dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-64 p-3 shadow-xl border border-base-200"
          >
            {/* User Info Section */}
            <li className="mb-2">
              <div className="flex items-center gap-3 p-3 bg-base-200 rounded-lg hover:bg-base-200 cursor-default">
                <div className="avatar">
                  <div className="w-12 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                    {user.picture_url && !imageError ? (
                      <Image
                        alt={`Profile picture of ${user.name}`}
                        src={user.picture_url}
                        width={48}
                        height={48}
                        className="rounded-full object-cover"
                        onError={() => setImageError(true)}
                      />
                    ) : (
                      <div className="bg-primary text-primary-content w-full h-full flex items-center justify-center">
                        <span className="text-xl font-bold">
                          {(user.name || user.email || 'U')
                            .charAt(0)
                            .toUpperCase()}
                        </span>
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-semibold truncate">{user.name}</p>
                  <p className="text-xs text-base-content/60 truncate">
                    {user.email}
                  </p>
                </div>
              </div>
            </li>

            <div className="divider my-1"></div>

            {/* Profile Section */}
            <li>
              <Link href="/profile" className="flex items-center gap-3 py-3">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
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
                <span>{t('profile')}</span>
              </Link>
            </li>
            <li>
              <Link
                href="/user-contacts"
                className="flex items-center gap-3 py-3"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                  />
                </svg>
                <span>{t('myContacts')}</span>
              </Link>
            </li>
            <li>
              <Link href="/favorites" className="flex items-center gap-3 py-3">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                  />
                </svg>
                <span>{t('myFavorites')}</span>
              </Link>
            </li>

            <div className="divider my-1"></div>

            {/* Business Section */}
            <li>
              <Link
                href="/profile/listings"
                className="flex items-center gap-3 py-3"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                  />
                </svg>
                <span>{t('myListings')}</span>
              </Link>
            </li>
            <li>
              <Link
                href="/profile/storefronts"
                className="flex items-center gap-3 py-3"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                  />
                </svg>
                <span>{t('myStorefronts')}</span>
              </Link>
            </li>
            <li>
              <Link
                href="/profile/orders"
                className="flex items-center gap-3 py-3"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-base-content/60"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                  />
                </svg>
                <span>{t('myOrders')}</span>
              </Link>
            </li>

            {user.is_admin && (
              <>
                <div className="divider my-1"></div>

                {/* Admin Section */}
                <li>
                  <Link href="/admin" className="flex items-center gap-3 py-3">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5 text-base-content/60"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                    </svg>
                    <span>{t('adminPanel')}</span>
                  </Link>
                </li>
                <li>
                  <Link href="/docs" className="flex items-center gap-3 py-3">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5 text-base-content/60"
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
                    <span>
                      {tCommon.has('documentation')
                        ? tCommon('documentation')
                        : 'Documentation'}
                    </span>
                  </Link>
                </li>
              </>
            )}

            <div className="divider my-1"></div>

            {/* Logout */}
            <li>
              <button
                onClick={handleLogout}
                disabled={isLoggingOut}
                className={`flex items-center gap-3 py-3 text-error hover:bg-error/10 ${isLoggingOut ? 'loading' : ''}`}
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
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                  />
                </svg>
                <span>{isLoggingOut ? t('loggingOut') : t('logout')}</span>
              </button>
            </li>
          </ul>
        </div>
      </>
    );
  }

  return (
    <button
      onClick={() =>
        onLoginClick ? onLoginClick() : router.push('/auth/login')
      }
      className="btn btn-primary btn-sm"
      aria-label={t('login.enter')}
    >
      {t('login.enter')}
    </button>
  );
}
