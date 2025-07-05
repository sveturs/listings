'use client';

import { useLocale } from 'next-intl';
import { usePathname, useRouter } from '@/i18n/routing';
import { useTransition } from 'react';

export default function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();
  const [isPending, startTransition] = useTransition();

  const handleLocaleChange = (nextLocale: string) => {
    startTransition(() => {
      router.push(pathname, { locale: nextLocale });
    });
  };

  return (
    <div className="dropdown dropdown-end">
      <button tabIndex={0} className="btn btn-ghost" disabled={isPending}>
        {locale.toUpperCase()}
        <svg
          className="ml-1 h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
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
        tabIndex={0}
        className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-24"
      >
        <li>
          <button
            onClick={() => handleLocaleChange('ru')}
            className={locale === 'ru' ? 'active' : ''}
          >
            РУС
          </button>
        </li>
        <li>
          <button
            onClick={() => handleLocaleChange('en')}
            className={locale === 'en' ? 'active' : ''}
          >
            ENG
          </button>
        </li>
        <li>
          <button
            onClick={() => handleLocaleChange('sr')}
            className={locale === 'sr' ? 'active' : ''}
          >
            SRP
          </button>
        </li>
      </ul>
    </div>
  );
}
