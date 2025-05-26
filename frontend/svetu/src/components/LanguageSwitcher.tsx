'use client';

import { usePathname, useRouter } from 'next/navigation';
import { locales, type Locale } from '@/i18n/request';

export default function LanguageSwitcher({ currentLocale }: { currentLocale: string }) {
  const pathname = usePathname();
  const router = useRouter();

  const handleLocaleChange = (newLocale: Locale) => {
    const pathnameWithoutLocale = pathname.split('/').slice(2).join('/');
    const newPath = `/${newLocale}${pathnameWithoutLocale ? `/${pathnameWithoutLocale}` : ''}`;
    
    document.cookie = `locale=${newLocale};path=/;max-age=31536000`;
    router.push(newPath);
  };

  return (
    <div className="flex gap-2">
      {locales.map((locale) => (
        <button
          key={locale}
          onClick={() => handleLocaleChange(locale)}
          className={`px-2 py-1 rounded ${
            currentLocale === locale 
              ? 'bg-blue-600 text-white' 
              : 'hover:text-gray-300'
          }`}
        >
          {locale.toUpperCase()}
        </button>
      ))}
    </div>
  );
}