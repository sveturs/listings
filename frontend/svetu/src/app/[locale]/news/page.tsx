'use client';

import { useTranslations } from 'next-intl';

export default function NewsPage() {
  const t = useTranslations('news');

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-6">{t('title')}</h1>
      <p className="text-lg text-gray-600 mb-8">{t('description')}</p>

      <div className="grid gap-4">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-2">
            {t('latestNews.title')}
          </h2>
          <p className="text-gray-600 mb-2">{t('latestNews.date')}</p>
          <p className="text-gray-700">{t('latestNews.content')}</p>
          <a
            href="#"
            className="text-blue-600 hover:underline mt-2 inline-block"
          >
            {t('readMore')}
          </a>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-2">
            {t('secondNews.title')}
          </h2>
          <p className="text-gray-600 mb-2">{t('secondNews.date')}</p>
          <p className="text-gray-700">{t('secondNews.content')}</p>
          <a
            href="#"
            className="text-blue-600 hover:underline mt-2 inline-block"
          >
            {t('readMore')}
          </a>
        </div>
      </div>
    </div>
  );
}
