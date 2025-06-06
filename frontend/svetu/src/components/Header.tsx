'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import LanguageSwitcher from './LanguageSwitcher';
import { AuthButton } from './AuthButton';
import LoginModal from './LoginModal';

export default function Header() {
  const t = useTranslations('header');
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

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
