'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import Image from 'next/image';

export function AuthButton() {
  const { user, isAuthenticated, isLoading, login, logout } = useAuth();
  const t = useTranslations('auth');

  if (isLoading) {
    return <div className="skeleton h-10 w-24"></div>;
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
            {user.picture_url ? (
              <div className="relative w-10 h-10">
                <Image
                  alt={user.name}
                  src={user.picture_url}
                  fill
                  className="rounded-full object-cover"
                />
              </div>
            ) : (
              <div className="bg-neutral text-neutral-content w-full h-full flex items-center justify-center">
                {user.name.charAt(0).toUpperCase()}
              </div>
            )}
          </div>
        </div>
        <ul
          tabIndex={0}
          className="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
        >
          <li className="menu-title">
            <span>{user.name}</span>
            <span className="text-xs font-normal">{user.email}</span>
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
            <button onClick={logout}>{t('logout')}</button>
          </li>
        </ul>
      </div>
    );
  }

  return (
    <button onClick={() => login()} className="btn btn-primary btn-sm">
      {t('login')}
    </button>
  );
}
