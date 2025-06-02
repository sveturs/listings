'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import Image from 'next/image';

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
  const [imageError, setImageError] = useState(false);
  const [mounted, setMounted] = useState(false);

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
          {user.is_admin && (
            <li>
              <Link href="/admin">{t('adminPanel')}</Link>
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
    );
  }

  return (
    <button
      onClick={() => (onLoginClick ? onLoginClick() : login())}
      className="btn btn-primary btn-sm"
      aria-label={t('login')}
    >
      {t('login')}
    </button>
  );
}
