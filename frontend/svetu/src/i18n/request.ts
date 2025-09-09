import { getRequestConfig } from 'next-intl/server';
import { routing } from './routing';

export default getRequestConfig(async ({ requestLocale }) => {
  // This typically corresponds to the `[locale]` segment
  let locale = await requestLocale;

  // Ensure that a valid locale is used
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    locale = routing.defaultLocale;
  }

  // Для серверного рендеринга загружаем переводы напрямую
  // Это более надёжно работает при SSG
  const messages: Record<string, any> = {};

  // Загружаем базовые модули
  const modules = [
    'common',
    'marketplace',
    'auth',
    'misc',
    'cart',
    'map',
    'storefronts',
    'admin',
    'cars',
    'reviews',
  ];

  for (const mod of modules) {
    try {
      const moduleData = await import(`@/messages/${locale}/${mod}.json`);
      // JSON импорты могут иметь default в некоторых случаях
      const data = moduleData.default || moduleData;

      // Добавляем модуль как namespace
      messages[mod] = data;

      // Также добавляем все ключи на верхний уровень для обратной совместимости
      Object.keys(data).forEach((key) => {
        if (!messages[key]) {
          messages[key] = data[key];
        }
      });
    } catch (error) {
      console.error(`Failed to load ${mod} for ${locale}:`, error);
    }
  }

  return {
    locale,
    messages,
  };
});
