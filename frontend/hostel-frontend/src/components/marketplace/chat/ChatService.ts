import axios from '../../../api/axios';

// Определяем интерфейсы для сообщений чата
export interface ChatMessage {
  id?: number;
  chat_id?: number;
  listing_id: number;
  sender_id?: number;
  receiver_id: number;
  content: string;  // Основное поле для содержимого сообщения - используется бэкендом
  message?: string;  // Устаревшее поле, оставлено для совместимости со старым кодом
  created_at?: string;
  is_read?: boolean;
  client_message_id?: string; // Уникальный клиентский идентификатор для дедупликации
  [key: string]: any;
}

export interface WebSocketMessage {
  type: string;
  [key: string]: any;
}

export interface ChatItem {
  id: number;
  listing_id: number;
  title?: string;
  user_id: number;
  participant_id: number;
  [key: string]: any;
}

type MessageHandler = (message: any) => void;
type ChatListHandler = (chats: ChatItem[]) => void;

class ChatService {
  private userId: number;
  private ws: WebSocket | null;
  private messageHandlers: Set<MessageHandler>;
  private reconnectAttempts: number;
  private maxReconnectAttempts: number;
  private reconnectTimer: any | null;
  private isConnecting: boolean;
  private isActive: boolean;
  private chatListHandlers: Set<ChatListHandler>;
  private lastPingTime: number;
  private pingInterval: any | null;
  // Для отслеживания изменений в количестве непрочитанных сообщений
  private _lastUnreadCount: number | null = null;

  constructor(userId: number) {
    this.userId = userId;
    this.ws = null;
    this.messageHandlers = new Set();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10; // Увеличено максимальное количество попыток
    this.reconnectTimer = null;
    this.isConnecting = false;
    this.isActive = true;
    this.chatListHandlers = new Set();
    this.lastPingTime = 0;
    this.pingInterval = null;
  }

  connect(): void {
    // Проверяем, активен ли сервис и нет ли уже активного подключения
    if (!this.isActive || this.ws?.readyState === WebSocket.OPEN || this.ws?.readyState === WebSocket.CONNECTING || this.isConnecting) {
      return;
    }

    this.isConnecting = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    try {
      // Используем window.ENV для WebSocket URL если доступно
      let wsUrl: string;
      if ((window as any).ENV && (window as any).ENV.REACT_APP_WS_URL) {
        wsUrl = (window as any).ENV.REACT_APP_WS_URL;
      } else {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = process.env.NODE_ENV === 'development' ? 'localhost:3000' : window.location.host;
        wsUrl = `${protocol}//${host}/ws/chat`;
      }

      console.log('Попытка подключения к WebSocket:', wsUrl);
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        if (!this.isActive) {
          this.ws?.close();
          return;
        }

        console.log('WebSocket соединение установлено');
        this.isConnecting = false;
        this.reconnectAttempts = 0;

        // Проверяем состояние ws перед отправкой, чтобы избежать ошибки
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
          this.ws.send(JSON.stringify({
            type: 'auth',
            user_id: this.userId
          }));
        } else {
          console.warn('WebSocket not ready for auth message, will retry');
          // Планируем повторную попытку через небольшой интервал
          setTimeout(() => {
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
              console.log('Retry sending auth message');
              this.ws.send(JSON.stringify({
                type: 'auth',
                user_id: this.userId
              }));
            }
          }, 500);
        }
        
        // Устанавливаем интервал для отправки ping-сообщений
        this.startPingInterval();
      };

      this.ws.onmessage = (event: MessageEvent) => {
        if (!this.isActive || !this.ws) return;

        try {
          const message = JSON.parse(event.data);
          
          // Обработка pong-сообщений
          if (message.type === 'pong') {
            console.log('Получен pong от сервера');
            this.lastPingTime = Date.now(); // Обновляем время последнего ping/pong
            return;
          }
          
          this.messageHandlers.forEach(handler => handler(message));
        } catch (error) {
          console.error('Ошибка обработки сообщения:', error);
        }
      };

      this.ws.onerror = (error: Event) => {
        if (!this.isActive) return;
        console.error('WebSocket ошибка:', error);
        this.isConnecting = false;

        // Переподключаемся только если соединение было разорвано
        if (this.ws?.readyState === WebSocket.CLOSED) {
          this.handleReconnect();
        }
      };

      this.ws.onclose = (event: CloseEvent) => {
        if (!this.isActive) return;
        console.log('WebSocket соединение закрыто:', event.code, event.reason);
        this.isConnecting = false;
        this.ws = null;

        // Переподключаемся только при определенных кодах закрытия
        if (event.code === 1006 || event.code === 1001 || event.code === 1011) {
          this.handleReconnect();
        }
      };

    } catch (error) {
      console.error('Ошибка при создании WebSocket:', error);
      this.isConnecting = false;
      if (this.isActive) {
        this.handleReconnect();
      }
    }
  }

  private handleReconnect(): void {
    if (!this.isActive || this.reconnectAttempts >= this.maxReconnectAttempts || this.isConnecting) {
      return;
    }

    // Остановим ping-интервал при переподключении
    this.stopPingInterval();

    // Увеличиваем задержку только до определенного предела и линейный рост вместо экспоненциального
    const delay = Math.min(2000 + (this.reconnectAttempts * 1000), 15000);
    console.log(`Попытка переподключения через ${delay}ms (попытка ${this.reconnectAttempts + 1})`);

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++;
      
      // Если WebSocket все еще существует, закрываем его перед новым подключением
      if (this.ws) {
        try {
          this.ws.close();
        } catch (e) {
          console.error('Ошибка при закрытии старого WebSocket:', e);
        }
        this.ws = null;
      }
      
      this.connect();
    }, delay);
  }

  private startPingInterval(): void {
    // Очищаем предыдущий интервал, если он существует
    this.stopPingInterval();
    
    // Установка интервала для отправки ping каждые 30 секунд
    this.pingInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          console.log('Отправка ping');
          // Проверяем состояние ws перед отправкой пинга
          if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ type: 'ping', timestamp: Date.now() }));
            this.lastPingTime = Date.now();
          } else {
            console.debug('WebSocket not ready for ping, will reconnect');
            this.handleReconnect();
          }
          
          // Проверяем, получен ли ответ на предыдущий ping в течение 15 секунд
          setTimeout(() => {
            // Если последний ping был давно и не получен ответ
            const pingTimeout = 15000; // 15 секунд
            if (Date.now() - this.lastPingTime > pingTimeout && this.ws?.readyState === WebSocket.OPEN) {
              console.log('Не получен ответ на ping, переподключение...');
              // Закрываем соединение и пытаемся переподключиться
              try {
                this.ws.close(1000);
              } catch (err) {
                console.error('Ошибка при закрытии неотвечающего WebSocket:', err);
              }
              this.ws = null;
              this.handleReconnect();
            }
          }, 15000);
          
        } catch (e) {
          console.error('Ошибка при отправке ping:', e);
          this.handleReconnect();
        }
      } else {
        // Если соединение закрыто, пытаемся переподключиться
        this.handleReconnect();
      }
    }, 30000); // 30 секунд
  }
  
  private stopPingInterval(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  disconnect(): void {
    console.log('Отключение WebSocket...');
    this.isActive = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.stopPingInterval();

    if (this.ws) {
      try {
        // Отправляем "чистое" закрытие
        this.ws.close(1000, 'Нормальное закрытие');
      } catch (e) {
        console.error('Ошибка при закрытии WebSocket:', e);
      }
      this.ws = null;
    }

    this.messageHandlers.clear();
    this.reconnectAttempts = 0;
    this.isConnecting = false;
  }

  async sendMessage(message: ChatMessage): Promise<any> {
    try {
      // Проверяем наличие обязательных полей
      if (!message.chat_id || Number(message.chat_id) <= 0) {
        console.error(`Invalid chat_id in message: ${message.chat_id}`);
        throw new Error(`Некорректный ID чата: ${message.chat_id}`);
      }

      if (!message.listing_id || Number(message.listing_id) <= 0) {
        console.error(`Invalid listing_id in message: ${message.listing_id}`);
        throw new Error(`Некорректный ID объявления: ${message.listing_id}`);
      }

      if (!message.receiver_id || Number(message.receiver_id) <= 0) {
        console.error(`Invalid receiver_id in message: ${message.receiver_id}`);
        throw new Error(`Некорректный ID получателя: ${message.receiver_id}`);
      }

      // Проверяем наличие текста сообщения
      if (!message.content || message.content.trim() === '') {
        console.error('Message content is empty');
        throw new Error('Текст сообщения не может быть пустым');
      }

      // Генерируем уникальный клиентский ID для сообщения, если его еще нет
      const clientMessageId = message.client_message_id ||
        `client_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

      // Подготавливаем сообщение для отправки на сервер
      const messageToSend = {
        chat_id: Number(message.chat_id),
        listing_id: Number(message.listing_id),
        receiver_id: Number(message.receiver_id),
        content: message.content.trim(),
        client_message_id: clientMessageId  // Добавляем клиентский ID
      };

      // Отладочное сообщение только в режиме разработки
      if (process.env.NODE_ENV === 'development') {
        console.debug('Preparing to send message');
      }

      // ИЗМЕНЕНО: Используем WebSocket как основной канал, HTTP как резервный
      let response;

      // Сначала проверяем, что WebSocket существует и соединение полностью установлено
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          // Отладочное сообщение
          if (process.env.NODE_ENV === 'development') {
            console.debug('Sending message via WebSocket');
          }
          this.ws.send(JSON.stringify({
            type: 'message',
            ...messageToSend
          }));

          // Создаем базовый объект ответа для оптимистичного UI
          response = {
            data: {
              success: true,
              message: "Message sent via WebSocket",
              client_message_id: clientMessageId
            },
            method: 'websocket'
          };

          // Не отправляем HTTP запрос, если WebSocket доступен и отправка прошла успешно
          return response;
        } catch (wsError) {
          // Если отправка через WebSocket не удалась, логируем и продолжаем с HTTP
          console.error('WebSocket send failed, falling back to HTTP:', wsError);
        }
      }

      // Отправляем через HTTP только если WebSocket недоступен или произошла ошибка
      // Отладочное сообщение
      if (process.env.NODE_ENV === 'development') {
        console.debug('Sending message via HTTP (fallback)');
      }
      const httpResponse = await axios.post('/api/v1/marketplace/chat/messages', messageToSend);

      // Если WebSocket был недоступен, пытаемся переподключиться
      if (this.isActive && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.connect();
      }

      // Отладочное сообщение
      if (process.env.NODE_ENV === 'development') {
        console.debug('HTTP response received');
      }
      return {
        ...httpResponse.data,
        client_message_id: clientMessageId,
        method: 'http'
      };

    } catch (error) {
      console.error('Ошибка отправки сообщения:', error);
      throw error;
    }
  }

  onChatListUpdate(handler: ChatListHandler): () => void {
    this.chatListHandlers.add(handler);
    return () => this.chatListHandlers.delete(handler);
  }

  // Вспомогательный метод для получения общего количества непрочитанных сообщений
  async getTotalUnreadCount(): Promise<number> {
    try {
      const response = await axios.get('/api/v1/marketplace/chat');
      const chats: ChatItem[] = response.data?.data || [];
      return chats.reduce((total, chat) => total + (chat.unread_count || 0), 0);
    } catch (error) {
      console.error('Error getting total unread count:', error);
      return 0;
    }
  }

  async updateChatsList(): Promise<void> {
    try {
      // Полностью отключаем это сообщение - вызывается слишком часто

      const response = await axios.get('/api/v1/marketplace/chat');
      const chats: ChatItem[] = response.data?.data || [];

      // Подсчет общего количества непрочитанных сообщений
      const totalUnread = chats.reduce((sum, chat) => sum + (chat.unread_count || 0), 0);

      // Только обновляем внутреннее состояние, без логирования
      this._lastUnreadCount = totalUnread;

      // Уведомляем всех подписчиков об обновлении списка чатов
      this.chatListHandlers.forEach(handler => handler(chats));
    } catch (error) {
      console.error('Error updating chats list:', error);
      throw error; // Пробрасываем ошибку для возможности обработки выше
    }
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => {
      this.messageHandlers.delete(handler);
    };
  }

  async getMessageHistory(chatId: number, listingId: number, limit: number = 100): Promise<ChatMessage[]> {
    if (!chatId || chatId <= 0) {
      console.error(`Invalid chat ID for message history: ${chatId}`);
      throw new Error(`Некорректный ID чата: ${chatId}`);
    }

    if (!listingId || listingId <= 0) {
      console.error(`Invalid listing ID for message history: ${listingId}`);
      throw new Error(`Некорректный ID объявления: ${listingId}`);
    }

    try {
      // Отладочное сообщение
      if (process.env.NODE_ENV === 'development') {
        console.debug(`Fetching message history with limit=${limit}`);
      }
      const response = await axios.get('/api/v1/marketplace/chat/messages', {
        params: {
          chat_id: chatId,
          listing_id: listingId,
          limit: limit // Добавляем параметр limit для получения большего количества сообщений
        }
      });

      if (!response.data) {
        console.error('Empty response when fetching message history');
        return [];
      }

      if (!response.data.data) {
        console.error('Missing data field in response', response.data);
        return [];
      }

      if (!Array.isArray(response.data.data)) {
        console.error('Response data is not an array', response.data.data);
        return [];
      }

      // Фильтруем сообщения с пустым содержимым
      let messagesArray = response.data.data
        .filter((msg: any) => msg && (msg.content || msg.message)); // Убедимся, что у сообщения есть хоть какой-то контент

      // УЛУЧШЕННАЯ ДЕДУПЛИКАЦИЯ: используем Map для хранения уникальных сообщений
      // Сначала дедуплицируем по ID и client_message_id
      const uniqueMessages = new Map<string, any>();

      messagesArray.forEach((msg: any) => {
        // Генерируем ключи для различных стратегий дедупликации
        const serverIdKey = msg.id ? `id:${msg.id}` : null;
        const clientIdKey = msg.client_message_id ? `client:${msg.client_message_id}` : null;

        // Сначала проверяем серверный ID
        if (serverIdKey) {
          // Если уже есть сообщение с таким ID, берем более новое
          if (!uniqueMessages.has(serverIdKey) ||
              new Date(msg.created_at || '').getTime() >
              new Date(uniqueMessages.get(serverIdKey).created_at || '').getTime()) {
            uniqueMessages.set(serverIdKey, msg);
          }
        }

        // Затем проверяем клиентский ID, если серверного ID нет
        else if (clientIdKey) {
          if (!uniqueMessages.has(clientIdKey)) {
            uniqueMessages.set(clientIdKey, msg);
          }
        }

        // Если нет ни серверного, ни клиентского ID, используем комбинацию полей
        else {
          const contentTimeKey = `${msg.sender_id}:${msg.content || msg.message}:${new Date(msg.created_at || '').getTime()}`;
          if (!uniqueMessages.has(contentTimeKey)) {
            uniqueMessages.set(contentTimeKey, msg);
          }
        }
      });

      // Получаем уникальные сообщения из Map
      const messagesAfterIdDedup = Array.from(uniqueMessages.values());

      // Теперь дедуплицируем по содержанию и отправителю с расширенным временным окном в 10 секунд
      const messagesByContent = new Map<string, any[]>();

      // Группируем сообщения по содержанию и отправителю
      messagesAfterIdDedup.forEach((msg: any) => {
        const contentSenderKey = `${msg.sender_id}:${msg.content || msg.message}`;
        if (!messagesByContent.has(contentSenderKey)) {
          messagesByContent.set(contentSenderKey, []);
        }
        messagesByContent.get(contentSenderKey)!.push(msg);
      });

      // Дедуплицируем сообщения в пределах временного окна
      const dedupedMessages: any[] = [];

      messagesByContent.forEach((messages: any[]) => {
        // Сортируем сообщения по времени
        messages.sort((a, b) =>
          new Date(a.created_at || '').getTime() - new Date(b.created_at || '').getTime()
        );

        // Добавляем только одно сообщение из группы, если они находятся в пределах 10 секунд друг от друга
        let lastAddedTime = 0;

        messages.forEach((msg: any) => {
          const msgTime = new Date(msg.created_at || '').getTime();

          if (lastAddedTime === 0 || msgTime - lastAddedTime > 10000) {
            dedupedMessages.push(msg);
            lastAddedTime = msgTime;
          } else {
            console.log('Removing duplicate message from history (time window):', msg);
          }
        });
      });

      // Преобразуем сообщения и сортируем их
      const messages: ChatMessage[] = dedupedMessages
        .map((msg: any) => ({
          ...msg,
          content: msg.content || msg.message || '' // Убедимся, что у всех сообщений есть поле content
        }))
        .sort((a: ChatMessage, b: ChatMessage) => {
          const timeA = new Date(a.created_at || '').getTime();
          const timeB = new Date(b.created_at || '').getTime();
          return timeA - timeB;
        });

      // Выводим информацию о дедупликации
      const originalCount = response.data.data.length;
      const afterIdDedup = messagesAfterIdDedup.length;
      const finalCount = messages.length;

      console.log(`Message history stats:
        - Original count: ${originalCount}
        - After ID/client_id deduplication: ${afterIdDedup}
        - After content/time deduplication: ${finalCount}
        - Removed ${originalCount - finalCount} duplicate messages
      `);

      return messages;
    } catch (error) {
      console.error('Ошибка получения истории сообщений:', error);
      throw error; // Выбрасываем ошибку для обработки на верхнем уровне
    }
  }

  async markMessagesAsRead(messageIds: number[]): Promise<boolean> {
    if (!messageIds || !Array.isArray(messageIds) || messageIds.length === 0) {
      console.warn('No message IDs provided for marking as read');
      return false;
    }

    // Проверяем валидность всех ID
    const validIds = messageIds.filter(id => id && id > 0);
    if (validIds.length === 0) {
      console.warn('No valid message IDs found for marking as read', messageIds);
      return false;
    }

    try {
      // Полностью отключаем это сообщение - вызывается слишком часто
      // Логирование будет только в методе handleSelectChat в ChatPage.tsx
      const response = await axios.put('/api/v1/marketplace/chat/messages/read', {
        message_ids: validIds
      });

      const success = response.data?.success === true;

      // Полностью отключаем это сообщение - вызывается слишком часто

      // После успешной отметки сообщений как прочитанных, запрашиваем
      // текущее количество непрочитанных сообщений для диагностики
      if (success) {
        try {
          // Вместо запроса к API используем информацию из chats
          // Получаем текущее количество непрочитанных сообщений, но не логируем
          const totalUnread = await this.getTotalUnreadCount();
          this._lastUnreadCount = totalUnread;
        } catch (err) {
          // Только отладочное логирование для диагностического действия
          console.debug('Info: Could not calculate unread count, continuing anyway');
        }
      }

      // Если успешно отметили сообщения как прочитанные, делаем дополнительное
      // обновление списка чатов для актуализации счетчиков во всем приложении
      if (success) {
        try {
          // Обновляем список чатов трижды с разными задержками,
          // чтобы гарантировать обновление счетчика после обработки на сервере
          await this.updateChatsList();

          setTimeout(() => {
            this.updateChatsList().catch(err =>
              console.error('Error updating chats list (500ms delay):', err)
            );
          }, 500);

          setTimeout(() => {
            this.updateChatsList().catch(err =>
              console.error('Error updating chats list (1500ms delay):', err)
            );
          }, 1500);

          setTimeout(() => {
            this.updateChatsList().catch(err =>
              console.error('Error updating chats list (3000ms delay):', err)
            );
          }, 3000);
        } catch (updateError) {
          console.error('Error updating chats list after marking as read:', updateError);
          // Игнорируем ошибку обновления, т.к. сообщения уже отмечены как прочитанные
        }
      }

      return success;
    } catch (error) {
      console.error('Ошибка отметки сообщений как прочитанных:', error);
      return false;
    }
  }

  handleNewMessage = async (message: ChatMessage): Promise<void> => {
    // Сначала вызываем обработчики для обработки сообщения
    this.messageHandlers.forEach(handler => handler(message));

    // Затем обновляем список чатов для актуализации счетчиков
    try {
      // Первое обновление - немедленно
      await this.updateChatsList();

      // Отметка сообщения как прочитанное, если оно открыто (реализуется в компоненте чата)

      // Обновления с задержкой для гарантии актуальных счетчиков
      // Несколько попыток с разными задержками для надежности
      const delays = [500, 1500, 3000]; // задержки в мс

      for (const delay of delays) {
        setTimeout(() => {
          this.updateChatsList().catch(err =>
            console.error(`Error updating chats list in delayed update (${delay}ms):`, err)
          );
        }, delay);
      }
    } catch (error) {
      console.error('Error updating chats list after new message:', error);

      // Всё равно пытаемся обновить, даже при ошибке первого обновления
      setTimeout(() => {
        this.updateChatsList().catch(err =>
          console.error('Error in recovery update of chats list:', err)
        );
      }, 2000);
    }
  }
}

export default ChatService;