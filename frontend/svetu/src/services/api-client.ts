import configManager from '@/config';
import { logger } from '@/utils/logger';

export interface ApiClientOptions extends RequestInit {
  // Использовать внутренний URL для серверных запросов
  internal?: boolean;
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
   * Получает CSRF токен через AuthService
   */
  private async getCsrfToken(): Promise<string | null> {
    if (typeof window === 'undefined') return null;
    try {
      const { AuthService } = await import('./auth');
      return await AuthService.getCsrfToken();
    } catch {
      return null;
    }
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

    // Добавляем CSRF токен для небезопасных методов
    const method = fetchOptions.method?.toUpperCase();
    if (method && ['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
      const csrfToken = await this.getCsrfToken();
      if (csrfToken) {
        headers.set('X-CSRF-Token', csrfToken);
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
    // ВРЕМЕННО ОТКЛЮЧЕНО: Блокировка радиусного поиска в API клиенте
    // Блокировка теперь работает на уровне backend
    if (
      endpoint.includes('/api/v1/gis/search/radius') &&
      typeof window !== 'undefined' &&
      (window.location.pathname.includes('/districts') ||
        localStorage.getItem('blockRadiusSearch') === 'true' ||
        (window as any).__BLOCK_RADIUS_SEARCH__)
    ) {
      logger.api.debug(
        'ℹ️ API CLIENT: Radius search request detected but allowing (backend will block)'
      );
    }

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

// Для обратной совместимости с существующим кодом
export abstract class ApiClientLegacy {
  protected baseURL: string;

  constructor(baseURL?: string) {
    // В development режиме и в браузере используем пустую строку для proxy
    if (
      !baseURL &&
      typeof window !== 'undefined' &&
      process.env.NODE_ENV === 'development'
    ) {
      this.baseURL = '';
    } else {
      this.baseURL = baseURL || configManager.getConfig().api.url;
    }
  }

  /**
   * Преобразует старый формат ответа в новый
   */
  private convertResponse<T>(response: ApiResponse<T>): {
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  } {
    if (response.data !== undefined) {
      return {
        success: true,
        data: response.data,
      };
    }

    return {
      success: false,
      data: null,
      error: response.error?.message,
      message: response.error?.message,
    };
  }

  protected async get<T>(
    endpoint: string,
    params?: Record<string, any>,
    options?: any
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const queryString = params
      ? '?' + new URLSearchParams(params).toString()
      : '';
    const response = await apiClient.get<T>(
      `${endpoint}${queryString}`,
      options
    );
    return this.convertResponse(response);
  }

  protected async post<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: any
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const response = await apiClient.post<T>(endpoint, data, options);
    return this.convertResponse(response);
  }

  protected async put<T, D = unknown>(
    endpoint: string,
    data?: D,
    options?: any
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const response = await apiClient.put<T>(endpoint, data, options);
    return this.convertResponse(response);
  }

  protected async delete<T>(
    endpoint: string,
    options?: any
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const response = await apiClient.delete<T>(endpoint, options);
    return this.convertResponse(response);
  }

  protected async upload<T>(
    endpoint: string,
    formData: FormData,
    options?: any
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const response = await apiClient.upload<T>(endpoint, formData, options);
    return this.convertResponse(response);
  }

  protected async request<T>(
    url: string,
    options?: RequestInit
  ): Promise<{
    success: boolean;
    data: T | null;
    error?: string;
    message?: string;
  }> {
    const endpoint = url.startsWith(this.baseURL)
      ? url.replace(this.baseURL, '')
      : url;
    const response = await apiClient.request<T>(endpoint, options);
    return this.convertResponse(response);
  }

  protected unwrap<T>(response: {
    success: boolean;
    data: T | null;
    error?: string;
  }): T {
    if (response.success && response.data !== null) {
      return response.data;
    }
    throw new Error(response.error || 'Request failed');
  }
}

// Экспортируем ApiClientLegacy как ApiClient для обратной совместимости
export { ApiClientLegacy as ApiClient };
