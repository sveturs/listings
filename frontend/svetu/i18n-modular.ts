import { notFound } from 'next/navigation';
import { getRequestConfig } from 'next-intl/server';
import { routing } from '@/i18n/routing';
import { loadMessages, getRequiredModules } from '@/lib/i18n/loadMessages';

export default getRequestConfig(async ({ locale, requestLocale }) => {
  // Validate that the incoming `locale` parameter is valid
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    notFound();
  }

  // Получаем путь запроса для определения необходимых модулей
  const pathname = (await requestLocale)?.pathname || '/';
  
  // Определяем какие модули нужны для текущей страницы
  const requiredModules = getRequiredModules(pathname);
  
  // Загружаем только необходимые модули
  const messages = await loadMessages(
    locale as 'ru' | 'en' | 'sr', 
    requiredModules
  );

  return {
    locale,
    messages,
  };
});