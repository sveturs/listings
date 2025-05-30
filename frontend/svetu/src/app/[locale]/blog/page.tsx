'use client';

import { useTranslations } from 'next-intl';

export default function BlogPage() {
  const t = useTranslations('blog');

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-6">{t('title')}</h1>
      <p className="text-lg text-gray-600 mb-8">{t('description')}</p>

      <div className="grid gap-6">
        <article className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-2">
            {t('latestPost.title')}
          </h2>
          <p className="text-gray-600 mb-2">{t('latestPost.date')}</p>
          <p className="text-gray-700">{t('latestPost.excerpt')}</p>
          <a
            href="#"
            className="text-blue-600 hover:underline mt-2 inline-block"
          >
            {t('readMore')}
          </a>
        </article>
      </div>
    </div>
  );
}
