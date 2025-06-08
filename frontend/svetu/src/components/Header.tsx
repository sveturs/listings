'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import LanguageSwitcher from './LanguageSwitcher';
import { AuthButton } from './AuthButton';
import LoginModal from './LoginModal';
import { useAuthContext } from '@/contexts/AuthContext';

export default function Header() {
  const t = useTranslations('header');
  const { user, isAuthenticated } = useAuthContext();
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  // Проверяем, что компонент смонтирован на клиенте
  useEffect(() => {
    setMounted(true);
  }, []);

  const navItems = [
    { href: '/blog', label: t('nav.blog') },
    { href: '/news', label: t('nav.news') },
    { href: '/contacts', label: t('nav.contacts') },
    { href: '/user-contacts', label: t('nav.userContacts') },
  ];

  return (
    <>
      <header className="navbar bg-base-100 shadow-lg fixed top-0 left-0 right-0 z-[100]">
        <div className="navbar-start">
          <Link href="/" className="btn btn-ghost text-xl">
            SveTu
          </Link>
        </div>

        <div className="navbar-center hidden lg:flex">
          <ul className="menu menu-horizontal px-1">
            {navItems.map((item) => (
              <li key={item.href}>
                <Link href={item.href} className="btn btn-ghost">
                  {item.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>

        <div className="navbar-end space-x-2">
          {mounted && isAuthenticated && user && (
            <Link href="/create-listing" className="btn btn-primary btn-sm">
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
              {t('nav.createListing')}
            </Link>
          )}
          <LanguageSwitcher />
          <AuthButton onLoginClick={() => setIsLoginModalOpen(true)} />
        </div>
      </header>

      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => setIsLoginModalOpen(false)}
      />
    </>
  );
}
