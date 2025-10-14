export interface ApiClientOptions extends RequestInit {
  // Timeout в миллисекундах
  timeout?: number;
  // Повторные попытки при ошибках сети
  retries?: number;
  // Управление Next.js кешированием (по умолчанию 'no-store' для client-side)
  // 'force-cache' - максимальное кеширование
  // 'no-store' - без кеширования (по умолчанию)
  // number - revalidate через N секунд
  nextCache?: RequestCache | number;
  // Теги для Next.js revalidation
  nextTags?: string[];
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
   *
   * Все запросы идут через Next.js BFF прокси /api/v2
   * который автоматически добавляет JWT токены из httpOnly cookies.
   *
   * Маппинг: /api/v2/* → backend /api/v1/*
   */
  private getBaseUrl(): string {
    // В browser это будет относительный путь: /api/v2
    // В SSR это будет полный URL: http://localhost:3001/api/v2
    if (typeof window === 'undefined') {
      // SSR: используем полный URL
      const frontendUrl =
        process.env.NEXT_PUBLIC_FRONTEND_URL || 'http://localhost:32341';
      return `${frontendUrl}/api/v2`;
    }

    // Browser: используем относительный путь
    return '/api/v2';
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
      ),
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
          await new Promise((resolve) => setTimeout(resolve, 1000 * (i + 1))); // Экспоненциальная задержка
          continue;
        }

        return response;
      } catch (error) {
        lastError = error as Error;

        // Если это не последняя попытка, ждем и повторяем
        if (i < retries) {
          await new Promise((resolve) => setTimeout(resolve, 1000 * (i + 1)));
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
      timeout = this.defaultTimeout,
      retries = this.defaultRetries,
      nextCache,
      nextTags,
      ...fetchOptions
    } = options;

    // Получаем базовый URL (всегда BFF прокси)
    const baseUrl = this.getBaseUrl();

    // Для BFF прокси нужно убрать /api/v1 из endpoint если он есть
    let finalEndpoint = endpoint;
    if (endpoint.startsWith('/api/v1/')) {
      // BFF прокси ожидает путь без /api/v1
      // Например: /api/v1/users/me → /users/me → BFF → backend /api/v1/users/me
      finalEndpoint = endpoint.replace('/api/v1/', '/');
    }

    const url = `${baseUrl}${finalEndpoint}`;

    // Подготавливаем заголовки
    const headers = new Headers(fetchOptions.headers);

    // Добавляем заголовки по умолчанию
    if (
      !headers.has('Content-Type') &&
      fetchOptions.body &&
      !(fetchOptions.body instanceof FormData)
    ) {
      headers.set('Content-Type', 'application/json');
    }

    // Добавляем заголовок Accept-Language
    if (!headers.has('Accept-Language') && typeof window !== 'undefined') {
      const locale = localStorage.getItem('locale') || 'sr';
      headers.set('Accept-Language', locale);
    }

    // Cookies (включая JWT токены) отправляются автоматически через credentials: 'include'
    // BFF proxy обрабатывает CSRF защиту через SameSite cookies

    // Финальные опции запроса
    const finalOptions: RequestInit = {
      ...fetchOptions,
      headers,
      // Добавляем credentials для CORS
      credentials: fetchOptions.credentials || 'include',
      // Настройки Next.js кеширования
      cache:
        typeof nextCache === 'number'
          ? undefined
          : nextCache || fetchOptions.cache || 'no-store',
      next: {
        ...(fetchOptions as any).next,
        revalidate: typeof nextCache === 'number' ? nextCache : undefined,
        tags: nextTags,
      },
    };

    try {
      // Выполняем запрос
      const response = await this.fetchWithRetries(
        url,
        finalOptions,
        timeout,
        retries
      );

      // Парсим ответ
      let data: T | undefined;
      const contentType = response.headers.get('content-type');

      if (response.status === 204) {
        // No content
        data = undefined;
      } else if (contentType?.includes('application/json')) {
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
            message:
              (data as any)?.['message'] ||
              (data as any)?.['error'] ||
              `Request failed with status ${response.status}`,
            code: (data as any)?.['code'] || `HTTP_${response.status}`,
            details: (data as any)?.['details'],
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
  async get<T = any>(
    endpoint: string,
    options?: ApiClientOptions
  ): Promise<ApiResponse<T>> {
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
      body:
        data instanceof FormData
          ? data
          : data
            ? JSON.stringify(data)
            : undefined,
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
      body:
        data instanceof FormData
          ? data
          : data
            ? JSON.stringify(data)
            : undefined,
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
      body:
        data instanceof FormData
          ? data
          : data
            ? JSON.stringify(data)
            : undefined,
    });
  }

  async delete<T = any>(
    endpoint: string,
    options?: ApiClientOptions & { data?: any }
  ): Promise<ApiResponse<T>> {
    const { data, ...restOptions } = options || {};
    return this.request<T>(endpoint, {
      ...restOptions,
      method: 'DELETE',
      body: data ? JSON.stringify(data) : undefined,
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
