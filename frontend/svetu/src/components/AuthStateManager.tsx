'use client';

import { useEffect, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useChat } from '@/hooks/useChat';

export default function AuthStateManager() {
  const { user } = useAuth();
  const { clearChatData } = useChat();
  const previousUserId = useRef<number | null>(null);

  useEffect(() => {
    // Если пользователь изменился (включая logout)
    if (
      previousUserId.current !== null &&
      previousUserId.current !== user?.id
    ) {
      console.log('[AuthStateManager] User changed, clearing chat data');
      clearChatData();

      // Очищаем только данные чата из sessionStorage, оставляя токены
      try {
        // Сохраняем важные данные, которые не должны быть удалены
        const locale = localStorage.getItem('NEXT_LOCALE');

        // Очищаем всё кроме токенов и локали
        const keysToRemove: string[] = [];
        for (let i = 0; i < sessionStorage.length; i++) {
          const key = sessionStorage.key(i);
          if (
            key &&
            !['svetu_access_token', 'svetu_user', 'client_id'].includes(key)
          ) {
            keysToRemove.push(key);
          }
        }

        keysToRemove.forEach((key) => sessionStorage.removeItem(key));

        // Для localStorage очищаем всё кроме локали и корзины
        const cart = localStorage.getItem('svetu_cart');
        localStorage.clear();
        if (locale) {
          localStorage.setItem('NEXT_LOCALE', locale);
        }
        if (cart) {
          localStorage.setItem('svetu_cart', cart);
        }
      } catch (error) {
        console.error('Error clearing storage:', error);
      }
    }

    // Обновляем предыдущего пользователя
    previousUserId.current = user?.id || null;
  }, [user, clearChatData]);

  return null;
}
