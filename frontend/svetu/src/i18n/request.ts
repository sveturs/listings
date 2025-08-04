import { getRequestConfig } from 'next-intl/server';
import { routing } from './routing';
import { loadMessages } from '@/lib/i18n/loadMessages';

export default getRequestConfig(async ({ requestLocale }) => {
  // This typically corresponds to the `[locale]` segment
  let locale = await requestLocale;

  // Ensure that a valid locale is used
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    locale = routing.defaultLocale;
  }

  // Загружаем базовые модули для серверного рендеринга
  const messages = await loadMessages(locale as any, [
    'common',
    'marketplace',
    'auth',
    'misc',
    'cart',
    'map',
  ]);

  return {
    locale,
    messages,
  };
});
