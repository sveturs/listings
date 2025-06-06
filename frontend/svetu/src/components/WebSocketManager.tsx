'use client';

import { useEffect, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useChat } from '@/hooks/useChat';
import { tokenManager } from '@/utils/tokenManager';

export default function WebSocketManager() {
  const { user, isAuthenticated } = useAuth();
  const { initWebSocket, closeWebSocket, loadChats } = useChat();
  const isInitialized = useRef(false);
  const initTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Очищаем предыдущий таймаут
    if (initTimeoutRef.current) {
      clearTimeout(initTimeoutRef.current);
      initTimeoutRef.current = null;
    }

    // Инициализируем WebSocket только если пользователь авторизован
    if (isAuthenticated && user && !isInitialized.current) {
      // Небольшая задержка для того, чтобы токен успел сохраниться
      initTimeoutRef.current = setTimeout(() => {
        const hasToken = tokenManager.getAccessToken() !== null;

        if (hasToken && !isInitialized.current) {
          console.log(
            '[WebSocketManager] Initializing WebSocket for user:',
            user.id,
            'Token available:',
            hasToken
          );
          isInitialized.current = true;
          // Инициализируем WebSocket с функцией получения ID пользователя
          initWebSocket(() => user.id);
          // Загружаем чаты только если пользователь авторизован
          loadChats();
        } else {
          console.log('[WebSocketManager] Waiting for token, will retry...');
        }
      }, 500); // Задержка 500мс
    }

    // Отключаем WebSocket при выходе
    if (!isAuthenticated && isInitialized.current) {
      console.log('[WebSocketManager] User logged out, cleaning up WebSocket');
      isInitialized.current = false;
      closeWebSocket();
    }

    // Cleanup при размонтировании компонента
    return () => {
      if (initTimeoutRef.current) {
        clearTimeout(initTimeoutRef.current);
        initTimeoutRef.current = null;
      }
      if (isInitialized.current) {
        console.log(
          '[WebSocketManager] Component unmounting, cleaning up WebSocket'
        );
        closeWebSocket();
      }
    };
  }, [isAuthenticated, user, initWebSocket, closeWebSocket, loadChats]);

  // Компонент не рендерит ничего в UI
  return null;
}
