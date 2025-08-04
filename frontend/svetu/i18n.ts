import { notFound } from 'next/navigation';
import { getRequestConfig } from 'next-intl/server';
import { routing } from '@/i18n/routing';

// Временное решение: пока используем старую систему
// TODO: Полностью перейти на модульную систему после миграции всех переводов
const USE_MODULAR_SYSTEM = process.env.USE_MODULAR_I18N === 'true';

export default getRequestConfig(async ({ locale }) => {
  // Validate that the incoming `locale` parameter is valid
  if (!locale || !routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    notFound();
  }

  if (USE_MODULAR_SYSTEM) {
    // Новая модульная система
    const { loadMessages, getRequiredModules } = await import('@/lib/i18n/loadMessages');
    
    // Для серверного рендеринга загружаем базовый набор модулей
    // На клиенте будут подгружены дополнительные модули по необходимости
    const baseModules = ['common', 'marketplace', 'auth'];
    const messages = await loadMessages(locale as any, baseModules as any);
    
    return {
      locale,
      messages,
    };
  } else {
    // Старая система - загружаем весь файл
    return {
      locale,
      messages: (await import(`./src/messages/${locale}.json`)).default,
    };
  }
});