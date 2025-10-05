'use client';

import { useEffect } from 'react';
import { useLocale } from 'next-intl';
import { chatService } from '@/services/chat';

/**
 * LocaleSync component - синхронизирует локаль приложения с настройками чата пользователя
 *
 * При смене языка приложения автоматически обновляет preferred_language в БД
 * через Auth Service metadata.
 *
 * Этот компонент должен быть добавлен в корневой layout для автоматической синхронизации.
 */
export default function LocaleSync() {
  const locale = useLocale();

  useEffect(() => {
    const syncLocale = async () => {
      try {
        // Получаем текущие настройки чата
        const currentSettings = await chatService.getChatSettings();

        // Синхронизируем локаль с сервером если она изменилась
        if (currentSettings.preferred_language !== locale) {
          await chatService.updateChatSettings({
            ...currentSettings,
            preferred_language: locale as 'ru' | 'en' | 'sr',
          });

          console.log(`[LocaleSync] Locale synced to server: ${locale}`);
        }
      } catch (error) {
        // Не показываем ошибку пользователю, просто логируем
        // Синхронизация локали не критична для работы приложения
        console.warn('[LocaleSync] Failed to sync locale:', error);
      }
    };

    syncLocale();
  }, [locale]);

  // Этот компонент не рендерит ничего
  return null;
}
