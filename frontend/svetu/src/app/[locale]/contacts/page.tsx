'use client';

import { useTranslations } from 'next-intl';

export default function ContactsPage() {
  const t = useTranslations('contacts');

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-6">{t('title')}</h1>
      <p className="text-lg text-gray-600 mb-8">{t('description')}</p>

      <div className="grid md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-2xl font-semibold mb-4">
            {t('contactInfo.title')}
          </h2>
          <div className="space-y-3">
            <div>
              <strong className="block text-gray-700">
                {t('contactInfo.phone')}:
              </strong>
              <a
                href="tel:+74951234567"
                className="text-blue-600 hover:underline"
              >
                +7 (495) 123-45-67
              </a>
            </div>
            <div>
              <strong className="block text-gray-700">
                {t('contactInfo.email')}:
              </strong>
              <a
                href="mailto:info@example.com"
                className="text-blue-600 hover:underline"
              >
                info@example.com
              </a>
            </div>
            <div>
              <strong className="block text-gray-700">
                {t('contactInfo.address')}:
              </strong>
              <p className="text-gray-600">{t('contactInfo.addressText')}</p>
            </div>
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-4">
            {t('contactForm.title')}
          </h2>
          <form className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-gray-700 mb-1">
                {t('contactForm.name')}
              </label>
              <input
                type="text"
                id="name"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('contactForm.namePlaceholder')}
              />
            </div>
            <div>
              <label htmlFor="email" className="block text-gray-700 mb-1">
                {t('contactForm.email')}
              </label>
              <input
                type="email"
                id="email"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('contactForm.emailPlaceholder')}
              />
            </div>
            <div>
              <label htmlFor="message" className="block text-gray-700 mb-1">
                {t('contactForm.message')}
              </label>
              <textarea
                id="message"
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder={t('contactForm.messagePlaceholder')}
              ></textarea>
            </div>
            <button
              type="submit"
              className="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 transition-colors"
            >
              {t('contactForm.submit')}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
