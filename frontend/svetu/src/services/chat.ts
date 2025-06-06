import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';
import {
  MarketplaceChat,
  MarketplaceMessage,
  SendMessagePayload,
  GetMessagesParams,
  MarkMessagesReadPayload,
  ChatListResponse,
  MessagesResponse,
  ChatAttachment,
} from '@/types/chat';

class ChatService {
  private baseUrl: string;
  private csrfToken: string | null = null;
  private reconnectAttempts = 0;

  constructor() {
    this.baseUrl = `${configManager.getApiUrl()}/api/v1/marketplace/chat`;
  }

  private async getCsrfToken(): Promise<string> {
    if (this.csrfToken) {
      return this.csrfToken;
    }

    try {
      const response = await fetch(
        `${configManager.getApiUrl()}/api/v1/csrf-token`,
        {
          method: 'GET',
          credentials: 'include',
        }
      );

      if (response.ok) {
        const data = await response.json();
        this.csrfToken = data.csrf_token;
        return this.csrfToken || '';
      }
    } catch (error) {
      console.warn('Failed to fetch CSRF token:', error);
    }

    // Fallback: generate client-side token for basic protection
    this.csrfToken = `client-${Date.now()}-${Math.random().toString(36).substring(2)}`;
    return this.csrfToken;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    console.log('Chat API request:', url, options);

    // Получаем JWT токен
    const accessToken = tokenManager.getAccessToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options?.headers as Record<string, string>),
    };

    // Добавляем JWT токен если есть
    if (accessToken) {
      (headers as any)['Authorization'] = `Bearer ${accessToken}`;
    }

    // Добавляем CSRF токен для изменяющих запросов
    const method = options?.method || 'GET';
    if (['POST', 'PUT', 'DELETE', 'PATCH'].includes(method.toUpperCase())) {
      const csrfToken = await this.getCsrfToken();
      headers['X-CSRF-Token'] = csrfToken;
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: 'include',
    });

    console.log('Chat API response status:', response.status);

    if (!response.ok) {
      let error: { message?: string; error?: string } = { message: '' };
      const contentType = response.headers.get('content-type');

      if (contentType?.includes('application/json')) {
        error = await response
          .json()
          .catch(() => ({ message: `HTTP ${response.status}` }));
      } else {
        const text = await response.text();
        error = { message: text || `HTTP ${response.status}` };
      }

      // Обрабатываем 429 (Too Many Requests) как предупреждение, а не ошибку
      if (response.status === 429) {
        console.warn('Chat API rate limit:', error);
        const rateLimitError = new Error(
          (error as any).message ||
            (error as any).error ||
            'Слишком много запросов. Подождите немного.'
        );
        (rateLimitError as Error & { status: number }).status = 429;
        throw rateLimitError;
      }

      console.error('Chat API error:', {
        status: response.status,
        error: error,
        url: url,
      });
      throw new Error(
        (error as any).message ||
          (error as any).error ||
          `Error: ${response.status}`
      );
    }

    const data = await response.json();
    console.log('Chat API response data:', data);
    return data;
  }

  // Получить список чатов
  async getChats(page = 1, limit = 20): Promise<ChatListResponse> {
    const response = await this.request<{
      data: MarketplaceChat[];
      success: boolean;
    }>(`/?page=${page}&limit=${limit}`);
    // Если сервер возвращает обернутый ответ, разворачиваем его
    if (response.data && response.success) {
      // Backend возвращает просто массив чатов, нужно преобразовать в ожидаемый формат
      const chats = Array.isArray(response.data) ? response.data : [];
      return {
        chats: chats,
        total: chats.length,
        page: page,
        limit: limit,
      };
    }
    return {
      chats: [],
      total: 0,
      page: 1,
      limit: limit,
    };
  }

  // Получить сообщения
  async getMessages(params: GetMessagesParams): Promise<MessagesResponse> {
    const query = new URLSearchParams();
    if (params.listing_id)
      query.append('listing_id', params.listing_id.toString());
    if (params.chat_id) query.append('chat_id', params.chat_id.toString());
    if (params.page) query.append('page', params.page.toString());
    if (params.limit) query.append('limit', params.limit.toString());

    const response = await this.request<{
      data:
        | MarketplaceMessage[]
        | {
            messages: MarketplaceMessage[];
            total: number;
            page: number;
            limit: number;
          };
      success: boolean;
    }>(`/messages?${query.toString()}`, {
      signal: params.signal,
    });

    console.log('Chat API response data:', response);

    // Если сервер возвращает обернутый ответ, разворачиваем его
    if (response && response.success && response.data) {
      const data = response.data;

      // Проверяем формат ответа
      if (Array.isArray(data)) {
        // Старый формат - просто массив сообщений
        return {
          messages: data,
          total: data.length,
          page: params.page || 1,
          limit: params.limit || 20,
        };
      } else if (data && typeof data === 'object' && 'messages' in data) {
        // Новый формат со структурированным ответом
        console.log('New format detected, messages:', data.messages?.length);
        return {
          messages: Array.isArray(data.messages) ? data.messages : [],
          total: data.total || -1,
          page: data.page || params.page || 1,
          limit: data.limit || params.limit || 20,
        };
      }
    }
    return {
      messages: [],
      total: 0,
      page: 1,
      limit: params.limit || 20,
    };
  }

  // Отправить сообщение
  async sendMessage(payload: SendMessagePayload): Promise<MarketplaceMessage> {
    const response = await this.request<{
      data: MarketplaceMessage;
      success: boolean;
    }>('/messages', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    // Если сервер возвращает {data: ..., success: true}, извлекаем data
    return response.data || (response as unknown as MarketplaceMessage);
  }

  // Пометить сообщения как прочитанные
  async markMessagesAsRead(payload: MarkMessagesReadPayload): Promise<void> {
    await this.request<void>('/messages/read', {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  }

  // Архивировать чат
  async archiveChat(chatId: number): Promise<void> {
    await this.request<void>(`/${chatId}/archive`, {
      method: 'POST',
    });
  }

  // Получить количество непрочитанных
  async getUnreadCount(): Promise<number> {
    const response = await this.request<{
      data: { count: number };
      success: boolean;
    }>('/unread-count');
    // Если сервер возвращает обернутый ответ, разворачиваем его
    if (response.data && response.success) {
      return response.data.count || 0;
    }
    return (response as unknown as { count: number }).count || 0;
  }

  // Загрузить вложения для сообщения
  async uploadAttachments(
    messageId: number,
    files: File[],
    onProgress?: (progress: number) => void
  ): Promise<ChatAttachment[]> {
    const formData = new FormData();
    files.forEach((file) => formData.append('files', file));

    const xhr = new XMLHttpRequest();

    return new Promise(async (resolve, reject) => {
      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable && onProgress) {
          const progress = Math.round((e.loaded * 100) / e.total);
          onProgress(progress);
        }
      });

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText);
            resolve(response.data || response);
          } catch {
            reject(new Error('Invalid response format'));
          }
        } else {
          // Обрабатываем ошибки с телом ответа
          try {
            const errorResponse = JSON.parse(xhr.responseText);
            const error = new Error(
              errorResponse.message ||
                errorResponse.error ||
                `Upload failed: ${xhr.status}`
            );
            (error as Error & { status: number }).status = xhr.status;
            reject(error);
          } catch {
            const error = new Error(`Upload failed: ${xhr.status}`);
            (error as Error & { status: number }).status = xhr.status;
            reject(error);
          }
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error('Network error'));
      });

      xhr.withCredentials = true; // Включаем отправку куки

      const uploadUrl = `${this.baseUrl}/messages/${messageId}/attachments`;
      console.log('Uploading attachments to:', uploadUrl);

      xhr.open('POST', uploadUrl);

      // Добавляем JWT токен
      const accessToken = tokenManager.getAccessToken();
      if (accessToken) {
        xhr.setRequestHeader('Authorization', `Bearer ${accessToken}`);
      }

      // Добавляем CSRF токен
      const csrfToken = await this.getCsrfToken();
      xhr.setRequestHeader('X-CSRF-Token', csrfToken);

      xhr.send(formData);
    });
  }

  // Получить информацию о вложении
  async getAttachment(attachmentId: number): Promise<ChatAttachment> {
    const response = await this.request<{
      data: ChatAttachment;
      success: boolean;
    }>(`/attachments/${attachmentId}`);
    return response.data || (response as unknown as ChatAttachment);
  }

  // Загрузить файлы к сообщению
  async uploadFiles(
    messageId: number,
    files: File[],
    onProgress?: (progress: number) => void
  ): Promise<{ attachments: ChatAttachment[] }> {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append('files', file);
    });

    // Получаем JWT токен
    const accessToken = tokenManager.getAccessToken();
    const headers: Record<string, string> = {};

    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    const xhr = new XMLHttpRequest();

    return new Promise(async (resolve, reject) => {
      // Отслеживание прогресса
      if (onProgress) {
        xhr.upload.addEventListener('progress', (e) => {
          if (e.lengthComputable) {
            const progress = (e.loaded / e.total) * 100;
            onProgress(progress);
          }
        });
      }

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText);
            // Сервер возвращает { data: [...], success: true }
            // Преобразуем в ожидаемый формат { attachments: [...] }
            if (response.data && response.success) {
              resolve({ attachments: response.data });
            } else {
              // Fallback на случай, если формат ответа изменится
              resolve({
                attachments: response.attachments || response.data || [],
              });
            }
          } catch {
            reject(new Error('Invalid response'));
          }
        } else {
          reject(new Error(`Upload failed: ${xhr.status}`));
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error('Network error'));
      });

      xhr.open('POST', `${this.baseUrl}/messages/${messageId}/attachments`);

      // Устанавливаем заголовки
      Object.entries(headers).forEach(([key, value]) => {
        xhr.setRequestHeader(key, value);
      });

      // Добавляем CSRF токен
      const csrfToken = await this.getCsrfToken();
      xhr.setRequestHeader('X-CSRF-Token', csrfToken);

      xhr.withCredentials = true;
      xhr.send(formData);
    });
  }

  // Удалить вложение
  async deleteAttachment(attachmentId: number): Promise<void> {
    await this.request<void>(`/attachments/${attachmentId}`, {
      method: 'DELETE',
    });
  }

  // WebSocket соединение
  connectWebSocket(onMessage: (event: MessageEvent) => void): WebSocket | null {
    const wsUrl = configManager.getApiUrl().replace(/^http/, 'ws');

    // Получаем JWT токен из tokenManager
    const accessToken = tokenManager.getAccessToken();

    // Если нет токена, не подключаемся
    if (!accessToken) {
      console.warn(
        '[ChatService] No access token available, skipping WebSocket connection'
      );
      return null;
    }

    // Добавляем токен как query параметр для WebSocket
    const wsUrlWithAuth = `${wsUrl}/ws/chat?token=${accessToken}`;

    console.log(
      '[ChatService] Connecting WebSocket:',
      wsUrlWithAuth.replace(/token=.*/, 'token=***')
    );
    const ws = new WebSocket(wsUrlWithAuth);

    ws.onopen = () => {
      console.log('WebSocket connected');
      // Сбрасываем счетчик попыток при успешном подключении
      this.reconnectAttempts = 0;
    };

    ws.onmessage = onMessage;

    ws.onerror = (error) => {
      console.warn('WebSocket error:', error);
    };

    ws.onclose = (event) => {
      console.log('WebSocket disconnected', event.code, event.reason);

      // Если отключение из-за отсутствия аутентификации
      if (event.code === 1008 || event.reason === 'Unauthorized') {
        console.error('WebSocket: Authentication required');
        // Можно перенаправить на страницу входа или показать уведомление
        return;
      }

      // Если ошибка 429 (Too Many Requests), увеличиваем задержку
      if (event.code === 1006 && this.reconnectAttempts > 0) {
        // 1006 может быть из-за 429, увеличиваем задержку
        this.reconnectAttempts = Math.max(this.reconnectAttempts, 3);
      }

      // Проверяем максимальное количество попыток
      if (this.reconnectAttempts >= 10) {
        console.error('WebSocket: Max reconnection attempts reached');
        return;
      }

      // Автоматическое переподключение с экспоненциальной задержкой
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      this.reconnectAttempts++;

      console.log(
        `WebSocket: Will reconnect in ${delay / 1000}s (attempt ${this.reconnectAttempts})`
      );

      setTimeout(() => {
        console.log(
          `WebSocket: Reconnecting... (attempt ${this.reconnectAttempts})`
        );
        this.connectWebSocket(onMessage);
      }, delay);
    };

    return ws;
  }
}

export const chatService = new ChatService();
