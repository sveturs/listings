'use client';

import { ReactNode, useEffect } from 'react';
import { NextIntlClientProvider } from 'next-intl';
import { usePathname } from 'next/navigation';
import { preloadModules } from '@/lib/i18n/loadMessages';

interface ModularIntlProviderProps {
  locale: string;
  messages: Record<string, any>;
  children: ReactNode;
}

/**
 * Провайдер для модульной системы переводов
 * Автоматически предзагружает модули для следующих страниц
 */
export function ModularIntlProvider({
  locale,
  messages,
  children,
}: ModularIntlProviderProps) {
  const pathname = usePathname();

  // Debug: проверяем наличие chat переводов
  if (typeof window !== 'undefined' && pathname.includes('/chat')) {
    console.log('[ModularIntlProvider] Chat translations check:', {
      hasChatNamespace: !!messages.chat,
      hasTranslationKey: !!messages.chat?.translation,
      hasTranslatingKey: !!messages.chat?.translation?.translating,
      translatingValue: messages.chat?.translation?.translating,
      chatKeys: messages.chat ? Object.keys(messages.chat).slice(0, 10) : [],
    });
  }

  useEffect(() => {
    // Предзагружаем модули для возможных переходов
    const preloadNextModules = async () => {
      // Определяем вероятные следующие модули на основе текущей страницы
      const nextModules = [];

      if (pathname.includes('/marketplace')) {
        nextModules.push('cart', 'chat');
      } else if (pathname.includes('/chat')) {
        // Для страницы чата загружаем модуль chat
        nextModules.push('chat');
      } else if (pathname.includes('/admin')) {
        nextModules.push('marketplace', 'storefront');
      } else if (pathname.includes('/map')) {
        nextModules.push('map');
      } else if (pathname === '/') {
        nextModules.push('marketplace', 'auth', 'chat');
      }

      if (nextModules.length > 0) {
        await preloadModules(locale as any, nextModules as any);
      }
    };

    // Предзагружаем с задержкой, чтобы не блокировать основную загрузку
    const timer = setTimeout(preloadNextModules, 1000);

    return () => clearTimeout(timer);
  }, [pathname, locale]);

  return (
    <NextIntlClientProvider
      locale={locale}
      messages={messages}
      timeZone="Europe/Belgrade"
      now={new Date()}
    >
      {children}
    </NextIntlClientProvider>
  );
}
