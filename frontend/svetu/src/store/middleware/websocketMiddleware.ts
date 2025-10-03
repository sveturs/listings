import { Middleware } from '@reduxjs/toolkit';
import { chatService } from '@/services/chat';
import {
  setWebSocket,
  handleNewMessage,
  handleMessageRead,
  handleUserOnline,
  handleUserOffline,
  setUserTyping,
  setCurrentUserId,
  updateMessageAttachments,
  refreshMessageWithAttachments,
} from '../slices/chatSlice';
import type { RootState } from '../index';

let ws: WebSocket | null = null;
let heartbeatInterval: NodeJS.Timeout | null = null;
let reconnectTimeout: NodeJS.Timeout | null = null;
let isConnecting = false;

export const websocketMiddleware: Middleware =
  (store: any) => (next) => (action: any) => {
    if (action.type === 'chat/initWebSocket') {
      const { getCurrentUserId } = action.payload;

      // Сохраняем ID текущего пользователя в store
      const currentUserId = getCurrentUserId();
      store.dispatch(setCurrentUserId(currentUserId));

      // Предотвращаем множественные соединения
      if (isConnecting) {
        console.log('[WebSocket] Connection already in progress, skipping...');
        return next(action);
      }

      // Если уже есть активное соединение, проверяем его состояние
      if (ws && ws.readyState === WebSocket.OPEN) {
        console.log('[WebSocket] Already connected, skipping initialization');
        return next(action);
      }

      // Закрываем существующее соединение если оно не активно
      if (ws && ws.readyState !== WebSocket.OPEN) {
        ws.close();
        ws = null;
      }

      // Очищаем таймауты
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        reconnectTimeout = null;
      }

      // Устанавливаем флаг подключения
      isConnecting = true;

      // Создаем новое подключение (async)
      chatService
        .connectWebSocket((event) => {
          try {
            const data = JSON.parse(event.data);
            if (process.env.NODE_ENV === 'development') {
              console.log('WebSocket message:', data);
            }

            switch (data.type) {
              case 'new_message':
                // Сообщение уже должно содержать вложения от сервера
                store.dispatch(handleNewMessage(data.payload));

                // Если сообщение помечено как имеющее вложения, но вложения не пришли
                // это может означать, что файлы еще загружаются
                if (data.payload.has_attachments && !data.payload.attachments) {
                  // Ждем немного и пытаемся загрузить сообщение с вложениями
                  setTimeout(() => {
                    store.dispatch(
                      refreshMessageWithAttachments({
                        chatId: data.payload.chat_id,
                        messageId: data.payload.id,
                      })
                    );
                  }, 3000); // Ждем 3 секунды, чтобы дать время файлам загрузиться
                }
                break;

              case 'message_read':
                store.dispatch(handleMessageRead(data.payload));
                break;

              case 'user_typing':
                store.dispatch(
                  setUserTyping({
                    chatId: data.payload.chat_id,
                    userId: data.payload.user_id,
                    isTyping: data.payload.is_typing,
                  })
                );
                break;

              case 'user_online':
                store.dispatch(
                  handleUserOnline({ user_id: data.payload.user_id })
                );
                break;

              case 'user_offline':
                store.dispatch(
                  handleUserOffline({
                    user_id: data.payload.user_id,
                    last_seen: data.payload.last_seen,
                  })
                );
                break;

              case 'attachment_upload':
                // Когда другой пользователь загрузил вложение
                // Это событие должно содержать полную информацию о сообщении с вложениями
                if (
                  data.payload &&
                  data.payload.message_id &&
                  data.payload.attachments
                ) {
                  store.dispatch(
                    updateMessageAttachments({
                      messageId: data.payload.message_id,
                      chatId: data.payload.chat_id,
                      attachments: data.payload.attachments,
                    })
                  );
                }
                break;

              case 'attachment_delete':
                // Обрабатываем удаление вложений если нужно
                break;

              case 'online_users_list':
                // Обработка списка онлайн пользователей
                if (data.payload && data.payload.users) {
                  data.payload.users.forEach((userId: number) => {
                    store.dispatch(handleUserOnline({ user_id: userId }));
                  });
                }
                break;

              default:
                if (process.env.NODE_ENV === 'development') {
                  console.log('Unknown WebSocket message type:', data.type);
                }
            }
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
          }
        })
        .then((websocket) => {
          ws = websocket;

          if (!ws) {
            console.warn(
              '[WebSocket] Connection not created, possibly no auth token'
            );
            isConnecting = false;
            return;
          }

          store.dispatch(setWebSocket(ws));

          // Запрашиваем статус всех пользователей при подключении
          ws.addEventListener('open', () => {
            isConnecting = false; // Сбрасываем флаг подключения
            if (process.env.NODE_ENV === 'development') {
              console.log('WebSocket connected, requesting user statuses');
            }

            // Получаем ID всех пользователей из чатов
            const state = store.getState() as RootState;
            const userIds = new Set<number>();

            state.chat.chats.forEach((chat) => {
              if (chat.buyer) userIds.add(chat.buyer.id);
              if (chat.seller) userIds.add(chat.seller.id);
              if (chat.other_user) userIds.add(chat.other_user.id);
            });

            // Запрашиваем статус для каждого пользователя
            userIds.forEach((userId) => {
              if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(
                  JSON.stringify({
                    type: 'get_user_status',
                    payload: { user_id: userId },
                  })
                );
              }
            });

            // Heartbeat для поддержания соединения
            heartbeatInterval = setInterval(() => {
              if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({ type: 'ping' }));
              }
            }, 30000);
          });

          ws.addEventListener('close', () => {
            isConnecting = false; // Сбрасываем флаг при закрытии
            if (heartbeatInterval) {
              clearInterval(heartbeatInterval);
              heartbeatInterval = null;
            }
          });

          ws.addEventListener('error', () => {
            isConnecting = false; // Сбрасываем флаг при ошибке
            console.error('[WebSocket] Connection error occurred');
          });
        })
        .catch((error) => {
          console.error('[WebSocket] Failed to connect:', error);
          isConnecting = false;
        });
    }

    if (action.type === 'chat/closeWebSocket') {
      isConnecting = false; // Сбрасываем флаг
      if (ws) {
        ws.close();
        ws = null;
      }
      if (heartbeatInterval) {
        clearInterval(heartbeatInterval);
        heartbeatInterval = null;
      }
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        reconnectTimeout = null;
      }
      store.dispatch(setWebSocket(null));
    }

    return next(action);
  };
