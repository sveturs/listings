import { notFound } from 'next/navigation';
import { getRequestConfig } from 'next-intl/server';
import { routing } from '@/i18n/routing';
import { loadMessages } from '@/lib/i18n/loadMessages';

export default getRequestConfig(async ({ locale }) => {
  // Validate that the incoming `locale` parameter is valid
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    notFound();
  }

  // Используем модульную систему
  // Для серверного рендеринга загружаем базовый набор модулей
  // На клиенте будут подгружены дополнительные модули по необходимости
  const baseModules = ['common', 'marketplace', 'auth', 'misc', 'cart'];
  const messages = await loadMessages(locale as any, baseModules as any);

  return {
    locale,
    messages,
  };
});
