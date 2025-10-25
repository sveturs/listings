'use client';

import { useLocale, useTranslations } from 'next-intl';
import { usePathname, useRouter } from '@/i18n/routing';
import { useTransition } from 'react';
import { i18n } from '@/i18n/config';

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();
  const [isPending, startTransition] = useTransition();
  const t = useTranslations('common');

  const handleLocaleChange = (nextLocale: string) => {
    // Сохраняем выбор пользователя в cookie
    document.cookie = `${i18n.localeDetection.cookieName}=${nextLocale}; max-age=${i18n.localeDetection.cookieMaxAge}; path=/; samesite=lax`;

    startTransition(() => {
      router.push(pathname, { locale: nextLocale });
    });
  };

  const languages = [
    { code: 'ru', label: 'РУС' },
    { code: 'en', label: 'ENG' },
    { code: 'sr', label: 'SRP' },
  ];

  return (
    <div className="dropdown dropdown-end">
      <button
        tabIndex={0}
        className="btn btn-ghost"
        disabled={isPending}
        aria-label={t('language.switchLanguage')}
        aria-expanded="false"
        aria-haspopup="listbox"
        aria-controls="language-menu"
      >
        {locale.toUpperCase()}
        <svg
          className="ml-1 h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>
      <ul
        id="language-menu"
        tabIndex={0}
        className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-24"
        role="listbox"
        aria-label={t('language.selectLanguage')}
      >
        {languages.map((lang) => (
          <li key={lang.code} role="none">
            <button
              onClick={() => handleLocaleChange(lang.code)}
              className={locale === lang.code ? 'active' : ''}
              role="option"
              aria-selected={locale === lang.code}
            >
              {lang.label}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
