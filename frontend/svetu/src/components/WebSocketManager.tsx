'use client';

import { useEffect, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useChat } from '@/hooks/useChat';
import { tokenManager } from '@/utils/tokenManager';

// Глобальный флаг для предотвращения множественных инициализаций
let globalWebSocketInitialized = false;

export default function WebSocketManager() {
  const { user, isAuthenticated } = useAuth();
  const { initWebSocket, closeWebSocket, loadChats } = useChat();
  const localInitRef = useRef(false);
  const initTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const userIdRef = useRef<number | null>(null);

  useEffect(() => {
    // Очищаем предыдущий таймаут
    if (initTimeoutRef.current) {
      clearTimeout(initTimeoutRef.current);
      initTimeoutRef.current = null;
    }

    // Инициализируем WebSocket только если пользователь авторизован
    if (isAuthenticated && user) {
      // Проверяем, изменился ли пользователь
      const userChanged =
        userIdRef.current !== null && userIdRef.current !== user.id;
      userIdRef.current = user.id;

      // Если пользователь изменился, закрываем старое соединение
      if (userChanged && localInitRef.current) {
        globalWebSocketInitialized = false;
        localInitRef.current = false;
        closeWebSocket();
      }

      // Инициализируем только если еще не инициализировано
      if (!globalWebSocketInitialized && !localInitRef.current) {
        // Небольшая задержка для того, чтобы токен успел сохраниться
        initTimeoutRef.current = setTimeout(() => {
          const hasToken = tokenManager.getAccessToken() !== null;

          if (
            hasToken &&
            !globalWebSocketInitialized &&
            !localInitRef.current
          ) {
            globalWebSocketInitialized = true;
            localInitRef.current = true;
            // Инициализируем WebSocket с функцией получения ID пользователя
            initWebSocket(() => user.id);
            // Загружаем чаты только если пользователь авторизован
            loadChats();
          }
        }, 500); // Задержка 500мс
      }
    }

    // Отключаем WebSocket при выходе
    if (!isAuthenticated && localInitRef.current) {
      globalWebSocketInitialized = false;
      localInitRef.current = false;
      userIdRef.current = null;
      closeWebSocket();
    }

    // Cleanup при размонтировании компонента
    return () => {
      if (initTimeoutRef.current) {
        clearTimeout(initTimeoutRef.current);
        initTimeoutRef.current = null;
      }
      // При размонтировании сбрасываем только локальный флаг
      // глобальный флаг сбрасывается только при выходе или смене пользователя
      if (localInitRef.current) {
        localInitRef.current = false;
        // Не закрываем WebSocket при размонтировании, если компонент может быть перемонтирован
      }
    };
  }, [isAuthenticated, user, initWebSocket, closeWebSocket, loadChats]);

  // Компонент не рендерит ничего в UI
  return null;
}
