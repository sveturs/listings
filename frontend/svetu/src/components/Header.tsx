'use client';

import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import LanguageSwitcher from './LanguageSwitcher';

export default function Header() {
  const t = useTranslations('header');

  const navItems = [
    { href: '/blog', label: t('nav.blog') },
    { href: '/news', label: t('nav.news') },
    { href: '/contacts', label: t('nav.contacts') },
  ];

  return (
    <header className="navbar bg-base-100 shadow-lg">
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
        <button className="btn btn-primary">{t('login')}</button>
      </div>
    </header>
  );
}
