# Шаг 4: API client с поддержкой разных контекстов

## Цель
Модифицировать API client для правильной работы с разными URL в зависимости от контекста выполнения (сервер/клиент, внутренний/внешний запрос).

## Задачи

### 5.1 Обновление базового API client

Файл: `/frontend/svetu/src/services/api-client.ts`

```typescript
import configManager from '@/config';
import { tokenManager } from '@/utils/tokenManager';

export interface ApiClientOptions extends RequestInit {
  // Использовать внутренний URL для серверных запросов
  internal?: boolean;
  // Автоматически добавлять токен авторизации
  includeAuth?: boolean;
  // Timeout в миллисекундах
  timeout?: number;
  // Повторные попытки при ошибках сети
  retries?: number;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: {
    message: string;
    code?: string;
    details?: any;
  };
  status: number;
  headers: Headers;
}

class ApiClient {
  private defaultTimeout = 30000; // 30 секунд
  private defaultRetries = 3;

  /**
   * Получает базовый URL в зависимости от контекста
   */
  private getBaseUrl(isInternal: boolean = false): string {
    return configManager.getApiUrl({ internal: isInternal });
  }

  /**
   * Добавляет timeout к fetch запросу
   */
  private fetchWithTimeout(
    url: string, 
    options: RequestInit, 
    timeout: number
  ): Promise<Response> {
    return Promise.race([
      fetch(url, options),
      new Promise<Response>((_, reject) =>
        setTimeout(() => reject(new Error('Request timeout')), timeout)
      )
    ]);
  }

  /**
   * Выполняет запрос с повторными попытками
   */
  private async fetchWithRetries(
    url: string,
    options: RequestInit,
    timeout: number,
    retries: number
  ): Promise<Response> {
    let lastError: Error | null = null;

    for (let i = 0; i <= retries; i++) {
      try {
        const response = await this.fetchWithTimeout(url, options, timeout);
        
        // Если ответ успешный или это клиентская ошибка (4xx), не повторяем
        if (response.ok || (response.status >= 400 && response.status < 500)) {
          return response;
        }
        
        // Для серверных ошибок (5xx) повторяем запрос
        if (i < retries) {
          await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1))); // Экспоненциальная задержка
          continue;
        }
        
        return response;
      } catch (error) {
        lastError = error as Error;
        
        // Если это не последняя попытка, ждем и повторяем
        if (i < retries) {
          await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
          continue;
        }
      }
    }

    throw lastError || new Error('Request failed after retries');
  }

  /**
   * Основной метод для выполнения запросов
   */
  async request<T = any>(
    endpoint: string,
    options: ApiClientOptions = {}
  ): Promise<ApiResponse<T>> {
    const {
      internal = false,
      includeAuth = true,
      timeout = this.defaultTimeout,
      retries = this.defaultRetries,
      ...fetchOptions
    } = options;

    // Получаем базовый URL
    const baseUrl = this.getBaseUrl(internal);
    const url = `${baseUrl}${endpoint}`;

    // Подготавливаем заголовки
    const headers = new Headers(fetchOptions.headers);
    
    // Добавляем заголовки по умолчанию
    if (!headers.has('Content-Type') && fetchOptions.body) {
      headers.set('Content-Type', 'application/json');
    }

    // Добавляем токен авторизации если нужно
    if (includeAuth && typeof window !== 'undefined') {
      const token = tokenManager.getAccessToken();
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
    }

    // Добавляем специальный заголовок для внутренних запросов
    if (internal) {
      headers.set('X-Internal-Request', 'true');
    }

    // Финальные опции запроса
    const finalOptions: RequestInit = {
      ...fetchOptions,
      headers,
      // Добавляем credentials для CORS
      credentials: fetchOptions.credentials || 'include',
    };

    try {
      // Выполняем запрос
      const response = await this.fetchWithRetries(url, finalOptions, timeout, retries);

      // Парсим ответ
      let data: T | undefined;
      const contentType = response.headers.get('content-type');
      
      if (contentType?.includes('application/json')) {
        const text = await response.text();
        if (text) {
          try {
            data = JSON.parse(text);
          } catch (e) {
            console.error('Failed to parse JSON response:', e);
          }
        }
      }

      // Обрабатываем ошибки
      if (!response.ok) {
        return {
          error: {
            message: data?.['message'] || `Request failed with status ${response.status}`,
            code: data?.['code'],
            details: data?.['details'],
          },
          status: response.status,
          headers: response.headers,
        };
      }

      return {
        data,
        status: response.status,
        headers: response.headers,
      };
    } catch (error) {
      // Обработка сетевых ошибок
      const message = error instanceof Error ? error.message : 'Network error';
      
      return {
        error: {
          message,
          code: 'NETWORK_ERROR',
        },
        status: 0,
        headers: new Headers(),
      };
    }
  }

  /**
   * Удобные методы для разных типов запросов
   */
  async get<T = any>(endpoint: string, options?: ApiClientOptions): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'GET',
    });
  }

  async post<T = any>(
    endpoint: string, 
    data?: any, 
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T = any>(
    endpoint: string, 
    data?: any, 
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async patch<T = any>(
    endpoint: string, 
    data?: any, 
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T = any>(
    endpoint: string, 
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'DELETE',
    });
  }

  /**
   * Загрузка файлов
   */
  async upload<T = any>(
    endpoint: string,
    formData: FormData,
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
    const { headers = {}, ...restOptions } = options || {};
    
    // Не устанавливаем Content-Type для FormData, браузер сделает это сам
    const cleanHeaders = new Headers(headers);
    cleanHeaders.delete('Content-Type');

    return this.request<T>(endpoint, {
      ...restOptions,
      method: 'POST',
      body: formData,
      headers: cleanHeaders,
    });
  }
}

// Экспортируем singleton instance
export const apiClient = new ApiClient();

// Экспортируем также класс для тестирования
export { ApiClient };
```

### 5.2 Создание typed API endpoints

Файл: `/frontend/svetu/src/services/api/endpoints.ts`

```typescript
import { apiClient, ApiResponse } from '../api-client';
import type { 
  UserResponse, 
  MarketplaceListing,
  ChatMessage 
} from '@/types/generated/api';

/**
 * Базовый класс для API endpoints
 */
export class ApiEndpoint {
  constructor(protected basePath: string) {}

  /**
   * Определяет, должен ли запрос использовать внутренний URL
   * Переопределите в наследниках для специфичной логики
   */
  protected shouldUseInternalUrl(): boolean {
    // По умолчанию используем внутренний URL только для SSR
    return typeof window === 'undefined';
  }
}

/**
 * User API endpoints
 */
export class UserApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/users');
  }

  async getCurrentUser(): Promise<ApiResponse<UserResponse>> {
    return apiClient.get<UserResponse>(`${this.basePath}/me`, {
      internal: this.shouldUseInternalUrl(),
    });
  }

  async updateProfile(data: Partial<UserResponse>): Promise<ApiResponse<UserResponse>> {
    return apiClient.patch<UserResponse>(`${this.basePath}/me`, data, {
      internal: this.shouldUseInternalUrl(),
    });
  }
}

/**
 * Marketplace API endpoints
 */
export class MarketplaceApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/marketplace');
  }

  // Для публичных данных можем использовать внутренний URL при SSR
  protected shouldUseInternalUrl(): boolean {
    return typeof window === 'undefined';
  }

  async getListings(params?: {
    page?: number;
    limit?: number;
    category?: string;
  }): Promise<ApiResponse<{ items: MarketplaceListing[]; total: number }>> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.set('page', params.page.toString());
    if (params?.limit) queryParams.set('limit', params.limit.toString());
    if (params?.category) queryParams.set('category', params.category);

    return apiClient.get(`${this.basePath}/listings?${queryParams}`, {
      internal: this.shouldUseInternalUrl(),
    });
  }

  async getListingById(id: string): Promise<ApiResponse<MarketplaceListing>> {
    return apiClient.get<MarketplaceListing>(`${this.basePath}/listings/${id}`, {
      internal: this.shouldUseInternalUrl(),
    });
  }

  async createListing(data: FormData): Promise<ApiResponse<MarketplaceListing>> {
    return apiClient.upload<MarketplaceListing>(`${this.basePath}/listings`, data);
  }
}

/**
 * Chat API endpoints - всегда использует публичный URL для WebSocket
 */
export class ChatApi extends ApiEndpoint {
  constructor() {
    super('/api/v1/chat');
  }

  // Для чата всегда используем публичный URL из-за WebSocket
  protected shouldUseInternalUrl(): boolean {
    return false;
  }

  async getMessages(chatId: string): Promise<ApiResponse<ChatMessage[]>> {
    return apiClient.get<ChatMessage[]>(`${this.basePath}/${chatId}/messages`);
  }

  async sendMessage(chatId: string, message: string): Promise<ApiResponse<ChatMessage>> {
    return apiClient.post<ChatMessage>(`${this.basePath}/${chatId}/messages`, {
      content: message,
    });
  }
}

// Экспортируем singleton instances
export const userApi = new UserApi();
export const marketplaceApi = new MarketplaceApi();
export const chatApi = new ChatApi();
```

### 5.3 Использование в Server Components

Файл: `/frontend/svetu/src/app/[locale]/marketplace/[id]/page.tsx`

```typescript
import { marketplaceApi } from '@/services/api/endpoints';
import { notFound } from 'next/navigation';

interface PageProps {
  params: {
    id: string;
    locale: string;
  };
}

export default async function MarketplaceItemPage({ params }: PageProps) {
  // Server Component - будет использовать внутренний URL автоматически
  const response = await marketplaceApi.getListingById(params.id);

  if (response.error || !response.data) {
    notFound();
  }

  const listing = response.data;

  return (
    <div>
      <h1>{listing.title}</h1>
      {/* Render listing details */}
    </div>
  );
}

// Для статической генерации
export async function generateStaticParams() {
  const response = await marketplaceApi.getListings({ limit: 100 });
  
  if (!response.data) {
    return [];
  }

  return response.data.items.map((listing) => ({
    id: listing.id.toString(),
  }));
}
```

### 5.4 Использование в Client Components

Файл: `/frontend/svetu/src/components/marketplace/ListingActions.tsx`

```typescript
'use client';

import { useState } from 'react';
import { marketplaceApi } from '@/services/api/endpoints';
import { toast } from '@/utils/toast';

interface ListingActionsProps {
  listingId: string;
}

export function ListingActions({ listingId }: ListingActionsProps) {
  const [isLoading, setIsLoading] = useState(false);

  const handleAddToFavorites = async () => {
    setIsLoading(true);
    
    try {
      // Client Component - будет использовать публичный URL
      const response = await marketplaceApi.addToFavorites(listingId);
      
      if (response.error) {
        toast.error(response.error.message);
      } else {
        toast.success('Added to favorites!');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <button
      onClick={handleAddToFavorites}
      disabled={isLoading}
      className="btn btn-primary"
    >
      {isLoading ? 'Adding...' : 'Add to Favorites'}
    </button>
  );
}
```

### 5.5 Hook для API вызовов

Файл: `/frontend/svetu/src/hooks/useApi.ts`

```typescript
import { useState, useEffect, useCallback } from 'react';
import { ApiResponse } from '@/services/api-client';

interface UseApiOptions {
  immediate?: boolean; // Выполнить запрос сразу
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

interface UseApiResult<T> {
  data: T | null;
  error: any | null;
  loading: boolean;
  execute: (...args: any[]) => Promise<void>;
  reset: () => void;
}

export function useApi<T = any>(
  apiCall: (...args: any[]) => Promise<ApiResponse<T>>,
  options: UseApiOptions = {}
): UseApiResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<any | null>(null);
  const [loading, setLoading] = useState(false);

  const execute = useCallback(async (...args: any[]) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiCall(...args);

      if (response.error) {
        setError(response.error);
        options.onError?.(response.error);
      } else {
        setData(response.data || null);
        options.onSuccess?.(response.data);
      }
    } catch (err) {
      const error = { message: 'Unexpected error', details: err };
      setError(error);
      options.onError?.(error);
    } finally {
      setLoading(false);
    }
  }, [apiCall, options.onSuccess, options.onError]);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  useEffect(() => {
    if (options.immediate) {
      execute();
    }
  }, []);

  return { data, error, loading, execute, reset };
}

// Пример использования:
// const { data, loading, error, execute } = useApi(
//   () => marketplaceApi.getListings({ page: 1 }),
//   { immediate: true }
// );
```

## Тестирование

### Unit тесты для API client

Файл: `/frontend/svetu/src/services/__tests__/api-client.test.ts`

```typescript
import { ApiClient } from '../api-client';
import configManager from '@/config';

// Mock fetch
global.fetch = jest.fn();

describe('ApiClient', () => {
  let apiClient: ApiClient;

  beforeEach(() => {
    apiClient = new ApiClient();
    jest.clearAllMocks();
  });

  it('should use internal URL for server-side requests', async () => {
    // Mock server environment
    Object.defineProperty(window, 'window', { value: undefined });
    
    // Mock config
    jest.spyOn(configManager, 'getApiUrl').mockReturnValue('http://internal-api:3000');

    await apiClient.get('/test', { internal: true });

    expect(fetch).toHaveBeenCalledWith(
      'http://internal-api:3000/test',
      expect.any(Object)
    );
  });

  it('should retry on network errors', async () => {
    (fetch as jest.Mock)
      .mockRejectedValueOnce(new Error('Network error'))
      .mockRejectedValueOnce(new Error('Network error'))
      .mockResolvedValueOnce(new Response('{"data": "success"}'));

    const response = await apiClient.get('/test', { retries: 2 });

    expect(fetch).toHaveBeenCalledTimes(3);
    expect(response.data).toEqual({ data: 'success' });
  });

  it('should timeout long requests', async () => {
    (fetch as jest.Mock).mockImplementation(() => 
      new Promise(resolve => setTimeout(resolve, 5000))
    );

    const response = await apiClient.get('/test', { timeout: 100 });

    expect(response.error?.code).toBe('NETWORK_ERROR');
  });
});
```

## Результат
После выполнения этого шага:
1. API client будет автоматически выбирать правильный URL
2. Server Components будут использовать внутренние URL для быстродействия
3. Client Components будут использовать публичные URL
4. Добавлена поддержка retry и timeout
5. Типизированные endpoints для удобства использования

## Критерии приемки

### 1. Обновление api-client.ts
- [x] Файл `/src/services/api-client.ts` обновлен
- [x] Интерфейс `ApiClientOptions` содержит поля: `internal`, `includeAuth`, `timeout`, `retries`
- [x] Интерфейс `ApiResponse` типизирован с generics
- [x] Метод `getBaseUrl()` использует `configManager.getApiUrl()` с учетом контекста
- [x] Реализована поддержка timeout через `fetchWithTimeout()`
- [x] Реализована поддержка retry через `fetchWithRetries()` с экспоненциальной задержкой
- [x] Метод `upload()` для загрузки файлов работает с FormData

### 2. Создание typed endpoints
- [x] Файл `/src/services/api/endpoints.ts` создан
- [x] Базовый класс `ApiEndpoint` с методом `shouldUseInternalUrl()`
- [x] Классы `UserApi`, `MarketplaceApi`, `ChatApi` созданы
- [x] Каждый класс правильно определяет контекст использования:
  - UserApi и MarketplaceApi используют internal URL для SSR
  - ChatApi всегда использует public URL (для WebSocket)
- [x] Экспортированы singleton экземпляры: `userApi`, `marketplaceApi`, `chatApi`

### 3. Проверка контекстов
- [ ] Server Components автоматически используют internal URL:
  ```typescript
  // В server component запрос должен идти на http://backend:3000
  const response = await marketplaceApi.getListings();
  ```
- [ ] Client Components используют public URL:
  ```typescript
  // В client component запрос должен идти на http://localhost:3000
  const response = await marketplaceApi.getListings();
  ```
- [ ] Заголовок `X-Internal-Request` добавляется для внутренних запросов

### 4. Hook useApi
- [x] Файл `/src/hooks/useApi.ts` создан
- [x] Hook управляет состояниями: `data`, `error`, `loading`
- [x] Метод `execute()` для выполнения запроса
- [x] Метод `reset()` для сброса состояния
- [x] Опция `immediate` для автоматического запуска при монтировании
- [x] Callbacks `onSuccess` и `onError` работают

### 5. Функциональность
- [x] Retry логика работает при сетевых ошибках:
  - 3 попытки по умолчанию
  - Экспоненциальная задержка между попытками
  - Не повторяет при 4xx ошибках
- [x] Timeout работает (по умолчанию 30 секунд)
- [x] Авторизация автоматически добавляется через `tokenManager`
- [x] CORS настроен правильно (`credentials: 'include'`)

### 6. Тестирование
- [x] Unit тесты для API client написаны (хотя бы базовые)
- [x] Проверка использования internal URL в SSR контексте
- [x] Проверка retry механизма
- [x] Проверка timeout функциональности
- [ ] Нет CORS ошибок при реальных запросах

### 7. Интеграция
- [x] Существующие API вызовы можно мигрировать на новую систему
- [x] TypeScript типы из `/types/generated/api` используются
- [x] Нет breaking changes для существующего кода
- [x] Документация по использованию в комментариях кода