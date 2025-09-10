import { getRequestConfig } from 'next-intl/server';
import { routing } from './routing';
import { getMessages } from './messages';

export default getRequestConfig(async ({ requestLocale }) => {
  // This typically corresponds to the `[locale]` segment
  let locale = await requestLocale;

  // Ensure that a valid locale is used
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    locale = routing.defaultLocale;
  }

  // Используем статически импортированные сообщения для надёжной работы SSG
  const messages = getMessages(locale as 'en' | 'ru' | 'sr');

  return {
    locale,
    messages,
  };
});
